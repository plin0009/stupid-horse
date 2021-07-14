package main

import "strings"

type Game struct {
	moveTree *MoveTree
	players  map[PieceColour]Player
	timers   map[PieceColour]TimeLeft
}

type Player struct {
	name string
	me   bool
}

type TimeLeft int

func NewGame(fen string, players map[PieceColour]Player) Game {
	pos := LoadInitialPosition(fen)
	mt := new(MoveTree)
	mt.position = pos
	g := Game{
		moveTree: mt,
		players:  players,
	}
	return g
}

// update moveTree based on input algebraic notation
func (g *Game) Move(moveString string) {
	oldPosition := g.moveTree.position
	move := oldPosition.StringToMove(moveString)
	child := new(MoveTree)
	child.move = move
	child.position = oldPosition.ProcessMove(move)
	g.moveTree = child
	g.moveTree.Peek()
}

func (g *Game) Moves(movesString string) {
	for _, moveString := range strings.Fields(movesString) {
		g.Move(moveString)
	}
}
