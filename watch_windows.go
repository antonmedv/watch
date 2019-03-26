// +build windows

package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	pty "github.com/iamacarpet/go-winpty"
)

const defaultShell = "cmd /c"

func cmdOutput(cmd *exec.Cmd, buf *bytes.Buffer) error {
	ptmx, err := pty.OpenWithOptions(pty.Options{
		Command: cmd.String(),
		Env:     os.Environ(),
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(buf, ptmx.StdOut)
	if err != nil {
		return err
	}

	ptmx.Close()

	return nil
}
