package Console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	Lib "github.com/distribRuted/framework/library"
	Connection "github.com/distribRuted/framework/library/connection"
	Log "github.com/distribRuted/framework/library/log"
	Parameters "github.com/distribRuted/framework/library/parameters"
	Modules "github.com/distribRuted/framework/modules"
	"github.com/fatih/color"
)

func Open() {
	for {
		var userCommand string = getUserCommand()
		evaluateCommand(userCommand)
	}
}

func getUserCommand() string {
	color.Set(color.FgHiBlue)
	fmt.Print("$ ")
	reader := bufio.NewReader(os.Stdin)
	variable, _ := reader.ReadString('\n')
	color.Unset()
	if len(variable) > 0 {
		variable = variable[:len(variable)-1]
	}
	return variable
}

func evaluateCommand(userCommand string) string {
	var words []string = strings.Fields(userCommand)

	if len(words) > 0 {
		var yellowStr func(a ...interface{}) string = color.New(color.FgYellow).Add(color.Underline).SprintFunc()

		if words[0] == "set" {
			// COMMAND = set [PARAMETER] [VALUE]
			if len(words) >= 2 {
				var isParamExist bool = false
				for k, v := range Parameters.CollectParameters {
					if v.Name == words[1] {
						isParamExist = true
						Parameters.CollectParameters[k].Value = strings.Join(words[2:], " ")
					}
				}
				if !isParamExist {
					fmt.Println("No such parameter. You can use the 'show parameters' command to list all available parameters.")
				}
			}

		} else if words[0] == "connect" {
			var isValid bool = true

			if len(words) >= 3 {
				if len(words) == 3 {
					// COMMAND = connect [IP] [PORT]
					if Parameters.IsClient {
						if !Connection.IsIPAddrValid(words[1]) {
							isValid = false
							fmt.Println("The given IP address is invalid.")
						}
						port, _ := strconv.Atoi(words[2])
						if port > 65535 || port < 0 {
							isValid = false
							fmt.Println("The given port number is invalid.")
						}
						if isValid {
							Parameters.ConnHost = words[1]
							Parameters.ConnPort = words[2]
							go Connection.Dial()
						}

					} else {
						fmt.Println("To connect to a server, you must start the application in the client role.\nThe server (you) can only listen for incoming connections.")
					}
				} else {
					var msgToSend string = "CUSTOM_MSG:" + strings.Join(words[3:], " ")
					if words[1] == "all" {
						// COMMAND = connect all send [CMD]
						if len(Parameters.Connections) > 0 {
							for nodeID, connection := range Parameters.Connections {
								if connection.IsActive {
									if Parameters.IsServer {
										Parameters.SendMsgAsServer(nodeID, msgToSend)
									} else {
										Parameters.SendMsgAsClient(nodeID, msgToSend)
									}
								}
							}

						} else {
							fmt.Println("Message could not be sent as no connected nodes were found.")
						}

					} else {
						// COMMAND = connect [NODE_ID] send [CMD]
						key, _ := strconv.Atoi(words[1])
						if len(Parameters.Connections) > key {
							if Parameters.Connections[key].IsActive {
								key, _ := strconv.Atoi(words[1])
								if Parameters.IsServer {
									Parameters.SendMsgAsServer(key, msgToSend)
								} else {
									Parameters.SendMsgAsClient(key, msgToSend)
								}
							} else {
								fmt.Println("Message could not be sent as the connection to the target node has been lost.")
							}
						} else {
							fmt.Println("No such node.")
						}

					}
				}
			} else {
				fmt.Println("Missing parameter. You can use the 'help' command to list all available commands.")
			}

		} else if words[0] == "disconnect" {
			if len(words) == 2 && words[1] == "all" {
				// COMMAND = disconnect all
				for key := range Parameters.Connections {
					if Parameters.Connections[key].IsActive {
						if Parameters.IsServer {
							Parameters.SendMsgAsServer(key, "PROTOCOL_MSG:disconnect")
							Log.Create("The connection of ["+Parameters.Connections[key].DestAddr+"] was terminated by the server.", "Node", "CRITICAL")
						} else {
							Parameters.SendMsgAsClient(key, "PROTOCOL_MSG:disconnect")
							Log.Create("The connection of ["+Parameters.Connections[key].DestAddr+"] was terminated by the client.", "Node", "CRITICAL")
						}
						Parameters.Connections[key].IsActive = false
					}
				}

			} else if len(words) == 2 {
				// COMMAND = disconnect [NODE_ID]
				key, _ := strconv.Atoi(words[1])
				if len(Parameters.Connections) > key {
					if Parameters.Connections[key].IsActive {
						if Parameters.IsServer {
							Parameters.SendMsgAsServer(key, "PROTOCOL_MSG:disconnect")
							Log.Create("The connection of ["+Parameters.Connections[key].DestAddr+"] was terminated by the server.", "Node", "CRITICAL")
						} else {
							Parameters.SendMsgAsClient(key, "PROTOCOL_MSG:disconnect")
							Log.Create("The connection of ["+Parameters.Connections[key].DestAddr+"] was terminated by the client.", "Node", "CRITICAL")
						}
						Parameters.Connections[key].IsActive = false
					}
				} else {
					fmt.Println("No such node.")
				}
			}
		} else if words[0] == "deploy" && len(words) == 2 {
			server_count, err := strconv.Atoi(words[1])
			if err == nil {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("An error occurred: ", r)
					}
				}()
				// Connection.DO_Servers(server_count)
				Connection.AWS_Servers(server_count)
			}
		} else {
			switch userCommand {
			case "help":
				fmt.Printf("%-44s%s \n", yellowStr("COMMAND"), yellowStr("DESCRIPTION"))
				fmt.Printf("%-30s%s \n", "show parameters", "List all the parameters of the module.")
				fmt.Printf("%-30s%s \n", "show nodes", "List all nodes connected to this server.")
				fmt.Printf("%-30s%s \n", "show command", "Show the CLI command to be executed on other nodes to start the attack.")
				fmt.Printf("%-30s%s \n", "show logs", "Show application, network and attack logs.")
				fmt.Printf("%-30s%s \n", "show output", "Show the attack output.")
				fmt.Printf("%-30s%s \n", "deploy [COUNT]", "Deploy [COUNT] clients and connect them to this server.")
				fmt.Printf("%-30s%s \n", "connect [IP] [PORT]", "Connect to the server with the specified IP:PORT address, as a client.")
				fmt.Printf("%-30s%s \n", "connect [NODE_ID] send [CMD]", "Send a command to the node with the specified IP:PORT address.")
				fmt.Printf("%-30s%s \n", "connect all send [CMD]", "Send a command to the all connected nodes.")
				fmt.Printf("%-30s%s \n", "disconnect [NODE_ID]", "Disconnect from the selected node.")
				fmt.Printf("%-30s%s \n", "disconnect all", "Disconnect from all nodes.")
				fmt.Printf("%-30s%s \n", "set [PARAMETER] [VALUE]", "Assigns [VALUE] to the [PARAMETER].")
				fmt.Printf("%-30s%s \n", "exclude", "Exclude this device from the attack.")
				fmt.Printf("%-30s%s \n", "include", "Include this device in the attack again.")
				fmt.Printf("%-30s%s \n", "status", "Show the status of the node and the attack.")
				fmt.Printf("%-30s%s \n", "start", "Start a simultaneous attack on all nodes.")
				fmt.Printf("%-30s%s \n", "pause", "Pause the simultaneous attack on all nodes.")
				fmt.Printf("%-30s%s \n", "continue", "Continue the paused simultaneous attack on all nodes.")
				fmt.Printf("%-30s%s \n", "stop", "Stop the simultaneous attack on all nodes.")
				fmt.Printf("%-30s%s \n", "help", "List all console commands.")
				fmt.Printf("%-30s%s \n", "exit", "Terminate the program.")
			case "show parameters":
				fmt.Printf("%-44s%-44v%s \n", yellowStr("PARAMETER"), yellowStr("VALUE"), yellowStr("DESCRIPTION"))
				for _, p := range Parameters.CollectParameters {
					fmt.Printf("%-30s%-30s%s \n", p.Name, p.Value, p.Description)
					fmt.Println(strings.Repeat("-", 150))
				}
			case "show nodes":
				fmt.Printf("%-44s%-44v%s \n", yellowStr("ID"), yellowStr("DESTINATION"), yellowStr("IS ACTIVE"))
				for k, n := range Parameters.Connections {
					fmt.Printf("%-30d%-30s%t \n", k, n.DestAddr, n.IsActive)
				}
			case "show command":
				fmt.Println(Parameters.ShowCommand())
			case "show logs":
				for _, getLog := range Log.AllLogs {
					var epochToDate time.Time = time.Unix(getLog.Epoch, 0)
					var currentDate string = epochToDate.Format("02/01/2006 15:04:05")
					var printLog string
					switch getLog.Level {
					case "INFO":
						printLog = color.CyanString("[INFO] " + getLog.Log)
					case "LOW":
						printLog = color.MagentaString("[LOW] " + getLog.Log)
					case "MEDIUM":
						printLog = color.BlueString("[MEDIUM] " + getLog.Log)
					case "HIGH":
						printLog = color.YellowString("[HIGH] " + getLog.Log)
					case "CRITICAL":
						printLog = color.RedString("[CRITICAL] " + getLog.Log)
					default:
						printLog = getLog.Log
					}
					printSource := color.New(color.FgWhite).SprintFunc()
					fmt.Println(currentDate, printSource("["+getLog.Source+"]"), printLog)
				}
			case "show output":
				fmt.Println(color.HiWhiteString(Parameters.AttackOutput))
			case "status":
				if Parameters.IsServer {
					fmt.Println("TotalActiveNodes =", Parameters.TotalActiveNodes)
					fmt.Println("AllNodesCompleted =", Parameters.AllNodesCompleted)
				}
				fmt.Println("IsThisDeviceExcluded =", Parameters.IsThisDeviceExcluded)
				fmt.Println("IsAttackStarted =", Parameters.IsAttackStarted)
				fmt.Println("IsAttackOngoing =", Parameters.IsAttackOngoing)
				fmt.Println("IsAttackPaused =", Parameters.IsAttackPaused)
				fmt.Println("IsAttackStopped =", Parameters.IsAttackStopped)
			case "include":
				if !Parameters.IsClient {
					if Parameters.IsThisDeviceExcluded {
						Parameters.TotalNodes += 1
						Parameters.TotalActiveNodes += 1
					}
					Parameters.IsThisDeviceExcluded = false
				} else {
					Log.Create("\"include\" command can only be executed if the application is started in the server role.", "Application", "INFO")
					Log.PrintMsg(color.YellowString("\"include\" command can only be executed if the application is started in the server role."))
				}
			case "exclude":
				if !Parameters.IsClient {
					if !Parameters.IsThisDeviceExcluded {
						Parameters.TotalNodes -= 1
						Parameters.TotalActiveNodes -= 1
					}
					Parameters.IsThisDeviceExcluded = true
				} else {
					Log.Create("\"exclude\" command can only be executed if the application is started in the server role.", "Application", "INFO")
					Log.PrintMsg(color.YellowString("\"exclude\" command can only be executed if the application is started in the server role."))
				}
			case "start":
				if (!Parameters.IsAttackOngoing && !Parameters.IsAttackPaused) || Parameters.IsAttackStopped {

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

						if Parameters.IsServer {
							if len(Parameters.Connections) > 0 {
								//go func() {
								for nodeID, connection := range Parameters.Connections {
									if connection.IsActive {
										Parameters.SendMsgAsServer(nodeID, "PROTOCOL_MSG:start")
									}
								}
								//}()
								var logMsg string
								if Parameters.IsThisDeviceExcluded {
									logMsg = "The attack has been started on all connected clients, but it has not been started on this server."
								} else {
									logMsg = "The attack has been started on all connected clients and this server."
								}
								Log.Create(logMsg, "Application", "HIGH")
								Log.PrintMsg(color.YellowString(logMsg))
							} else {
								var logMsg string
								if Parameters.IsClient {
									logMsg = "Since there are no clients connected to this server, the attack has been started only on this server in both client and server roles."
								} else {
									logMsg = "Since there are no clients connected to this server, the attack has been started only on this server."
								}
								Log.Create(logMsg, "Node", "HIGH")
								Log.PrintMsg(color.YellowString(logMsg))
							}
						} else {
							Log.Create("The attack has been started.", "Application", "HIGH")
							Log.PrintMsg(color.YellowString("The attack has been started."))
						}
						Parameters.IsAttackStarted = true
						Parameters.IsAttackOngoing = true
						Parameters.IsAttackStopped = false
						Parameters.IsAttackPaused = false
						Parameters.AllNodesCompleted = false
						go Modules.Start_Attack()
					}
				} else {
					Log.PrintMsg(color.YellowString("There is already an ongoing attack still in progress."))
				}
			case "continue":
				if Parameters.IsAttackStarted {
					if !Parameters.IsAttackOngoing && Parameters.IsAttackPaused && !Parameters.IsAttackStopped {
						if Parameters.IsServer {
							if len(Parameters.Connections) > 0 {
								go func() {
									for nodeID, connection := range Parameters.Connections {
										if connection.IsActive {
											Parameters.SendMsgAsServer(nodeID, "PROTOCOL_MSG:continue")
										}
									}
								}()
								Log.Create("This server and all connected clients have started to continue the attack.", "Node", "HIGH")
								Log.PrintMsg(color.YellowString("This server and all connected clients have started to continue the attack."))
							} else {

								Log.Create("Since there are no clients connected to this server, only this server has started to continue the attack.", "Node", "HIGH")
								Log.PrintMsg(color.YellowString("Since there are no clients connected to this server, only this server has started to continue the attack."))

							}
						} else {
							Log.Create("The attack has started to continue.", "Node", "HIGH")
							Log.PrintMsg(color.YellowString("The attack has started to continue."))
						}

						Parameters.IsAttackStarted = true
						Parameters.IsAttackOngoing = true
						Parameters.IsAttackStopped = false
						Parameters.IsAttackPaused = false
					} else {
						Log.PrintMsg(color.YellowString("There is no paused attack."))
					}
				} else {
					Log.PrintMsg(color.YellowString("No attack has been started yet."))
				}
			case "pause":
				if Parameters.IsAttackStarted {
					if !Parameters.IsAttackPaused {
						if Parameters.IsServer {
							if len(Parameters.Connections) > 0 {
								go func() {
									for nodeID, connection := range Parameters.Connections {
										if connection.IsActive {
											Parameters.SendMsgAsServer(nodeID, "PROTOCOL_MSG:pause")
										}
									}
								}()
								Log.Create("The attack has been paused on all connected clients and this server.", "Application", "HIGH")
								Log.PrintMsg(color.YellowString("The attack has been paused on all connected clients and this server."))
							} else {
								Log.Create("Since there are no clients connected to this server, the attack has been paused only on this server.", "Node", "HIGH")
								Log.PrintMsg(color.YellowString("Since there are no clients connected to this server, the attack has been paused only on this server."))
							}
						} else {
							Log.Create("The attack has been paused.", "Application", "HIGH")
							Log.PrintMsg(color.YellowString("The attack has been paused."))
						}

						Parameters.IsAttackStarted = true
						Parameters.IsAttackOngoing = false
						Parameters.IsAttackStopped = false
						Parameters.IsAttackPaused = true
					} else {
						Log.PrintMsg(color.YellowString("The attack has already been paused."))
					}
				} else {
					Log.PrintMsg(color.YellowString("No attack has been started yet."))
				}
			case "stop":
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

						Parameters.IsAttackStopped = true
						Parameters.IsAttackStarted = false
						Parameters.IsAttackOngoing = false
						Parameters.IsAttackPaused = false
						go Modules.Stop_Attack()
					} else {
						Log.PrintMsg(color.YellowString("There is no ongoing or paused attack. First, start a new one."))
					}
				} else {
					Log.PrintMsg(color.YellowString("No attack has been started yet."))
				}
			case "exit":
				Lib.Exit()

			default:
				fmt.Println("No such command. You can use the 'help' command to list all available commands.")
			}
		}
	}
	return userCommand
}
