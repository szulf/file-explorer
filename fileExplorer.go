package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

// Maybe change the name
type Events struct {
	Quit     string
	MoveUp   string
	MoveDown string
	EnterDir string

	CmdMode string
	Touch   string
	MkDir   string
	Remove  string
	Move    string
	Copy    string
	Rename  string

	Editor    string
	SupEditor string
}

type FileEntry struct {
	name  string
	isDir bool
}

type FileExplorer struct {
	win     *gc.Window
	entries []FileEntry
	chosen  int
	err     error
	errLoc  string
}

func NewFileExplorer() FileExplorer {
	fe := FileExplorer{chosen: 0}
	fe.Init()
	return fe
}

func (fe *FileExplorer) Init() {
	var err error
	fe.win, err = gc.Init()
	if err != nil {
		fe.err = err
		fe.errLoc = "Init() Error creation"
		return
	}

	if !gc.HasColors() {
		gc.End()
		fe.err = fmt.Errorf("your terminal has to support colors in order to run the application")
		fe.errLoc = "Init() Color check"
		return
	}

	err = gc.StartColor()
	if err != nil {
		fe.err = err
		fe.errLoc = "Init() Starting colors"
		return
	}

	gc.InitPair(1, gc.C_GREEN, gc.C_BLACK)

	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	fe.win.Clear()
	fe.win.Keypad(true)

	fe.getEntries()
	fe.render()
}

func (fe *FileExplorer) getEntries() {
	file, err := os.Open(".")
	if err != nil {
		fe.err = err
		fe.errLoc = "getEntries() Opening dir"
		return
	}

	contents, err := file.Readdir(0)
	if err != nil {
		fe.err = err
		fe.errLoc = "getEntries() Reading dir"
		return
	}
	file.Close()

	entries := make([]FileEntry, len(contents)+1)
	entries[0] = FileEntry{name: "..", isDir: true}
	for i, val := range contents {
		entries[i+1] = FileEntry{name: val.Name(), isDir: val.IsDir()}
	}
	fe.entries = entries
}

func (fe FileExplorer) renderLine(entry FileEntry, i int) {
	if i == fe.chosen {
		fe.win.AttrOn(gc.A_REVERSE)
	}
	if entry.isDir {
		fe.win.AttrOn(gc.ColorPair(1))
	}
	fe.win.Println(entry.name)
	if entry.isDir {
		fe.win.AttrOff(gc.ColorPair(1))
	}
	if i == fe.chosen {
		fe.win.AttrOff(gc.A_REVERSE)
	}
}

func (fe FileExplorer) render() {
	// Maybe clear?
	fe.win.Move(0, 0)
	for i, item := range fe.entries {
		fe.renderLine(item, i)
	}
}

// Maybe change the name later
// Also dont know if this is really necessary
func (fe FileExplorer) reRenderLine(id int) {
	fe.win.Move(id, 0)
	fe.win.ClearToEOL()
	fe.renderLine(fe.entries[id], id)
}

func (fe *FileExplorer) MoveDown() {
	if fe.chosen < len(fe.entries)-1 {
		fe.chosen++
		fe.reRenderLine(fe.chosen - 1)
		fe.reRenderLine(fe.chosen)
	}
}

func (fe *FileExplorer) MoveUp() {
	if fe.chosen > 0 {
		fe.chosen--
		fe.reRenderLine(fe.chosen + 1)
		fe.reRenderLine(fe.chosen)
	}
}

func (fe *FileExplorer) EnterDir() {
	if fe.entries[fe.chosen].isDir {
		err := os.Chdir(fe.entries[fe.chosen].name)
		if err != nil {
			fe.err = err
			fe.errLoc = "EnterDir() Changing dir"
			return
		}
		fe.chosen = 0
		fe.getEntries()
		fe.win.Clear()
		fe.render()
	}
}

func (fe *FileExplorer) EventLoop(events Events) {
	for {
		if fe.err != nil {
			gc.End()
			log.Fatal("Error at: ", fe.errLoc, "\n", fe.err)
		}

		ch := fe.win.GetChar()

		switch gc.KeyString(ch) {
		case events.Quit:
			gc.End()
			os.Exit(0)

		case events.EnterDir:
			fe.EnterDir()

		case events.MoveUp:
			fe.MoveUp()

		case events.MoveDown:
			fe.MoveDown()

		case events.CmdMode:
			maxY, _ := fe.win.MaxYX()
			gc.Echo(true)
			fe.win.Move(maxY-1, 0)
			gc.Cursor(1)
			str, err := fe.win.GetString(128)
			if err != nil {
				log.Fatal("Error getting command string: ", err)
			}

			strArr := strings.Split(str, " ")
			cmd := exec.Command(strArr[0], strArr[1:]...)
			cmdOutput, err := cmd.Output()
			if err != nil {
				cmdOutput = []byte("Invalid Command: " + str)
			}

			gc.Echo(false)
			gc.Cursor(0)

			gc.End()
			fmt.Println(string(cmdOutput), maxY)
			os.Exit(0)

		case events.Touch:

		case events.MkDir:

		case events.Remove:

		case events.Move:

		case events.Copy:

		case events.Rename:

		case events.Editor:

		case events.SupEditor:
		}
	}
}
