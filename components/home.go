package components

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
	"github.com/lancekrogers/lazysql/commands"
	"github.com/lancekrogers/lazysql/drivers"
	"github.com/lancekrogers/lazysql/helpers/logger"
	"github.com/lancekrogers/lazysql/internal/history"
	"github.com/lancekrogers/lazysql/internal/ui"
	"github.com/lancekrogers/lazysql/internal/vim/buffer"
	"github.com/lancekrogers/lazysql/internal/vim/cmdline"
	cmdlinecommands "github.com/lancekrogers/lazysql/internal/vim/cmdline/commands"
	"github.com/lancekrogers/lazysql/internal/vim/leader"
	"github.com/lancekrogers/lazysql/internal/vim/modes"
	"github.com/lancekrogers/lazysql/internal/vim/namespace"
	"github.com/lancekrogers/lazysql/internal/vim/whichkey"
	"github.com/lancekrogers/lazysql/models"
)

type Home struct {
	*tview.Flex
	Tree                 *Tree
	TabbedPane           *TabbedPane
	LeftWrapper          *tview.Flex
	RightWrapper         *tview.Flex
	MainContent          *tview.Flex
	leftWrapperVisible   bool
	treePinned           bool
	HelpStatus           HelpStatus
	ModeManager          *modes.ModeManager
	ModeIndicator        *modes.ModeIndicator
	CommandLine          *cmdline.CommandLine
	StatusPages          *tview.Pages
	StatusLine           *StatusLine
	BufferManager        *buffer.Manager
	BufferController     *BufferController
	BufferNavigator      *buffer.Navigator
	BufferPicker         *BufferPicker
	ShellPane            *ui.ShellPane
	ShellHistoryModal    *ShellHistoryModal
	ContentPages         *tview.Pages
	StatusBar            *tview.Flex
	LeaderRegistry       *leader.Registry
	NamespaceRegistry    *namespace.Registry
	LeaderManager        *leader.Manager
	LeaderOverlay        *whichkey.WhichKeyOverlay
	HelpModal            *HelpModal
	QueryHistoryModal    *QueryHistoryModal
	DBDriver             drivers.Driver
	FocusedWrapper       string
	ListOfDBChanges      []models.DBDMLChange
	ConnectionIdentifier string
	ConnectionURL        string
	ReadOnly             bool
}

func NewHomePage(connection models.Connection, dbdriver drivers.Driver) *Home {
	tree := NewTree(connection.DBName, dbdriver)
	leftWrapper := tview.NewFlex()
	rightWrapper := tview.NewFlex()

	maincontent := tview.NewFlex()

	connectionIdentifier := connection.Name
	if connectionIdentifier == "" {
		parsedURL, err := url.Parse(connection.URL)
		if err == nil {
			connectionIdentifier = history.SanitizeFilename(parsedURL.Host + strings.ReplaceAll(parsedURL.Path, "/", "_"))
		} else {
			connectionIdentifier = "unnamed_or_invalid_url_connection"
		}
	}

	modeManager := modes.NewManager()
	modeIndicator := modes.NewModeIndicator(modeManager)
	commandLine := cmdline.NewCommandLine(app.App.Application, cmdline.NewCommandHistory(100))
	commandLine.SetContext(app.App.Context())
	bufferManager := buffer.NewManager()
	bufferNavigator := buffer.NewNavigator(bufferManager)
	bufferPicker := NewBufferPicker(bufferManager)
	shellPane := ui.NewShellPane()
	statusLine := NewStatusLine()
	leaderRegistry := leader.NewRegistry()
	leaderOverlay := whichkey.NewOverlay(app.App.Application)
	leaderOverlay.SetPosition(whichkey.PositionBottomRight)

	home := &Home{
		Flex:               tview.NewFlex().SetDirection(tview.FlexRow),
		Tree:               tree,
		LeftWrapper:        leftWrapper,
		RightWrapper:       rightWrapper,
		MainContent:        maincontent,
		leftWrapperVisible: true,
		treePinned:         true,
		HelpStatus:         NewHelpStatus(),
		ModeManager:        modeManager,
		ModeIndicator:      modeIndicator,
		CommandLine:        commandLine,
		StatusLine:         statusLine,
		BufferManager:      bufferManager,
		BufferNavigator:    bufferNavigator,
		BufferPicker:       bufferPicker,
		ShellPane:          shellPane,
		LeaderRegistry:     leaderRegistry,
		LeaderOverlay:      leaderOverlay,
		HelpModal:          NewHelpModal(),

		DBDriver:             dbdriver,
		ListOfDBChanges:      []models.DBDMLChange{},
		ConnectionIdentifier: connectionIdentifier,
		ConnectionURL:        connection.URL,
		ReadOnly:             connection.ReadOnly,
	}

	shellPane.SetOnExecute(home.executeShellQuery)

	shellHistoryModal := NewShellHistoryModal(connectionIdentifier, func(selectedQuery string) {
		if err := home.showShellPane(); err != nil {
			home.showStatusError(err.Error())
			return
		}
		if home.ShellPane != nil {
			home.ShellPane.SetInputText(selectedQuery)
			app.App.SetFocus(home.ShellPane.Input())
		}
	})
	home.ShellHistoryModal = shellHistoryModal

	leaderManager := leader.NewManager(leaderRegistry, leaderOverlay, leader.ConfigFromApp(app.App.Config()).Timeout, homeStatusReporter{home: home})
	home.LeaderManager = leaderManager

	bufferNamespace := namespace.NewBufferNamespace(bufferManager, bufferNavigator, bufferPicker)
	shellNamespace := namespace.NewShellNamespace(home)
	runNamespace := namespace.NewRunNamespace(home)
	describeNamespace := namespace.NewDescribeNamespace(home)
	namespaceRegistry, err := namespace.Initialize(leaderRegistry, bufferNamespace, shellNamespace, runNamespace, describeNamespace)
	if err != nil {
		logger.Error("Failed to initialize namespaces", map[string]any{"error": err})
	} else {
		home.NamespaceRegistry = namespaceRegistry
	}

	tabbedPane := NewTabbedPane()

	home.TabbedPane = tabbedPane

	qhm := NewQueryHistoryModal(connectionIdentifier, func(selectedQuery string) {
		home.createOrFocusEditorTab()

		currentTab := home.TabbedPane.GetCurrentTab()
		if currentTab != nil {
			table := currentTab.Content.(*ResultsTable)
			table.Editor.SetText(selectedQuery, true)
		}
	})

	home.QueryHistoryModal = qhm

	go home.subscribeToTreeChanges()

	leftWrapper.SetBorderColor(app.Styles.InverseTextColor)
	leftWrapper.AddItem(tree.Wrapper, 0, 1, true)

	if connection.ReadOnly {
		leftWrapper.SetTitle(" [READ-ONLY] ")
		leftWrapper.SetTitleColor(tcell.ColorLightBlue)
		leftWrapper.SetBorder(true)
	}

	rightWrapper.SetBorderColor(app.Styles.InverseTextColor)
	rightWrapper.SetBorder(true)
	rightWrapper.SetDirection(tview.FlexColumnCSS)
	rightWrapper.SetInputCapture(home.rightWrapperInputCapture)
	rightWrapper.AddItem(tabbedPane.HeaderContainer, 1, 0, false)
	rightWrapper.AddItem(tabbedPane.Pages, 0, 1, false)

	maincontent.AddItem(leftWrapper, 30, 1, false)
	maincontent.AddItem(rightWrapper, 0, 5, false)

	contentPages := tview.NewPages()
	contentPages.AddPage(pageNameHomeContent, maincontent, true, true)
	contentPages.AddPage(pageNameWhichKeyOverlay, leaderOverlay, true, true)
	home.ContentPages = contentPages

	statusBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	statusPages := tview.NewPages()
	statusPages.AddPage(pageNameStatusHelp, home.HelpStatus, true, true)
	statusPages.AddPage(pageNameStatusCommandLine, commandLine, true, false)
	statusPages.AddPage(pageNameStatusMessage, statusLine, true, false)
	commandLine.SetOnShow(func() {
		statusPages.SwitchToPage(pageNameStatusCommandLine)
	})
	commandLine.SetOnHide(func() {
		current, _ := statusPages.GetFrontPage()
		if current == pageNameStatusCommandLine || current == "" {
			statusPages.SwitchToPage(pageNameStatusHelp)
		}
	})
	home.StatusPages = statusPages

	statusBar.AddItem(statusPages, 0, 1, false)
	statusBar.AddItem(home.ModeIndicator, 14, 0, false)
	home.StatusBar = statusBar

	commandLine.SetOnError(func(err error) {
		if err != nil {
			home.showStatusError(err.Error())
		}
	})

	registry := cmdline.NewRegistry()
	registry.Register(&cmdlinecommands.WriteCommand{})
	registry.Register(&cmdlinecommands.EditCommand{})
	registry.Register(&cmdlinecommands.QuitCommand{})
	registry.Register(&cmdlinecommands.WriteQuitCommand{})
	registry.Register(&cmdlinecommands.QuitAllCommand{})
	registry.RegisterAlias("x", "wq")
	registry.RegisterAlias("exit", "qa")

	executor := cmdline.NewCommandExecutor(registry, &cmdline.CommandContext{
		Buffers: bufferManager,
		Status:  homeStatusReporter{home: home},
		App:     app.App,
	})
	commandLine.SetExecutor(executor)

	home.layoutMain()

	home.SetInputCapture(home.homeInputCapture)

	home.SetFocusFunc(func() {
		if home.FocusedWrapper == focusedWrapperLeft || home.FocusedWrapper == "" {
			home.focusLeftWrapper()
		} else {
			home.focusRightWrapper()
		}
	})

	mainPages.AddPage(connection.URL, home, true, false)
	return home
}

func (home *Home) subscribeToTreeChanges() {
	ch := home.Tree.Subscribe()

	for stateChange := range ch {
		switch stateChange.Key {
		case eventTreeSelectedTable:
			databaseName := home.Tree.GetSelectedDatabase()
			tableName := stateChange.Value.(string)

			home.showTable(databaseName, tableName)
		case eventTreeIsFiltering:
			isFiltering := stateChange.Value.(bool)
			if isFiltering {
				home.SetInputCapture(nil)
			} else {
				home.SetInputCapture(home.homeInputCapture)
			}
		case eventTreeSelectedFunction:
			home.createOrFocusEditorTab()
			currentTab := home.TabbedPane.GetCurrentTab()
			if currentTab != nil {
				table := currentTab.Content.(*ResultsTable)
				databaseName := home.Tree.GetSelectedDatabase()
				functionName := stateChange.Value.(string)
				functionDefinition, err := home.Tree.DBDriver.GetFunctionDefinition(databaseName, functionName)
				if err != nil {
					logger.Error(err.Error(), nil)
					continue
				}
				table.Editor.SetText(functionDefinition, false)
				App.ForceDraw()
			}
		case eventTreeSelectedProcedure:
			home.createOrFocusEditorTab()
			currentTab := home.TabbedPane.GetCurrentTab()
			if currentTab != nil {
				table := currentTab.Content.(*ResultsTable)
				databaseName := home.Tree.GetSelectedDatabase()
				procedureName := stateChange.Value.(string)
				procedureDefinition, err := home.Tree.DBDriver.GetProcedureDefinition(databaseName, procedureName)
				if err != nil {
					logger.Error(err.Error(), nil)
					continue
				}
				table.Editor.SetText(procedureDefinition, false)
				App.ForceDraw()
			}
		case eventTreeSelectedView:
			home.createOrFocusEditorTab()
			currentTab := home.TabbedPane.GetCurrentTab()
			if currentTab != nil {
				table := currentTab.Content.(*ResultsTable)
				databaseName := home.Tree.GetSelectedDatabase()
				viewName := stateChange.Value.(string)
				viewDefinition, err := home.Tree.DBDriver.GetViewDefinition(databaseName, viewName)
				if err != nil {
					logger.Error(err.Error(), nil)
					continue
				}
				table.Editor.SetText(viewDefinition, false)
				App.ForceDraw()
			}
		}
	}
}

func (home *Home) showTable(databaseName, tableName string) {
	if tableName == "" {
		return
	}
	tabReference := fmt.Sprintf("%s.%s", databaseName, tableName)

	tab := home.TabbedPane.GetTabByReference(tabReference)

	var table *ResultsTable

	if tab != nil {
		table = tab.Content.(*ResultsTable)
		home.TabbedPane.SwitchToTabByReference(tab.Reference)
	} else {
		table = NewResultsTable(&home.ListOfDBChanges, home.Tree, home.DBDriver, home.ConnectionIdentifier, home.ConnectionURL, home.ReadOnly).WithFilter()
		table.SetDatabaseName(databaseName)
		table.SetTableName(tableName)

		home.TabbedPane.AppendTab(tableName, table, tabReference)
	}

	results := table.FetchRecords(func() {
		home.focusLeftWrapper()
	})

	// Show sidebar if there is more then 1 row (row 0 are
	// the column names) and the sidebar is not disabled.
	if !app.App.Config().DisableSidebar && len(results) > 1 && !table.GetShowSidebar() {
		table.ShowSidebar(true)
	}

	if table.state.error == "" {
		if !home.treePinned && home.leftWrapperVisible {
			home.toggleLeftWrapper()
		}

		home.focusRightWrapper()
	}

	app.App.ForceDraw()
}

func (home *Home) focusRightWrapper() {
	home.Tree.RemoveHighlight()

	home.RightWrapper.SetBorderColor(app.Styles.PrimaryTextColor)
	home.LeftWrapper.SetBorderColor(app.Styles.InverseTextColor)
	home.TabbedPane.Highlight()
	tab := home.TabbedPane.GetCurrentTab()

	if tab != nil {
		home.focusTab(tab)
	}

	home.FocusedWrapper = focusedWrapperRight
}

func (home *Home) focusTab(tab *Tab) {
	if tab != nil {
		table := tab.Content.(*ResultsTable)
		table.HighlightAll()

		if table.GetIsFiltering() {
			go func() {
				if table.Filter != nil {
					app.App.SetFocus(table.Filter.Input)
					table.Filter.HighlightLocal()
				} else if table.Editor != nil {
					app.App.SetFocus(table.Editor)
					table.Editor.Highlight()
				}

				table.RemoveHighlightTable()
				app.App.Draw()
			}()
		} else {
			table.SetInputCapture(table.tableInputCapture)
			app.App.SetFocus(table)
		}

		if tab.Name == tabNameEditor {
			home.HelpStatus.SetStatusOnEditorView()
		} else {
			home.HelpStatus.SetStatusOnTableView()
		}
		home.showHelpStatus()
	}
}

func (home *Home) focusLeftWrapper() {
	home.Tree.Highlight()

	home.RightWrapper.SetBorderColor(app.Styles.InverseTextColor)
	home.LeftWrapper.SetBorderColor(app.Styles.PrimaryTextColor)

	tab := home.TabbedPane.GetCurrentTab()

	if tab != nil {
		table := tab.Content.(*ResultsTable)

		table.RemoveHighlightAll()

	}

	home.TabbedPane.SetBlur()

	app.App.SetFocus(home.Tree)

	home.FocusedWrapper = focusedWrapperLeft
}

func (home *Home) isCurrentTabFiltering() bool {
	tab := home.TabbedPane.GetCurrentTab()

	if tab != nil {
		table := tab.Content.(*ResultsTable)
		return table.GetIsFiltering()
	}

	return false
}

func (home *Home) rightWrapperInputCapture(event *tcell.EventKey) *tcell.EventKey {
	var tab *Tab

	command := app.Keymaps.Group(app.TableGroup).Resolve(event)

	switch command {
	case commands.TabPrev:

		tab := home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)
			if !table.GetIsEditing() && !table.GetIsFiltering() {
				home.TabbedPane.SwitchToPreviousTab()
				// home.focusTab(home.TabbedPane.SwitchToPreviousTab())
				return nil
			}

		}

		return event
	case commands.TabNext:
		tab := home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)
			if !table.GetIsEditing() && !table.GetIsFiltering() {
				home.TabbedPane.SwitchToNextTab()
				// home.focusTab(home.TabbedPane.SwitchToNextTab())
				return nil
			}
		}

		return event
	case commands.TabFirst:
		if home.isCurrentTabFiltering() {
			return event
		}

		home.TabbedPane.SwitchToFirstTab()
		// home.focusTab(home.TabbedPane.SwitchToFirstTab())
		return nil
	case commands.TabLast:
		if home.isCurrentTabFiltering() {
			return event
		}

		home.TabbedPane.SwitchToLastTab()
		// home.focusTab(home.TabbedPane.SwitchToLastTab())
		return nil
	case commands.TabClose:
		tab = home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)

			if !table.GetIsFiltering() && !table.GetIsEditing() && !table.GetIsLoading() {
				home.TabbedPane.RemoveCurrentTab()

				if home.TabbedPane.GetLength() == 0 {
					home.focusLeftWrapper()
					return nil
				}
			}
		}
	case commands.PagePrev:
		if home.isCurrentTabFiltering() {
			return event
		}

		tab = home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)

			if ((table.Menu != nil && table.Menu.GetSelectedOption() == 1) ||
				table.Menu == nil) && !table.Pagination.GetIsFirstPage() && !table.GetIsLoading() {
				table.Pagination.SetOffset(table.Pagination.GetOffset() - table.Pagination.GetLimit())
				table.FetchRecords(nil)
			}
		}

	case commands.PageNext:
		if home.isCurrentTabFiltering() {
			return event
		}

		tab = home.TabbedPane.GetCurrentTab()

		if tab != nil {
			table := tab.Content.(*ResultsTable)

			if ((table.Menu != nil && table.Menu.GetSelectedOption() == 1) ||
				table.Menu == nil) && !table.Pagination.GetIsLastPage() && !table.GetIsLoading() {
				table.Pagination.SetOffset(table.Pagination.GetOffset() + table.Pagination.GetLimit())
				table.FetchRecords(nil)
			}
		}
	}

	return event
}

func (home *Home) homeInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if home.LeaderManager != nil && home.shouldHandleLeader(event) {
		if home.LeaderManager.HandleEvent(event) {
			return nil
		}
	}

	tab := home.TabbedPane.GetCurrentTab()

	var table *ResultsTable

	if tab != nil {
		table = tab.Content.(*ResultsTable)
	}

	command := app.Keymaps.Group(app.HomeGroup).Resolve(event)

	switch command {
	case commands.MoveLeft:
		if table != nil && !table.GetIsEditing() && !table.GetIsFiltering() && home.FocusedWrapper == focusedWrapperRight {
			if !home.leftWrapperVisible {
				home.toggleLeftWrapper()
			}

			home.focusLeftWrapper()
		}
	case commands.MoveRight:
		if table != nil && !table.GetIsEditing() && !table.GetIsFiltering() && home.FocusedWrapper == focusedWrapperLeft {
			if home.leftWrapperVisible && !home.treePinned {
				home.toggleLeftWrapper()
			}

			home.focusRightWrapper()
		}
	case commands.SwitchToEditorView:
		home.createOrFocusEditorTab()
	case commands.SwitchToConnectionsView:
		if (table != nil && !table.GetIsEditing() && !table.GetIsFiltering() && !table.GetIsLoading()) || table == nil {
			mainPages.SwitchToPage(pageNameConnections)
		}
	case commands.Quit:
		if tab == nil || (!table.GetIsEditing() && !table.GetIsFiltering()) {
			app.App.Stop()
		}
	case commands.Save:
		if home.ReadOnly {
			errorModal := tview.NewModal().
				SetText("Cannot save changes: Connection is in read-only mode").
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(_ int, _ string) {
					mainPages.RemovePage(pageNameReadOnlyError)
				})
			mainPages.AddPage(pageNameReadOnlyError, errorModal, true, true)
			return event
		}
		if (len(home.ListOfDBChanges) > 0) && !table.GetIsEditing() {
			queryPreviewModal := NewQueryPreviewModal(&home.ListOfDBChanges, home.DBDriver, func() {
				for _, change := range home.ListOfDBChanges {
					queryString, err := home.DBDriver.DMLChangeToQueryString(change)
					if err != nil {
						logger.Error("Failed to convert DML change to query string", map[string]any{"error": err})
						continue
					}
					err = history.AddQueryToHistory(home.ConnectionIdentifier, queryString)
					if err != nil {
						logger.Error("Failed to add query to history", map[string]any{"error": err})
					}
				}
				home.ListOfDBChanges = []models.DBDMLChange{}
				table.FetchRecords(nil)
				home.Tree.ForceRemoveHighlight()
			})

			mainPages.AddPage(pageNameDMLPreview, queryPreviewModal, true, true)
		}
	case commands.HelpPopup:
		if table == nil || !table.GetIsEditing() {
			mainPages.AddPage(pageNameHelp, home.HelpModal, true, true)
		}
	case commands.SearchGlobal:
		if !home.leftWrapperVisible {
			home.toggleLeftWrapper()
		}

		if table != nil && !table.GetIsEditing() && !table.GetIsFiltering() && !table.GetIsLoading() && home.FocusedWrapper == focusedWrapperRight {
			home.focusLeftWrapper()
		}

		home.Tree.ForceRemoveHighlight()
		home.Tree.ClearSearch()
		app.App.SetFocus(home.Tree.Filter)
		home.Tree.SetIsFiltering(true)
	case commands.ToggleQueryHistory:
		if mainPages.HasPage(pageNameQueryHistory) {
			mainPages.SwitchToPage(pageNameQueryHistory)
		} else {
			mainPages.AddPage(pageNameQueryHistory, home.QueryHistoryModal, true, true)
		}

		home.QueryHistoryModal.queryHistoryComponent.LoadHistory(home.ConnectionIdentifier)
		return nil
	case commands.ToggleTree:
		home.toggleLeftWrapper()
		home.treePinned = home.leftWrapperVisible
		return nil
	}

	return event
}

func (home *Home) shouldHandleLeader(event *tcell.EventKey) bool {
	if event == nil {
		return false
	}

	focus := app.App.GetFocus()
	if focus == nil {
		return true
	}

	if focus == home.CommandLine {
		return false
	}

	if _, ok := focus.(*tview.InputField); ok {
		return false
	}

	if home.ModeManager != nil && home.ModeManager.Mode() == modes.ModeInsert {
		tab := home.TabbedPane.GetCurrentTab()
		if tab != nil {
			table, ok := tab.Content.(*ResultsTable)
			if ok && table.Editor != nil && focus == table.Editor {
				return false
			}
		}
	}

	return true
}

func (home *Home) createOrFocusEditorTab() {
	tab := home.TabbedPane.GetTabByName(tabNameEditor)

	if tab != nil {
		home.TabbedPane.SwitchToTabByName(tabNameEditor)
		table := tab.Content.(*ResultsTable)
		table.SetIsFiltering(true)
		if table.Editor != nil {
			table.Editor.EnableVim(home.ModeManager, home.LeaderManager, home.CommandLine)
			if home.BufferController == nil {
				home.BufferController = NewBufferController(home.BufferManager, table.Editor)
			}
		}
	} else {
		tableWithEditor := NewResultsTable(&home.ListOfDBChanges, home.Tree, home.DBDriver, home.ConnectionIdentifier, home.ConnectionURL, home.ReadOnly).WithEditor()
		home.TabbedPane.AppendTab(tabNameEditor, tableWithEditor, tabNameEditor)
		tableWithEditor.SetIsFiltering(true)
		if tableWithEditor.Editor != nil {
			tableWithEditor.Editor.EnableVim(home.ModeManager, home.LeaderManager, home.CommandLine)
			if home.BufferController == nil {
				home.BufferController = NewBufferController(home.BufferManager, tableWithEditor.Editor)
			}
		}
		home.TabbedPane.GetCurrentTab()
	}

	home.HelpStatus.SetStatusOnEditorView()
	home.showHelpStatus()
	home.focusRightWrapper()
	App.ForceDraw()
}

func (home *Home) showHelpStatus() {
	if home.StatusPages == nil {
		return
	}
	home.StatusPages.SwitchToPage(pageNameStatusHelp)
}

func (home *Home) showStatusInfo(message string) {
	if home.StatusLine == nil || home.StatusPages == nil {
		return
	}
	home.StatusLine.Info(message)
	home.StatusPages.SwitchToPage(pageNameStatusMessage)
}

func (home *Home) showStatusError(message string) {
	if home.StatusLine == nil || home.StatusPages == nil {
		return
	}
	home.StatusLine.Error(message)
	home.StatusPages.SwitchToPage(pageNameStatusMessage)
}

type homeStatusReporter struct {
	home *Home
}

func (r homeStatusReporter) Info(message string) {
	if r.home == nil {
		return
	}
	r.home.showStatusInfo(message)
}

func (r homeStatusReporter) Error(message string) {
	if r.home == nil {
		return
	}
	r.home.showStatusError(message)
}

func (home *Home) toggleLeftWrapper() {
	if home.leftWrapperVisible {
		home.MainContent.Clear()
		home.MainContent.AddItem(home.RightWrapper, 0, 5, false)
		home.leftWrapperVisible = false
		home.focusRightWrapper()
	} else {
		home.MainContent.Clear()
		home.MainContent.AddItem(home.LeftWrapper, 30, 1, false)
		home.MainContent.AddItem(home.RightWrapper, 0, 5, false)
		home.leftWrapperVisible = true
		home.focusLeftWrapper()
	}
	app.App.ForceDraw()
}

func (home *Home) layoutMain() {
	home.Clear()
	if home.ContentPages != nil {
		home.AddItem(home.ContentPages, 0, 1, false)
	}
	if home.ShellPane != nil && home.ShellPane.IsVisible() {
		home.AddItem(home.ShellPane, home.ShellPane.Height(), 0, false)
	}
	if home.StatusBar != nil {
		home.AddItem(home.StatusBar, 1, 0, false)
	}
}

func (home *Home) shellPane() (*ui.ShellPane, error) {
	if home == nil || home.ShellPane == nil {
		return nil, errors.New("shell pane not configured")
	}
	return home.ShellPane, nil
}

func (home *Home) isShellFocused() bool {
	if home == nil || home.ShellPane == nil {
		return false
	}
	focus := app.App.GetFocus()
	if focus == nil {
		return false
	}
	return focus == home.ShellPane || focus == home.ShellPane.Input() || focus == home.ShellPane.Output()
}

func (home *Home) setShellVisible(visible bool) error {
	pane, err := home.shellPane()
	if err != nil {
		return err
	}
	wasFocused := home.isShellFocused()
	if visible {
		pane.Show()
	} else {
		pane.Hide()
	}
	home.layoutMain()
	if !visible && wasFocused {
		home.focusRightWrapper()
	}
	app.App.ForceDraw()
	return nil
}

func (home *Home) showShellPane() error {
	return home.setShellVisible(true)
}

func (home *Home) hideShellPane() error {
	return home.setShellVisible(false)
}

func (home *Home) ToggleShell() error {
	pane, err := home.shellPane()
	if err != nil {
		return err
	}
	if pane.IsVisible() {
		return home.hideShellPane()
	}
	return home.showShellPane()
}

func (home *Home) CloseShell() error {
	return home.hideShellPane()
}

func (home *Home) FocusShellInput() error {
	if err := home.showShellPane(); err != nil {
		return err
	}
	if home.ShellPane != nil {
		app.App.SetFocus(home.ShellPane.Input())
	}
	return nil
}

func (home *Home) ClearShellOutput() error {
	pane, err := home.shellPane()
	if err != nil {
		return err
	}
	pane.ClearOutput()
	app.App.ForceDraw()
	return nil
}

func (home *Home) ResizeShell() error {
	pane, err := home.shellPane()
	if err != nil {
		return err
	}
	if !pane.IsVisible() {
		if err := home.showShellPane(); err != nil {
			return err
		}
	}
	heights := []int{10, 15, 25}
	current := pane.Height()
	next := heights[0]
	for i, height := range heights {
		if height == current {
			next = heights[(i+1)%len(heights)]
			break
		}
	}
	pane.SetHeight(next)
	home.layoutMain()
	app.App.ForceDraw()
	return nil
}

func (home *Home) ShowShellHistory() error {
	if home.ShellHistoryModal == nil {
		return errors.New("shell history modal not configured")
	}
	if mainPages == nil {
		return errors.New("pages not configured")
	}
	if mainPages.HasPage(pageNameShellHistory) {
		mainPages.SwitchToPage(pageNameShellHistory)
	} else {
		mainPages.AddPage(pageNameShellHistory, home.ShellHistoryModal, true, true)
	}
	home.ShellHistoryModal.LoadHistory(home.ConnectionIdentifier)
	app.App.SetFocus(home.ShellHistoryModal.GetPrimitive())
	return nil
}

func (home *Home) executeShellQuery(query string) {
	query = strings.TrimSpace(query)
	if query == "" {
		return
	}
	pane, err := home.shellPane()
	if err != nil {
		home.showStatusError(err.Error())
		return
	}
	if home.DBDriver == nil {
		home.appendShellError(pane, errors.New("database driver not configured"))
		return
	}

	if home.ReadOnly {
		if err := drivers.ValidateQueryForReadOnly(query); err != nil {
			home.appendShellError(pane, err)
			return
		}
	}

	if isSelectQuery(query) {
		rows, _, err := home.DBDriver.ExecuteQuery(query)
		if err != nil {
			home.appendShellError(pane, err)
			return
		}
		pane.AppendOutput(formatShellRows(rows))
	} else {
		result, err := home.DBDriver.ExecuteDMLStatement(query)
		if err != nil {
			home.appendShellError(pane, err)
			return
		}
		pane.AppendOutput(result)
	}

	if err := history.AddQueryToHistory(home.ConnectionIdentifier, query); err != nil {
		logger.Error("Failed to add shell query to history", map[string]any{"error": err, "query": query, "connection": home.ConnectionIdentifier})
	}
	app.App.ForceDraw()
}

func (home *Home) appendShellError(pane *ui.ShellPane, err error) {
	if err == nil {
		return
	}
	if pane != nil {
		pane.AppendOutput(fmt.Sprintf("Error: %s", err.Error()))
	}
	home.showStatusError(err.Error())
	app.App.ForceDraw()
}

func isSelectQuery(query string) bool {
	queryLower := strings.ToLower(strings.TrimSpace(query))
	return strings.HasPrefix(queryLower, "select") ||
		strings.HasPrefix(queryLower, "with") ||
		strings.HasPrefix(queryLower, "explain") ||
		strings.HasPrefix(queryLower, "show") ||
		strings.HasPrefix(queryLower, "describe") ||
		strings.HasPrefix(queryLower, "desc")
}

func formatShellRows(rows [][]string) string {
	if len(rows) == 0 {
		return "No results."
	}
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if i >= len(widths) {
				continue
			}
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var builder strings.Builder
	for rowIndex, row := range rows {
		for colIndex, cell := range row {
			if colIndex > 0 {
				builder.WriteString(" | ")
			}
			builder.WriteString(padRight(cell, widths[colIndex]))
		}
		if rowIndex == 0 {
			builder.WriteString("\n")
			for colIndex, width := range widths {
				if colIndex > 0 {
					builder.WriteString("-+-")
				}
				builder.WriteString(strings.Repeat("-", width))
			}
		}
		if rowIndex < len(rows)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func padRight(text string, width int) string {
	if width <= len(text) {
		return text
	}
	return text + strings.Repeat(" ", width-len(text))
}
