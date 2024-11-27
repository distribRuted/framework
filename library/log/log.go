package Log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	AllLogs []Logs
)

type Logs struct {
	Log    string // Any message
	Source string // Application / Destination / Node / Module / OS ...
	Level  string // INFO / LOW / MEDIUM / HIGH / CRITICAL
	Epoch  int64  // Timestamp
}

func Create(log, source string, level string) {
	var lcSource string = strings.ToLower(source)
	var lcLevel string = strings.ToLower(level)
	if (lcSource == "application" || lcSource == "destination" || lcSource == "node" || lcSource == "module" || lcSource == "os") &&
		(lcLevel == "info" || lcLevel == "low" || lcLevel == "medium" || lcLevel == "high" || lcLevel == "critical") {
		var createLog Logs = Logs{Log: log, Source: source, Level: level, Epoch: time.Now().Unix()}
		AllLogs = append(AllLogs, createLog)
	}
}

func PrintMsg(message string) {
	var epochToDate time.Time = time.Now()
	var currentDate string = epochToDate.Format("02/01/2006 15:04:05")
	var printLog string = message
	fmt.Println(currentDate, printLog)
}

func OutputToFile(dir_name, file_name, file_content string) {
	var outputsDir string = "outputs"
	if _, err := os.Stat(outputsDir); os.IsNotExist(err) {
		err := os.Mkdir(outputsDir, os.ModePerm)
		if err != nil {
			// fmt.Println("Error while creating directory:", err)
			return
		}
	}

	var nodeDir string = filepath.Join(outputsDir, dir_name)
	if _, err := os.Stat(nodeDir); os.IsNotExist(err) {
		err := os.Mkdir(nodeDir, os.ModePerm)
		if err != nil {
			// fmt.Println("Error while creating directory:", err)
			return
		}
	}

	var fileName string = time.Now().Format("02_01_2006_15_04_05") + ".txt"
	var filePath string = filepath.Join(nodeDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		// fmt.Println("Error while creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(file_content)
	if err != nil {
		// fmt.Println("Error while writing to the file:", err)
		return
	}
}

func ScanOutputToFile(scan_output string) {
	var outputsDir string = "outputs"
	if _, err := os.Stat(outputsDir); os.IsNotExist(err) {
		err := os.Mkdir(outputsDir, os.ModePerm)
		if err != nil {
			// fmt.Println("Error while creating directory:", err)
			return
		}
	}

	var logFileName string = "scan_output_" + time.Now().Format("02_01_2006_15_04_05") + ".txt"

	filePath := filepath.Join(outputsDir, logFileName)

	file, err := os.Create(filePath)
	if err != nil {
		// fmt.Println("Error while creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(scan_output)
	if err != nil {
		// fmt.Println("Error while writing to the file:", err)
		return
	}
}
