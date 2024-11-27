package Connection

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	Lib "github.com/distribRuted/framework/library"
	Log "github.com/distribRuted/framework/library/log"
	Parameters "github.com/distribRuted/framework/library/parameters"
)

func Dial() {

	if !(Parameters.ListenDialBoth && Parameters.FirstDialCall) {

		var connectionString string
		if net.ParseIP(Parameters.ConnHost).To4() != nil {
			connectionString = Parameters.ConnHost + ":" + Parameters.ConnPort
		} else {
			// If the IP address is IPv6
			connectionString = "[" + Parameters.ConnHost + "]:" + Parameters.ConnPort
		}

		c, err := net.Dial("tcp", connectionString)
		if err != nil {
			Log.Create(err.Error(), "Node", "CRITICAL")
			return
		}

		Parameters.ServerC = c

		time.Sleep(1 * time.Second)

		var nodeIPAddr string = c.RemoteAddr().String()
		var getNodeIndex int = Parameters.GetNodeKey(nodeIPAddr)

		Parameters.Connections = append(Parameters.Connections, Parameters.Node{DestAddr: nodeIPAddr, NetConn: c, IsActive: true})
		Parameters.NodeIndexes[nodeIPAddr] = len(Parameters.Connections) - 1
		getNodeIndex = len(Parameters.Connections) - 1

		Parameters.TotalNodes += 1
		Parameters.TotalActiveNodes += 1
		Log.Create("The central server ["+nodeIPAddr+"] accepted the connection.", "Node", "MEDIUM")

		defer c.Close()

		reader := bufio.NewReader(c)

		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				Parameters.Connections[getNodeIndex].IsActive = false
				Log.Create("The connection to the server ["+nodeIPAddr+"] has been lost: "+err.Error(), "Node", "CRITICAL")
				Parameters.TotalActiveNodes -= 1
				break
			}

			var msgToEvaluate string = strings.TrimSpace(message)
			evaluateMessage(false, true, nodeIPAddr, msgToEvaluate)

			if msgToEvaluate == "STOP" {
				fmt.Println("TCP client exiting...")
				Parameters.TotalActiveNodes -= -1
				break
			}
		}
	} else {
		Parameters.FirstDialCall = false
	}
}

func Listen() {

	var connectionString string
	if net.ParseIP(Parameters.ConnHost).To4() != nil {
		connectionString = Parameters.ConnHost + ":" + Parameters.ConnPort
	} else {
		// If the IP address is IPv6
		connectionString = "[" + Parameters.ConnHost + "]:" + Parameters.ConnPort
	}

	l, err := net.Listen("tcp", connectionString)
	if err != nil {
		fmt.Println()
		Lib.ExitWithError("ERROR! " + err.Error())
		return
	}
	time.Sleep(1 * time.Second)
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			Log.Create(err.Error(), "Node", "HIGH")
			break
		} else {
			var nodeIPAddr string = c.RemoteAddr().String()
			var getNodeIndex int = Parameters.GetNodeKey(nodeIPAddr)
			if getNodeIndex == -1 {
				Parameters.Connections = append(Parameters.Connections, Parameters.Node{DestAddr: nodeIPAddr, NetConn: c, IsActive: true})
				Parameters.NodeIndexes[nodeIPAddr] = len(Parameters.Connections) - 1
				getNodeIndex = len(Parameters.Connections) - 1
				Parameters.TotalNodes += 1
				Parameters.TotalActiveNodes += 1
			} else {
				Parameters.Connections[getNodeIndex].NetConn = c
				Parameters.Connections[getNodeIndex].IsActive = true
				Parameters.TotalActiveNodes += 1
			}
			go handleConnection(Parameters.Connections[getNodeIndex].NetConn)
		}
	}

}

func handleConnection(c net.Conn) {
	Log.Create(c.RemoteAddr().String()+" has been connected.", "Node", "HIGH")

	reader := bufio.NewReader(c)

	for {
		netData, err := reader.ReadString('\n')
		if err != nil {
			Log.Create(c.RemoteAddr().String()+" has been disconnected: "+err.Error(), "Node", "CRITICAL")
			Parameters.Connections[Parameters.GetNodeKey(c.RemoteAddr().String())].IsActive = false
			Parameters.TotalActiveNodes -= 1
			c.Close()
			break
		}

		var msgToEvaluate string = strings.TrimSpace(netData)
		evaluateMessage(true, false, c.RemoteAddr().String(), msgToEvaluate)

		if msgToEvaluate == "PROTOCOL_MSG:disconnect" {
			c.Close()
			Parameters.TotalActiveNodes -= 1
			break
		}

	}
}

func IsIPAddrValid(addr string) bool {
	if net.ParseIP(addr) != nil || addr == "localhost" {
		return true
	}
	return false
}
