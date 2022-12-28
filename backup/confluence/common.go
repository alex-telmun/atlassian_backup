// Package confluence implenments Confluence Cloud backup
package confluence

const (
	backupBasePath   string = "/wiki/rest/obm/1.0/runbackup"
	progressBasePath string = "/wiki/rest/obm/1.0/getprogress"
	downloadBasePath string = "/wiki/download/"
	percentageRegex  string = "[0-9]{1,3}"
)

var backupReqData = []byte(`{"cbAttachments":"true","exportToCloud":"true"}`)

type Backup struct {
	Name               string
	atlassianAccount   string
	atlassianWorkspace string
	atlassianToken     string
}

type progressResponse struct {
	Progress string `json:"alternativePercentage"`
	Result   string `json:"fileName"`
}
