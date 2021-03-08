package main

import (
	"fmt"
)


const(
	ORDER = 3
	MAX_KEYS = ORDER - 1
)

type Item struct{
	item int
}

type Node struct{
	m			int // meaning number of keys
	keys	 	[]int
	isLeaf		bool // Identifies a node as being a leaf or not
	nodes 		[]*Node //List of pointers to other nodes 
}

func (n *Node) ToStr(){
	fmt.Printf("Number of keys: %d\n", n.m)
	fmt.Printf("Keys: %v\n", n.keys)
	fmt.Printf("Leaf: %v\n", n.isLeaf)
	fmt.Printf("Nodes: %v\n", n.nodes)
	fmt.Println("-----------------------")
}

func AllocateNode() *Node{
	return &Node{
		m: MAX_KEYS,
		keys: make([]int, MAX_KEYS),
		isLeaf: true,
		nodes: make([]*Node, ORDER),
	}
}

func FindPlace(root* Node, key int) (*Node, int){
	var i = 0
	//Finding suitable index for choosing the correct branch
	for  i < len(root.keys) && key > root.keys[i] && root.keys[i] != 0{
		i++
	}

	if root.isLeaf{
		return root, i
	}else if key > root.keys[i-1]{
		return FindPlace(root.nodes[i], key)
	}else if i == 0{
		if key < root.keys[i]{
			return FindPlace(root.nodes[i], key)
		} 
	}

	return nil, -1
}

func InsertKey(root *Node, key int) {
	node,idx := FindPlace(root,key)
	fmt.Printf("node: %v, index: %d\n", root, idx)
	
	if len(root.keys) == root.m{
		Split(node)
	}else{
		node.keys[idx] = key
		return
	}
}

func PrintTree(root* Node){
	if root != nil{
		root.ToStr()
		if root.nodes != nil{
			for _, node := range root.nodes{
				PrintTree(node)
			}
		}
	}
}

func Split(node* Node){
		
}

func SearchKey(root *Node, key int) *Node{
	return nil
}

func SampleTree() *Node{
	var cn1 = &Node{
		m: MAX_KEYS,
		keys: make([]int, MAX_KEYS),
		isLeaf: true,
		nodes: make([]*Node, MAX_KEYS + 1),
	}
	cn1.keys[0] = 1

	var cn2 = &Node{
		m: MAX_KEYS,
		keys: make([]int, MAX_KEYS),
		isLeaf: true,
		nodes: make([]*Node, MAX_KEYS + 1),
	}
	cn2.keys[0] = 3

	var root = &Node{
		m: MAX_KEYS,
		keys: make([]int, MAX_KEYS),
		isLeaf: false,
		nodes: []*Node{cn1, cn2},
	}
	root.keys[0] = 2

	return root
}