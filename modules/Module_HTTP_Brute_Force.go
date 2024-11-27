package Modules

import (
	Parameters "github.com/distribRuted/framework/library/parameters"
)

func (F *Module) Module_HTTP_Brute_Force_Init() {

	var (
		moduleShortName   string = "HTTP_Brute_Force"
		moduleName        string = "HTTP Brute Force"
		moduleDescription string = "This module concurrently performs HTTP based brute force attacks."
		moduleAuthor      string = "@numanozdemircom <root@numanozdemir.com>"
	)

	// ATTENTION! Do not edit or remove the line below:
	Set_Module(moduleShortName, moduleName, moduleDescription, moduleAuthor)
}

// Variables of user inputs to pass to the Start() function.
var (
/*
targetHost       string
wordlistPaths    []string
clientCount      int
inputCount       int
stopOnFirstMatch bool
*/
)

// Questions to ask the user are defined here.
func (F *Module) Module_HTTP_Brute_Force_Parameters() {

	// add_parameter(PARAMETER_NAME, DESCRIPTION, DEFAULT_VALUE, IS_REQUIRED)
	Parameters.Add_Parameter("TARGET_HOST", "Enter the target host to attack.", "", true)
	Parameters.Add_Parameter("WORDLIST_PATHS", "Enter the paths of the wordlist files containing the combinations to be tested.", "", true)
	Parameters.Add_Parameter("CLIENT_COUNT", "Enter the number of clients (except this computer).", "-1", true)
	Parameters.Add_Parameter("INPUT_COUNT", "Enter the number of inputs (i.e. enter 2 for username & password).", "-1", true)
	Parameters.Add_Parameter("SINGLE_MATCH", "Stop the attack on the first match.", "true", false)

}

// This function will start the attack.
func (F *Module) Module_HTTP_Brute_Force_Start() {

	/* To call a function within a module, you can use one of the following syntaxes:

	-> 			Call_Client_Role_Function()

	OR -> 		Call_Custom_Function("Module_HTTP_Brute_Force_Client_Role")

	OR -> 		(&initModule{}).Module_HTTP_Brute_Force_Client_Role()

	*/
}

func (F *Module) Module_HTTP_Brute_Force_Server_Role() {
	// your very own codes
}

func (F *Module) Module_HTTP_Brute_Force_Client_Role() {
	// your very own codes

}

func (F *Module) Module_HTTP_Brute_Force_Success_Indicator() string {
	return "Custom_Success_Message_in_Response"
}

func (F *Module) Module_HTTP_Brute_Force_Stop() {
	// Goodbye.
}
