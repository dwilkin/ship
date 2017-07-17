package main

import (
	"net"
	"fmt"
	"bufio"
	"strings"
	"os"
	"io"
	"sync"
)


type Index struct {
	sync.Mutex
	dex map[string][]string
}

func initIndex() Index {
	return Index{dex: make(map[string][]string)}
}

var d = initIndex() //This is the map in which we'll be storing the index

//Ensures that all dependencies of a package are indexed prior to indexing a new package
func dependenciesSatisfied(dependencies []string) bool{
	for i := 0; i < len(dependencies); i++ {
		_, exists := d.dex[dependencies[i]]
		if exists == false {
			return false
		}
	}
	return true
}


func index(pkg string, dependencies string) string {

	var dep_parsed []string

	//If there are dependencies, break them into an array and check them against currently indexed packages
	//If all of those packages are already indexed, index the new package.  Otherwise return an error
	if dependencies != "" {
		dep_parsed = strings.Split(dependencies, ",")
		if dependenciesSatisfied(dep_parsed) == true {
			d.dex[pkg] = dep_parsed
		} else {
			return "FAIL"
		}
	} else {
		//If there are no dependencies, we can simply install the package with no listed dependencies
		d.dex[pkg] = nil
	}
	return "OK"
}

//Just checks to see if the package is already indexed or not
func query(pkg string) string {
	_, exists := d.dex[pkg]
	if exists == false {
		return "FAIL"
	}
	return "OK"
}

//This is the tricky one because I didn't use a tree, I just used a map.
//I know enough to know that this is very bad from a big O perspective
//However I don't know a better way to do it without switching away from maps
func remove(pkg string) string {
	_, exists := d.dex[pkg]
	if exists == true {
		for _, value := range d.dex {
			for i := 0; i < len(value); i++ {
				if value[i] == pkg {
					return "FAIL"
				}
			}
		}

		delete(d.dex, pkg)
	}
	return "OK"
}

//This function parses the data from the client and sends it to the appropriate function to process the request
func parseData(conn net.Conn) {

	var result string

	defer conn.Close()

	for {
		//Take the raw connection and turn it into a string
		message, err := bufio.NewReader(conn).ReadString('\n')
		//Learned about io.EOF here: https://appliedgo.net/networking/
		switch {
		case err == io.EOF:
			e := "Reached EOF"
			fmt.Println(e)
			return
		case err != nil:
			e := "Error reading command"
			fmt.Println(e)
			return
		}
		fmt.Println("Message Received:", string(message))

		//Trim the newline character from the end so we can parse the message
		message = strings.TrimSuffix(message, "\n")

		//Next we separate the string into three parts
		data := strings.Split(message, "|")
		if len(data) < 3 {
			e := "Not enough pipe-separated fields in the message"
			fmt.Println(e)
			result = "ERROR"
		} else {

			//Put the pieces into their appropriate buckets
			command := data[0]
			pkg := data[1]
			dependencies := data[2]

			fmt.Println("Command given is:", string(command))
			fmt.Println("Package given is:", string(pkg))
			fmt.Println("That package depends on:", string(dependencies))

			//based on the command sent
			switch command {
			case "INDEX":
				result = index(pkg, dependencies)
			case "QUERY":
				result = query(pkg)
			case "REMOVE":
				result = remove(pkg)
			default:
				e := "Invalid command: "
				fmt.Println(e, string(command))
				result = "ERROR"
			}
		}

		fmt.Println("Sending response", result+"\n")
		conn.Write([]byte(result + "\n"))
	}
}


func main() {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		e := "Error creating listener"
		fmt.Println(e)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			e := "Error accepting connection"
			fmt.Println(e)
		}

		go parseData(conn)
		fmt.Println("Got through parseData")
	}
}
