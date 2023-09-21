package main

import (
	"log"
	"time"

	"github.com/afmireski/alp-backup-system/internal/backup_system"
)

func main() {
	src := "/home/afmireski/Área de trabalho/tmp"
	dst := "/home/afmireski/Documentos/BCC/p4/aspectos-linguagens-programacao/trabalhos/trabalho01/alp-backup-system/backup/"

	bs := backup_system.InitBackupSystem(dst)
	bs.SetBackupSrc(src)
	bs.Print()
	err := bs.Sync()
	if err != nil {
		log.Fatal(err)
	}
	bs.Print()
	time.Sleep(1 * time.Minute)
	err = bs.Sync()
	if err != nil {
		log.Fatal(err)
	}
	bs.Print()
}
