package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type (
	HostInfo struct {
		ServerHost  string `json:"serverhost" bson:"serverhost"`
		Port        string `json:"port" bson:"port"`
		replicahost string `json:"replicahost" bson:"replicahost"`
		tenantcode  string `json:"tenantcode" bson:"tenantcode"`
	}
)

func main() {

	fmt.Println("---------------------------Checking Configuration For Replica Server Please Wait---------------------------")

	command := "cat /etc/os-release"
	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		fmt.Println(err.Error())
		log.Println(err)
	}
	osName := string(output)
	name := osName[:strings.IndexByte(osName, '\n')]
	if name == `NAME="CentOS Linux"` {
		checkFileExists := "test -f /opt/zona/sslkey " + " && echo file exists "
		checkFileExistsResponse, err := exec.Command("sh", "-c", checkFileExists).Output()
		if err != nil {
			fmt.Println("sslkey file missing!!")
		}
		if string(checkFileExistsResponse) != "" {
			fmt.Println("Stop Installation")
		}
		if string(checkFileExistsResponse) == "" {
			fmt.Println("---------------------------Initializing Installation For Primary Server Please Wait---------------------------")
			//1. Create Log file
			logfilecommand := "sudo mkdir /opt/mongosetuplog"
			logfileResponse, err := exec.Command("sh", "-c", logfilecommand).Output()
			if err != nil {
				fmt.Println("Error in creating /opt/mongosetuplog folder command or /opt/mongosetuplog folder already exists" + err.Error())
			}
			if string(logfileResponse) != "" {
				fmt.Println("Mongosetuplog Log Folder Created")
				setupLog("mongosetuplog Log Folder Created")
			}

			//1.1. log File Read Permission
			logfilePermissioncommand := "sudo chmod -R 777 /opt/mongosetuplog"
			logfilePermisssionResponse, err := exec.Command("sh", "-c", logfilePermissioncommand).Output()
			if err != nil {
				fmt.Println("Error in Giving permission for log file Command:" + err.Error())
			}
			if string(logfilePermisssionResponse) != "" {
				fmt.Println("mongosetuplog Log File Permission")
				setupLog("mongosetuplog Log Folder Permission")
			}

			//2. Create /opt/zona/ folder
			createZonafolder := "sudo mkdir /opt/zona"
			createZonafolderResponse, err := exec.Command("sh", "-c", createZonafolder).Output()
			if err != nil {
				fmt.Println("Error in creating /opt/zona/ folder command or the folder already exists" + err.Error())
				setupLog("Error in creating Zona folder command")
			}
			if string(createZonafolderResponse) != "" {
				fmt.Println("Zona folder Created")
				setupLog("Zona folder Created")
			}

			//2.1. Give Permission to /opt/zona folder
			optZonaFolderPermissionCommand := "sudo chmod 777 -R /opt/zona"
			createOptZonaPermissionResponse, err := exec.Command("sh", "-c", optZonaFolderPermissionCommand).Output()
			if err != nil {
				fmt.Println("Error in Giving Permission to /opt/zona Folder Command:" + err.Error())
				setupLog("Error in Giving Permission to /opt/zona Folder Command:" + err.Error())
			}
			if string(createOptZonaPermissionResponse) != "" {
				fmt.Println("/opt/zona Folder Permission Set")
				setupLog("/opt/zona Folder Permission Set")
			}

			//3 Create /opt/zona/logs folder
			createLogsfolder := "sudo mkdir /opt/zona/logs"
			createLogsfolderResponse, err := exec.Command("sh", "-c", createLogsfolder).Output()
			if err != nil {
				fmt.Println("Error in creating /opt/zona/logs folder command or the folder already exists" + err.Error())
				setupLog("Error in creating Zona folder command")
			}
			if string(createLogsfolderResponse) != "" {
				fmt.Println("Logs folder Created")
				setupLog("Logs folder Created")
			}

			//3.1 Give Permission to /opt/zona/logs folder
			optLogsFolderPermissionCommand := "sudo chmod 777 -R /opt/zona/logs"
			createOptZonaLogsPermissionResponse, err := exec.Command("sh", "-c", optLogsFolderPermissionCommand).Output()
			if err != nil {
				fmt.Println("Error in Giving Permission to /opt/zona/logs Folder Command:" + err.Error())
				setupLog("Error in Giving Permission to /opt/zona/logs Folder Command:" + err.Error())
			}
			if string(createOptZonaLogsPermissionResponse) != "" {
				fmt.Println("/opt/zona/logs Folder Permission Set")
				setupLog("/opt/zona/logs Folder Permission Set")
			}
			// DO PWD and copy whole folder(mongosetup) to /opt/zona -extra setup

			//4 Create mongosetup in /opt/zona
			createmongosetupfolder := "sudo mkdir /opt/zona/mongosetup"
			createmongosetupfolderResponse, err := exec.Command("sh", "-c", createmongosetupfolder).Output()
			if err != nil {
				fmt.Println("Error in creating /opt/zona/mongosetup folder command or the folder already exists" + err.Error())
				setupLog("Error in creating Zona folder command")
			}
			if string(createmongosetupfolderResponse) != "" {
				fmt.Println("mongosetup folder Created")
				setupLog("mongosetup folder Created")
			}

			//4.1 Give Permission to /opt/zona/logs folder
			optmongosetupFolderPermissionCommand := "sudo chmod 777 -R /opt/zona/mongosetup"
			createOptZonamongosetupPermissionResponse, err := exec.Command("sh", "-c", optmongosetupFolderPermissionCommand).Output()
			if err != nil {
				fmt.Println("Error in Giving Permission to /opt/zona/logs Folder Command:" + err.Error())
				setupLog("Error in Giving Permission to /opt/zona/logs Folder Command:" + err.Error())
			}
			if string(createOptZonamongosetupPermissionResponse) != "" {
				fmt.Println("/opt/zona/logs Folder Permission Set")
				setupLog("/opt/zona/logs Folder Permission Set")
			}

			//5. Move All files to /opt/zona/mongosetup
			mongoFilesMoveCommand := "sudo mv mongod.conf sslkey /opt/zona/mongosetup"
			mongoFileMoveCommand, err := exec.Command("sh", "-c", mongoFilesMoveCommand).Output()
			if err != nil {
				fmt.Println("Error in Executing moving conf files Command:" + err.Error())
				setupLog("Error in Executing moving conf files Command:" + err.Error())
			} else {
				if string(mongoFileMoveCommand) == "" {
					fmt.Println("Mongo All Files Moved ")
					setupLog("Mongo All Files Moved ")
				}
			}

			//6. Move mongod.conf file
			mongoConfMoveCommand := "sudo mv /opt/zona/mongosetup/mongod.conf /etc/mongod.conf"
			mongoConfMoveFile, err := exec.Command("sh", "-c", mongoConfMoveCommand).Output()
			if err != nil {
				fmt.Println("Error in Executing moving conf Command:" + err.Error())
				setupLog("Error in Executing moving conf Command:" + err.Error())
			} else {
				if string(mongoConfMoveFile) == "" {
					fmt.Println("Mongo Conf File Moved ")
					setupLog("Mongo Conf File Moved ")
				}
			}

			//7. Move sslkey file
			mongoSSlMoveCommand := "sudo mv /opt/zona/mongosetup/sslkey /opt/zona/sslkey"
			mongoSSlMoveFile, err := exec.Command("sh", "-c", mongoSSlMoveCommand).Output()
			if err != nil {
				fmt.Println("Error in Executing moving ssl file Command:" + err.Error())
				setupLog("Error in Executing moving ssl file Command:" + err.Error())
			} else {
				if string(mongoSSlMoveFile) == "" {
					fmt.Println("Mongo sslkey File Moved ")
					setupLog("Mongo sslkey File Moved ")
				}
			}

			//8. Set Permission to ssl key file
			mongoSetPermissionCommand := "sudo chmod 400 /opt/zona/sslkey"
			mongoSetPermissionSSLCommand, err := exec.Command("sh", "-c", mongoSetPermissionCommand).Output()
			if err != nil {
				fmt.Println("Error in Setting Permission to sslkey File:" + err.Error())
				setupLog("Error in Setting Permission to sslkey File:" + err.Error())
			} else {
				if string(mongoSetPermissionSSLCommand) == "" {
					fmt.Println("/opt/zona/sslkey Permission Set ")
					setupLog("/opt/zona/sslkey Permission Set ")
				}
			}

			//9. Set Ownership to ssl key file
			mongoSetOwnershipCommand := "sudo chown mongod:root /opt/zona/sslkey"
			mongoSetOwnershipSSLCommand, err := exec.Command("sh", "-c", mongoSetOwnershipCommand).Output()
			if err != nil {
				fmt.Println("Error in Setting Ownership to sslkey File:" + err.Error())
				setupLog("Error in Setting Ownership to sslkey File:" + err.Error())
			} else {
				if string(mongoSetOwnershipSSLCommand) == "" {
					fmt.Println("/opt/zona/sslkey Ownership Set ")
					setupLog("/opt/zona/sslkey Ownership Set ")
				}
			}

			//10. RestartMongoCommand
			RestartMongoCommand := "sudo service mongod restart"
			// RestartMongoCommand := "sudo mongod --config /etc/mongod.conf"
			RestartMongo, err := exec.Command("sh", "-c", RestartMongoCommand).Output()
			if err != nil {
				fmt.Println("Error in Executing restart Mongodb:" + err.Error())
				setupLog("Error in Executing restart Mongodb:" + err.Error())
			} else {
				if string(RestartMongo) == "" {
					fmt.Println("Mongodb ReStarted")
					setupLog("Mongodb ReStarted")
				}
			}

			//11. Check Mongo service is UP or not -extra setup
			sess, err := mgo.Dial("localhost")
			if err != nil {
				fmt.Println(err)
				fmt.Println("Mongo Service not Restarted!!")
				fmt.Println("Trying to Restart Mongod Service Forcefully!!")

				//11.1 Forcefully start mongo service
				RestartMongoForcefullyCommand := "sudo mongod --config /etc/mongod.conf"
				RestartMongoForcefully, err := exec.Command("sh", "-c", RestartMongoForcefullyCommand).Output()
				if err != nil {
					fmt.Println("Error in Executing restart Mongodb ForceFully:" + err.Error())
					setupLog("Error in Executing restart Mongodb ForceFully:" + err.Error())
				} else {
					if string(RestartMongoForcefully) == "" {
						fmt.Println("Mongodb ReStarted ForceFully")
						setupLog("Mongodb ReStarted ForceFully")

						sess1, err1 := mgo.Dial("localhost")
						if err1 != nil {
							fmt.Println(err)
							fmt.Println("Mongo Service not Restarted ForceFully ")
							fmt.Println("Mongo Service not Restarted ForceFully ")
						} else {
							defer sess1.Close()
							fmt.Println("MongoDB server is healthy.")
							setupLog("MongoDB server is healthy. ")
						}
					}
				}
			} else {
				defer sess.Close()
				fmt.Println("MongoDB server is healthy.")
				setupLog("MongoDB server is healthy. ")
			}

			time.Sleep(10 * time.Second)

			//12.Remove Mongo file from yum.repos.d
			removeMongofile := "sudo rm -rf /opt/zona/mongosetup"
			removeMongoFileResponse, err := exec.Command("sh", "-c", removeMongofile).Output()
			if err != nil {
				fmt.Println("Error in Removing Setup File:" + err.Error())
				setupLog("Error in Removing Setup File:" + err.Error())
			}
			if string(removeMongoFileResponse) == "" {
				fmt.Println("Removed Setup file")
				setupLog("Removed Setup file")
			}

			fmt.Println("Installation Completed!!")
			setupLog("Installation Completed!!")
		} else {
			fmt.Println("Mongo is installed already")
			//	setupLog("Mongo is installed already")
		}
	} else {
		fmt.Println("Please run in the centOS machine")
		//	setupLog("Please run in the centOS machine")
	}
}

func setupLog(msg string) {

	date := time.Now().Local().Format("2006-01-02 15")
	logpath := "/opt/mongosetuplog/"

	err := os.Chmod("/opt/mongosetuplog/", 0777)
	if err != nil {
		fmt.Println("Error 1, Function : setupLog, File : mongoController.go")
		fmt.Println(err.Error())
	}

	fileName := logpath + date + "-mongosetup.log"
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Error 1, Function : setupLog, File : mongoController.go")
		fmt.Println(err.Error())
		log.Println("error opening file: ", err.Error())
	} else {
		defer f.Close()
		w := bufio.NewWriter(f)
		currentDateAndTime := time.Now()
		LogFileName := "mongosetup"
		LogFunctionName := "main.go"

		logString := "" + currentDateAndTime.Format("2006/01/02 15:04:05") + "->" + LogFileName + "->" + LogFunctionName + " : " + msg

		_, err = fmt.Fprintf(w, "%v\n", logString)
		if err != nil {
			fmt.Println("Error 2, Function : setupLog, File : mongoController.go")
			fmt.Println(err.Error())
		}
		w.Flush()
	}
}
