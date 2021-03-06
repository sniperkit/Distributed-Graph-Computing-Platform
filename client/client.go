package main

import (
	"bytes"
	"cs425_mp4/protocol-buffer/master-client"
	"cs425_mp4/utility"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	clientPort = ":1234"
	masterID   = 10
	nodeName   = "fa17-cs425-g28-%02d.cs.illinois.edu%s"
	APP1       = "PageRank"
	APP2       = "SSSP"
)

var (
	myID            int
	clientID        int
	standbyMasterID int
	output          masterclient.MasterClient
)

func main() {
	myID = util.GetIDFromHostname()
	var app, data string
	var src int

	for {
		// handle input
		fmt.Println("Please enter command like: <Application> <Dataset filename>\nApplication includes PageRank and SSSP")
		fmt.Scanln(&app, &data)

		if app != APP1 && app != APP2 {
			fmt.Println("Invalid command, please enter correct command")
			continue
		}

		if app == APP2 {
			fmt.Println("SSSP Source?")
			fmt.Scanln(&src)
		}

		start := time.Now()
		dataset, err := ioutil.ReadFile(data)
		if err != nil {
			fmt.Println("unable to open the file, please enter correct command")
			continue
		}

		// send msg to master
		msg := &masterclient.MasterClient{ClientID: uint32(myID)}
		msg.Application = app
		msg.Dataset = dataset
		msg.Src = uint32(src)

		pb, err := proto.Marshal(msg)
		if err != nil {
			fmt.Println("error occured!")
			return
		}

		conn, err := net.Dial("tcp", fmt.Sprintf(nodeName, masterID, clientPort))
		//conn, err := net.Dial("tcp", "localhost"+masterPort)
		if err != nil {
			fmt.Printf("error has occured! %s\n", err)
			return
		}
		conn.Write(pb)
		conn.Close()

		// listen on response
		fmt.Println("Command and dataset have been sent to master, waiting for response!")
		ln, err := net.Listen("tcp", clientPort)
		if err != nil {
			fmt.Printf("error has occured! %s\n", err)
			return
		}
		connln, err := ln.Accept()
		if err != nil {
			fmt.Printf("error has occured! %s\n", err)
			return
		}
		var buf bytes.Buffer

		_, err = io.Copy(&buf, connln)
		if err != nil {
			fmt.Printf("error has occured! %s\n", err)
			return
		}

		proto.Unmarshal(buf.Bytes(), &output)

		fmt.Println(time.Since(start))
		fmt.Println("get output!")
		fmt.Println(string(output.GetResult()))
		connln.Close()
		ln.Close()
	}
}
