# Go-RevShell

<img src="https://go.dev/images/gophers/pilot-bust.svg" width="200" height="200" align="right">

## Installation & Usage

1. Clone the repo `git clone git@github.com:Lucaskrell/go_revshell.git` and get into it `cd go_revshell`
2. Modify line 19 of file `go_revshell.go` and set `localhost` to client IP and `1111` to client port
3. Compile the code using `go build`
4. Execute the file produced at step 3 server side
5. Listen to the given port at step 2 client side using `nc -lp port`

## Screenshot

![Go-RevShell](img/go_revshell.PNG)

## To-Do

- Generate Reverse Shell file using user arguments (client IP, client port)
