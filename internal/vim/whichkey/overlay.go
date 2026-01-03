package whichkey

import (
	"sort"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Position uint8

const (
	PositionCenter Position = iota
	PositionBottomRight
)

type KeyHint struct {
	Key         string
	Description string
	IsGroup     bool
}

type rect struct {
	x      int
	y      int
	width  int
	height int
}

type WhichKeyOverlay struct {
	*tview.Box
	app           *tview.Application
	content       []KeyHint
	visible       bool
	position      Position
	title         string
	styles        Styles
	containerRect rect
	tree          *KeyTree
	leaderPrefix  rune
	focusSetter   func(p tview.Primitive)
	returnFocus   tview.Primitive
}

func NewOverlay(app *tview.Application) *WhichKeyOverlay {
	box := tview.NewBox()
	box.SetBorder(true)
	box.SetTitleAlign(tview.AlignLeft)

	return &WhichKeyOverlay{
		Box:          box,
		app:          app,
		position:     PositionCenter,
		styles:       DefaultStyles,
		leaderPrefix: '\\',
	}
}

func (w *WhichKeyOverlay) SetContent(content []KeyHint) {
	w.content = content
}

func (w *WhichKeyOverlay) SetKeyTree(tree *KeyTree) {
	w.tree = tree
	w.syncFromTree()
}

func (w *WhichKeyOverlay) SetLeaderPrefix(prefix rune) {
	w.leaderPrefix = prefix
}

func (w *WhichKeyOverlay) SetPosition(position Position) {
	w.position = position
}

func (w *WhichKeyOverlay) SetTitle(title string) {
	w.title = title
}

func (w *WhichKeyOverlay) SetStyles(styles Styles) {
	w.styles = styles
}

func (w *WhichKeyOverlay) Show() {
	w.visible = true
	if w.tree != nil {
		w.syncFromTree()
	}
	if w.focusSetter != nil {
		w.focusSetter(w)
	}
	if w.app != nil {
		w.app.Draw()
	}
}

func (w *WhichKeyOverlay) Hide() {
	w.visible = false
	if w.focusSetter != nil && w.returnFocus != nil {
		w.focusSetter(w.returnFocus)
	}
	if w.app != nil {
		w.app.Draw()
	}
}

func (w *WhichKeyOverlay) IsVisible() bool {
	return w.visible
}

func (w *WhichKeyOverlay) SetRect(x, y, width, height int) {
	w.containerRect = rect{x: x, y: y, width: width, height: height}
	w.Box.SetRect(x, y, width, height)
}

func (w *WhichKeyOverlay) SetFocusSetter(setFocus func(p tview.Primitive)) {
	w.focusSetter = setFocus
}

func (w *WhichKeyOverlay) SetReturnFocus(primitive tview.Primitive) {
	w.returnFocus = primitive
}

func (w *WhichKeyOverlay) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return w.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if !w.visible || w.tree == nil {
			return
		}

		switch event.Key() {
		case tcell.KeyEscape:
			w.handleEscape()
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			w.handleBack()
		default:
			if key := event.Rune(); key != 0 {
				w.handleKeyPress(key)
			}
		}
	})
}

func (w *WhichKeyOverlay) Draw(screen tcell.Screen) {
	if !w.visible {
		return
	}

	if w.tree != nil {
		w.syncFromTree()
	}

	container := w.containerRect
	if container.width == 0 || container.height == 0 {
		container.x, container.y, container.width, container.height = w.GetRect()
	}

	overlayWidth, overlayHeight := w.calculateSize(container.width, container.height)
	overlayX, overlayY := w.positionRect(container, overlayWidth, overlayHeight)

	w.Box.SetRect(overlayX, overlayY, overlayWidth, overlayHeight)
	w.Box.SetBorderColor(w.styles.BorderColor)
	w.Box.SetTitleColor(w.styles.TitleColor)
	w.Box.SetTitle(w.title)
	w.Box.SetBackgroundColor(w.styles.BackgroundColor)
	w.Box.DrawForSubclass(screen, w)

	innerX, innerY, innerW, innerH := w.GetInnerRect()
	w.drawContent(screen, innerX, innerY, innerW, innerH)
}

func (w *WhichKeyOverlay) calculateSize(maxWidth, maxHeight int) (int, int) {
	const (
		minWidth   = 20
		minHeight  = 3
		paddingX   = 1
		paddingY   = 0
		columnGap  = 2
		borderSize = 2
	)

	maxKey := 0
	maxDesc := 0
	for _, hint := range w.content {
		if length := utf8.RuneCountInString(hint.Key); length > maxKey {
			maxKey = length
		}
		if length := utf8.RuneCountInString(hint.Description); length > maxDesc {
			maxDesc = length
		}
	}

	if maxKey == 0 {
		maxKey = 1
	}
	if maxDesc == 0 {
		maxDesc = 1
	}

	contentWidth := maxKey + columnGap + maxDesc
	innerWidth := contentWidth + paddingX*2
	width := innerWidth + borderSize
	if width < minWidth {
		width = minWidth
	}
	if width > maxWidth {
		width = maxWidth
	}

	contentHeight := len(w.content)
	if contentHeight == 0 {
		contentHeight = 1
	}
	innerHeight := contentHeight + paddingY*2
	height := innerHeight + borderSize
	if height < minHeight {
		height = minHeight
	}
	if height > maxHeight {
		height = maxHeight
	}

	return width, height
}

func (w *WhichKeyOverlay) positionRect(container rect, width, height int) (int, int) {
	const margin = 1

	x := container.x
	y := container.y

	switch w.position {
	case PositionBottomRight:
		x = container.x + container.width - width - margin
		y = container.y + container.height - height - margin
	default:
		x = container.x + (container.width-width)/2
		y = container.y + (container.height-height)/2
	}

	if x < container.x {
		x = container.x
	}
	if y < container.y {
		y = container.y
	}

	return x, y
}

func (w *WhichKeyOverlay) drawContent(screen tcell.Screen, x, y, width, height int) {
	if width <= 0 || height <= 0 {
		return
	}

	maxKey := 0
	for _, hint := range w.content {
		if length := utf8.RuneCountInString(hint.Key); length > maxKey {
			maxKey = length
		}
	}
	if maxKey == 0 {
		maxKey = 1
	}

	columnGap := 2
	if maxKey+columnGap >= width {
		columnGap = 1
	}
	keyWidth := maxKey
	if keyWidth > width {
		keyWidth = width
	}
	descWidth := width - keyWidth - columnGap
	if descWidth < 0 {
		descWidth = 0
	}

	rowCount := len(w.content)
	if rowCount > height {
		rowCount = height
	}

	for i := 0; i < rowCount; i++ {
		hint := w.content[i]
		color := w.styles.DescriptionColor
		if hint.IsGroup {
			color = w.styles.GroupColor
		}

		tview.Print(screen, hint.Key, x, y+i, keyWidth, tview.AlignLeft, w.styles.KeyColor)
		if descWidth > 0 {
			tview.Print(screen, hint.Description, x+keyWidth+columnGap, y+i, descWidth, tview.AlignLeft, color)
		}
	}
}

func (w *WhichKeyOverlay) handleKeyPress(key rune) {
	if w.tree == nil {
		return
	}
	if !w.tree.Enter(key) {
		return
	}

	if w.tree.Current != nil && w.tree.Current.Action != nil {
		w.tree.Current.Action()
		w.tree.Reset()
		w.Hide()
		return
	}

	w.syncFromTree()
}

func (w *WhichKeyOverlay) handleEscape() {
	if w.tree == nil {
		w.Hide()
		return
	}
	if w.tree.Current != w.tree.Root {
		w.tree.Back()
		w.syncFromTree()
		return
	}

	w.tree.Reset()
	w.Hide()
}

func (w *WhichKeyOverlay) handleBack() {
	if w.tree == nil {
		return
	}
	if w.tree.Current != w.tree.Root {
		w.tree.Back()
		w.syncFromTree()
	}
}

func (w *WhichKeyOverlay) syncFromTree() {
	if w.tree == nil || w.tree.Current == nil {
		w.content = []KeyHint{{Description: "No mappings"}}
		w.title = string(w.leaderPrefix)
		return
	}

	children := w.tree.Current.Children
	if len(children) == 0 {
		w.content = []KeyHint{{Description: "No mappings"}}
		w.title = w.tree.Breadcrumb(w.leaderPrefix)
		return
	}

	keys := make([]rune, 0, len(children))
	for key := range children {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	hints := make([]KeyHint, 0, len(children))
	for _, key := range keys {
		child := children[key]
		description := child.Description
		isGroup := len(child.Children) > 0 || child.Action == nil
		if isGroup {
			if description == "" {
				description = "+ group"
			} else {
				description = "+ " + description
			}
		}
		hints = append(hints, KeyHint{
			Key:         string(key),
			Description: description,
			IsGroup:     isGroup,
		})
	}

	w.content = hints
	w.title = w.tree.Breadcrumb(w.leaderPrefix)
}
