// Package utils contains common functions for use in the app
package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

/*
Timestamp get current time and returns string with timestamp.
Timestamp format: YY_MM_DD_HH_m
More about time formatting string: https://pkg.go.dev/time#Time.Format
*/
func Timestamp() string {
	return time.Now().Format("2006_01_02_15_04")
}

/*
Request do web request, returns data, as an array bytes or error if
request fail

Arguments:

	method string
	url *url.URL
	reqData io.Reader

Returns:

	data []bytes
	err error
*/
func Request(
	method string,
	url *url.URL,
	reqData io.Reader,
) (data []byte, err error) {
	defer func() { err = WrapIfErr("can't do request", err) }()

	c := http.Client{
		Timeout: 0,
	}

	req, err := http.NewRequest(method, url.String(), reqData)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

/*
NiceSize convert bytes count to human-readable format, e.q. 4.6 GiB

Arguments:

	b int64: bytes count
*/
func NiceSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// Wrap returns formatted error message
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapIfErr returns formatted error message if error is
func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}
