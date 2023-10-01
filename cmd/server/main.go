package main

import (
	"crypto/sha256"
	"fmt"
	"data_structures"
)

func main() {
	hashTable := data_structures.CreateHashTable[int](100)

	key := "exemplo"
	value := 42

	hash := hashTable.Hash(key)
	insertedKey := hashTable.Insert(key, value)
	fmt.Printf("Chave inserida: %s\n", insertedKey)

	foundHash, foundValue, exists := hashTable.Search(key)
	if exists {
		fmt.Printf("Chave encontrada: %s, Valor: %v\n", foundHash, foundValue)
	} else {
		fmt.Println("Chave não encontrada.")
	}

	backupTable := data_structures.CreateBackupTable[int](50) 

	key2 := "exemplo2"
	value2 := 24

	hash2 := backupTable.Hash(key2)
	insertedKey2 := backupTable.Insert(key2, value2)
	fmt.Printf("Chave inserida na tabela de backup: %s\n", insertedKey2)

	http.HandleFunc("/api/requisicao", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			requestData := map[string]interface{}{
				"chave": "exemplo",
				"valor": 42,
			}

			chave := requestData["chave"].(string)
			valor := int(requestData["valor"].(float64))

			hash := hashTable.Hash(chave)

			insertedKey := hashTable.Insert(chave, valor)
			fmt.Printf("Chave inserida: %s\n", insertedKey)

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Chave inserida com sucesso: %s\n", insertedKey)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, "Método não permitido. Utilize POST.")
		}
	})

	http.ListenAndServe(":8080", nil)
}

func (bt *data_structures.BackupTable[int]) Hash(key string) string {
	shaBytes := sha256.Sum256([]byte(key))
	return string(shaBytes[:])
}