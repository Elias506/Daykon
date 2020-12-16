package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// scaning your request
func scanStr() string {
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	if err := in.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error -", err)
	}
	return in.Text()
}

func main() {
	fmt.Println("Try to connect...")
	var (
		conn net.Conn
		err  error
	)
	for {
		conn, err = net.Dial("tcp", "127.0.0.1:8080")
		if err == nil {
			break
		}
	}
	fmt.Println("Connection succeeded")
	defer conn.Close()
	for {
		fmt.Print("daykon> ")
		source := scanStr()
		// Write request to server
		n, err := conn.Write([]byte(source))
		if n == 0 || err != nil {
			fmt.Println(err)
			return
		}
		// Read response
		buff := make([]byte, 1024)
		n, err = conn.Read(buff)
		if err != nil {
			break
		}
		fmt.Println(string(buff[0:n]))
	}
}
