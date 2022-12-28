// Package jira implements Jira Cloud backup
package jira

const (
	backupBasePath     string = "/rest/backup/1/export/runbackup"
	lastTaskIdBasePath string = "/rest/backup/1/export/lastTaskId"
	progressBasePath   string = "/rest/backup/1/export/getProgress"
	downloadBasePath   string = "/plugins/servlet/export/download/"
	fileIdRegex        string = "([a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12})"
)

var backupReqData = []byte(`{"cbAttachments":"true","exportToCloud":"true"}`)

type Backup struct {
	Name               string
	atlassianAccount   string
	atlassianWorkspace string
	atlassianToken     string
}

type progressResponse struct {
	Progress int    `json:"progress"`
	Result   string `json:"result"`
}
