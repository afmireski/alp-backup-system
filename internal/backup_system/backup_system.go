package backup_system

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"time"

	"github.com/afmireski/alp-backup-system/internal/data_structures"
)

type BackupSystem interface {
	Init() BackupSystem
	SetBackupSrc(path string) error
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
	syncedAt      time.Time
	backupSrc     string
	srcDir        string
	backupDst     string
	backupHistory data_structures.BackupTable[FileMetadata]
}

func (mbs MyBackupSystem) Init(dst string) MyBackupSystem {
	return MyBackupSystem{
		syncedAt:      time.Now(),
		backupSrc:     "",
		srcDir:        "",
		backupDst:     dst,
		backupHistory: data_structures.CreateBackupTable[FileMetadata](50),
	}
}

func (mbs *MyBackupSystem) SetBackupSrc(path string) error {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return err
	} else if !fileInfo.IsDir() {
		return errors.New("O seu caminho de backup deve apontar para um diretório")
	}

	mbs.backupSrc = path
	mbs.srcDir = fileInfo.Name() // Obtém o nome útil do diretório de backup

	return nil
}

func (mbs *MyBackupSystem) Sync() error {
	defer mbs.setSyncedAt()

	if len(mbs.backupSrc) > 0 {
		mbs.sync(mbs.backupSrc)
	}
	return nil
}

func (mbs *MyBackupSystem) sync(path string) error {
	pathInfo, pathErr := os.Lstat(path)

	if pathErr == nil { // Existe o caminho para verificar
		// Verifica se o caminho já existe no histórico
		_, value, existsPath := mbs.backupHistory.Search(path)

		if pathInfo.IsDir() { // Verifica se o caminho aponta para um diretório
			dirContent, dirErr := os.ReadDir(path) // Tenta ler o diretório

			if dirErr == nil { // Verifica se conseguiu ler
				if len(dirContent) > 0 { // Verifica se o diretório não está vazio

					if !existsPath {
						mbs.mkdir(path, pathInfo.Mode())
					}

					for _, dirItem := range dirContent { // Percorre os conteúdos do diretório
						mbs.sync(path + string(os.PathSeparator) + dirItem.Name()) // Sincroniza esses conteúdos
					}
				}
			} else {
				return dirErr // Retorna o erro
			}
		} else { // Se não é um diretório, então, é um arquivo

			if existsPath {
				// Verifica se o arquivo não foi modificado
				if pathInfo.ModTime().Equal(value.ModifiedAt) {
					return nil // Não faz nada
				}
			}

			copyErr := mbs.copy(path, pathInfo.Mode())

			if copyErr != nil {
				return copyErr
			}

			metadata := FileMetadata{
				Path:       path,
				Filename:   pathInfo.Name(),
				IsDir:      pathInfo.IsDir(),
				ModifiedAt: time.Now(),
			}

			mbs.backupHistory.Insert(path, metadata) // Atualiza o histórico de Backup
		}
	}
	return nil
}

func (mbs *MyBackupSystem) mkdir(path string, fileMode os.FileMode) error {
	after, _ := strings.CutPrefix(path, mbs.backupSrc) // Remove todo caminho até o diretório de backup
	dirPath := mbs.backupDst + mbs.srcDir + after      // Monta o caminho correto para a cópia
	err := os.Mkdir(dirPath, fileMode)

	return err
}

func (mbs *MyBackupSystem) copy(path string, fileMode os.FileMode) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return errors.New("Falha ao ler o conteúdo de " + path)
	}
	after, _ := strings.CutPrefix(path, mbs.backupSrc) // Remove todo caminho até o diretório de backup
	bPath := mbs.backupDst + mbs.srcDir + after  // Monta o caminho correto para a cópia

	err = os.WriteFile(bPath, fileData, fileMode)
	return err
}

func (mbs* MyBackupSystem) setSyncedAt() {
	mbs.syncedAt = time.Now()
}