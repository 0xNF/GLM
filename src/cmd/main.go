package main

import (
	"fmt"
	"log"
	"os"
	conf "github.com/0xNF/glm/src/internal/conf"
	monitor "github.com/0xNF/glm/src/internal/monitor"
	notify "github.com/0xNF/glm/src/internal/notify"
)

func main() {
	config := checkLoadConfig()
	shouldNotify := config.NotifyIfNothing
	for _, trigger := range config.Triggers {
		moved, err := monitor.Monitor(&trigger)
		if err != nil {
			log.Fatal(err)
		}
		msg := ""
		subj := ""
		if !moved {
			/* system executed fine, but nothing trigger wasnt detected, so just deleted some files */
			subj = "GLM - No Report"
			msg = fmt.Sprintf("glm finished executing. No trigger (%s) was found. Any matching files were deleted.", trigger.TriggerFile)
			log.Println(msg)
		} else {
			/* send notifications to the administrator */
			shouldNotify = true /* override any ShouldNotify setting - this is entire purpose for the program */
			subj = "GLM - Found A Trigger"
			log.Println(fmt.Sprintf("Trigger file (%s) was found. Moved affected files to safe loc (%s). Notifying administrator.", trigger.TriggerFile, trigger.SaveTo))
			msg = fmt.Sprintf("GLM: Trigger file (%s) was found. Moved affected files to safe location (%s).", trigger.TriggerFile, trigger.SaveTo)
		}

		if shouldNotify {
			if config.HasSlack {
				err = notify.SendSlack(&config.Slack, msg)
				if err != nil {
					fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to send to slack: %s", err))
				}
			}
			if config.HasEmail {
				err = notify.SendMail(&config.Email, subj, msg)
				if err != nil {
					fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to send to email: %s", err))
				}
			}
		}
	}
}


func checkLoadConfig() *conf.GLMConfig {

	var config *conf.GLMConfig = &conf.GLMConfig{}

	ok, err := conf.DoesConfigExist()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load config file: %s", err))
	}
	if !ok {
		log.Println("local config file did not exist. Creating new one.")
		config, err = conf.MakeDefaultConfigFile()
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to save new config file.: %s", err))
		}
		ok, err = conf.DoesConfigExist()
		if err != nil {
			log.Fatal(fmt.Sprintf("Created new config file, but it wasn't available to read: %s", err))
		}
	} else {
		config, err = conf.ReadConfigFile()
		if err != nil {
			log.Fatal(fmt.Sprintf("Found config file, but failed to load it: %s", err))
		}
	}

	return config
}
