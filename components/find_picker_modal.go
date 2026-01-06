package components

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
)

type FindPickerItem struct {
	Label     string
	Secondary string
	OnSelect  func()
}

type FindPickerConfig struct {
	Title       string
	Placeholder string
	Items       []FindPickerItem
	OnCancel    func()
}

type FindPicker struct {
	container     *tview.Flex
	input         *tview.InputField
	list          *tview.List
	items         []FindPickerItem
	filteredItems []FindPickerItem
	loading       bool
	onCancel      func()
	previousFocus tview.Primitive
}

func NewFindPicker() *FindPicker {
	input := tview.NewInputField()
	list := tview.NewList()
	container := tview.NewFlex().SetDirection(tview.FlexRow)

	input.SetLabel("Find: ")
	input.SetLabelColor(app.Styles.InverseTextColor)
	input.SetFieldStyle(tcell.StyleDefault.Background(app.Styles.PrimitiveBackgroundColor).Foreground(app.Styles.PrimaryTextColor))
	input.SetPlaceholderStyle(tcell.StyleDefault.Background(app.Styles.PrimitiveBackgroundColor).Foreground(app.Styles.InverseTextColor))
	input.SetFieldTextColor(app.Styles.PrimaryTextColor)

	list.SetMainTextColor(app.Styles.PrimaryTextColor)
	list.SetSecondaryTextColor(app.Styles.InverseTextColor)
	list.ShowSecondaryText(true)

	container.SetBorder(true)
	container.SetTitleAlign(tview.AlignLeft)
	container.AddItem(input, 1, 0, false)
	container.AddItem(list, 0, 1, true)

	picker := &FindPicker{
		container: container,
		input:     input,
		list:      list,
	}

	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			picker.selectCurrent()
		case tcell.KeyEscape:
			picker.dismiss(true)
		}
	})

	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event == nil {
			return nil
		}
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyUp, tcell.KeyTab:
			App.SetFocus(picker.list)
			return nil
		}
		return event
	})

	list.SetDoneFunc(func() {
		picker.dismiss(true)
	})

	return picker
}

func (p *FindPicker) Show(config FindPickerConfig) error {
	if mainPages == nil {
		return errors.New("main pages not configured")
	}

	p.onCancel = config.OnCancel
	p.loading = false
	p.items = append([]FindPickerItem(nil), config.Items...)
	p.filteredItems = append([]FindPickerItem(nil), config.Items...)

	p.container.SetTitle(fmt.Sprintf(" %s ", config.Title))
	p.input.SetPlaceholder(config.Placeholder)
	p.input.SetText("")
	p.input.SetChangedFunc(func(text string) {
		p.applyFilter(text)
	})

	p.refreshList()

	p.previousFocus = App.GetFocus()
	if mainPages.HasPage(pageNameFindPicker) {
		mainPages.RemovePage(pageNameFindPicker)
	}
	mainPages.AddPage(pageNameFindPicker, p.container, true, true)
	App.SetFocus(p.input)
	return nil
}

func (p *FindPicker) applyFilter(query string) {
	if p.loading {
		p.refreshList()
		return
	}
	p.filteredItems = filterFindItems(p.items, query)
	p.refreshList()
}

func (p *FindPicker) refreshList() {
	p.list.Clear()
	if p.loading {
		p.list.AddItem("Loading...", "", 0, nil)
		return
	}
	if len(p.filteredItems) == 0 {
		p.list.AddItem("No matches", "", 0, nil)
		return
	}
	for _, item := range p.filteredItems {
		item := item
		p.list.AddItem(item.Label, item.Secondary, 0, func() {
			p.selectItem(item)
		})
	}
}

func (p *FindPicker) selectCurrent() {
	if len(p.filteredItems) == 0 {
		return
	}
	index := p.list.GetCurrentItem()
	if index < 0 || index >= len(p.filteredItems) {
		index = 0
	}
	p.selectItem(p.filteredItems[index])
}

func (p *FindPicker) selectItem(item FindPickerItem) {
	p.dismiss(false)
	if item.OnSelect != nil {
		item.OnSelect()
	}
}

func (p *FindPicker) dismiss(callCancel bool) {
	if mainPages != nil {
		mainPages.RemovePage(pageNameFindPicker)
	}
	if p.previousFocus != nil {
		App.SetFocus(p.previousFocus)
	}
	if callCancel && p.onCancel != nil {
		p.onCancel()
	}
}

func (p *FindPicker) ShowLoading(config FindPickerConfig) error {
	if mainPages == nil {
		return errors.New("main pages not configured")
	}

	p.onCancel = config.OnCancel
	p.loading = true
	p.items = nil
	p.filteredItems = nil

	p.container.SetTitle(fmt.Sprintf(" %s ", config.Title))
	p.input.SetPlaceholder(config.Placeholder)
	p.input.SetText("")
	p.input.SetChangedFunc(func(text string) {
		p.applyFilter(text)
	})

	p.refreshList()

	p.previousFocus = App.GetFocus()
	if mainPages.HasPage(pageNameFindPicker) {
		mainPages.RemovePage(pageNameFindPicker)
	}
	mainPages.AddPage(pageNameFindPicker, p.container, true, true)
	App.SetFocus(p.input)
	return nil
}

func (p *FindPicker) SetItems(items []FindPickerItem) {
	p.loading = false
	p.items = append([]FindPickerItem(nil), items...)
	p.filteredItems = append([]FindPickerItem(nil), items...)
	p.applyFilter(p.input.GetText())
}

func filterFindItems(items []FindPickerItem, query string) []FindPickerItem {
	query = strings.TrimSpace(query)
	if query == "" {
		return append([]FindPickerItem(nil), items...)
	}

	type rankedItem struct {
		item FindPickerItem
		rank int
	}

	ranked := make([]rankedItem, 0, len(items))
	for _, item := range items {
		searchText := item.Label
		if item.Secondary != "" {
			searchText = searchText + " " + item.Secondary
		}
		rank := fuzzy.RankMatchNormalizedFold(query, searchText)
		if rank >= 0 {
			ranked = append(ranked, rankedItem{item: item, rank: rank})
		}
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].rank == ranked[j].rank {
			return ranked[i].item.Label < ranked[j].item.Label
		}
		return ranked[i].rank < ranked[j].rank
	})

	filtered := make([]FindPickerItem, 0, len(ranked))
	for _, entry := range ranked {
		filtered = append(filtered, entry.item)
	}
	return filtered
}
