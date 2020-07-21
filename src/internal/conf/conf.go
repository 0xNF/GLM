package conf

import (
	fsops "github.com/0xNF/glm/src/internal/fsops"
	"gopkg.in/ini.v1"
	"errors"
)

const iniLoc string = "./glmconfig.ini"

type GLMConfig struct {
	Triggers []GLMTrigger
	Email    GLMEmail
	Slack    GLMSlack
	HasSlack bool
	HasEmail bool
	NotifyIfNothing bool
}
type GLMEmail struct {
	User       string
	SenderAddress string
	RecipientAddress string
	Password   string
	SmtpServer string
}
type GLMSlack struct {
	Token   string
	Channel string
}
type GLMTrigger struct {
	TriggerFile    string
	SaveFromFolder string
	SavePattern    string
	SaveTo         string
}

// DoesConfigExist Checks whether the local config file exists
func DoesConfigExist() (bool, error) {
	return fsops.CheckExists(iniLoc)
}

// MakeDefaultConfigFile writes default config file to "./glmconfig.ini"
func MakeDefaultConfigFile() (*GLMConfig, error) {
	cfg := ini.Empty()

	/* write the General section */
	cfg.Section("General").Key("NotifyIfNothing").SetValue("True")

	/* write the GLM Trigger section */
	trigger := &GLMTrigger{
		TriggerFile:    "/your/trigger/file.txt",
		SaveFromFolder: "/save/this/folder/",
		SavePattern:    "[file(.*)regex]",
		SaveTo:         "/save/to/this/folder",
	}

	cfg.Section("Trigger").Key("TriggerFile").SetValue(trigger.TriggerFile)
	cfg.Section("Trigger").Key("SaveFromFolder").SetValue(trigger.SaveFromFolder)
	cfg.Section("Trigger").Key("SavePattern").SetValue(trigger.SavePattern)
	cfg.Section("Trigger").Key("SaveTo").SetValue(trigger.SaveTo)

	/* write the GLM Email section */
	email := &GLMEmail{
		User:       "email@localhost",
		Password:   "YourPassword",
		SmtpServer: "smtp.localhost:443",
	}
	cfg.Section("Email").Key("User").SetValue(email.User)
	cfg.Section("Email").Key("SenderAddress").SetValue(email.SenderAddress)
	cfg.Section("Email").Key("RecipientAddress").SetValue(email.RecipientAddress)
	cfg.Section("Email").Key("Password").SetValue(email.Password)
	cfg.Section("Email").Key("SmtpServer").SetValue(email.SmtpServer)

	/* Write the GLM Slack section */
	slack := &GLMSlack{
		Token:   "YourSlackToken",
		Channel: "#you",
	}
	cfg.Section("Slack").Key("Token").SetValue(slack.Token)
	cfg.Section("Slack").Key("Channel").SetValue(slack.Channel)

	err := cfg.SaveToIndent(iniLoc, "\t")
	if err != nil {
		return &GLMConfig{}, err
	}

	conf := &GLMConfig{
		Email: *email,
		Slack: *slack,
		Triggers: []GLMTrigger{
			*trigger,
		},
	}

	return conf, nil

}

// Reads the config file located at "./glmconfig.ini"
func ReadConfigFile() (*GLMConfig, error) {

	/* create Conf */
	conf := &GLMConfig{
		HasEmail: false,
		HasSlack: false,
	}

	/* load Config from disk */
	inif, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}, iniLoc)

	if err != nil {
		return &GLMConfig{}, nil
	}

	/* load General section */
	generalSec, err := inif.GetSection("General")
	if err != nil {
		return conf, errors.New("config file did not have a General section")
	}
	notifyIfNothing, err := generalSec.Key("NotifyIfNothing").Bool()
	if err != nil {
		return conf, errors.New("Failed to parse NotifyIfNothing config field. Must be either [True,False]")
	}
	conf.NotifyIfNothing = notifyIfNothing

	/* load Trigger Section */
	triggerSec, err := inif.GetSection("Trigger")
	if err != nil {
		return conf, errors.New("config file did not have a Trigger section")
	}
	trigger := &GLMTrigger{
		TriggerFile:    triggerSec.Key("TriggerFile").String(),
		SaveFromFolder: triggerSec.Key("SaveFromFolder").String(),
		SavePattern:    triggerSec.Key("SavePattern").String(),
		SaveTo:         triggerSec.Key("SaveTo").String(),
	}
	conf.Triggers = append(conf.Triggers, *trigger)

	/* load Email Section */
	emailSec, err := inif.GetSection("Email")
	if err == nil {
		email := &GLMEmail{
			User:       emailSec.Key("User").String(),
			Password:   emailSec.Key("Password").String(),
			SmtpServer: emailSec.Key("SmtpServer").String(),
			SenderAddress: emailSec.Key("SenderAddress").String(),
			RecipientAddress: emailSec.Key("RecipientAddress").String(),
		}
		conf.Email = *email
		conf.HasEmail = true
	}

	/* Load Slack Section */
	slackSec, err := inif.GetSection("Slack")
	if err == nil {
		slack := &GLMSlack{
			Token:   slackSec.Key("Token").String(),
			Channel: slackSec.Key("Channel").String(),
		}
		conf.Slack = *slack
		conf.HasSlack = true
	}


	return conf, nil
}
