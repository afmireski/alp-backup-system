package main

import (
	"fmt"

	"github.com/afmireski/alp-backup-system/internal/data_structures"
)

func main() {
	backupTable := data_structures.CreateBackupTable[uint](11)

	paths := []string{
		"/home/afmireski/Documentos/",
		"/home/afmireski/Documentos/BCC/p3/ed2/BCC33A/hashing/ea/hashtable_ea.c",
		"/home/afmireski/Documentos/BCC/seinfos",
		"/home/afmireski/Documentos/isos",
		"/home/afmireski/Documentos/isos/manjaro-kde-22.1.3-23.iso",
		"/home/afmireski/Documentos/training/go/fullcycle/curso-go",
		"/home/afmireski/Documentos/isos/manjaro-kde-22.1.3-23.iso",
		"/home/afmireski/Documentos/isos/manjaro-gnome-22.1.3.iso",
	}

	for i, p := range paths {
		backupTable.Insert(p, uint(i))
	}
	backupTable.Print()
	x, y, z := backupTable.Search("/home/afmireski/Documentos/BCC/seinfos")
	fmt.Printf("%x -> %d | %t\n", x, y, z)
}
