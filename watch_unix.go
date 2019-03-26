// +build !windows

package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/kr/pty"
)

const defaultShell = "bash -cli"

func cmdOutput(cmd *exec.Cmd, buf *bytes.Buffer) error {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	err = pty.InheritSize(os.Stdin, ptmx)
	if err != nil {
		return err
	}

	_, err = io.Copy(buf, ptmx)
	if err != nil {
		return err
	}

	return nil
}
