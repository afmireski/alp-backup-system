package backup_system

import (
	"fmt"

	"time"

	"github.com/afmireski/alp-backup-system/internal/data_structures"
)

type BackupSystem interface {
	Init() BackupSystem
	SetBackup(path string) error
	Sync() error
}

type FileMetadata struct {
	Path       string
	Filename   string
	IsDir      bool
	ModifiedAt time.Time
}

func (fm FileMetadata) String() string {
	return fmt.Sprintf("{\n\t path=%s\nfilename=%s\nisDir=%t\nmodifiedAt=%s\t\n}", fm.Path, fm.Filename, fm.IsDir, fm.ModifiedAt)
}

type MyBackupSystem struct {
	modifiedAt time.Time
	tableData  data_structures.BackupTable[FileMetadata]
}

func (mbs MyBackupSystem) Init() MyBackupSystem {
	return MyBackupSystem{
		modifiedAt: time.Now(),
		tableData:  data_structures.CreateBackupTable[FileMetadata](50),
	}
}
