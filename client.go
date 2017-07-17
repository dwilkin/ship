package main

import "net"
import "fmt"
import "bufio"
import "os"

//Used for testing, got from https://systembash.com/a-simple-go-tcp-server-and-tcp-client/
func main() {

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text + "\n")
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
	}
}
