package ui

func Run(list []string) string {
	sampleData = list
	t := NewTerminal()
	return t.Loop()
}
