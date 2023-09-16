package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/afmireski/alp-backup-system/internal/data_structures"
)

func main() {
	src := "/home/afmireski/Área de trabalho/tmp"
	srcDir := "tmp"
	path := "/home/afmireski/Área de trabalho/tmp/x"
	dst := "/home/afmireski/Área de trabalho/tmp/bkp/"

	fileInfo, _ := os.Lstat(path)
	after, _ := strings.CutPrefix(path, src) 
	dirPath := dst + srcDir + after;
	fmt.Println(dirPath)

	os.Mkdir(dirPath, fileInfo.Mode())
}

func test2() {
	target := "y.txt"
	path := "/home/afmireski/Área de trabalho/tmp/"

	fileInfo, err := os.Lstat(path+target)

	if err != nil {
		log.Fatal(err);
	}

	fileData, err := os.ReadFile(path+target)
	if err != nil {
		log.Fatal(err);
	}
	bPath := "/home/afmireski/Área de trabalho/tmp/bkp/z/" + target 

	err = os.WriteFile(bPath, fileData, fileInfo.Mode())
	if err != nil {
		log.Fatal(err);
	}
}

func test1() {
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

	fileInfo, err := os.Lstat("/home/afmireski/Documentos/BCC/seinfos")

	// cannot handle symbolic link
	//fileInfo, err := os.Lstat("file.txt")

	if err != nil {
		panic(err)
	}

	fmt.Println("Name : ", fileInfo.Name())

	fmt.Println("Size : ", fileInfo.Size())

	fmt.Println("Mode/permission : ", fileInfo.Mode())

	// --- check if file is a symlink

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		fmt.Println("File is a symbolic link")
	}

	fmt.Println("Modification Time : ", fileInfo.ModTime())

	fmt.Println("Is a directory? : ", fileInfo.IsDir())

	fmt.Println("Is a regular file? : ", fileInfo.Mode().IsRegular())

	fmt.Println("Unix permission bits? : ", fileInfo.Mode().Perm())

	fmt.Println("Permission in string : ", fileInfo.Mode().String())

	fmt.Println("What else underneath? : ", fileInfo.Sys())

	dirs, e := os.ReadDir("/home/afmireski/Documentos/isos/")

	if e == nil {
		fmt.Println(" Len ", len(dirs))
		for _, d := range dirs {
			fmt.Println(" -> ", d.Name())
		}
	} else {
		log.Fatal(e)
	}

	for i, p := range paths {
		backupTable.Insert(p, uint(i))
	}
	backupTable.Print()
	x, y, z := backupTable.Search("/home/afmireski/Documentos/BCC/seinfos")
	fmt.Printf("%x -> %d | %t\n", x, y, z)
}
