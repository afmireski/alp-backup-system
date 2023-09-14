package data_structures

type BTreeNode struct {
	Len     int
	M 		int
	Key     []int
	Value   []int
	Chields []*BTreeNode
	IsLeaf bool
}

func (btn *BTreeNode) Search(key int) (*BTreeNode) {
	if btn.Len == 0 {
		return nil
	}
	i := 0
	for j, k := range btn.Key {
		// Verifica se a chave atual é menor que cada cada chave já cadastrada
		if key <= k {
			i = j
			break
		}
	}	

	if btn.Key[i] == key {
		return btn
	} else if btn.IsLeaf {
		return nil
	}

	return btn.Chields[i].Search(key)
}

func (btn *BTreeNode) Insert(key, value int) (*BTreeNode) {
	if btn == nil {
		node := BTreeNode{
			Len: 0,
			Key: make([]int, 0, btn.M),
			Value: make([]int, 0, btn.M),
			Chields: nil,
			IsLeaf: true,
		}
		node.Key = append(node.Key, key)
		node.Value = append(node.Key, value)
		node.Len++

		return &node
	} 

	i := 0
	for j, k := range btn.Key {
		// Verifica se a chave atual é menor que cada cada chave já cadastrada
		if key <= k {
			i = j
			break
		}
	}
	if btn.Key[i] == key {
		btn.Value[i] = value;
	} else if i == btn.Len-1 {
		btn.Key = append(btn.Key, key)
		btn.Value = append(btn.Value, value)
		btn.Len++
	}

	if btn.Len == btn.M {
		// split
		
	}

	return btn
}

type BTree struct {
	M    int32
	Root *BTreeNode
}

func (bt BTree) Search(key int) (*BTreeNode) {
	if bt.Root == nil {
		return nil
	}
	return bt.Root.Search(key)
}
