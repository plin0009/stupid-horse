package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	pos := LoadInitialPosition("bnrnqkrb/pppppppp/8/8/8/8/PPPPPPPP/BNRNQKRB w KQkq -")
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

func startBot() {
	err := godotenv.Load()
	if err != nil {
		panic("could not load .env file")
	}
	b := Bot{token: os.Getenv("lichess_key")}
	b.Listen()
}
