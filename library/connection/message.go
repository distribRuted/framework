package Connection

import (
	"strconv"
	"strings"
	"time"

	b64 "encoding/base64"

	Log "github.com/distribRuted/framework/library/log"
	Parameters "github.com/distribRuted/framework/library/parameters"
	Modules "github.com/distribRuted/framework/modules"
	"github.com/fatih/color"
)

func evaluateMessage(isServer, isClient bool, sourceIP, message string) {
	if len(message) > 11 && message[:11] == "CUSTOM_MSG:" {
		// var getMessage string = message[11:]

		// Create a log for the incoming custom message
		if isServer {
			Log.Create("Message from ["+sourceIP+"] client: ["+message+"]", "Node", "LOW")
		} else {
			Log.Create("Message from ["+sourceIP+"] server: ["+message+"]", "Node", "LOW")
		}

	} else if len(message) > 13 && message[:13] == "PROTOCOL_MSG:" {
		var getMessage string = message[13:]

		switch getMessage {
		case "disconnect":
			for key, conn := range Parameters.Connections {
				if conn.DestAddr == sourceIP && conn.IsActive {
					Parameters.Connections[key].IsActive = false
					Parameters.TotalActiveNodes -= -1
					if isServer {
						Log.Create("The connection of ["+sourceIP+"] was terminated by the client.", "Node", "CRITICAL")
					} else {
						Log.Create("The connection of ["+sourceIP+"] was terminated by the server.", "Node", "CRITICAL")
					}
					break
				}
			}
		case "start":
			if isClient {
				var canStartAttack bool = true
				for _, parameter := range Parameters.CollectParameters {
					if parameter.Required && parameter.Value == "" {
						canStartAttack = false
						var logMsg string = "Required module parameters cannot be left blank. The attack could not be started."
						Log.Create(logMsg, "Module", "HIGH")
						Log.PrintMsg(color.YellowString(logMsg))
						break
					}
				}

				if canStartAttack {
					if !Parameters.IsAttackStarted {
						Log.Create("The attack has been started by the command of the server: ["+sourceIP+"]", "Application", "HIGH")
						Parameters.IsAttackStarted = true
						Parameters.IsAttackOngoing = true
						Parameters.IsAttackStopped = false
						Parameters.IsAttackPaused = false
						Parameters.AllNodesCompleted = false
						go Modules.Start_Attack()
					}
				}
			} else {
				Log.Create("The attack couldn't be started. The client ["+sourceIP+"] doesn't have the authority to command this server.", "Node", "LOW")
			}
		case "stop":
			// Node verification has intentionally not been implemented here.
			// To prevent the application from being used as a botnet, clients are allowed to send commands to servers to stop attacks.
			if Parameters.IsAttackStarted && !Parameters.IsAttackStopped {
				Log.Create("The attack has been stopped by the command of the server: ["+sourceIP+"]", "Application", "HIGH")
				Parameters.IsAttackStopped = true
				Parameters.IsAttackStarted = false
				Parameters.IsAttackOngoing = false
				Parameters.IsAttackPaused = false
				go Modules.Stop_Attack()
			}
		case "pause":
			// Node verification has intentionally not been implemented here.
			// To prevent the application from being used as a botnet, clients are allowed to send commands to servers to pause attacks.
			if Parameters.IsAttackOngoing && !Parameters.IsAttackPaused {
				Log.Create("The attack has been paused by the command of the server: ["+sourceIP+"]", "Application", "HIGH")
				Parameters.IsAttackStarted = true
				Parameters.IsAttackOngoing = false
				Parameters.IsAttackStopped = false
				Parameters.IsAttackPaused = true
			}
		case "continue":
			// Node verification has intentionally not been implemented here.
			// To prevent the application from being used as a botnet, clients are allowed to send commands to servers to continue attacks.
			if !Parameters.IsAttackOngoing && Parameters.IsAttackPaused {
				Log.Create("The attack has been started to continue by the command of the server: ["+sourceIP+"]", "Application", "HIGH")
				Parameters.IsAttackStarted = true
				Parameters.IsAttackOngoing = true
				Parameters.IsAttackStopped = false
				Parameters.IsAttackPaused = false
			}
		default:
			if getMessage[:13] == "DISTRIBUTION:" {
				if Parameters.IsClient {
					getMessage = getMessage[13:]
					var parseMsg []string = strings.Split(getMessage, ";")
					var taskNames []string = strings.Split(parseMsg[0], "=")
					var tasksFrom []string = strings.Split(parseMsg[1], "=")
					var tasksTo []string = strings.Split(parseMsg[2], "=")

					var taskName, taskFromStr, taskToStr string

					if len(taskNames) == 2 && len(tasksFrom) == 2 && len(tasksTo) == 2 {
						if taskNames[0] == "NAME" && tasksFrom[0] == "START" && tasksTo[0] == "STOP" {
							taskName = taskNames[1]
							taskFromStr = tasksFrom[1]
							taskToStr = tasksTo[1]
							if strings.TrimSpace(taskName) == "" {
								return
							}
						} else {
							return
						}
					} else {
						return
					}

					taskFrom, err := strconv.Atoi(taskFromStr)
					if err != nil {
						return
					}
					taskTo, err := strconv.Atoi(taskToStr)
					if err != nil {
						return
					}

					var createTask Parameters.Task = Parameters.Task{ID: Parameters.TotalTaskCount, Name: taskName, From: taskFrom, To: taskTo, State: "Queued"}
					Parameters.Tasks = append(Parameters.Tasks, createTask)
					Parameters.TotalTaskCount++
					Log.Create("The server ["+sourceIP+"] distributed a task: [NAME:"+taskName+" - START:"+taskFromStr+" - STOP:"+taskToStr+"].", "Node", "INFO")
				}
			} else if getMessage[:13] == "SHARE_OUTPUT:" {
				getMessage = getMessage[13:]
				var anyNodeLeft bool = false
				var isThisNodeCompleted bool = false
				var getTaskName string
				var splitMsg []string = strings.Split(getMessage, ":")
				if len(splitMsg) == 2 {
					getTaskName = splitMsg[0]
					getMessage = splitMsg[1]
				}

				if Parameters.IsServer {
					if node, ok := Parameters.SharedTasks[getTaskName]; ok {
						for nodeID, node := range node {
							if node.NodeIP != "THIS_NODE" {
								if node.NodeIP == sourceIP {
									Parameters.SharedTasks[getTaskName][nodeID] = Parameters.ToNode{NodeID: node.NodeID, NodeIP: node.NodeIP, From: node.From, To: node.To, Completed: true}
								}
								var isCompleted bool = Parameters.SharedTasks[getTaskName][nodeID].Completed
								if !isCompleted {
									anyNodeLeft = true
								}
							} else {
								if Parameters.SharedTasks[getTaskName][nodeID].Completed {
									isThisNodeCompleted = true
								}
							}
						}
					}

					if !anyNodeLeft {
						var logMsg string
						if isThisNodeCompleted {
							logMsg = "The attack has been completed on all connected clients and this server."
						} else {
							logMsg = "The attack has been completed on all connected clients."
						}
						Parameters.AllNodesCompleted = true
						Log.Create(logMsg, "Node", "HIGH")
						Log.PrintMsg(color.YellowString(logMsg))
					}
				}

				decodedOutput, _ := b64.StdEncoding.DecodeString(getMessage)
				getMessage = "[ID: " + strconv.Itoa(Parameters.GetNodeKey(sourceIP)) + " - NODE: " + sourceIP + "] Output: \n\n" + string(decodedOutput) + "\n" + time.Now().Format("02/01/2006 15:04:05") + "\n------------------------------" + "\n"
				Parameters.AttackOutput += getMessage
				Log.OutputToFile(sourceIP, time.Now().Format("02_01_2006_15_04_05"), getMessage)
				// TODO: Instead of grouping all incoming values into a single variable called IncomingOutput, consider designing a more comprehensive structure.
			} else {
				Log.Create("Unsupported protocol message from ["+sourceIP+"] address: ["+message+"]", "Node", "MEDIUM")
			}
		}

	} else {
		// If the incoming message is undefined
		Log.Create("Unsupported message from ["+sourceIP+"] address: ["+message+"]", "Node", "MEDIUM")
	}
}
