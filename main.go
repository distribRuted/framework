package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	Lib "github.com/distribRuted/framework/library"
	Connection "github.com/distribRuted/framework/library/connection"
	Console "github.com/distribRuted/framework/library/console"
	Log "github.com/distribRuted/framework/library/log"
	Parameters "github.com/distribRuted/framework/library/parameters"
	Modules "github.com/distribRuted/framework/modules"
	"github.com/fatih/color"
)

func main() {

	// Create a startup log
	Log.Create("The application has been started.", "Application", "INFO")

	// Print welcome message
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	web := color.New(color.FgWhite, color.BgRed, color.Bold, color.Underline).SprintFunc()

	fmt.Println(fmt.Sprintf(`%c     %s     %s        %s %c  %s %s %s%c %s %s%s %s%c%s %s%s%c %s%s%s%c`, '\n', white("_ _     _        _ _"), red("____"), white("_           _"), '\n', white("__| (_)___| |_ _ __(_) |__"), red("|  _ \\"), white("_   _| |_ ___  __| |"), '\n', white("/ _`"), white("| / __| __| '__| | '_ \\"), red("| |_)"), white("| | | | __/ _ \\/ _` |"), '\n', white("| (_| | \\__ \\ |_| |  | | |_)"), red("|  _ <"), white("| |_| | ||  __/ (_| |"), '\n', white("\\__,_|_|___/\\__|_|  |_|_.__/"), red("|_| \\_\\"), white("\\__,_|\\__\\___|\\__,_|  v1"), '\n'))
	fmt.Println(fmt.Sprintf(`                %s %s                     %s`, white("Distributed attack framework."), red("λ\n\n"), web("www.distribRuted.com")))
	fmt.Print("\n\n")

	// Set CLI arguments
	if len(os.Args) > 1 {
		flag.BoolVar(&Parameters.IsServer, "server", false, "Start application as a server, not as a client.")
		flag.BoolVar(&Parameters.IsClient, "client", false, "Start application as a client, not as a server.")
		flag.StringVar(&Parameters.DestIPAddr, "ip", "0.0.0.0", "IP address of the central server.")
		flag.StringVar(&Parameters.SelectedModuleName, "module", "", "Start application with the selected module.")
		flag.StringVar(&Parameters.CLIParameters, "parameters", "", "Start application module with the defined parameters.")
		flag.UintVar(&Parameters.DefaultPort, "port", 1337, "Port number for the TCP communication. Default value is 1337.")
		flag.Parse()
	} else {
		fmt.Println(color.HiWhiteString("Use --help parameter to see commands.\n"))
	}

	// Check if the given IP address and port is valid or not
	if Parameters.IsClient && !Connection.IsIPAddrValid(Parameters.DestIPAddr) {
		if Parameters.DestIPAddr == "" {
			Lib.ExitWithError("Enter an IP address to connect by using --ip parameter.")
		}
		Lib.ExitWithError("The given IP address is invalid.")
	}
	if Parameters.DefaultPort > 65535 {
		Lib.ExitWithError("The given port number is invalid.")
	}

	// Print application mode
	Parameters.ConnHost, Parameters.ConnPort = Parameters.DestIPAddr, strconv.FormatUint(uint64(Parameters.DefaultPort), 10)
	if Parameters.IsClient && Parameters.IsServer {
		Parameters.ListenDialBoth = true
		go Connection.Listen()
		go Connection.Dial()
		Log.Create("The application is running in both client and server modes.", "Application", "INFO")
		Log.PrintMsg(color.YellowString("The application is running in both client and server modes."))
	} else {
		if Parameters.IsClient {
			go Connection.Dial()
			Log.Create("The application has been running in client mode.", "Application", "INFO")
			Log.PrintMsg(color.YellowString("The application has been running in client mode."))
		} else {
			// Listen for incoming connections
			go Connection.Listen()
			Log.Create("The application has been running in server mode.", "Application", "INFO")
			Log.PrintMsg(color.YellowString("The application has been running in server mode."))
		}
	}

	// Load modules
	Modules.Load()

	// Print questions
	printQuestions()

	// Execute the functions of the selected module
	Modules.ExecuteFunctions()

	// Open the console to receive and evaluate user commands.
	Console.Open()

}

func colorTest() {
	log.Println(color.RedString("Chrome path could not found."))
	log.Println(color.GreenString("Chrome path could not found."))
	log.Println(color.YellowString("Chrome path could not found."))
	log.Println(color.BlueString("Chrome path could not found."))
	log.Println(color.MagentaString("Chrome path could not found."))
	log.Println(color.CyanString("Chrome path could not found."))
	log.Println(color.WhiteString("Chrome path could not found."))

	log.Println(color.HiRedString("Chrome path could not found."))
	log.Println(color.HiGreenString("Chrome path could not found."))
	log.Println(color.HiYellowString("Chrome path could not found."))
	log.Println(color.HiBlueString("Chrome path could not found."))
	log.Println(color.HiMagentaString("Chrome path could not found."))
	log.Println(color.HiCyanString("Chrome path could not found."))
	log.Println(color.HiWhiteString("Chrome path could not found."))

	redX := color.New(color.FgRed, color.Bold).SprintFunc()
	greenX := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellowX := color.New(color.FgYellow, color.Bold).SprintFunc()
	blueX := color.New(color.FgBlue, color.Bold).SprintFunc()
	magentaX := color.New(color.FgMagenta, color.Bold).SprintFunc()
	cyanX := color.New(color.FgCyan, color.Bold).SprintFunc()
	whiteX := color.New(color.FgWhite, color.Bold).SprintFunc()

	fmt.Println()

	log.Printf("%s", redX("Chrome path could not found."))
	log.Printf("%s", greenX("Chrome path could not found."))
	log.Printf("%s", yellowX("Chrome path could not found."))
	log.Printf("%s", blueX("Chrome path could not found."))
	log.Printf("%s", magentaX("Chrome path could not found."))
	log.Printf("%s", cyanX("Chrome path could not found."))
	log.Printf("%s", whiteX("Chrome path could not found."))
}

func printQuestions() {
	var currentQuestion uint8 = 1
	var userChoice string

	if currentQuestion == 1 {

		Log.PrintMsg(color.GreenString(strconv.Itoa(len(Parameters.AllModules))) + color.YellowString(" modules loaded. Choose any from below:\n"))
		var modulesCount int = len(Parameters.ModuleNames)
		var moduleIndex int = 1
		var moduleMatched bool = false
		for _, v := range Parameters.ModuleNames {
			if Parameters.SelectedModuleName == v {
				moduleMatched = true
				Parameters.SelectedModule = moduleIndex
			}
			fmt.Println(color.HiWhiteString(strconv.Itoa(moduleIndex) + ") " + Parameters.AllModules[v].Name))
			moduleIndex++
		}
		fmt.Println()

		if moduleMatched {
			fmt.Println(color.HiMagentaString("Selected module:"), strconv.Itoa(Parameters.SelectedModule)+") "+Parameters.SelectedModuleName)
			currentQuestion = 2
		} else {
			for {
				fmt.Print(color.HiMagentaString("Your choice (1 - " + strconv.Itoa(modulesCount) + "): "))
				reader := bufio.NewReader(os.Stdin)
				userChoice, _ = reader.ReadString('\n')
				userChoice = userChoice[:len(userChoice)-1]
				intVal, err := strconv.Atoi(userChoice)
				if err == nil && intVal > 0 && intVal <= modulesCount {
					Parameters.SelectedModule, _ = strconv.Atoi(userChoice)
					currentQuestion = 2
					break
				}
			}
		}
	}

	// Import details and questions from the module
	if currentQuestion == 2 {
		printDetails := color.New(color.FgWhite).Add(color.Underline).SprintFunc()
		var moduleDetails = Parameters.AllModules[Parameters.ModuleNames[Parameters.SelectedModule-1]]
		Parameters.SelectedModuleName = Parameters.ModuleNames[Parameters.SelectedModule-1]
		fmt.Println(fmt.Sprintf(printDetails("\nModule:")) + " " + color.HiWhiteString(moduleDetails.Name))
		fmt.Println(fmt.Sprintf(printDetails("Description:")) + " " + color.HiWhiteString(moduleDetails.Description))
		fmt.Println(fmt.Sprintf(printDetails("Author:")) + " " + color.HiWhiteString(moduleDetails.Author) + "\n")

		// currentQuestion = 3
	}

}
