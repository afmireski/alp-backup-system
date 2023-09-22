package main

import (
	"fmt"
	// "log"
	"os"

	"github.com/afmireski/alp-backup-system/internal/backup_system"
)

const src = "/home/afmireski/Área de trabalho/tmp"
const dst = "/home/afmireski/Documentos/BCC/p4/aspectos-linguagens-programacao/trabalhos/trabalho01/alp-backup-system/backup/"
const configPath = "/home/afmireski/Documentos/BCC/p4/aspectos-linguagens-programacao/trabalhos/trabalho01/alp-backup-system/config"

func main() {
	wdir, _ := os.Getwd()

	fmt.Println(wdir)

	bs := backup_system.InitBackupSystem(dst, configPath)
	// bs.SetBackupSrc(src)
	// bs.Print()
	// err := bs.Sync()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	bs.Print()
}
