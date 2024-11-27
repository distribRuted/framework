package Modules

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"

	Lib "github.com/distribRuted/framework/library"
	Log "github.com/distribRuted/framework/library/log"
	Parameters "github.com/distribRuted/framework/library/parameters"
	"github.com/fatih/color"
)

var F Module

type Module []string

func Load() {

	Parameters.ModuleNames = getModuleNames()
	if len(Parameters.ModuleNames) > 0 {
		for _, moduleName := range Parameters.ModuleNames {
			Parameters.AllModules[moduleName] = Parameters.Module{}
			defer func() {
				if err := recover(); err != nil {
					Log.Create("One of the module files could not parse correctly. Make sure your module contains the \"Init()\" function.", "Module", "CRITICAL")
					Lib.ExitWithError("One of the module files could not parse correctly. Make sure your module contains the \"Init()\" function.")
				}
			}()
			reflect.ValueOf(&F).MethodByName("Module_" + moduleName + "_Init").Call([]reflect.Value{})
		}
	} else {
		Log.Create("Could not load any module. Check \"./modules/\" directory.", "Application", "CRITICAL")
		Lib.ExitWithError("Could not load any module. Check \"./modules/\" directory.")
	}
}

func getModuleNames() []string {
	var moduleNames []string
	listDir, err := ioutil.ReadDir("./modules/")
	if err != nil {
		Log.Create("Could not read \"./modules/\" directory. Error message: "+err.Error(), "Application", "CRITICAL")
		Lib.ExitWithError("Could not read \"./modules/\" directory. Error message: " + err.Error())
	} else {
		for _, item := range listDir {
			if !item.IsDir() && strings.HasPrefix(item.Name(), "Module_") && strings.HasSuffix(item.Name(), ".go") {
				moduleNames = append(moduleNames, item.Name()[7:len(item.Name())-3])
			}
		}
	}
	return moduleNames
}

// set_module(MODULE_SHORT_NAME, MODULE_NAME, MODULE_DESCRIPTION, MODULE_AUTHOR)
func Set_Module(moduleShortName, moduleName, moduleDescription, moduleAuthor string) {
	addModule := new(Parameters.Module)
	addModule.Name = moduleName
	addModule.Description = moduleDescription
	addModule.Author = moduleAuthor
	if _, ok := Parameters.AllModules[moduleShortName]; ok {
		Parameters.AllModules[moduleShortName] = *addModule
	}
}

func ExecuteFunctions() {
	func() {
		defer func() {
			if err := recover(); err != nil {
				Log.Create("Module files could not parse correctly. Make sure your module contains \"Parameters()\" function.", "Module", "CRITICAL")
				Lib.ExitWithError("Module files could not parse correctly. Make sure your module contains \"Parameters()\" function.")
			}
		}()
		// Execute the Parameters() function of the module to collect user inputs.
		reflect.ValueOf(&F).MethodByName("Module_" + Parameters.ModuleNames[Parameters.SelectedModule-1] + "_Parameters").Call([]reflect.Value{})

	}()
}

func Distribute(taskName string, totalTasks int) {

	if Parameters.IsServer {

		if strings.TrimSpace(taskName) == "" {
			Log.Create("The \"taskName\" parameter of the function cannot be left blank. Distribution failed.", "Module", "HIGH")
			return
		}

		if totalTasks > 0 {

			for _, tName := range Parameters.Tasks {
				if tName.Name == taskName {
					Log.Create("The task ["+taskName+"] already exists. Please choose a new task name.", "Node", "INFO")
					return
				}
			}

			var distributionArr [][2]int = CalculateDistribution(totalTasks, Parameters.TotalActiveNodes)

			var latestNodeOrder, stopTaskInt, startTaskInt int = 0, 0, 0
			var startTask, stopTask, msgToSend string

			if !Parameters.IsThisDeviceExcluded {
				startTaskInt = distributionArr[latestNodeOrder][0]
				startTask = strconv.Itoa(startTaskInt)
				stopTaskInt = distributionArr[latestNodeOrder][1]
				stopTask = strconv.Itoa(stopTaskInt)

				if startTaskInt != -1 && stopTaskInt != -1 {
					// msgToSend = "PROTOCOL_MSG:" + "DISTRIBUTION:NAME=" + taskName + ";START=" + startTask + ";STOP=" + stopTask
					var createTask Parameters.Task = Parameters.Task{ID: Parameters.TotalTaskCount, Name: taskName, From: startTaskInt, To: stopTaskInt, State: "Queued"}
					Parameters.Tasks = append(Parameters.Tasks, createTask)
					Parameters.SharedTasks[taskName] = append(Parameters.SharedTasks[taskName], Parameters.ToNode{NodeID: -1, NodeIP: "THIS_NODE", From: 0, To: stopTaskInt, Completed: false})
					Parameters.TotalTaskCount++
					Log.Create("A task was created: [NAME:"+taskName+" - START:"+startTask+" - STOP:"+stopTask+"].", "Node", "INFO")

				}

				latestNodeOrder++
			}

			if len(Parameters.Connections) > 0 {
				for nodeID, connection := range Parameters.Connections {
					if connection.IsActive {
						var startTaskInt int = distributionArr[latestNodeOrder][0]
						startTask = strconv.Itoa(startTaskInt)
						var stopTaskInt int = distributionArr[latestNodeOrder][1]
						stopTask = strconv.Itoa(stopTaskInt)

						if startTaskInt != -1 && stopTaskInt != -1 {
							// TODO: Prevent the start command from being sent to servers that do not receive any tasks if the number of tasks is less than the number of servers.
							msgToSend = "PROTOCOL_MSG:" + "DISTRIBUTION:NAME=" + taskName + ";START=" + startTask + ";STOP=" + stopTask
							Parameters.SendMsgAsServer(nodeID, msgToSend)
							Parameters.SharedTasks[taskName] = append(Parameters.SharedTasks[taskName], Parameters.ToNode{NodeID: nodeID, NodeIP: connection.DestAddr, From: startTaskInt, To: stopTaskInt, Completed: false})
						}

						latestNodeOrder++
					}
				}
			}
			Log.Create("Distribution was performed among the nodes.", "Module", "LOW")
		} else {
			fmt.Println("No tasks to distribute were found. Did you set the required parameters?")
		}
	} else {
		Log.Create("Only nodes in the server role can run the \"Distribute()\" function.", "Module", "LOW")
	}
}

func CalculateDistribution(totalTasks, TotalActiveNodeCount int) [][2]int {
	var distribution [][2]int = make([][2]int, TotalActiveNodeCount)

	// We assign an equal number of tasks to each server and calculate the remaining tasks.
	var tasksPerServer int = totalTasks / TotalActiveNodeCount
	var remainder int = totalTasks % TotalActiveNodeCount

	var start int = 0

	for i := 0; i < TotalActiveNodeCount; i++ {
		var end int = start + tasksPerServer - 1

		if i < remainder {
			end++
		}

		if start <= end {
			distribution[i] = [2]int{start, end}
		} else {
			distribution[i] = [2]int{-1, -1}
		}

		start = end + 1
	}

	return distribution
}

func Run_Distributed_Task(task_name string, function_name string, run_async, loop bool) {
	for {
		if Parameters.IsAttackOngoing {
			if len(Parameters.Tasks) > 0 {
				for taskID, currentTask := range Parameters.Tasks {
					if currentTask.Name == task_name {
						if currentTask.State == "Queued" {
							Parameters.Tasks[taskID].State = "Progressing"
							Log.Create("Queued task was taken into processing: [NAME:"+currentTask.Name+" - START:"+strconv.Itoa(currentTask.From)+" - STOP:"+strconv.Itoa(currentTask.To)+"].", "Module", "INFO")
							if run_async {
								defer func() {
									if err := recover(); err != nil {
										Log.Create("Make sure your module contains the \" "+function_name+"()\" function.", "Module", "CRITICAL")
										Lib.ExitWithError("Make sure your module contains the \" " + function_name + "()\" function.")
									}
								}()
								go reflect.ValueOf(&F).MethodByName(function_name).Call([]reflect.Value{reflect.ValueOf(currentTask.ID), reflect.ValueOf(currentTask.Name), reflect.ValueOf(currentTask.From), reflect.ValueOf(currentTask.To)})
							} else {
								defer func() {
									if err := recover(); err != nil {
										Log.Create("Make sure your module contains the \" "+function_name+"()\" function.", "Module", "CRITICAL")
										Lib.ExitWithError("Make sure your module contains the \" " + function_name + "()\" function.")
									}
								}()
								reflect.ValueOf(&F).MethodByName(function_name).Call([]reflect.Value{})
							}
						}
					}
				}
			}
			if !loop {
				break
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func Start_Attack() {
	defer func() {
		if err := recover(); err != nil {
			Log.Create("Module files could not parse correctly. Make sure your module contains \"Start()\" function.", "Module", "CRITICAL")
			Lib.ExitWithError("Module files could not parse correctly. Make sure your module contains \"Start()\" function.")
		}
	}()
	// Execute the Start() function of the module to start the attack.
	reflect.ValueOf(&F).MethodByName("Module_" + Parameters.ModuleNames[Parameters.SelectedModule-1] + "_Start").Call([]reflect.Value{})

	if !Parameters.IsThisDeviceExcluded {
		// Next, execute the Server_Role() function of the module if the application is started in the server role.
		if Parameters.IsServer {
			go Call_Server_Role_Function()
		}

		// Next, execute the Client_Role() function of the module if the application is started in the client role.
		if Parameters.IsClient {
			go Call_Client_Role_Function()
		}
	}
}

func Stop_Attack() {
	if Parameters.IsAttackStarted {
		if Parameters.IsAttackOngoing || Parameters.IsAttackPaused {
			if Parameters.IsServer {
				for {
					if Parameters.AllNodesCompleted || Parameters.TotalActiveNodes == 1 {
						break
					}
					time.Sleep(5 * time.Second)
				}
			}
			Log.Create("The attack has been stopped.", "Application", "HIGH")
			Log.PrintMsg(color.YellowString("The attack has been stopped."))
			Parameters.IsAttackStopped = true
			Parameters.IsAttackStarted = false
			Parameters.IsAttackOngoing = false
			Parameters.IsAttackPaused = false
			defer func() {
				if err := recover(); err != nil {
					Log.Create("Module files could not parse correctly. Make sure your module contains \"Stop()\" function.", "Module", "CRITICAL")
					Lib.ExitWithError("Module files could not parse correctly. Make sure your module contains \"Stop()\" function.")
				}
			}()
			// Execute the Start() function of the module to start the attack.
			reflect.ValueOf(&F).MethodByName("Module_" + Parameters.ModuleNames[Parameters.SelectedModule-1] + "_Stop").Call([]reflect.Value{})
		}
	}
}

func Stop_Attack_On_All_Nodes() {
	if Parameters.IsAttackStarted {
		if Parameters.IsAttackOngoing || Parameters.IsAttackPaused {
			if Parameters.IsServer {
				if len(Parameters.Connections) > 0 {
					go func() {
						for nodeID, connection := range Parameters.Connections {
							if connection.IsActive {
								Parameters.SendMsgAsServer(nodeID, "PROTOCOL_MSG:stop")
							}
						}
					}()
					Log.Create("The attack has been stopped on all connected clients and this server.", "Application", "HIGH")
					Log.PrintMsg(color.YellowString("The attack has been stopped on all connected clients and this server."))
				} else {
					Log.Create("Since there are no clients connected to this server, the attack has been stopped only on this server.", "Node", "HIGH")
					Log.PrintMsg(color.YellowString("Since there are no clients connected to this server, the attack has been stopped only on this server."))
				}
			} else {
				Log.Create("The attack has been stopped.", "Application", "HIGH")
				Log.PrintMsg(color.YellowString("The attack has been stopped."))
			}
			go Stop_Attack()
		} else {
			Log.PrintMsg(color.YellowString("There is no ongoing or paused attack. First, start a new one."))
		}
	} else {
		Log.PrintMsg(color.YellowString("No attack has been started yet."))
	}
}

func Call_Server_Role_Function() {
	defer func() {
		if err := recover(); err != nil {
			Log.Create("Module files could not parse correctly. Make sure your module contains \"Server_Role()\" function.", "Module", "INFO")
			Lib.ExitWithError("Module files could not parse correctly. Make sure your module contains \"Server_Role()\" function.")
		}
	}()
	if Parameters.IsServer {
		reflect.ValueOf(&F).MethodByName("Module_" + Parameters.ModuleNames[Parameters.SelectedModule-1] + "_Server_Role").Call([]reflect.Value{})
	}
}

func Call_Client_Role_Function() {
	defer func() {
		if err := recover(); err != nil {
			Log.Create("Module files could not parse correctly. Make sure your module contains \"Client_Role()\" function.", "Module", "CRITICAL")
			Lib.ExitWithError("Module files could not parse correctly. Make sure your module contains \"Client_Role()\" function.")
		}
	}()
	if Parameters.IsClient {
		reflect.ValueOf(&F).MethodByName("Module_" + Parameters.ModuleNames[Parameters.SelectedModule-1] + "_Client_Role").Call([]reflect.Value{})
	}
}

func Call_Custom_Function(function_name string) {
	defer func() {
		if err := recover(); err != nil {
			Log.Create("Module files could not parse correctly. Make sure your module contains \""+function_name+"()\" function.", "Module", "CRITICAL")
			Lib.ExitWithError("Module files could not parse correctly. Make sure your module contains \"" + function_name + "()\" function.")
		}
	}()
	// Execute the specified function of the module.
	reflect.ValueOf(&F).MethodByName(function_name).Call([]reflect.Value{})
}
