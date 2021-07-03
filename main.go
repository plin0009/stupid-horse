package main

import "fmt"

func main() {
	pos := LoadInitialPosition("nbqrknbr/pppppppp/8/8/8/8/PPPPPPPP/NBQRKNBR w KQkq - 0 1")
	tree := MoveTree{
		position: pos,
	}
	fmt.Println(tree.FindMoves(4))
	for i, child := range tree.children {
		fmt.Printf("%v %v\n", i, child.move)
		fmt.Println(child.FindMoves(3))
	}
}
