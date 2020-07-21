# Gravio Log Monitor
Small tool purpose-built to monitor iflab for ConigAccess errors

This tool will monitor for the presence of the `nlog-internal-configurationfile-errors.log`, and will preserve the last current `ConfigAccess` file, and then send an email to the registered administrator. It will delete any ConfigAccess file by default unless it meets the above criteria.


# Raisen D'etre
Because the ConfigAccessDebug logs are so heavy, we can't keep too many of them around at once. But because they will contain critical information for debugging Gorilla License errors, we need to have them in case a license error occurs. 

# Usage


## Configuration
```ini
[Trigger]
	TriggerFile = /your/trigger/file.txt
	SavePattern = /save/this/[file(.*)regex]
	SaveTo      = /save/to/this/folder

[Email]
	User       = email@localhost
	Password   = YourPassword
	SmtpServer = smtp.localhost:443

[Slack]
	Token   = YourSlackToken
	Channel = `#you`
```
# Installation

