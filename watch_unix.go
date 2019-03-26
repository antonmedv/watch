// +build !windows

package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"

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
		// Linux kernel return EIO when attempting to read from a master pseudo
		// terminal which no longer has an open slave. So ignore error here.
		// See https://github.com/kr/pty/issues/21
		if pathErr, ok := err.(*os.PathError); !ok || pathErr.Err != syscall.EIO {
			return err
		}
	}

	return nil
}
