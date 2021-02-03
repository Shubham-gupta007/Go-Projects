package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"zona_tip_batch/tip_batch_abuse.ch/constants"
	"zona_tip_batch/tip_batch_abuse.ch/controllers"
	"zona_tip_batch/tip_batch_lib/common"
	"zona_tip_batch/tip_batch_lib/config"
	"zona_tip_batch/tip_batch_lib/elastic"
	"zona_tip_batch/tip_batch_lib/logger"
)

type (
	App struct {
		ConfigObject config.ConfigStruct
	}
)

func (a *App) InitConfig(logObj logger.ErrorLog) {
	logObj.LogFileName = "main.go"
	logObj.LogFunctionName = "InitConfig"

	configfilepath := "/opt/zona/config.json"
	jsonFile, err := os.Open(configfilepath)
	if err != nil {
		logObj.BatchLog("Error 1 : " + err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logObj.BatchLog("Error 2 : " + err.Error())
	}

	err = json.Unmarshal(byteValue, &a.ConfigObject)
	if err != nil {
		logObj.BatchLog("Error 3 : " + err.Error())
	}
}

func main() {
	logObj := logger.ErrorLog{}
	logObj.Source = "abusech"
	logObj.LogFileName = "main.go"
	logObj.LogFunctionName = "main"

	a := App{}
	a.InitConfig(logObj)

	// checking for TIPLicense
	hasLicense := common.CheckTIPLicense(a.ConfigObject, logObj)
	if hasLicense {
		elasticsearchclient, err := elastic.GetESClient(a.ConfigObject.ElasticSearchHost, a.ConfigObject.ElasticSearchPort)
		if err != nil {
			logObj.BatchLog("Error 1 : " + err.Error())
		}

		go func() {
			for {
				logObj := logger.ErrorLog{}
				logObj.Source = "abusech"
				logObj.LogFileName = "main.go"
				logObj.LogFunctionName = "main"

				// Source - Feodo Tracker
				sourceFT := "abuse.ch Feodo Tracker"
				feedConfigFT := elastic.GetSourceConfig(sourceFT, logObj)
				if feedConfigFT.Enable {
					logObj.BatchLog("Processing abuse.ch Feodo Tracker")
					ip := controllers.NewIPDataController()
					ip.GetIPInformation(a.ConfigObject, elasticsearchclient)

					logObj.BatchLog("Processing Abuse.ch Botnet IP BlockList")
					botnetip := controllers.NewBotnetIPDataController()
					botnetip.GetBotnetIPInformation(a.ConfigObject, elasticsearchclient)

					logObj.BatchLog("Processing Abuse.ch Malware hashes")
					hash := controllers.NewHASHDataController()
					hash.GetHASHInformation(a.ConfigObject, elasticsearchclient)

					if feedConfigFT.Interval != 0 {
						logObj.BatchLog(sourceFT + " - Sleeping for " + fmt.Sprintf("%d", feedConfigFT.Interval) + " min")
						time.Sleep(time.Duration(feedConfigFT.Interval) * time.Minute)
					} else {
						logObj.BatchLog(sourceFT + " - Sleeping for " + fmt.Sprintf("%d", constants.FrequencyInSec) + " min")
						time.Sleep(time.Duration(constants.FrequencyInSec) * time.Second)
					}
				} else {
					time.Sleep(time.Duration(1) * time.Minute)
				}
			}
		}()

		go func() {
			for {
				// Source - abuse.ch SSL Blacklist
				sourceBL := "abuse.ch SSL Blacklist"
				feedConfigBL := elastic.GetSourceConfig(sourceBL, logObj)
				if feedConfigBL.Enable {
					logObj.BatchLog("Processing " + sourceBL)

					ssl := controllers.NewSSLDataController()
					ssl.GetSSLInformation(a.ConfigObject, elasticsearchclient)

					if feedConfigBL.Interval != 0 {
						logObj.BatchLog(sourceBL + " - Sleeping for " + fmt.Sprintf("%d", feedConfigBL.Interval) + " min")
						time.Sleep(time.Duration(feedConfigBL.Interval) * time.Minute)
					} else {
						logObj.BatchLog(sourceBL + " - Sleeping for " + fmt.Sprintf("%d", constants.FrequencyInSec) + " min")
						time.Sleep(time.Duration(constants.FrequencyInSec) * time.Second)
					}
				} else {
					time.Sleep(time.Duration(1) * time.Minute)
				}
			}
		}()

		go func() {
			for {
				// Source - abuse.ch URL Haus
				sourceUH := "abuse.ch URL Haus"
				feedConfigUH := elastic.GetSourceConfig(sourceUH, logObj)
				if feedConfigUH.Enable {
					logObj.BatchLog("Processing " + sourceUH)

					url := controllers.NewURLDataController()
					url.GetURLInformation(a.ConfigObject, elasticsearchclient)

					if feedConfigUH.Interval != 0 {
						logObj.BatchLog(sourceUH + " - Sleeping for " + fmt.Sprintf("%d", feedConfigUH.Interval) + " min")
						time.Sleep(time.Duration(feedConfigUH.Interval) * time.Minute)
					} else {
						logObj.BatchLog(sourceUH + " - Sleeping for " + fmt.Sprintf("%d", constants.FrequencyInSec) + " min")
						time.Sleep(time.Duration(constants.FrequencyInSec) * time.Second)
					}
				} else {
					time.Sleep(time.Duration(1) * time.Minute)
				}
			}
		}()

		for {
			time.Sleep(time.Duration(10) * time.Second)
		}
	} else {
		logObj.BatchLog("TIP License Not Found")
	}
}
