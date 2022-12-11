package ui

func Run(list []string) string {
	fileNames = list
	t := NewTerminal()
	return t.Loop()
}
