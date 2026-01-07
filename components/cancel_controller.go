package components

import "errors"

func (home *Home) CancelQuery() error {
	table, err := home.currentResultsTable()
	if err != nil {
		return err
	}
	if table == nil {
		return nil
	}
	if !table.GetIsLoading() {
		return nil
	}
	table.SetLoading(false)
	home.showStatusInfo("Canceled running query")
	return nil
}

func (home *Home) ClearResults() error {
	table, err := home.currentResultsTable()
	if err != nil {
		return err
	}
	if table == nil {
		return nil
	}
	table.state.error = ""
	if table.Page != nil {
		table.Page.HidePage(pageNameTableError)
	}
	table.SetResultsInfo("")
	table.SetRecords([][]string{})
	return nil
}

func (home *Home) currentResultsTable() (*ResultsTable, error) {
	if home == nil || home.TabbedPane == nil {
		return nil, errors.New("tabbed pane not configured")
	}
	tab := home.TabbedPane.GetCurrentTab()
	if tab == nil {
		return nil, nil
	}
	table, ok := tab.Content.(*ResultsTable)
	if !ok {
		return nil, nil
	}
	return table, nil
}
