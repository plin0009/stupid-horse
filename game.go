package main

import "fmt"

type Game struct {
	moveTree MoveTree
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
	g := Game{
		moveTree: MoveTree{position: pos},
		players:  players,
	}
	return g
}

// update moveTree to be a child
func (g *Game) Move(m Move) error {
	return fmt.Errorf("no child with move %v found", m)
}
