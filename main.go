package main

import "fmt"

func main() {
	pos := LoadInitialPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	tree := MoveTree{
		position: pos,
	}
	eval := Think(&tree)
	fmt.Println(eval)
	fmt.Println(tree.follow.move)

	lineLength := 50

	beginning := &tree

	for i := 0; i < lineLength; i++ {
		cur := beginning
		for cur.follow != nil {
			fmt.Printf("%v ", cur.move)
			cur = cur.follow
		}
		fmt.Printf("\n")
		beginning = beginning.follow
		if beginning == nil {
			break
		}
		Think(beginning)
	}
}
