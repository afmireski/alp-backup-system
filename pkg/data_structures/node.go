package data_stuctures

type Node interface {
	Insert()
	Search()
	Sort()
}

type NodeKey interface {
	~int32 | ~string
}

type NodeValue interface {} // any
