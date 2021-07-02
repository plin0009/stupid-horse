package main

import (
	"fmt"
)

func main() {
	pos := LoadInitialPosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	fmt.Println(pos.board)

	s := Square(52)
	fmt.Println(s)
}
