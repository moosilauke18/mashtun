package main

import (
	"os"
	"log"
	"strconv"
	"path"
	//"fmt"
	zmq "github.com/alecthomas/gozmq"
)
var basepath = "/home/core/.docker/"

func main() {
	context, _ := zmq.NewContext()
	socket, _ := context.NewSocket(zmq.REP)
	defer context.Close()
	defer socket.Close()
	socket.Bind("tcp://*:2222")


	// Wait for messages
	for {
		msg, _ := socket.Recv(0)
		if string(msg) == "START << EOF" {
			for {
				filepath, _ := socket.Recv(0)
				if string(filepath) == "EOF" {break}

				// Get Dir name
				dir := path.Dir(string(filepath))

				// Make directories
				os.MkdirAll(basepath+dir, 0755)

				// Make file
				file, err := os.Create(string(filepath))
				defer file.Close()
				if err != nil {
					log.Printf("%s",err)
				}
				packets, _ := socket.Recv(0)
				packNum, _ := strconv.Atoi(string(packets))
				for i := 0; i < packNum; i++  {
					msg, _ := socket.Recv(0)

					_, err = file.Write(msg)
				}
			}
		}
		// Received All Files
		socket.Send([]byte("Received"), 0)
	}
}
