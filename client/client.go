package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

const (
	clientPort = ":1234"
	APP1       = "PageRank"
	APP2       = "Another"
)

var (
	myID        int
	rpcRequest  service.rpcRequest
	rpcResponse service.rpcResponse
)

// TODO:change to util!
func getIDFromHostname() int {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Println("hostname:", hostname)
	list := strings.SplitN(hostname, ".", 2)
	if len(list) > 0 {
		tempStr := list[0]
		id, err := strconv.Atoi(tempStr[len(tempStr)-2:])
		if err != nil {
			// If not in the format of "fa17-cs425-g28-%02d.cs.illinois.edu"
			// just return 0 (to allow running in local developement)
			return 0
		}
		return id - 1
	}
	panic("No valid hostname!")
}

func main() {
	var app, data string

	for {
		// handle input
		fmt.Println("Please enter command like: <Application> <Dataset filename>\nApplication includes PageRank and Another\n")
		fmt.Scanln(&app, &data)

		if app != APP1 && app != APP2 {
			fmt.Println("Invalid command, please enter correct command\n")
			continue
		}

		dataset, err := ioutil.ReadFile(data)
		if err != nil {
			fmt.Println("unable to open the file, please enter correct command\n")
			continue
		}

		// build Request
		rpcRequest.ClientID = myID
		rpcRequest.Application = app
		rpcRequest.Dataset = dataset

		// rpc clent
		conn, err := net.Dial("tcp", "localhost"+clientPort)
		if err != nil {
			fmt.Println("cannot connect the master, enter command and try again.")
			continue
		}
		defer conn.Close()

		var response service.Response
		rpcClient := rpc.NewClient(conn)
		err = rpcClient.Call("Service.ComputeGraph", rpcRequest, &rpcResponse)
		if err != nil {
			fmt.Println("fail to call rpc, try again!\n")
			continue
		}

		fmt.Println(rpcResponse)
	}
}
