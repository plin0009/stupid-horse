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
