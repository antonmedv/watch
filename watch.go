package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var version = "1.0.0"

func main() {
	// Command to watch can come from stdin or arguments.
	var command string

	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stdout, "Welcome to watch %v.\nType command to watch.\n> ", version)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			command = scanner.Text()
			break
		}
	} else {
		command = strings.Join(os.Args[1:], " ")
	}

	startTime := time.Now()
	sleep := 1 * time.Second
	shell := defaultShell
	if s, ok := os.LookupEnv("WATCH_COMMAND"); ok {
		shell = s
	}
	sh := strings.Split(shell, " ")
	sh = append(sh, command)

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = screen.Init()
	if err != nil {
		panic(err)
	}

	app := tview.NewApplication()
	app.SetScreen(screen)
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
				screen.Clear()
				viewer.SetText(tview.TranslateANSI(buf.String()))
				elapsed.SetText(fmt.Sprintf("%v", time.Since(startTime).Round(time.Second)))
			})
			time.Sleep(sleep)
		}
	}()

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
