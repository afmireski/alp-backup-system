package backup_system

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/afmireski/alp-backup-system/internal/data_structures"
)

type BackupModeEnum int

type BackupSystemInterface interface {
	SetBackupSrc(path string) error
	Sync() error
	SetMode(newMode BackupModeEnum)
	saveConfigFile()
	loadConfigFile(path string)
}

type FileMetadata struct {
	Path       string
	Filename   string
	IsDir      bool
	ModifiedAt time.Time
}

const (
	MIRROR BackupModeEnum = iota // representa uma sequência de inteiros sem tipo
	PERSISTANCE
)

func (fm FileMetadata) String() string {
	return fmt.Sprintf("{\n\t path=%s\nfilename=%s\nisDir=%t\nmodifiedAt=%s\t\n}", fm.Path, fm.Filename, fm.IsDir, fm.ModifiedAt)
}

type BackupSystem struct {
	SyncedAt       time.Time `json:"synced_at"`
	ConfigFilePath string `json:"config_file_path"`
	BackupSrc      string `json:"backup_src"`
	SrcDir         string `json:"src_dir"`
	BackupDst      string `json:"backup_dst"`
	BackupHistory  data_structures.BackupTable [FileMetadata] `json:"-"`
	BackupMode     BackupModeEnum `json:"backup_mode"`
}

func (bs *BackupSystem) Print() {
	fmt.Println("\n\n-------------- BACKUP  SYSTEM --------------")
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -")
	switch mode := bs.BackupMode; mode {
	case MIRROR:
		fmt.Printf("| Mode=MIRROR |\n")
	default:
		fmt.Printf("| Mode=PERSISTANCE |\n")
	}
	fmt.Printf("| ConfigFilePath=%s |\n", bs.ConfigFilePath)
	fmt.Printf("| BackupDst=%s |\n", bs.BackupDst)
	fmt.Printf("| BackupSrc=%s |\n", bs.BackupSrc)
	fmt.Printf("| SrcDir=%s |\n", bs.SrcDir)
	fmt.Printf("| SyncedAt=%d/%d/%d  %d:%d:%d |\n", bs.SyncedAt.Year(), bs.SyncedAt.Month(), bs.SyncedAt.Day(), bs.SyncedAt.Hour(), bs.SyncedAt.Minute(), bs.SyncedAt.Second())
	bs.BackupHistory.Print()
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -\n\n")

}

func InitBackupSystem(dst, configPath string) BackupSystem {
	if !strings.HasSuffix(dst, "/") {
		dst = dst + "/"
	}
	if strings.HasSuffix(configPath, "/") {
		configPath, _ = strings.CutSuffix(configPath, "/")
	}
	bs, err := loadConfigFile(configPath)

	if err != nil {
		// Não conseguiu ler a config, então backup puro
		return BackupSystem{
			SyncedAt:       time.Now(),
			BackupSrc:      "",
			SrcDir:         "",
			ConfigFilePath: configPath,
			BackupDst:      dst,
			BackupHistory:  data_structures.CreateBackupTable[FileMetadata](50),
			BackupMode:     MIRROR,
		}
	}

	return *bs
}

func (bs *BackupSystem) SetMode(newMode BackupModeEnum) {
	switch newMode {
	case MIRROR:
		bs.BackupMode = MIRROR
	default:
		bs.BackupMode = PERSISTANCE
	}
	defer bs.saveConfigFile()
}

func (bs *BackupSystem) SetBackupSrc(path string) error {
	path, _ = strings.CutSuffix(path, "/")

	fileInfo, err := os.Stat(path)

	if err != nil {
		return err
	} else if !fileInfo.IsDir() {
		return errors.New("o seu caminho de backup deve apontar para um diretório")
	}
	
	bs.BackupSrc = path
	bs.SrcDir = fileInfo.Name() // Obtém o nome útil do diretório de backup

	return nil
}

func (bs *BackupSystem) Sync() error {
	// Executam depois que Sync sair da pilha
	defer bs.saveConfigFile()     // 3
	defer bs.removeDeletedFiles() // 2
	defer bs.setSyncedAt()        // 1
	// --

	if len(bs.BackupSrc) > 0 {
		return bs.sync(bs.BackupSrc)
	}
	return nil
}

func (bs *BackupSystem) sync(path string) error {
	pathInfo, pathErr := os.Lstat(path)

	if pathErr == nil { // Existe o caminho para verificar
		// Verifica se o caminho já existe no histórico
		_, value, existsPath := bs.BackupHistory.Search(path)

		if pathInfo.IsDir() { // Verifica se o caminho aponta para um diretório
			dirContent, dirErr := os.ReadDir(path) // Tenta ler o diretório

			if dirErr == nil { // Verifica se conseguiu ler
				if len(dirContent) > 0 { // Verifica se o diretório não está vazio

					if !existsPath {
						bs.mkdir(path, pathInfo.Mode())

						metadata := FileMetadata{
							Path:       path,
							Filename:   pathInfo.Name(),
							IsDir:      pathInfo.IsDir(),
							ModifiedAt: pathInfo.ModTime(),
						}
						bs.BackupHistory.Insert(path, metadata)
					}

					for _, dirItem := range dirContent { // Percorre os conteúdos do diretório
						syncErr := bs.sync(path + string(os.PathSeparator) + dirItem.Name()) // Sincroniza esses conteúdos
						if syncErr != nil {
							return syncErr
						}
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

			copyErr := bs.copy(path, pathInfo.Mode())

			if copyErr != nil {
				return copyErr
			}

			metadata := FileMetadata{
				Path:       path,
				Filename:   pathInfo.Name(),
				IsDir:      pathInfo.IsDir(),
				ModifiedAt: pathInfo.ModTime(),
			}

			bs.BackupHistory.Insert(path, metadata) // Atualiza o histórico de Backup
		}
	}
	return nil
}

func (bs *BackupSystem) removeDeletedFiles() {
	if bs.BackupMode == MIRROR && bs.BackupHistory.Len > 0 {
		// Percorre a tabela hash inteira se o MIRROR mod estiver ativo e tiver algo gravado no histórico.
		for key, value := range bs.BackupHistory.Data {
			_, err := os.Stat(value.Path)

			// Verifica se cada arquivo salvo na hash continua na origem
			if os.IsNotExist(err) {
				err := bs.del(value.Path) // Tenta excluir o arquivo no backup
				fmt.Println(err)
				if err == nil {
					// Caso tenha conseguido excluir, remove da hash
					bs.BackupHistory.RemoveByHash(key)
				}
			}
		}
	}
}

func (bs *BackupSystem) mkdir(path string, fileMode os.FileMode) error {
	after, _ := strings.CutPrefix(path, bs.BackupSrc) // Remove todo caminho até o diretório de backup
	dirPath := bs.BackupDst + bs.SrcDir + after       // Monta o caminho correto para a cópia
	err := os.Mkdir(dirPath, fileMode)

	fmt.Printf("--- \n mkdir: %s! \n ---", dirPath)

	return err
}

func (bs *BackupSystem) copy(path string, fileMode os.FileMode) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return errors.New("Falha ao ler o conteúdo de " + path)
	}
	after, _ := strings.CutPrefix(path, bs.BackupSrc) // Remove todo caminho até o diretório de backup
	bPath := bs.BackupDst + bs.SrcDir + after         // Monta o caminho correto para a cópia

	err = os.WriteFile(bPath, fileData, fileMode)

	fmt.Printf("--- \n copy: %s \n para \n %s \n ---", path, bPath)
	return err
}

func (bs *BackupSystem) del(path string) error {
	after, _ := strings.CutPrefix(path, bs.BackupSrc) // Remove todo caminho até o diretório de backup
	bPath := bs.BackupDst + bs.SrcDir + after         // Monta o caminho para o backup

	fmt.Printf("--- \n del: %s \n ---", bPath)

	return os.Remove(bPath)
}

func (bs *BackupSystem) setSyncedAt() {
	bs.SyncedAt = time.Now()
}

func (bs *BackupSystem) saveConfigFile() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(bs)
	if err != nil {
		fmt.Println("Houve uma falha ao salvar as configurações do sistema: ", err)
	}

	err = os.WriteFile(bs.ConfigFilePath, buffer.Bytes(), 0666)

	if err != nil {
		fmt.Println("Houve uma falha ao salvar o arquivo de configuração: ", err)
	}
}

func loadConfigFile(path string) (*BackupSystem, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var loadedBackup BackupSystem
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	err = decoder.Decode(&loadedBackup)

	if err != nil {
		return nil, err
	}

	return &loadedBackup, nil
}
