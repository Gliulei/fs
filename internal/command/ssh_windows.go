//go:build windows
// +build windows

package command

import "golang.org/x/crypto/ssh"

func initWindowListener(fd int, session *ssh.Session) {
    // Windows 不支持 SIGWINCH，什么都不做
}