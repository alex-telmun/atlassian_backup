/*
Package local implements function for saving file, which can be downloaded
from URL, to local filesystem

Required environment:

	LOCAL_FOLDER:
*/
package local

import (
	"atlassian_backup/lib/utils"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

/*
A LocalStorage is representation local filesystem for store backup files
Implements storage.Storage interface
*/
type LocalStorage struct {
	LocalPath string
}

/*
New returns new LocalStorage object or error if failure. Get local folder from
environment: LOCAL_FOLDER

Returns:

	*LocalStorage
	error
*/
func New() (*LocalStorage, error) {
	folder, ok := os.LookupEnv("LOCAL_FOLDER")
	if !ok {
		return nil, errors.New("Local folder is not specified")
	}

	return &LocalStorage{
		LocalPath: folder,
	}, nil
}

/*
Save download backup file from URL and save if to local filesystem. Returns
backup file size or error (if failure).

Arguments:

	downloadUrl *url.URL
	obj string

Returns:

	size string
	err error
*/
func (ls *LocalStorage) Save(
	downloadUrl *url.URL,
	obj string,
) (size string, err error) {
	defer func() {
		err = utils.WrapIfErr("can't save backup to folder", err)
	}()

	client := http.Client{
		Timeout: 0,
	}

	req, err := http.NewRequest(http.MethodGet, downloadUrl.String(), nil)
	if err != nil {
		return "", nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer func() { _ = resp.Body.Close() }()

	filename := filepath.Join(ls.LocalPath, obj)

	err = ifDirNotExists(filename)
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	nBytes, err := io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return utils.NiceSize(nBytes), nil

}

/*
ifDirNotExists create dir tree for stora backup file if dirs doesn't exists

Returns: error
*/
func ifDirNotExists(fPath string) error {
	if _, err := os.Stat(filepath.Dir(fPath)); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
