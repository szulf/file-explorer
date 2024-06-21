package main

func main() {
	fe := NewFileExplorer()

	fe.EventLoop(Events{
		Quit:     "q",
		EnterDir: "enter",
		MoveUp:   "k",
		MoveDown: "j",

		CmdMode: ":",
	})
}
