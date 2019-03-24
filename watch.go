package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/kr/pty"
	"github.com/rivo/tview"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()
	sleep := 1 * time.Second
	command := strings.Join(os.Args[1:], " ")

	shell := "bash -cli"
	if s, ok := os.LookupEnv("WATCH_COMMAND"); ok {
		shell = s
	}

	sh := strings.Split(shell, " ")
	sh = append(sh, command)

	app := tview.NewApplication()
	viewer := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetTextColor(tcell.ColorDefault)
	viewer.
		SetBackgroundColor(tcell.ColorDefault)
	elapsed := tview.NewTextView().
		SetTextColor(tcell.ColorBlack).
		SetTextAlign(tview.AlignRight).
		SetText("0s")
	elapsed.
		SetBackgroundColor(tcell.ColorLimeGreen)
	title := tview.NewTextView().
		SetTextColor(tcell.ColorBlack).
		SetText(command)
	title.
		SetBackgroundColor(tcell.ColorLimeGreen)
	statusBar := tview.NewFlex().
		AddItem(title, 0, 1, false).
		AddItem(elapsed, 7, 1, false)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	flex.AddItem(viewer, 0, 1, true)
	flex.AddItem(statusBar, 1, 1, false)
	app.SetRoot(flex, true)

	go func() {
		for {
			var err error

			cmd := exec.Command(sh[0], sh[1:]...)

			ptmx, err := pty.Start(cmd)
			if err != nil {
				panic(err)
			}

			err = pty.InheritSize(os.Stdin, ptmx)
			if err != nil {
				panic(err)
			}

			viewer.SetText("")

			w := tview.ANSIWriter(viewer)
			if _, err := io.Copy(w, ptmx); err != nil {
				panic(err)
			}

			err = ptmx.Close()
			if err != nil {
				panic(err)
			}

			app.QueueUpdateDraw(func() {
				d := time.Now().Sub(startTime).Round(time.Second)
				elapsed.SetText(fmt.Sprintf("%v", d))
			})
			time.Sleep(sleep)
		}
	}()

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
