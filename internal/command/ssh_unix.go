//go:build !windows
// +build !windows

package command

import (
    "os"
    "os/signal"
    "syscall"
    "golang.org/x/crypto/ssh"
    "golang.org/x/term"
)

func initWindowListener(fd int, session *ssh.Session) {
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGWINCH)
    go func() {
        for range sigc {
            if w, h, err := term.GetSize(fd); err == nil {
                _ = session.WindowChange(h, w)
            }
        }
    }()
}