// Package processor implements backup process pipeline
package processor

import (
	"atlassian_backup/backup"
	"atlassian_backup/backup/confluence"
	"atlassian_backup/backup/jira"
	"atlassian_backup/config"
	"atlassian_backup/lib/utils"
	"atlassian_backup/logger"
	"atlassian_backup/notifyer"
	"atlassian_backup/notifyer/slack"
	"atlassian_backup/storage"
	"atlassian_backup/storage/gs"
	"atlassian_backup/storage/local"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	errStartMsg    = "Start %s cloud backup process failure: %v\n"
	errFollowMsg   = "Follow backup %s process failure: %v\n"
	errGetUrlMsg   = "Can't get %s backup file download URL: %v\n"
	errInitStorage = "Can't initialize %s storage: %v\n"
	errSaveMsg     = "Saving backup to %s failure: %v\n"
	successMsg     = "Backup %s successfully saved! Backup size: %s\n"
)

// A Processor object
type Processor struct {
	config *config.Config
}

/*
New create new Processor object

Returns: *Processor
*/
func New() *Processor {
	return &Processor{
		config: config.MustLoad(),
	}
}

/*
Process handle backup pileline. In the beginning run cloud backup procedure,
then check backup progress and download backup file to storage
*/
func (p *Processor) Process() {
	logger.Init()

	b := p.backup()

	logger.Info.Printf(
		"Start %s cloud backup process",
		strings.Title(p.config.BackupType),
	)
	err := b.Run()
	if err != nil {
		p.handleErr(
			errStartMsg,
			p.config.BackupType,
			err,
		)
	}

	for {
		progress, err := b.Progress()
		if err != nil {
			p.handleErr(errFollowMsg, p.config.BackupType, err)
		}
		logger.Info.Printf(
			"Current backup %s progress is: %d%%\n",
			p.config.BackupType,
			progress,
		)

		if progress == int(100) {
			break
		}

		time.Sleep(time.Minute)
	}

	logger.Info.Printf(
		"%s cloud backup success. Downloading backup file to %s storage...\n",
		strings.Title(p.config.BackupType),
		p.config.StorageType,
	)

	fileUrl, err := b.File()
	if err != nil {
		p.handleErr(errGetUrlMsg, p.config.BackupType, err)
	}

	s, err := p.storage()
	if err != nil {
		p.handleErr(errInitStorage, p.config.StorageType, err)
	}

	size, err := s.Save(fileUrl, p.obj())
	if err != nil {
		p.handleErr(errSaveMsg, p.config.StorageType, err)
	}

	n, _ := slack.New()
	n.Send(fmt.Sprintf(
		"Backup %s successfully saved! Backup size is: %s",
		strings.Title(p.config.BackupType),
		size,
	),
	)

	logger.Info.Printf(
		"Backup %s successfully saved! Backup size is: %s",
		strings.Title(p.config.BackupType),
		size,
	)
}

/*
backup parse application config and create object, which implements
backup.Backup interface

Returns: backup.Backup
*/
func (p *Processor) backup() backup.Backup {
	switch p.config.BackupType {
	case "jira":
		return jira.New(
			p.config.AtlassianAccount,
			p.config.AtlassianWorkspace,
			p.config.AtlassianToken,
		)
	case "confluence":
		return confluence.New(
			p.config.AtlassianAccount,
			p.config.AtlassianWorkspace,
			p.config.AtlassianToken,
		)
	default:
		panic("Unsupported backup type parameter")
	}
}

/*
storage parse application config and create object, which implements
storage.Storage interface. Return error if failure

Returns:

	s storage.Storage
	err error
*/
func (p *Processor) storage() (s storage.Storage, err error) {
	switch p.config.StorageType {

	case "gs":
		s, err := gs.New()
		if err != nil {
			return nil, err
		}
		return s, nil

	case "local":
		s, err := local.New()
		if err != nil {
			return nil, err
		}
		return s, nil
	default:
		panic("Unsupported storage type parameter")
	}
}

/*
notifyer parse application config and create object, which implements
notifyer.Notifyer interface. Return error if failure

Returns:

	n notifyer.Notifyer
	err error
*/
func (p *Processor) notifyer() (n notifyer.Notifyer, err error) {
	switch p.config.NotifyType {

	case "slack":
		s, err := slack.New()
		if err != nil {
			return nil, err
		}
		return s, nil

	default:
		panic("Unsupported notify type parameter")
	}
}

/*
obj set backup file path or blob with filename

Returns: string
*/
func (p *Processor) obj() string {
	fName := fmt.Sprintf(
		"%s_%s_%s.tar.gz",
		p.config.BackupType,
		"cloud",
		utils.Timestamp(),
	)

	name := filepath.Join(
		strings.Title(p.config.BackupType),
		"Cloud",
		fName,
	)

	return name
}

/*
handleErr write error message to log and notify to some notifyer

Arguments:

	msg string: error message template
	ph string: placeholder
	e error
*/
func (p *Processor) handleErr(msg, ph string, e error) {
	n, err := p.notifyer()
	if err != nil {
		logger.Error.Fatalf(
			"Can't create %s notifyer object: %v",
			p.config.NotifyType,
			err,
		)
	}

	if err := n.Send(fmt.Sprintf(msg, strings.Title(ph), e)); err != nil {
		logger.Error.Fatalf(
			"Can't send notification to %s: %v\n",
			strings.Title(p.config.NotifyType),
			err,
		)
	}

	logger.Error.Fatalf(msg, strings.Title(ph), e)

}
