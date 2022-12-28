// Package confluence implements Confluence Cloud backup
package confluence

import (
	"atlassian_backup/lib/utils"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

func New(acc, workspace, token string) *Backup {
	return &Backup{
		Name:               "confluence_cloud_backup_" + utils.Timestamp() + ".zip",
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

	URL := url.URL{
		Scheme: "https",
		User: url.UserPassword(
			b.atlassianAccount,
			b.atlassianToken,
		),
		Host: b.atlassianWorkspace + ".atlassian.net",
		Path: progressBasePath,
	}

	data, err := utils.Request(http.MethodGet, &URL, nil)
	if err != nil {
		return 0, err
	}

	var resp progressResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return 0, err
	}

	return convertSize(resp.Progress), nil
}

func (b *Backup) File() (u *url.URL, err error) {
	defer func() { err = utils.WrapIfErr("can't get backup file URL", err) }()

	URL := url.URL{
		Scheme: "https",
		User: url.UserPassword(
			b.atlassianAccount,
			b.atlassianToken,
		),
		Host: b.atlassianWorkspace + ".atlassian.net",
		Path: progressBasePath,
	}

	data, err := utils.Request(http.MethodGet, &URL, nil)
	if err != nil {
		return nil, err
	}

	var resp progressResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &url.URL{
		Scheme: "https",
		User: url.UserPassword(
			b.atlassianAccount,
			b.atlassianToken,
		),
		Host: b.atlassianWorkspace + ".atlassian.net",
		Path: downloadBasePath + resp.Result,
	}, nil
}

func convertSize(sizeStr string) int {
	re := regexp.MustCompile(percentageRegex)

	p, err := strconv.Atoi(re.FindString(sizeStr))
	if err != nil {
		log.Fatal(err)
	}

	return p
}
