package tui

import (
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

var (
	_screen tcell.Screen
)

const (
	Rune EventType = iota
	CtrlA
	CtrlB
	CtrlC
	CtrlD
	CtrlE
	CtrlF
	CtrlG
	CtrlH
	Tab
	CtrlJ
	CtrlK
	CtrlL
	CtrlM
	CtrlN
	CtrlO
	CtrlP
	CtrlQ
	CtrlR
	CtrlS
	CtrlT
	CtrlU
	CtrlV
	CtrlW
	CtrlX
	CtrlY
	CtrlZ
	ESC
	CtrlSpace
	CtrlBackSlash
	CtrlRightBracket
	CtrlCaret
	CtrlSlash
	Invalid
	Resize
	Mouse
	DoubleClick
	LeftClick
	RightClick
	BTab
	BSpace
	Del
	PgUp
	PgDn
	Up
	Down
	Left
	Right
	Home
	End
	Insert
	SUp
	SDown
	SLeft
	SRight
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	Change
	BackwardEOF
	AltBS
	AltUp
	AltDown
	AltLeft
	AltRight
	AltSUp
	AltSDown
	AltSLeft
	AltSRight
	Alt
	CtrlAlt
)

type EventType int

type FillReturn int

type Event struct {
	Type EventType
	Char rune
}

const (
	FillContinue FillReturn = iota
	FillNextLine
	FillSuspend
)

func (et EventType) AsEvent() Event {
	return Event{et, 0}
}

func (e Event) Comparable() Event {
	return Event{e.Type, 0}
}

type FrameStyle struct {
	horizontal  rune
	vertical    rune
	topLeft     rune
	topRight    rune
	bottomLeft  rune
	bottomRight rune
}

func MakeFrameStyle() FrameStyle {
	return FrameStyle{
		horizontal:  '-',
		vertical:    '│',
		topLeft:     '┌',
		topRight:    '┐',
		bottomLeft:  '└',
		bottomRight: '┘',
	}
}

type Window interface {
	Top() int
	Left() int
	Width() int
	Height() int
	Refresh()
	X() int
	Y() int
	Move(y int, x int)
	MoveAndClear(y int, x int)

	Print(text string)
	Fill(text string) FillReturn
}

type Renderer interface {
	Init()
	Clear()
	RefreshWindows(windows []Window)
	Close()
	GetChar() Event
	MaxX() int
	MaxY() int

	NewWindow(top int, left int, width int, height int, frameStyle FrameStyle, isFrame bool) Window
}

type LightRenderer struct {
	width  int
	height int
	x      int
	y      int
}

type LightWindow struct {
	renderer *LightRenderer
	frame    FrameStyle
	top      int
	left     int
	width    int
	height   int
	lastX    int
	lastY    int
	is_frame bool
}

func NewLightRenderer() Renderer {
	return &LightRenderer{}
}

func (r *LightRenderer) screenInit() {
	s, e := tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	if e = s.Init(); e != nil {
		panic(e)
	}
	_screen = s
}

func (r *LightRenderer) Init() {
	r.screenInit()
}

func (r *LightRenderer) NewWindow(top int, left int, width int, height int, frame FrameStyle, isFrame bool) Window {
	w := LightWindow{
		renderer: r,
		frame:    frame,
		top:      top,
		left:     left,
		width:    width,
		height:   height,
		is_frame: isFrame,
	}
	return &w
}

func (r *LightRenderer) Clear() {
	_screen.Sync()
	_screen.Clear()
}
func (r *LightRenderer) MaxX() int {
	x, _ := _screen.Size()
	return int(x)
}

func (r *LightRenderer) MaxY() int {
	_, y := _screen.Size()
	return int(y)
}

func (w *LightWindow) Move(y int, x int) {
	if y != 0 {
		for i := 1; i < w.height; i++ {
			_screen.SetContent(w.left, i+w.top, rune(' '), nil, tcell.StyleDefault)
		}
	}
	w.lastX = x
	w.lastY = y
}

func (w *LightWindow) Top() int {
	return w.top
}

func (w *LightWindow) Left() int {
	return w.left
}

func (w *LightWindow) Width() int {
	return w.width
}

func (w *LightWindow) Height() int {
	return w.height
}

func (w *LightWindow) X() int {
	return w.lastX
}

func (w *LightWindow) Y() int {
	return w.lastY
}

func (w *LightWindow) drawFrame() {
	left := w.left
	right := left + w.width
	top := w.top
	bot := top + w.height
	style := tcell.StyleDefault
	for x := left; x < right; x++ {
		_screen.SetContent(x, top, w.frame.horizontal, nil, style)
	}
	for x := left; x < right; x++ {
		_screen.SetContent(x, bot-1, w.frame.horizontal, nil, style)
	}
	for y := top; y < bot; y++ {
		_screen.SetContent(left, y, w.frame.vertical, nil, style)
	}
	for y := top; y < bot; y++ {
		_screen.SetContent(right-1, y, w.frame.vertical, nil, style)
	}
	_screen.SetContent(left, top, w.frame.topLeft, nil, style)
	_screen.SetContent(right-1, top, w.frame.topRight, nil, style)
	_screen.SetContent(left, bot-1, w.frame.bottomLeft, nil, style)
	_screen.SetContent(right-1, bot-1, w.frame.bottomRight, nil, style)
}

func (w *LightWindow) MoveAndClear(y int, x int) {
	w.Move(y, x)
	for i := w.lastX; i < w.width; i++ {
		_screen.SetContent(i+w.left, w.lastY+w.top, rune(' '), nil, tcell.StyleDefault)
	}
	w.lastY = y
	w.lastX = x
}

func (w *LightWindow) fillString(text string) FillReturn {
	lx := 0

	var style tcell.Style
	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		rs := gr.Runes()
		if len(rs) == 1 && rs[0] == '\n' {
			w.lastY++
			w.lastX = 0
			lx = 0
			continue
		}

		xPos := w.left + w.lastX + lx
		if xPos >= (w.left + w.width) {
			w.lastY++
			w.lastX = 0
			lx = 0
			xPos = w.left
		}

		yPos := w.top + w.lastY
		if yPos >= (w.top + w.height) {
			return FillSuspend
		}

		_screen.SetContent(xPos, yPos, rs[0], rs[1:], style)
		lx += runewidth.StringWidth(string(rs))
	}
	w.lastX += lx
	if w.lastX == w.width {
		w.lastY++
		w.lastX = 0
		return FillNextLine
	}

	return FillContinue
}

func (w *LightWindow) Fill(str string) FillReturn {
	return w.fillString(str)
}

func (w *LightWindow) printString(text string) {
	lx := 0
	var style tcell.Style
	style = style.Normal()
	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		rs := gr.Runes()
		if len(rs) == 1 {
			r := rs[0]
			if r < rune(' ') {
				continue
			} else if r == '\n' {
				w.lastY++
				lx = 0
				continue
			} else if r == '\u000D' {
				continue
			}
		}
		var xPos = w.left + w.lastX + lx
		var yPos = w.top + w.lastY
		if xPos < (w.left+w.width) && yPos < (w.top+w.height) {
			_screen.SetContent(xPos, yPos, rs[0], rs[1:], style)
		}
		lx += runewidth.StringWidth(string(rs))
	}
	w.lastX += lx
}

func (w *LightWindow) Print(str string) {
	w.printString(str)
}

func (w *LightWindow) Refresh() {
	if !w.is_frame {
		_screen.ShowCursor(w.left+w.lastX, w.top+w.lastY)
	}
	w.lastX = 0
	w.lastY = 0
	if w.is_frame {
		w.drawFrame()
	}
}

func (r *LightRenderer) RefreshWindows(windows []Window) {
	for _, w := range windows {
		w.Refresh()
	}
	_screen.Show()
}

func (r *LightRenderer) GetChar() Event {
	ev := _screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			return Event{CtrlC, 0}
		case tcell.KeyCtrlH:
			return Event{CtrlH, 0}
		case tcell.KeyCtrlF:
			return Event{CtrlF, 0}
		case tcell.KeyCtrlB:
			return Event{CtrlB, 0}
		case tcell.KeyUp:
			return Event{Up, 0}
		case tcell.KeyDown:
			return Event{Down, 0}
		case tcell.KeyLeft:
			return Event{Left, 0}
		case tcell.KeyRight:
			return Event{Right, 0}
		case tcell.KeyDelete:
			return Event{Del, 0}
		case tcell.KeyBackspace2:
			return Event{BSpace, 0}
		case tcell.KeyEsc:
			return Event{ESC, 0}
		case tcell.KeyEnter:
			return Event{Insert, 0}
		case tcell.KeyRune:
			r := ev.Rune()
			switch {
			case r == ' ':
				return Event{CtrlSpace, 0}
			default:
				return Event{Rune, r}
			}
		}
	}
	return Event{Invalid, 0}
}

func (r *LightRenderer) Close() {
	_screen.Fini()
}
