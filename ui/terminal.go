package ui

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/dogerescat/vim-note/ui/algo"
	"github.com/dogerescat/vim-note/ui/tui"
	"github.com/dogerescat/vim-note/ui/util"
)

const (
	minWidth  = 4
	minHeight = 4
)

const (
	reqPrompt util.EventType = iota
	reqHeader
	reqList
	reqRefresh
	reqRedraw
	reqClose
	reqPrintQuery
	reqQuit
)

type action struct {
	t actionType
	a string
}

type actionType int

const (
	actIgnore actionType = iota
	actInvalid
	actRune
	actMouse
	actBeginningOfLine
	actAbort
	actAccept
	actAcceptNonEmpty
	actBackwardChar
	actBackwardDeleteChar
	actBackwardDeleteCharEOF
	actBackwardWord
	actCancel
	actChangePrompt
	actClearScreen
	actClearQuery
	actClearSelection
	actClose
	actDeleteChar
	actDeleteCharEOF
	actEndOfLine
	actForwardChar
	actForwardWord
	actYank
	actBackwardKillWord
	actSelectAll
	actDeselectAll
	actDown
	actUp
	actPrintQuery
	actExecute
	actExecuteSilent
	actSigStop
	actFirst
	actLast
	actReload
	actDisableSearch
	actEnableSearch
	actSelect
	actDeselect
	actUnbind
	actRebind
)

var sampleData []string

func toAction(types ...actionType) []*action {
	actions := make([]*action, len(types))
	for i, t := range types {
		actions[i] = &action{t: t, a: ""}
	}
	return actions
}

func defaultKeymap() map[tui.Event][]*action {
	keymap := make(map[tui.Event][]*action)
	add := func(e tui.EventType, a actionType) {
		keymap[e.AsEvent()] = toAction(a)
	}
	add(tui.Invalid, actInvalid)
	add(tui.CtrlC, actAbort)
	add(tui.ESC, actAbort)
	add(tui.CtrlH, actBackwardDeleteChar)
	add(tui.BSpace, actBackwardDeleteChar)
	add(tui.Del, actDeleteChar)
	add(tui.CtrlN, actDown)
	add(tui.CtrlP, actUp)
	add(tui.Insert, actExecute)
	add(tui.Up, actUp)
	add(tui.Down, actDown)
	add(tui.Left, actBackwardChar)
	add(tui.Right, actForwardChar)
	return keymap
}

func copySlice(slice []rune) []rune {
	ret := make([]rune, len(slice))
	copy(ret, slice)
	return ret
}

type Terminal struct {
	window   tui.Window
	frame    tui.Window
	mutex    sync.Mutex
	tui      tui.Renderer
	reqBox   *util.EventBox
	prompt   string
	pointer  string
	cx       int
	cy       int
	margin   [4]sizeSpec
	padding  [4]sizeSpec
	input    []rune
	initFunc func()
	printer  func(string)
	keymap   map[tui.Event][]*action
}

func trimQuery(query string) []rune {
	return []rune(strings.Replace(query, "\t", " ", -1))
}

func NewTerminal() *Terminal {
	opt := DefaultOption()
	input := trimQuery(opt.prompt)
	renderer := tui.NewLightRenderer()
	return &Terminal{
		mutex:    sync.Mutex{},
		tui:      renderer,
		reqBox:   util.NewEventBox(),
		prompt:   opt.prompt,
		pointer:  ">",
		initFunc: func() { renderer.Init() },
		margin:   defaultMargin(),
		padding:  defaultMargin(),
		printer:  func(str string) { fmt.Println(str) },
		keymap:   defaultKeymap(),
		input:    input,
		cx:       len(input),
		cy:       2,
	}
}

func (t *Terminal) refresh() {
	windows := make([]tui.Window, 0, 2)
	windows = append(windows, t.frame)
	windows = append(windows, t.window)
	t.tui.RefreshWindows(windows)
}

func (t *Terminal) resizeWindows() {
	screenWidth := t.tui.MaxX()
	screenHeight := t.tui.MaxY()
	marginInt := [4]int{}
	paddingInt := [4]int{}
	sizeSpecToInt := func(idx int, spec sizeSpec) int {
		if spec.percent {
			var max float64
			if idx%2 == 0 {
				max = float64(screenHeight)
			} else {
				max = float64(screenWidth)
			}
			return int(max * 0.3)
		}
		return int(spec.size)
	}

	for idx, sizeSpec := range t.padding {
		paddingInt[idx] = sizeSpecToInt(idx, sizeSpec)
	}

	extraMargin := [4]int{}
	for idx, sizeSpec := range t.margin {
		extraMargin[idx] += 1 + idx%2
		marginInt[idx] = sizeSpecToInt(idx, sizeSpec)
	}

	adjust := func(idx1 int, idx2 int, max int, min int) {
		if max >= min {
			margin := marginInt[idx1] + marginInt[idx2] + paddingInt[idx1] + paddingInt[idx2]
			if max-margin < min {
				desired := max - min
				paddingInt[idx1] = desired * paddingInt[idx1] / margin
				paddingInt[idx2] = desired * paddingInt[idx2] / margin
				marginInt[idx1] = util.Max(extraMargin[idx1], desired*marginInt[idx1]/margin)
				marginInt[idx2] = util.Max(extraMargin[idx2], desired*marginInt[idx2]/margin)
			}
		}
	}

	minAreaWidth := minWidth
	minAreaHeight := minHeight
	adjust(1, 3, screenWidth, minAreaWidth)
	adjust(0, 2, screenHeight, minAreaHeight)

	width := screenWidth - marginInt[1] - marginInt[3]
	height := screenHeight - marginInt[0] - marginInt[2]
	t.frame = t.tui.NewWindow(marginInt[0]-1, marginInt[3]-2, width+4, height+2, tui.MakeFrameStyle(), true)
	if t.window == nil {
		t.window = t.tui.NewWindow(marginInt[0], marginInt[3], width, height, tui.MakeFrameStyle(), false)
	}
	for i := 0; i < t.window.Height(); i++ {
		t.window.MoveAndClear(i, 0)
	}
}

func (t *Terminal) move(y int, x int, is_clear bool) {
	if is_clear {
		t.window.MoveAndClear(y, x)
	} else {
		t.window.Move(y, x)
	}
}

func (t *Terminal) delChar() bool {
	if len(t.input) > 0 && t.cx < len(t.input) {
		t.input = append(t.input[:t.cx], t.input[t.cx+1:]...)
		return true
	}
	return false
}

func (t *Terminal) printItem() {
	m := math.Min(float64(t.window.Height()-2), float64(len(sampleData)))
	for i := 0; i < int(m); i++ {
		t.move(i+2, 2, true)
		t.window.Print(sampleData[i])
	}
}

func (t *Terminal) printPrompt() {
	t.move(0, 0, true)
	t.window.Print(string(t.input))
	t.move(t.cy, 0, false)
	t.window.Print(t.pointer)
	t.window.Move(0, t.cx)
}

func (t *Terminal) printAll() {
	t.resizeWindows()
	t.printList()
	t.printPrompt()
	t.refresh()
}

func (t *Terminal) printList() {
	t.printItem()
}

func (t *Terminal) redraw() {
	t.printAll()
}

func (t *Terminal) sortList() {
	str := string(t.input[len(t.prompt):])
	sampleData = algo.MatchString(str, sampleData)
}

func Constrain(val int, min int, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (t *Terminal) Loop() string {
	t.mutex.Lock()
	t.initFunc()
	t.resizeWindows()
	t.printList()
	t.printPrompt()
	t.refresh()
	t.mutex.Unlock()

	var res string
	go func() {
		running := true
		for running {
			t.reqBox.Wait(func(events *util.Events) {
				defer events.Clear()
				t.mutex.Lock()
				for req := range *events {
					switch req {
					case reqList:
						t.sortList()
						t.printList()
						t.printPrompt()
					case reqPrompt:
						t.printPrompt()
					case reqClose:
						running = false
						return
					case reqQuit:
						running = false
						return
					}
				}
			})
			t.refresh()
			t.mutex.Unlock()
		}
	}()

	looping := true
	for looping {
		queryChanged := false
		previousCx := t.cx
		previousCy := t.cy
		previousInput := t.input
		event := t.tui.GetChar()
		t.mutex.Lock()
		events := []util.EventType{}
		req := func(evts ...util.EventType) {
			for _, event := range evts {
				events = append(events, event)
				if event == reqClose || event == reqQuit {
					looping = false
				}
			}
		}

		var doAction func(*action) bool
		doActions := func(actions []*action) bool {
			for _, action := range actions {
				if !doAction(action) {
					return false
				}
			}
			return true
		}

		doAction = func(action *action) bool {
			switch action.t {
			case actIgnore:
			case actChangePrompt:
				req(reqPrompt)
			case actAbort:
				req(reqQuit)
			case actRune:
				prefix := copySlice(t.input[:t.cx])
				t.input = append(append(prefix, event.Char), t.input[t.cx:]...)
				t.cx++
				req(reqList)
			case actUp:
				if t.cy > 2 {
					t.cy--
				}
				//req(reqList)
			case actDown:
				if t.cy < t.window.Height()-1 && t.cy <= len(sampleData) {
					t.cy++
				}
				//req(reqList)
			case actBackwardChar:
				if t.cx > len(t.prompt) {
					t.cx--
				}
			case actForwardChar:
				if t.cx < len(t.input) {
					t.cx++
				}
			case actDeleteChar:
				t.delChar()
				req(reqList)
			case actBackwardDeleteChar:
				if t.cx > len(t.prompt) {
					t.input = append(t.input[:t.cx-1], t.input[t.cx:]...)
					t.cx--
					req(reqList)
				}
			case actExecute:
				res = sampleData[t.cy-2]
				req(reqQuit)
			case actClose:
				req(reqQuit)
			}
			return true
		}

		actions := t.keymap[event.Comparable()]
		if len(actions) == 0 && event.Type == tui.Rune {
			doAction(&action{t: actRune})
		} else if !doActions(actions) {
			continue
		}
		queryChanged = string(previousInput) != string(t.input)
		if queryChanged || previousCx != t.cx || previousCy != t.cy {
			req(reqPrompt)
		}
		t.mutex.Unlock()
		for _, event := range events {
			t.reqBox.Set(event, nil)
		}
	}
	t.tui.Close()
	return res
}
