package Parameters

import "net"

// Global variables
var (
	IsServer             bool = true
	IsClient             bool = false
	DestIPAddr           string
	DefaultPort          uint   = 1337
	ConnHost             string = "0.0.0.0"
	ConnPort             string = "1337"
	NodeIndexes                 = make(map[string]int)
	Connections          []Node
	ServerC              net.Conn
	ListenDialBoth       bool = false
	FirstDialCall        bool = true
	SelectedModule       int  = 1
	SelectedModuleName   string
	ModuleNames          []string
	AllModules           = make(map[string]Module)
	CLIParameters        string
	IsAttackStarted      bool = false
	IsAttackOngoing      bool = false
	IsAttackPaused       bool = false
	IsAttackStopped      bool = false
	IsThisDeviceExcluded bool = false
	AllNodesCompleted    bool = false
	CompletedNodeCount   int  = 0
	TotalNodes           int  = 1
	TotalActiveNodes     int  = 1
	CollectParameters    []Parameter
	Tasks                []Task
	SharedTasks              = make(map[string][]ToNode)
	TotalTaskCount       int = 0
	AttackOutput         string
)
