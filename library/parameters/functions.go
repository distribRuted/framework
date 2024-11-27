package Parameters

import (
	"errors"
	"net"
	"strconv"
	"strings"

	b64 "encoding/base64"

	Log "github.com/distribRuted/framework/library/log"
)

// add_parameter(PARAMETER_NAME, DESCRIPTION, DEFAULT_VALUE, IS_REQUIRED)
func Add_Parameter(parameter_name, desc string, default_value string, is_required bool) {
	if len(strings.Fields(parameter_name)) == 1 {

		var createParam Parameter = Parameter{Name: parameter_name, Description: desc, Value: default_value, Required: is_required}

		if CLIParameters != "" {
			var predefinedParams []string = strings.Split(CLIParameters, ",")
			for _, currentParam := range predefinedParams {
				parameter := strings.Split(currentParam, "=")
				if len(parameter) == 2 {
					var parameterName string = parameter[0]
					var parameterValue string = parameter[1]
					if parameterName == parameter_name {
						createParam = Parameter{Name: parameter_name, Description: desc, Value: parameterValue, Required: is_required}
					}
				}
			}
		}

		CollectParameters = append(CollectParameters, createParam)

	}
}

func Read_Parameter_Str(parameter_name string) (string, error) {
	for _, v := range CollectParameters {
		if v.Name == parameter_name {
			return v.Value, nil
		}
	}
	Log.Create("Parameter could not read. No such parameter: "+parameter_name, "Module", "MEDIUM")
	return "", errors.New("ERR:NO_SUCH_PARAMETER")
}

func Read_Parameter_Int(parameter_name string) (int, error) {
	for _, v := range CollectParameters {
		if v.Name == parameter_name {
			var intValue, err = strconv.Atoi(v.Value)
			if err != nil {
				Log.Create("["+parameter_name+"] parameter could not be converted to an integer. Try reading it as a string.", "Module", "MEDIUM")
				return -1, errors.New("ERR:TYPE_CONVERSION")
			}
			return intValue, nil
		}
	}
	Log.Create("Parameter could not read. No such parameter: "+parameter_name, "Module", "MEDIUM")
	return -1, errors.New("ERR:NO_SUCH_PARAMETER")
}

func SendMsgAsServer(destNode int, msg string) {
	if len(Connections) > destNode {
		if Connections[destNode].IsActive {

			Connections[destNode].NetConn.Write([]byte(string(msg + "\n")))
			Log.Create("Message sent to ["+Connections[destNode].DestAddr+"] client: ["+msg+"]", "Destination", "INFO")

		} else {
			Log.Create("Message could not sent to ["+ConnHost+"] address (client is offline): ["+msg+"]", "Node", "MEDIUM")
		}
	} else {
		Log.Create("Message could not sent to ["+ConnHost+"] address (no such client): ["+msg+"]", "Node", "LOW")
	}
}

func SendMsgAsClient(destNode int, msg string) {
	if len(Connections) > destNode {
		if Connections[destNode].IsActive {

			Connections[destNode].NetConn.Write([]byte(string(msg + "\n")))
			Log.Create("Message sent to ["+Connections[destNode].DestAddr+"] server: ["+msg+"]", "Destination", "INFO")

		} else {
			Log.Create("Message could not sent to ["+ConnHost+"] address (server is offline): ["+msg+"]", "Node", "MEDIUM")
		}
	} else {
		Log.Create("Message could not sent to ["+ConnHost+"] address (no such server): ["+msg+"]", "Node", "LOW")
	}
}

func UpdateTask(taskID int, newState string) {
	var isUpdated bool = false
	var taskIDStr string = strconv.Itoa(taskID)
	var taskName string
	if newState == "Queued" || newState == "Progressing" || newState == "Succeeded" || newState == "Failed" {
		for task_id, task := range Tasks {
			if taskID == task_id {
				Tasks[taskID] = Task{ID: taskID, Name: task.Name, From: task.From, To: task.To, State: newState}
				taskName = task.Name
				isUpdated = true
				Log.Create("The specified task ["+taskIDStr+":"+taskName+"] has been updated to the ["+newState+"] state.", "Module", "INFO")
				if IsServer && !IsThisDeviceExcluded && (newState == "Succeeded" || newState == "Failed") {
					for nodeID, node := range SharedTasks[taskName] {
						if node.NodeIP == "THIS_NODE" {
							SharedTasks[taskName][nodeID] = ToNode{NodeID: -1, NodeIP: "THIS_NODE", From: task.From, To: task.To, Completed: true}
						}
					}
				}
				return
			}
		}
	} else {
		Log.Create("The specified task ["+taskIDStr+":"+taskName+"] could not be updated. The new state can only be one of the following values: Queued / Progressing / Succeeded / Failed.", "Module", "HIGH")
		return
	}
	if !isUpdated {
		Log.Create("The specified task ["+taskIDStr+"] could not be found.", "Module", "HIGH")
	}
}

func ShareOutput(destNode int, task, output string) {
	output = b64.StdEncoding.EncodeToString([]byte(output))
	if IsClient {
		SendMsgAsClient(destNode, "PROTOCOL_MSG:SHARE_OUTPUT:"+task+":"+output)
	} else {
		SendMsgAsServer(destNode, "PROTOCOL_MSG:SHARE_OUTPUT:"+task+":"+output)
	}
}

func GetNodeKey(ipAddr string) int {
	if key, ok := NodeIndexes[ipAddr]; ok {
		return key
	}
	return -1
}

func GetOutboundIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func ShowCommand() string {
	var commandStr string = "./distribRuted --client "
	commandStr += "--module " + SelectedModuleName + " "
	commandStr += "--parameters "
	var parametersStr string = "\""
	for _, p := range CollectParameters {
		if p.Value != "" {
			parametersStr += p.Name + "=" + p.Value + ","
		}
	}
	if strings.HasSuffix(parametersStr, ",") {
		parametersStr = parametersStr[:len(parametersStr)-1]
	}
	parametersStr += "\""
	commandStr += parametersStr + " "
	if IsServer {
		commandStr += "--ip " + GetOutboundIP() + " --port " + ConnPort + " "
	} else {
		for _, conn := range Connections {
			if conn.IsActive {
				var parseConn []string = strings.Split(conn.DestAddr, ":")
				if len(parseConn) == 2 {
					var connIP string = parseConn[0]
					var connPort string = parseConn[1]
					commandStr += "--ip " + connIP + " --port " + connPort + " "
					break
				}
			}
		}
	}
	return commandStr
}
