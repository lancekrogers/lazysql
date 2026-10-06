package components

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
	"github.com/lancekrogers/lazysql/commands"
	"github.com/lancekrogers/lazysql/helpers/logger"
	"github.com/lancekrogers/lazysql/internal/vim/modes"
	"github.com/lancekrogers/lazysql/models"
)

type SQLEditorState struct {
	isFocused bool
}

type SQLEditor struct {
	*tview.TextArea
	state         *SQLEditorState
	subscribers   []chan models.StateChange
	ConnectionURL string
	modeManager   *modes.ModeManager
	normalHandler *modes.NormalHandler
	insertHandler *modes.InsertHandler
}

func NewSQLEditor(connectionURL string) *SQLEditor {
	textarea := tview.NewTextArea()
	textarea.SetBorder(true)
	textarea.SetTitleAlign(tview.AlignLeft)
	textarea.SetPlaceholder("Enter your SQL query here...")
	sqlEditor := &SQLEditor{
		TextArea: textarea,
		state: &SQLEditorState{
			isFocused: false,
		},
		ConnectionURL: connectionURL,
	}
	sqlEditor.SetInputCapture(sqlEditor.handleInput)

	return sqlEditor
}

func (s *SQLEditor) EnableVim(modeManager *modes.ModeManager, leader modes.LeaderActivator, commandLine modes.CommandLineActivator) {
	if modeManager == nil {
		return
	}
	adapter := modes.NewTextAreaAdapter(s.TextArea)
	s.modeManager = modeManager
	s.normalHandler = modes.NewNormalHandler(modeManager, leader, commandLine, adapter)
	s.insertHandler = modes.NewInsertHandler(modeManager, adapter)
}

func (s *SQLEditor) Subscribe() chan models.StateChange {
	subscriber := make(chan models.StateChange)
	s.subscribers = append(s.subscribers, subscriber)
	return subscriber
}

func (s *SQLEditor) Publish(key string, message string) {
	for _, sub := range s.subscribers {
		sub <- models.StateChange{
			Key:   key,
			Value: message,
		}
	}
}

func (s *SQLEditor) GetIsFocused() bool {
	return s.state.isFocused
}

func (s *SQLEditor) SetIsFocused(isFocused bool) {
	s.state.isFocused = isFocused
}

func (s *SQLEditor) Highlight() {
	s.SetBorderColor(app.Styles.PrimaryTextColor)
	s.SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.PrimaryTextColor))
}

func (s *SQLEditor) SetBlur() {
	if s.modeManager != nil {
		s.modeManager.EnterNormal()
	}
	s.SetBorderColor(app.Styles.InverseTextColor)
	s.SetTextStyle(tcell.StyleDefault.Foreground(app.Styles.InverseTextColor))
}

func (s *SQLEditor) handleInput(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return nil
	}

	if s.modeManager != nil {
		if s.modeManager.Mode() == modes.ModeInsert {
			if s.insertHandler != nil && s.insertHandler.HandleKey(event) {
				return nil
			}
		} else {
			if s.normalHandler != nil && s.normalHandler.HandleKey(event) {
				return nil
			}
			if event.Key() == tcell.KeyRune {
				return nil
			}
		}
	}

	command := app.Keymaps.Group(app.EditorGroup).Resolve(event)

	switch command {
	case commands.Execute:
		s.Publish(eventSQLEditorQuery, s.GetText())
		return nil
	case commands.UnfocusEditor:
		s.Publish(eventSQLEditorEscape, "")
	case commands.OpenInExternalEditor:
		// THIS IS A LINUX-ONLY FEATURE (for now)
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			var newText string
			app.App.Suspend(func() {
				newText = openExternalEditor(s.GetText(), s.ConnectionURL)
			})
			s.SetText(newText, true)
		}
	}

	return event
}

// openExternalEditor opens the user's preferred editor to edit the query.
// It should be called within app.Suspend() to ensure the TUI is properly restored.
func openExternalEditor(currentText string, connectionURL string) string {
	tmpFile, err := os.CreateTemp("", "lazysql-*.sql")
	if err != nil {
		logger.Error("Failed to create temporary file", map[string]any{"error": err.Error()})
		return currentText
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			logger.Error("Failed to remove temporary file", map[string]any{"error": err.Error()})
		}
	}()

	path := tmpFile.Name()
	content := []byte(currentText)

	if _, err := tmpFile.Write(content); err != nil {
		logger.Error("Failed to write to temporary file", map[string]any{"error": err.Error()})
		err := tmpFile.Close()
		if err != nil {
			logger.Error("Failed to close temporary file", map[string]any{"error": err.Error()})
		}
		return currentText
	}

	if err := tmpFile.Close(); err != nil {
		logger.Error("Failed to close temporary file", map[string]any{"error": err.Error()})
		return currentText
	}

	if connectionURL != "" {
		err := os.Setenv("LAZYSQL_CONNECTION_URL", connectionURL)
		if err != nil {
			logger.Error("Failed to set environment variable", map[string]any{"error": err.Error()})
			return currentText
		}
		// Defer unsetting the environment variable to ensure it's cleaned up
		defer func() {
			if err := os.Unsetenv("LAZYSQL_CONNECTION_URL"); err != nil {
				logger.Error("Failed to unset environment variable", map[string]any{"error": err.Error()})
			}
		}()
	}

	editor := editorCommand(true)

	// #nosec G204 -- user-configured editor command
	cmd := exec.Command(editor[0], append(editor[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.Error("Error executing command", map[string]any{"error": err.Error(), "command": cmd.String()})
	}

	// #nosec G304 -- reading temporary file created by this function
	updatedContent, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Failed to read from temporary file", map[string]any{"error": err.Error()})
		return currentText
	}

	return string(updatedContent)
}

func editorCommand(preferSQL bool) []string {
	editor := ""
	if preferSQL {
		editor = os.Getenv("SQL_EDITOR")
	}
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return []string{"vi"}
	}
	return parts
}

// openCellInExternalEditor opens the user's preferred editor to edit a cell value.
// It should be called within app.Suspend() to ensure the TUI is properly restored.
func openCellInExternalEditor(currentText string) string {
	tmpFile, err := os.CreateTemp("", "lazysql-cell-*.txt")
	if err != nil {
		logger.Error("Failed to create temporary file", map[string]any{"error": err.Error()})
		return currentText
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			logger.Error("Failed to remove temporary file", map[string]any{"error": err.Error()})
		}
	}()

	if _, err := tmpFile.WriteString(currentText); err != nil {
		logger.Error("Failed to write to temporary file", map[string]any{"error": err.Error()})
		if closeErr := tmpFile.Close(); closeErr != nil {
			logger.Error("Failed to close temporary file", map[string]any{"error": closeErr.Error()})
		}
		return currentText
	}
	if err := tmpFile.Close(); err != nil {
		logger.Error("Failed to close temporary file", map[string]any{"error": err.Error()})
		return currentText
	}

	editor := editorCommand(false)
	// #nosec G204 -- user-configured editor command
	cmd := exec.Command(editor[0], append(editor[1:], tmpFile.Name())...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.Error("Error executing editor", map[string]any{"error": err.Error(), "command": cmd.String()})
		return currentText
	}

	// #nosec G304 -- reading temporary file created by this function
	updatedContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		logger.Error("Failed to read from temporary file", map[string]any{"error": err.Error()})
		return currentText
	}

	return string(updatedContent)
}
