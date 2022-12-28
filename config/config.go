// Package config implements an application configuration values
package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	AtlassianAccount   string
	AtlassianWorkspace string
	AtlassianToken     string
	BackupType         string
	StorageType        string
	NotifyType         string
}

/*
Supported types
*/
var (
	backupTypes  [2]string = [2]string{"jira", "confluence"}
	storageTypes [2]string = [2]string{"gs", "local"}
	notifyTypes  [1]string = [1]string{"slack"}
)

/*
MustLoad get application configuration from command-line arguments
or environment. Close programm with fatal message, if no arguments are set
or their values are incorrect. Some storage and notification types has specific
configurations.

Environment variables:

	ATLASSIAN_ACCOUNT
	ATLASSIAN_WORKSPACE
	ATLASSIAN_TOKEN
	BACKUP_TYPE
	STORAGE_TYPE
	NOTIFY_TYPE

Command-line flags:

	-atlassianAccount
	-atlassianWorkspace
	-atlassianToken
	-backupType
	-storageType
	-notifyType

Returns: Config
*/
func MustLoad() *Config {
	atlassianAccount := flag.String(
		"atlassianAccount",
		"",
		"Atlassian account name (email or username)",
	)

	atlassianWorkspace := flag.String(
		"atlassianWorkspace",
		"",
		"Atlassian workspace name",
	)

	atlassianToken := flag.String(
		"atlassianToken",
		"",
		"Atlassian account name (email or username)",
	)

	backupType := flag.String(
		"backupType",
		"",
		"What you want to backup(confluence or jira)",
	)

	storageType := flag.String(
		"storageType",
		"",
		"Where you want to save the backup (gs)",
	)

	notifyType := flag.String(
		"notifyType",
		"",
		"How you want to get notification",
	)

	flag.Parse()

	if *atlassianAccount == "" {
		acc, ok := os.LookupEnv("ATLASSIAN_ACCOUNT")
		if !ok {
			log.Fatal("Atlassian account is not specified")
		}
		*atlassianAccount = acc
	}

	if *atlassianWorkspace == "" {
		ws, ok := os.LookupEnv("ATLASSIAN_WORKSPACE")
		if !ok {
			log.Fatal("Atlassian workspace is not specified")
		}
		*atlassianWorkspace = ws
	}

	if *atlassianToken == "" {
		token, ok := os.LookupEnv("ATLASSIAN_TOKEN")
		if !ok {
			log.Fatal("Atlassian token is not specified")
		}
		*atlassianToken = token
	}

	if *backupType == "" {
		bType, ok := os.LookupEnv("BACKUP_TYPE")
		if !ok {
			log.Fatal("Backup type is not specified")
		}
		*backupType = bType
	}

	if *storageType == "" {
		sType, ok := os.LookupEnv("STORAGE_TYPE")
		if !ok {
			log.Fatal("Backup type is not specified")
		}
		*storageType = sType
	}

	if *notifyType == "" {
		nType, ok := os.LookupEnv("NOTIFY_TYPE")
		if !ok {
			log.Fatal("Notify type is not specified")
		}
		*notifyType = nType
	}

	if !validateType(backupTypes[:], *backupType) {
		log.Fatal("Backup type is incorrect")
	}

	if !validateType(storageTypes[:], *storageType) {
		log.Fatal("Storage type is incorrect")
	}

	if !validateType(notifyTypes[:], *notifyType) {
		log.Fatal("Notify type is incorrect")
	}

	return &Config{
		AtlassianAccount:   *atlassianAccount,
		AtlassianWorkspace: *atlassianWorkspace,
		AtlassianToken:     *atlassianToken,
		BackupType:         *backupType,
		StorageType:        *storageType,
		NotifyType:         *notifyType,
	}
}

/*
validateType checks if selected type in supported types.

Arguments:

	supTypes []string - supported types
	selType string - selected type

Returns: bool
*/
func validateType(supTypes []string, selType string) bool {
	for _, t := range supTypes {
		if selType == t {
			return true
		}
	}
	return false
}
