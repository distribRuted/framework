package Modules

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	Log "github.com/distribRuted/framework/library/log"
	Parameters "github.com/distribRuted/framework/library/parameters"
)

func (F *Module) Module_Anonymous_FTP_Scanner_Init() {

	var (
		moduleShortName   string = "Anonymous_FTP_Scanner"
		moduleName        string = "Anonymous FTP Scanner"
		moduleDescription string = "Distributed Nmap demonstration of anonymous FTP scanning for Black Hat"
		moduleAuthor      string = "@numanozdemircom <root@numanozdemir.com>"
	)

	Set_Module(moduleShortName, moduleName, moduleDescription, moduleAuthor)
}

var (
	targetHosts    string
	startTime_AFTP time.Time
	IPrange        []string
)

func (F *Module) Module_Anonymous_FTP_Scanner_Parameters() {
	Parameters.Add_Parameter("TARGET_HOSTS", "Enter the target hosts to be scanned, separated by space.", "", false)
}

func (F *Module) Module_Anonymous_FTP_Scanner_Start() {

	startTime_AFTP = time.Now()
	targetHosts, _ = Parameters.Read_Parameter_Str("TARGET_HOSTS")
	IPrange = strings.Split(targetHosts, " ")
}

func (F *Module) Module_Anonymous_FTP_Scanner_Server_Role() {
	var totalIPs int = len(IPrange)
	Distribute("IP_COUNT_DISTRIBUTION_AFTP", totalIPs)
	Run_Distributed_Task("IP_COUNT_DISTRIBUTION_AFTP", "Module_Anonymous_FTP_Scanner_Run_Nmap", true, true)
}

func (F *Module) Module_Anonymous_FTP_Scanner_Client_Role() {
	Run_Distributed_Task("IP_COUNT_DISTRIBUTION_AFTP", "Module_Anonymous_FTP_Scanner_Run_Nmap", true, true)
}

func (F *Module) Module_Anonymous_FTP_Scanner_Success_Indicator() string {
	return "Anonymous FTP login allowed"
}

func (F *Module) Module_Anonymous_FTP_Scanner_Stop() {
	// Terminate cloud instances, create a scan report, send notification, etc.
	fmt.Println("The attack on the server side is complete.")
}

func (F *Module) Module_Anonymous_FTP_Scanner_Run_Nmap(taskID int, taskName string, taskFrom, taskTo int) {
	var nmapOutput string
	var ipRange string = strings.Join(IPrange[taskFrom:taskTo+1], " ")

	var enteredParamter string = "nmap " + ipRange + " --script=ftp-anon -p 21 -Pn"
	log.Println("The command to be executed on this node =  ", enteredParamter)
	output, err := exec.Command("/bin/sh", "-c", enteredParamter).Output()
	if err != nil {
		Log.Create("Nmap could not be started: "+err.Error(), "Module", "HIGH")
		Log.PrintMsg("Nmap could not be started. Check your OS command.")
	}

	nmapOutput += string(output)

	if !strings.Contains(nmapOutput, (&Module{}).Module_Anonymous_FTP_Scanner_Success_Indicator()) {
		nmapOutput = "Anonymous FTP access could not be established for " + ipRange
	}

	Parameters.AttackOutput += nmapOutput

	if Parameters.IsClient {
		Parameters.ShareOutput(0, taskName, nmapOutput) // Send scan results to the server.
		fmt.Println("The attack on the client side is complete (" + strconv.FormatFloat(time.Since(startTime_AFTP).Seconds(), 'f', 3, 64) + "s)")
	}

	if Parameters.IsServer {
		fmt.Println("The attack on the server side is complete (" + strconv.FormatFloat(time.Since(startTime_AFTP).Seconds(), 'f', 3, 64) + "s)")
	}

	Parameters.UpdateTask(taskID, "Succeeded")
}

func ParseNmapReports(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var result []string
	var currentReport []string
	foundAnonFTP := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Nmap scan report for") {
			if foundAnonFTP {
				result = append(result, strings.Join(currentReport, "\n"))
			}
			currentReport = []string{}
			foundAnonFTP = false
		}

		currentReport = append(currentReport, line)

		if strings.Contains(line, "Anonymous FTP login allowed") {
			foundAnonFTP = true
		}
	}

	if foundAnonFTP {
		result = append(result, strings.Join(currentReport, "\n"))
	}

	return strings.Join(result, "\n\n")
}
