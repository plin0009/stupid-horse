package main

import "strings"

type Game struct {
	id         string
	moveTree   *MoveTree
	players    map[PieceColour]Player
	timers     map[PieceColour]uint
	initialFen string
	moves      string
}

type Player struct {
	id     string
	name   string
	rating int
	title  string
	me     bool
}

func NewGame(fen string, players map[PieceColour]Player, clock map[PieceColour]uint) *Game {
	pos := LoadInitialPosition(fen)
	mt := new(MoveTree)
	mt.position = pos
	g := &Game{
		moveTree:   mt,
		players:    players,
		timers:     clock,
		initialFen: fen,
		moves:      "",
	}
	return g
}

// update moveTree based on input algebraic notation
func (g *Game) AddMoves(movesString string) {
	child := new(MoveTree)
	child.position = g.moveTree.position
	for _, moveString := range strings.Fields(movesString) {
		child.move = child.position.StringToMove(moveString)
		child.position = child.position.ProcessMove(child.move)
		if g.moves != "" {
			g.moves += " "
		}
		g.moves += moveString
	}
	g.moveTree = child
	g.moveTree.Peek()
}
