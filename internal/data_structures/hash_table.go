package data_structures

import (
	"crypto/sha256"
	"fmt"
)

type HashTable[T any] interface {
	Hash(key string) string
	Alpha() float32
	Insert(key string, value T) string
	Search(key string) (string, T, bool)
	Resize()
	Print()
	GetM() uint
	GetN() uint
}

type BackupTable[T any] struct {
	size uint // M
	len  uint // N
	data map[string]T
}

func CreateBackupTable[T any](m uint) BackupTable[T] {
	return BackupTable[T]{
		size: m,
		len:  0,
		data: make(map[string]T, m),
	}
}

func (bt BackupTable[T]) GetM() uint {
	return bt.size
}

func (bt BackupTable[T]) GetN() uint {
	return bt.len
}

func (bt *BackupTable[T]) Alpha() float32 {
	return float32(bt.len) / float32(bt.size)
}

func (bt *BackupTable[T]) Hash(key string) string {
	shaBytes := sha256.Sum256([]byte(key))
	return string(shaBytes[:])
}

func (bt *BackupTable[T]) Search(key string) (string, T, bool) {
	hashKey := bt.Hash(key)

	value, exists := bt.data[hashKey]

	return hashKey, value, exists
}

func (bt *BackupTable[T]) Resize(newM uint) {
	newData := make(map[string]T, newM)
	for key, value := range bt.data {
		newData[key] = value
	}

	bt.size = newM
	bt.data = newData
}

func (bt *BackupTable[T]) Insert(key string, value T) string {
	hash, _, exists := bt.Search(key)

	if !exists {
		if bt.Alpha() > 0.5 {
			bt.Resize(bt.size * 2)
		}

		bt.len++
	}

	bt.data[hash] = value

	return hash
}

func (bt *BackupTable[T]) Print() {
	fmt.Println("-------------- BACKUP  TABLE --------------")
	fmt.Printf("| m = %d --------------------------- n = %d |\n", bt.size, bt.len)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -")
	for key, value := range bt.data {
		fmt.Printf("| %x | %s |\n", key, value)
	}
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -")
}
