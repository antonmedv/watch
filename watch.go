package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/kr/pty"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "usage: watch [command]")
		os.Exit(2)
	}

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

			out, err := ioutil.ReadAll(ptmx)
			if err != nil {
				panic(err)
			}

			err = ptmx.Close()
			if err != nil {
				panic(err)
			}

			app.QueueUpdateDraw(func() {
				viewer.Clear()
				viewer.SetText(tview.TranslateANSI(string(out)))
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
