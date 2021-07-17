package main

import (
	"fmt"
	"testing"
)

func TestFindAllMoves(t *testing.T) {
	normalpos := LoadInitialPosition("startpos")
	t.Run("normal d=4", func(t *testing.T) {
		tree := MoveTree{position: normalpos}
		nodes := tree.FindAllMoves(4)
		if nodes != 197281 {
			t.Errorf("FindAllMoves(4) = %d; want 197281", nodes)
		}
	})
	t.Run("normal d=5", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping normal d=5 in short mode.")
		}
		tree := MoveTree{position: normalpos}
		nodes := tree.FindAllMoves(5)
		if nodes != 4865609 {
			t.Errorf("FindAllMoves(5) = %d; want 4865609", nodes)
		}
	})

	// too slow to test right now (> 10 min)
	/*
	 *t.Run("normal d=6", func(t *testing.T) {
	 *  if testing.Short() {
	 *    t.Skip("skipping normal d=6 in short mode.")
	 *  }
	 *  tree := MoveTree{position: normalpos}
	 *  nodes := tree.FindAllMoves(6)
	 *  if nodes != 119060324 {
	 *    t.Errorf("FindAllMoves(6) = %d; want 119060324", nodes)
	 *  }
	 *})
	 */

	// too slow to test right now (> 10 min)
	/*
	 *t.Run("normal d=7", func(t *testing.T) {
	 *  if testing.Short() {
	 *    t.Skip("skipping normal d=7 in short mode.")
	 *  }
	 *  tree := MoveTree{position: normalpos}
	 *  nodes := tree.FindAllMoves(7)
	 *  if nodes != 3195901860 {
	 *    t.Errorf("FindAllMoves(7) = %d; want 3195901860", nodes)
	 *  }
	 *})
	 */

	pos := LoadInitialPosition("nbqrknbr/pppppppp/8/8/8/8/PPPPPPPP/NBQRKNBR w KQkq - 0 1")
	t.Run("960 #1, d=4", func(t *testing.T) {
		tree := MoveTree{position: pos}
		nodes := tree.FindAllMoves(4)
		if nodes != 166056 {
			t.Errorf("FindAllMoves(4) = %d; want 166056", nodes)
		}
	})
	t.Run("960 #1, d=5", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping 960 #1 d=5 in short mode.")
		}
		tree := MoveTree{position: pos}
		nodes := tree.FindAllMoves(5)
		if nodes != 3995752 {
			t.Errorf("FindAllMoves(5) = %d; want 3995752", nodes)
		}

	})
}

func TestMoveTreeState(t *testing.T) {
	t.Run("Fastest checkmate", func(t *testing.T) {
		g := NewGame("startpos", nil, nil)
		g.AddMoves("f2f4 e7e5 g2g4")
		if g.moveTree.state != Active {
			t.Errorf("Fastest checkmate half-move 3: %v, want %v", g.moveTree.state, Active)
		}
		g.AddMoves("d8h4")
		if g.moveTree.state != BlackWon {
			t.Errorf("Fastest checkmate half-move 4: %v, want %v", g.moveTree.state, BlackWon)
		}
	})
	t.Run("Fastest stalemate", func(t *testing.T) {
		g := NewGame("startpos", nil, nil)
		g.AddMoves("c2c4 h7h5 h2h4 a7a5 d1a4 a8a6 a4a5 a6h6 a5c7 f7f6 c7d7 e8f7 d7b7 d8d3 b7b8 d3h7 b8c8 f7g6")
		if g.moveTree.state != Active {
			t.Errorf("Fastest stalemate half-move 18: %v, want %v", g.moveTree.state, Active)
		}
		g.AddMoves("c8e6")
		if g.moveTree.state != Stalemate {
			t.Errorf("Fastest stalemate half-move 19: %v, want %v", g.moveTree.state, Stalemate)
		}
	})
}

func BenchmarkFindAllMoves(b *testing.B) {
	pos := LoadInitialPosition("nbqrknbr/pppppppp/8/8/8/8/PPPPPPPP/NBQRKNBR w KQkq - 0 1")
	tree := MoveTree{
		position: pos,
	}
	b.Run("d=4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree.FindAllMoves(4)
		}
	})
	b.Run("d=5", func(b *testing.B) {
		if testing.Short() {
			b.Skip("skipping d=5 in short mode.")
		}
		for i := 0; i < b.N; i++ {
			tree.FindAllMoves(5)
		}
	})
}
func BenchmarkThink(b *testing.B) {
	pos := LoadInitialPosition("bnrnqkrb/pppppppp/8/8/8/8/PPPPPPPP/BNRNQKRB w KQkq -")
	tree := MoveTree{
		position: pos,
	}
	b.Run("l=4", func(b *testing.B) {
		beginning := &tree
		for i := 0; i < 4; i++ {
			Think(beginning)
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
		}
	})
	b.Run("l=10", func(b *testing.B) {
		beginning := &tree
		if testing.Short() {
			b.Skip("skipping l=10 in short mode.")
		}
		for i := 0; i < 10; i++ {
			Think(beginning)
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
		}
	})
}
