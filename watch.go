package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "1.1.0"

const friendlyAppName = "watch"

func main() {
	v := viper.New()
	rcmd := newRootCmd(v)
	configure(rcmd.Flags(), v)

	rcmd.SetArgs(os.Args[1:])
	err := rcmd.Execute()

	if err != nil {
		panic(err)
	}
}

type option string
const (
	intervalOption option = "interval"
)

func configure(f *flag.FlagSet, v *viper.Viper) {
	{
		key := string(intervalOption)
		f.IntP(key, "n", 1, "seconds to wait between updates")
		v.BindPFlag(key, f.Lookup(key))
	}
}

func newRootCmd(v *viper.Viper) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [options] command", friendlyAppName),
		Short: fmt.Sprintf("%s rewritten in go", friendlyAppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command to watch can come from stdin or arguments.
			var command string

			if len(args) == 0 {
				fmt.Fprintf(os.Stdout, "Welcome to watch %v.\nType command to watch.\n> ", version)
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					command = scanner.Text()
					break
				}
			} else {
				command = strings.Join(args, " ")
			}

			startTime := time.Now()
			sleep := time.Duration(v.GetInt(string(intervalOption))) * time.Second

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

			return app.Run()
		},
	}

	return rootCmd
}
