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
	Remove(key string) 
	RemoveByHash(hash string)
	Resize()
	Print()
}

type BackupTable[T any] struct {
	Size uint         `json:"size" gob:"size"`
	Len  uint         `json:"len" gob:"len"`
	Data map[string]T `json:"data" gob:"data"`
}

func CreateBackupTable[T any](m uint) BackupTable[T] {
	return BackupTable[T]{
		Size: m,
		Len:  0,
		Data: make(map[string]T, m),
	}
}

func (bt *BackupTable[T]) Alpha() float32 {
	return float32(bt.Len) / float32(bt.Size)
}

func (bt *BackupTable[T]) Hash(key string) string {
	shaBytes := sha256.Sum256([]byte(key))
	return string(shaBytes[:])
}

func (bt *BackupTable[T]) Search(key string) (string, T, bool) {
	hashKey := bt.Hash(key)

	value, exists := bt.Data[hashKey]

	return hashKey, value, exists
}

func (bt *BackupTable[T]) Resize(newM uint) {
	newData := make(map[string]T, newM)
	for key, value := range bt.Data {
		newData[key] = value
	}

	bt.Size = newM
	bt.Data = newData
}

func (bt *BackupTable[T]) Insert(key string, value T) string {
	hash, _, exists := bt.Search(key)

	if !exists {
		if bt.Alpha() > 0.5 {
			bt.Resize(bt.Size * 2)
		}

		bt.Len++
	}

	bt.Data[hash] = value

	return hash
}

func (bt *BackupTable[T]) Remove(key string) {
	hash, _, exists := bt.Search(key)

	if exists {
		delete(bt.Data, hash)
		bt.Len--
	}
}

func (bt *BackupTable[T]) RemoveByHash(hash string) {	
	delete(bt.Data, hash)
	bt.Len--
}

func (bt *BackupTable[T]) Print() {
	fmt.Println("-------------- BACKUP  TABLE --------------")
	fmt.Printf("| m = %d --------------------------- n = %d |\n", bt.Size, bt.Len)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -")
	for key, value := range bt.Data {
		fmt.Printf("| %x | %s |\n", key, value)
	}
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - -")
}
