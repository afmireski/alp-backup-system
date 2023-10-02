package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/afmireski/alp-backup-system/internal/backup_system"
	"github.com/gin-gonic/gin"
)

type SetBackupSrcInput struct {
	Path string `json:"path"`
}

type SetModeSrcInput struct {
	Mode backup_system.BackupModeEnum `json:"mode"`
}

var bs backup_system.BackupSystem

func setBackupSrc(c *gin.Context) {
	var input SetBackupSrcInput

	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Input inválido"})
		return
	} else if len(input.Path) <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "O caminho deve ser informado"})
		return
	}

	err := bs.SetBackupSrc(input.Path)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Houve uma falha ao definir a fonte do Backup"})
		return
	}

	c.IndentedJSON(http.StatusOK, bs)
}

func sync(c *gin.Context) {

	err := bs.Sync()

	fmt.Println(err)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Houve uma falha ao sincronizar o Backup, tente novamente mais tarde"})
		return
	}

	c.IndentedJSON(http.StatusOK, bs)
}

func changeMode(c *gin.Context) {

	var input SetModeSrcInput

	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Input inválido"})
		return
	} 

	bs.SetMode(input.Mode)	

	c.IndentedJSON(http.StatusOK, bs)
}

func main() {
	workDir, _ := os.Getwd()
	workDir, _ = strings.CutSuffix(workDir, "cmd/server")

	dst := workDir + "/backup"
	configPath := workDir + "/config"

	bs = backup_system.InitBackupSystem(dst, configPath)

	router := gin.New()
	router.Use(gin.Logger())
	router.POST("/alp-backup/setBackupSrc", setBackupSrc)
	router.POST("/alp-backup/sync", sync)
	router.POST("/alp-backup/changeMode", changeMode)
	router.Run("localhost:3000") // Inicia o servidor
}
