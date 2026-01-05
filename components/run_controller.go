package components

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lancekrogers/lazysql/app"
	"github.com/lancekrogers/lazysql/drivers"
)

func (home *Home) RunCurrentQuery() error {
	table, editor, err := home.editorContext()
	if err != nil {
		return err
	}
	fromRow, fromCol, _, _ := editor.GetCursor()
	query := extractQueryAtCursor(editor.GetText(), fromRow, fromCol)
	if strings.TrimSpace(query) == "" {
		return errors.New("no query found at cursor")
	}
	return home.runQueryWithTiming(table, query)
}

func (home *Home) RunAllQueries() error {
	table, editor, err := home.editorContext()
	if err != nil {
		return err
	}
	queries := splitQueries(editor.GetText())
	if len(queries) == 0 {
		return errors.New("no queries found")
	}
	return home.runQueriesWithTiming(table, queries)
}

func (home *Home) RunSelectedLines() error {
	table, editor, err := home.editorContext()
	if err != nil {
		return err
	}
	if !editor.HasSelection() {
		return errors.New("no selection to run")
	}
	selection, _, _ := editor.GetSelection()
	query := strings.TrimSpace(selection)
	if query == "" {
		return errors.New("selected query is empty")
	}
	return home.runQueryWithTiming(table, query)
}

func (home *Home) RunWithParams() error {
	table, editor, err := home.editorContext()
	if err != nil {
		return err
	}
	fromRow, fromCol, _, _ := editor.GetCursor()
	query := extractQueryAtCursor(editor.GetText(), fromRow, fromCol)
	if strings.TrimSpace(query) == "" {
		return errors.New("no query found at cursor")
	}
	if mainPages == nil {
		return errors.New("main pages not configured")
	}

	modal := NewRunParamsModal(func(params string) {
		substituted, err := substituteQueryParams(query, params, home.DBDriver)
		if err != nil {
			home.showStatusError(err.Error())
			return
		}
		_ = home.runQueryWithTiming(table, substituted)
	})

	if mainPages.HasPage(pageNameRunParams) {
		mainPages.RemovePage(pageNameRunParams)
	}
	mainPages.AddPage(pageNameRunParams, modal, true, true)
	app.App.SetFocus(modal.GetPrimitive())
	return nil
}

func (home *Home) ShowQueryHistory() error {
	if home.QueryHistoryModal == nil {
		return errors.New("query history modal not configured")
	}
	if mainPages == nil {
		return errors.New("main pages not configured")
	}
	if mainPages.HasPage(pageNameQueryHistory) {
		mainPages.SwitchToPage(pageNameQueryHistory)
	} else {
		mainPages.AddPage(pageNameQueryHistory, home.QueryHistoryModal, true, true)
	}
	home.QueryHistoryModal.queryHistoryComponent.LoadHistory(home.ConnectionIdentifier)
	app.App.SetFocus(home.QueryHistoryModal.GetPrimitive())
	return nil
}

func (home *Home) editorContext() (*ResultsTable, *SQLEditor, error) {
	if home == nil || home.TabbedPane == nil {
		return nil, nil, errors.New("editor not available")
	}
	tab := home.TabbedPane.GetTabByName(tabNameEditor)
	if tab == nil {
		home.createOrFocusEditorTab()
		tab = home.TabbedPane.GetTabByName(tabNameEditor)
	}
	if tab == nil {
		return nil, nil, errors.New("editor tab not available")
	}
	table, ok := tab.Content.(*ResultsTable)
	if !ok || table == nil {
		return nil, nil, errors.New("editor tab not initialized")
	}
	if table.Editor == nil {
		return nil, nil, errors.New("editor not initialized")
	}
	return table, table.Editor, nil
}

func (home *Home) runQueryWithTiming(table *ResultsTable, query string) error {
	start := time.Now()
	err := table.ExecuteEditorQuery(query)
	elapsed := time.Since(start)
	if err != nil {
		home.showStatusError(err.Error())
		return err
	}
	home.showStatusInfo(fmt.Sprintf("Query executed in %s", formatDuration(elapsed)))
	return nil
}

func (home *Home) runQueriesWithTiming(table *ResultsTable, queries []string) error {
	start := time.Now()
	for _, query := range queries {
		if err := table.ExecuteEditorQuery(query); err != nil {
			home.showStatusError(err.Error())
			return err
		}
	}
	home.showStatusInfo(fmt.Sprintf("Executed %d queries in %s", len(queries), formatDuration(time.Since(start))))
	return nil
}

func extractQueryAtCursor(text string, row int, col int) string {
	if text == "" {
		return ""
	}
	cursor := cursorIndex(text, row, col)
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(text) {
		cursor = len(text)
	}
	start := strings.LastIndex(text[:cursor], ";")
	if start == -1 {
		start = 0
	} else {
		start++
	}
	end := strings.Index(text[cursor:], ";")
	if end == -1 {
		end = len(text)
	} else {
		end = cursor + end
	}
	return strings.TrimSpace(text[start:end])
}

func cursorIndex(text string, row int, col int) int {
	if row < 0 {
		row = 0
	}
	if col < 0 {
		col = 0
	}
	lines := strings.Split(text, "\n")
	if row >= len(lines) {
		return len(text)
	}
	index := 0
	for i := 0; i < row; i++ {
		index += len(lines[i]) + 1
	}
	if col > len(lines[row]) {
		col = len(lines[row])
	}
	return index + col
}

func splitQueries(text string) []string {
	parts := strings.Split(text, ";")
	queries := make([]string, 0, len(parts))
	for _, part := range parts {
		query := strings.TrimSpace(part)
		if query == "" {
			continue
		}
		queries = append(queries, query)
	}
	return queries
}

func substituteQueryParams(query string, paramsInput string, driver drivers.Driver) (string, error) {
	if driver == nil {
		return "", errors.New("database driver not configured")
	}
	placeholderCount := countPlaceholders(query, driver)
	if placeholderCount == 0 {
		return "", errors.New("query has no placeholders")
	}
	args := parseParamValues(paramsInput)
	if len(args) == 0 {
		return "", errors.New("no parameters provided")
	}
	if len(args) != placeholderCount {
		return "", fmt.Errorf("expected %d parameters, got %d", placeholderCount, len(args))
	}

	placeholder := driver.FormatPlaceholder(1)
	if placeholder == "?" {
		for _, arg := range args {
			query = strings.Replace(query, "?", driver.FormatArgForQueryString(arg), 1)
		}
		return query, nil
	}

	for i := placeholderCount; i >= 1; i-- {
		token := driver.FormatPlaceholder(i)
		query = strings.ReplaceAll(query, token, driver.FormatArgForQueryString(args[i-1]))
	}
	return query, nil
}

func countPlaceholders(query string, driver drivers.Driver) int {
	placeholder := driver.FormatPlaceholder(1)
	if placeholder == "?" {
		return strings.Count(query, "?")
	}

	re := regexp.MustCompile(`\$(\d+)`)
	matches := re.FindAllStringSubmatch(query, -1)
	maxIndex := 0
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		index, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

func parseParamValues(input string) []any {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil
	}
	parts := strings.Split(trimmed, ",")
	args := make([]any, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			args = append(args, "")
			continue
		}
		upper := strings.ToUpper(value)
		if upper == "NULL" || upper == "DEFAULT" {
			args = append(args, upper)
			continue
		}
		if intValue, err := strconv.Atoi(value); err == nil {
			args = append(args, intValue)
			continue
		}
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			args = append(args, floatValue)
			continue
		}
		args = append(args, value)
	}
	return args
}

func formatDuration(duration time.Duration) string {
	rounded := duration.Round(time.Millisecond)
	if rounded < time.Second {
		return fmt.Sprintf("%dms", rounded.Milliseconds())
	}
	return rounded.String()
}
