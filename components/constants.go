package components

import (
	"github.com/gdamore/tcell/v2"

	"github.com/lancekrogers/lazysql/app"
)

var App = app.App

// Pages
const (
	// General
	pageNameHelp              string = "Help"
	pageNameConfirmation      string = "Confirmation"
	pageNameConnections       string = "Connections"
	pageNameDMLPreview        string = "DMLPreview"
	pageNameErrorModal        string = "ErrorModal"
	pageNameReadOnlyError     string = "readOnlyError"
	pageNameStatusHelp        string = "StatusHelp"
	pageNameStatusCommandLine string = "StatusCommandLine"
	pageNameStatusMessage     string = "StatusMessage"
	pageNameBufferPicker      string = "BufferPicker"
	pageNameHomeContent       string = "HomeContent"
	pageNameWhichKeyOverlay   string = "WhichKeyOverlay"
	pageNameDescribePicker    string = "DescribePicker"
	pageNameFindPicker        string = "FindPicker"

	// Results table
	pageNameTable                  string = "Table"
	pageNameTableError             string = "TableError"
	pageNameTableLoading           string = "TableLoading"
	pageNameTableEditorTable       string = "TableEditorTable"
	pageNameTableEditorResultsInfo string = "TableEditorResultsInfo"
	pageNameTableEditCell          string = "TableEditCell"
	pageNameQueryPreviewError      string = "QueryPreviewError"
	pageNameJSONViewer                    = "json_viewer"

	// Sidebar
	pageNameSidebar string = "Sidebar"

	// Connections
	pageNameConnectionSelection string = "ConnectionSelection"
	pageNameConnectionForm      string = "ConnectionForm"

	// SetValueList
	pageNameSetValue string = "SetValue"

	// Query History
	pageNameQueryHistory     string = "QueryHistoryModal"
	pageNameShellHistory     string = "ShellHistoryModal"
	pageNameSaveQuery        string = "SaveQueryModal"
	pageNameSavedQueryDelete string = "SavedQueryDeleteModal"
	pageNameRunParams        string = "RunParamsModal"
)

// Tabs
const (
	tabNameEditor string = "Editor"

	savedQueryTabReference   string = "saved_queries"
	queryHistoryTabReference string = "query_history"
)

// Events
const (
	eventSidebarEditing       string = "EditingSidebar"
	eventSidebarUnfocusing    string = "UnfocusingSidebar"
	eventSidebarToggling      string = "TogglingSidebar"
	eventSidebarCommitEditing string = "CommitEditingSidebar"
	eventSidebarError         string = "ErrorSidebar"

	eventSQLEditorQuery  string = "Query"
	eventSQLEditorEscape string = "Escape"

	eventResultsTableFiltering string = "FilteringResultsTable"

	eventTreeSelectedDatabase  string = "SelectedDatabase"
	eventTreeSelectedTable     string = "SelectedTable"
	eventTreeSelectedFunction  string = "SelectedFunction"
	eventTreeSelectedProcedure string = "SelectedProcedure"
	eventTreeSelectedView      string = "SelectedView"
	eventTreeIsFiltering       string = "IsFiltering"
)

// Results table menu items
const (
	menuRecords     string = "Records"
	menuColumns     string = "Columns"
	menuConstraints string = "Constraints"
	menuForeignKeys string = "Foreign Keys"
	menuIndexes     string = "Indexes"
)

// Actions
const (
	actionNewConnection  string = "NewConnection"
	actionEditConnection string = "EditConnection"
)

// Misc (until i find a better name)
const (
	focusedWrapperLeft  string = "left"
	focusedWrapperRight string = "right"

	colorTableChange = tcell.ColorOrange
	colorTableInsert = tcell.ColorDarkGreen
	colorTableDelete = tcell.ColorRed
)
