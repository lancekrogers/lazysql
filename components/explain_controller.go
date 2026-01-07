package components

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func (home *Home) ExplainCurrentQuery() error {
	table, editor, err := home.editorContext()
	if err != nil {
		return err
	}
	fromRow, fromCol, _, _ := editor.GetCursor()
	query := extractQueryAtCursor(editor.GetText(), fromRow, fromCol)
	if strings.TrimSpace(query) == "" {
		return errors.New("no query found at cursor")
	}
	return home.runExplain(table, "Explain", "EXPLAIN "+query)
}

func (home *Home) ExplainAnalyzeCurrentQuery() error {
	table, editor, err := home.editorContext()
	if err != nil {
		return err
	}
	fromRow, fromCol, _, _ := editor.GetCursor()
	query := extractQueryAtCursor(editor.GetText(), fromRow, fromCol)
	if strings.TrimSpace(query) == "" {
		return errors.New("no query found at cursor")
	}

	analyzeQuery := "EXPLAIN ANALYZE " + query
	start := time.Now()
	if err := table.ExecuteEditorQuery(analyzeQuery); err == nil {
		home.showStatusInfo(fmt.Sprintf("Explain analyze executed in %s", formatDuration(time.Since(start))))
		return nil
	}

	fallbackStart := time.Now()
	if err := table.ExecuteEditorQuery("EXPLAIN " + query); err != nil {
		home.showStatusError(err.Error())
		return err
	}
	home.showStatusInfo(fmt.Sprintf("Explain analyze not supported; ran EXPLAIN in %s", formatDuration(time.Since(fallbackStart))))
	return nil
}

func (home *Home) runExplain(table *ResultsTable, label string, query string) error {
	if table == nil {
		return errors.New("results table not configured")
	}
	start := time.Now()
	if err := table.ExecuteEditorQuery(query); err != nil {
		home.showStatusError(err.Error())
		return err
	}
	home.showStatusInfo(fmt.Sprintf("%s executed in %s", label, formatDuration(time.Since(start))))
	return nil
}
