package data_structures

type BTreeNode struct {
	Len     int
	Key     []int
	Value   []int
	Chields []BTreeNode
	IsLeaf bool
}

func (btn BTreeNode) Search(key int) (*BTreeNode) {
	if btn.Len == 0 {
		return nil
	}
	i := 0
	for i < int(btn.Len) && key > btn.Key[i] {
		
	}

	if btn.Key[i] == key {
		return &btn
	} else if btn.IsLeaf {
		return nil
	}

	return btn.Chields[i].Search(key)
}

type BTree struct {
	M    int32
	Root BTreeNode
}
