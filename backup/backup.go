// Package backup include interface for describe backup object
package backup

import (
	"net/url"
)

/*
A Backup presents backup object

Methods:

	Run() (err error)
	Progress() (progress int, err error)
	File() (u *url.URL, err error)
*/
type Backup interface {
	Run() (err error)
	Progress() (progress int, err error)
	File() (u *url.URL, err error)
}
