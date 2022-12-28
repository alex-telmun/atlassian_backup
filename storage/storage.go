// Package storage include interface for describe backup storage
package storage

import "net/url"

/*
A Storage presents backups storage object

Methods:

	Save(downloadUrl *url.URL, obj string) (size string, err error)
*/
type Storage interface {
	Save(downloadUrl *url.URL, obj string) (size string, err error)
}
