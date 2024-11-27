package Modules

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	Log "github.com/distribRuted/framework/library/log"
	Parameters "github.com/distribRuted/framework/library/parameters"
)

func (F *Module) Module_Distributed_Nmap_Init() {

	var (
		moduleShortName   string = "Distributed_Nmap"
		moduleName        string = "Distributed Nmap"
		moduleDescription string = "Distributed Nmap demonstration for DEF CON 32"
		moduleAuthor      string = "@numanozdemircom <root@numanozdemir.com>"
	)

	// ATTENTION! Do not edit or remove the line below:
	Set_Module(moduleShortName, moduleName, moduleDescription, moduleAuthor)
}

// Variables of user inputs to pass to the Start() function.
var (
	targetHost         string
	startPort, endPort int
	startTime          time.Time
)

// Questions to ask the user are defined here.
func (F *Module) Module_Distributed_Nmap_Parameters() {

	Parameters.Add_Parameter("TARGET_HOST", "Enter the target host to scan.", "", true)
	Parameters.Add_Parameter("START_PORT", "Enter the starting port number for scanning.", "", true)
	Parameters.Add_Parameter("END_PORT", "Enter the ending port number for scanning.", "", true)
}

// This function will start the attack.
func (F *Module) Module_Distributed_Nmap_Start() {

	startTime = time.Now()

	targetHost, _ = Parameters.Read_Parameter_Str("TARGET_HOST")
	startPort, _ = Parameters.Read_Parameter_Int("START_PORT")
	endPort, _ = Parameters.Read_Parameter_Int("END_PORT")

}

func (F *Module) Module_Distributed_Nmap_Server_Role() {

	var totalPorts int = (endPort - startPort) + 1

	Distribute("PORT_DISTRIBUTION", totalPorts)

	Run_Distributed_Task("PORT_DISTRIBUTION", "Module_Distributed_Nmap_Run_Nmap", true, true)

}

func (F *Module) Module_Distributed_Nmap_Client_Role() {

	Run_Distributed_Task("PORT_DISTRIBUTION", "Module_Distributed_Nmap_Run_Nmap", true, true)

}

func (F *Module) Module_Distributed_Nmap_Success_Indicator() string {
	return "Nmap done:"
}

func (F *Module) Module_Distributed_Nmap_Stop() {
	// Goodbye.
}

func (F *Module) Module_Distributed_Nmap_Run_Nmap(taskID int, taskName string, taskFrom, taskTo int) {
	taskFrom++
	taskTo++
	var nmapOutput string
	var portRange string = strconv.Itoa(taskFrom) + "-" + strconv.Itoa(taskTo)

	var enteredParamter string = "nmap " + targetHost + " -p " + portRange + " -v -A -sT -Pn"
	log.Println("girilen komut  =", enteredParamter)
	output, err := exec.Command("/bin/sh", "-c", enteredParamter).Output()
	if err != nil {
		Log.Create("Nmap could not be started: "+err.Error(), "Module", "HIGH")
		Log.PrintMsg("Nmap could not be started. Check your OS command.")
	}

	nmapOutput += string(output)

	if strings.Contains(nmapOutput, (&Module{}).Module_Distributed_Nmap_Success_Indicator()) {
		Parameters.AttackOutput += nmapOutput
		if Parameters.IsClient {
			Parameters.ShareOutput(0, taskName, nmapOutput) // Send scan results to the server.
		}
	}

	if Parameters.IsServer {
		fmt.Println("the attack took (in seconds) =", time.Since(startTime).Seconds())
	}

	// Update the state of the relevant task.
	Parameters.UpdateTask(taskID, "Succeeded")
}
