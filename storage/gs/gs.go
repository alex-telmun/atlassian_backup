/*
Package gs implements function for saving file, which can be downloaded
from URL, to Google Cloud Storage.

Required environment

	GOOGLE_APPLICATION_CREDENTIALS: path to google service-account json file
	GS_BUCKET_NAME: Google Storage bucket name
*/
package gs

import (
	"atlassian_backup/lib/utils"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

/*
A GoogleStorage is representation Google Cloud Storage for store backup files.
Implements storage.Storage interface.
*/
type GoogleStorage struct {
	bucketName string
}

/*
New returns new GoogleStorage object or error if failure. Get google storage
bucket name from environment: GS_BUCKET_NAME

Returns:

	*GoogleStorage
	error
*/
func New() (*GoogleStorage, error) {
	bucket, ok := os.LookupEnv("GS_BUCKET_NAME")
	if !ok {
		return nil, errors.New("GS bucker is not specified")
	}

	return &GoogleStorage{
		bucketName: bucket,
	}, nil
}

/*
Save download backup file from URL and save it to Google Storage. Returns
backup file size or error (if failure).

Arguments:

	downloadUrl *url.URL
	obj string

Returns:

	size string
	err error
*/
func (gs *GoogleStorage) Save(
	downloadUrl *url.URL,
	obj string,
) (size string, err error) {
	defer func() { err = utils.WrapIfErr("can't save backup to storage", err) }()

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

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Hour)
	defer cancel()

	gsClient, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func() { _ = gsClient.Close() }()

	object := gsClient.Bucket(gs.bucketName).Object(obj)
	writer := object.NewWriter(ctx)
	defer func() { _ = writer.Close() }()

	nBytes, err := io.Copy(writer, resp.Body)
	if err != nil {
		return "", err
	}

	return utils.NiceSize(nBytes), nil
}
