package data_structures

import (
	"crypto/sha256"
	"fmt"
)

type HashTable interface {
	Hash(key string) string 
	Alpha() float32
	Insert(key string, value uint) string
	Search(key string) (string, uint, bool)
	Resize()
	Print()
	GetM() uint
	GetN() uint
}

type BackupTable struct {
	size uint // M
	len  uint // N
	data map[string]uint
}

func CreateBackupTable(m uint) BackupTable {
	return BackupTable{
		size: m,
		len: 0,
		data: make(map[string]uint, m),
	}
}

func (bt BackupTable) GetM() (uint) {
	return bt.size
}

func (bt BackupTable) GetN() (uint) {
	return bt.len
}

func (bt* BackupTable) Alpha() float32 {
	return float32(bt.len)/float32(bt.size)
}

func (bt* BackupTable) Hash(key string) string {
	shaBytes := sha256.Sum256([]byte(key))
	return string(shaBytes[:])
}

func (bt* BackupTable) Search(key string) (string, uint, bool) {
	hashKey := bt.Hash(key)

	value, exists := bt.data[hashKey]
	
	return hashKey, value, exists
}

func (bt* BackupTable) Resize(newM uint) {
	newData := make(map[string]uint, newM)
	for key, value := range bt.data {
		newData[key] = value
	}

	bt.size = newM
	bt.data = newData
}

func (bt* BackupTable) Insert(key string, value uint) string {
	hash, _, exists := bt.Search(key)

	if !exists {
		if (bt.Alpha() > 0.5) { 
			bt.Resize(bt.size*2)
		}

		bt.len++
	}

	bt.data[hash] = value

	return hash
}

func (bt* BackupTable) Print() {
	fmt.Println("-------------- BACKUP  TABLE --------------")
	fmt.Printf("| m = %d --------------------------- n = %d |\n", bt.size, bt.len)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -")
	for key, value := range bt.data {
		fmt.Printf("| %x | %d |\n", key, value)
	}
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -")
}
