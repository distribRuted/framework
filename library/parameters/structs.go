package Parameters

import "net"

type Parameter struct {
	Name        string
	Description string
	Value       string
	Required    bool
}

type Node struct {
	DestAddr string
	IsActive bool
	NetConn  net.Conn
	// LastActive
}

type Module struct {
	Name        string
	Description string
	Author      string
}

type Task struct {
	ID    int
	Name  string
	From  int
	To    int
	State string // Queued / Progressing / Succeeded / Failed
}

type ToNode struct {
	NodeID    int
	NodeIP    string
	From      int
	To        int
	Completed bool // true / false
}
