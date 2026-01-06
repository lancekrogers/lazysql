package components

import "errors"

func (home *Home) Connect() error {
	if err := home.canSwitchToConnections(); err != nil {
		return err
	}
	return home.switchToConnections()
}

func (home *Home) SwitchConnection() error {
	if err := home.canSwitchToConnections(); err != nil {
		return err
	}
	return home.switchToConnections()
}

func (home *Home) ListConnections() error {
	if err := home.canSwitchToConnections(); err != nil {
		return err
	}
	return home.switchToConnections()
}

func (home *Home) Disconnect() error {
	if err := home.canSwitchToConnections(); err != nil {
		return err
	}
	if mainPages == nil {
		return errors.New("main pages not configured")
	}
	current, _ := mainPages.GetFrontPage()
	if current != "" && current != pageNameConnections {
		mainPages.RemovePage(current)
	}
	mainPages.SwitchToPage(pageNameConnections)
	if home.ConnectionIndicator != nil {
		home.ConnectionIndicator.SetDisconnected()
	}
	return nil
}

func (home *Home) switchToConnections() error {
	if mainPages == nil {
		return errors.New("main pages not configured")
	}
	mainPages.SwitchToPage(pageNameConnections)
	return nil
}

func (home *Home) canSwitchToConnections() error {
	if home == nil {
		return errors.New("home not configured")
	}
	if home.TabbedPane == nil {
		return errors.New("tabbed pane not configured")
	}
	tab := home.TabbedPane.GetCurrentTab()
	if tab == nil {
		return nil
	}
	table, ok := tab.Content.(*ResultsTable)
	if !ok {
		return nil
	}
	if table.GetIsEditing() || table.GetIsFiltering() || table.GetIsLoading() {
		return errors.New("finish editing or filtering before switching connections")
	}
	return nil
}
