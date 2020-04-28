package view

import (
	"log"
	"os"

	ui "github.com/cswank/gocui"
)

func Start() error {
	g, err := ui.NewGui(ui.Output256)
	if err != nil {
		return err
	}

	defer g.Close()

	//ui.DefaultEditor = s.footer
	g.SetManagerFunc(render)
	g.Cursor = true
	g.InputEsc = true

	err = g.MainLoop()
	log.SetOutput(os.Stderr)
	return err
}

func render(g *ui.Gui) error {
	return nil
}
