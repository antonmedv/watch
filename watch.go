package main

import (
	"bytes"
	"fmt"

	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell"

	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "usage: watch [command]")
		os.Exit(2)
	}

	startTime := time.Now()
	sleep := 1 * time.Second
	command := strings.Join(os.Args[1:], " ")

	shell := defaultShell
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
		SetBackgroundColor(tcell.ColorLightCyan)
	title := tview.NewTextView().
		SetTextColor(tcell.ColorBlack).
		SetText(command)
	title.
		SetBackgroundColor(tcell.ColorLightCyan)
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
			cmd := exec.Command(sh[0], sh[1:]...)

			var buf bytes.Buffer
			err := cmdOutput(cmd, &buf)
			if err != nil {
				panic(err)
			}

			app.QueueUpdateDraw(func() {
				viewer.Clear()
				viewer.SetText(tview.TranslateANSI(buf.String()))
				elapsed.SetText(fmt.Sprintf("%v", time.Since(startTime).Round(time.Second)))
			})
			time.Sleep(sleep)
		}
	}()

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
