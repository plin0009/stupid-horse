package main

import "fmt"

const TTMaxSize = 10000000

type TranspositionTable struct {
	lookup   map[Position]*MoveTree
	stack    []Position
	counter  int
	capacity int
}

func NewTranspositionTable(capacity int) *TranspositionTable {
	tt := new(TranspositionTable)
	tt.lookup = make(map[Position]*MoveTree)
	tt.stack = make([]Position, 0, capacity)
	tt.counter = 0
	tt.capacity = capacity
	return tt
}

func (tt *TranspositionTable) Get(pos Position) *MoveTree {
	if tt == nil {
		return nil
	}
	return tt.lookup[pos]
}

func (tt *TranspositionTable) Add(mt *MoveTree) {
	if tt == nil {
		return
	}
	if tt.lookup[mt.position] != nil {
		// already an entry
		return
	}
	if tt.counter == tt.capacity {
		// TODO: remove earliest entry
		fmt.Println("Encountered capacity")
		tt.counter = 0
		return
	}
	if len(tt.stack) > tt.counter {
		pos := tt.stack[tt.counter]
		delete(tt.lookup, pos)
		tt.stack[tt.counter] = mt.position
	} else {
		tt.stack = append(tt.stack, mt.position)
	}
	tt.lookup[mt.position] = mt
	tt.counter++
}
