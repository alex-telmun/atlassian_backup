// Package jira implements Jira Cloud backup
package jira

import (
	"atlassian_backup/lib/utils"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
)

func New(acc, workspace, token string) *Backup {
	return &Backup{
		Name:               "jira_cloud_backup_" + utils.Timestamp() + ".zip",
		atlassianAccount:   acc,
		atlassianWorkspace: workspace,
		atlassianToken:     token,
	}
}

func (b *Backup) Run() (err error) {
	defer func() { err = utils.WrapIfErr("can't run backup", err) }()

	URL := url.URL{
		Scheme: "https",
		User: url.UserPassword(
			b.atlassianAccount,
			b.atlassianToken,
		),
		Host: b.atlassianWorkspace + ".atlassian.net",
		Path: backupBasePath,
	}

	reqData := bytes.NewReader(backupReqData)

	data, err := utils.Request(http.MethodPost, &URL, reqData)
	if err != nil {
		return err
	}

	if len(data) != 0 {
		return errors.New(string(data))
	}

	return nil
}

func (b *Backup) Progress() (progress int, err error) {
	defer func() { err = utils.WrapIfErr("can't get backup progress", err) }()

	taskId, err := b.lastTaskId()
	if err != nil {
		return 0, err
	}

	query := url.Values{}
	query.Add("taskId", taskId)

	URL := url.URL{
		Scheme: "https",
		User: url.UserPassword(
			b.atlassianAccount,
			b.atlassianToken,
		),
		Host:     b.atlassianWorkspace + ".atlassian.net",
		Path:     progressBasePath,
		RawQuery: query.Encode(),
	}

	data, err := utils.Request(http.MethodGet, &URL, nil)
	if err != nil {
		return 0, err
	}

	var resp progressResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return 0, err
	}

	return resp.Progress, nil
}

func (b *Backup) File() (u *url.URL, err error) {
	defer func() { err = utils.WrapIfErr("can't get backup file URL", err) }()

	taskId, err := b.lastTaskId()
	if err != nil {
		return nil, err
	}

	progQuery := url.Values{}
	progQuery.Add("taskId", taskId)

	URL := url.URL{
		Scheme: "https",
		User: url.UserPassword(
			b.atlassianAccount,
			b.atlassianToken,
		),
		Host:     b.atlassianWorkspace + ".atlassian.net",
		Path:     progressBasePath,
		RawQuery: progQuery.Encode(),
	}

	data, err := utils.Request(http.MethodGet, &URL, nil)
	if err != nil {
		return nil, err
	}
	var resp progressResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	resQuery := url.Values{}
	resQuery.Add("fileId", b.fileId(resp.Result))

	u = &url.URL{
		Scheme: "https",
		User: url.UserPassword(
			b.atlassianAccount,
			b.atlassianToken,
		),
		Host:     b.atlassianWorkspace + ".atlassian.net",
		Path:     downloadBasePath,
		RawQuery: resQuery.Encode(),
	}

	return u, nil
}

func (b *Backup) lastTaskId() (string, error) {

	URL := url.URL{
		Scheme: "https",
		User:   url.UserPassword(b.atlassianAccount, b.atlassianToken),
		Host:   b.atlassianWorkspace + ".atlassian.net",
		Path:   lastTaskIdBasePath,
	}

	data, err := utils.Request(http.MethodGet, &URL, nil)
	if err != nil {
		return "", utils.Wrap("can't get last task ID", err)
	}

	return string(data), nil

}

func (b *Backup) fileId(result string) (id string) {
	re := regexp.MustCompile(fileIdRegex)

	return re.FindString(result)
}
