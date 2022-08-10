package ui

import "github.com/dogerescat/vim-note/ui/tui"

type sizeSpec struct {
	size    float64
	percent bool
}

//[4]sizeSpec trbl
func defaultMargin() [4]sizeSpec {
	var s [4]sizeSpec
	s[0].percent = true
	s[1].percent = true
	s[2].percent = true
	s[3].percent = true
	return s
}

type Option struct {
	prompt  string
	pointer string
	query   string
	tabstop int
	margin  [4]sizeSpec
	padding [4]sizeSpec
	keymap  map[tui.Event][]*action
}

func DefaultOption() *Option {
	return &Option{
		prompt:  "search> ",
		pointer: ">",
		query:   "",
		tabstop: 8,
		margin:  defaultMargin(),
		padding: defaultMargin(),
	}
}
