package main

import "testing"

func TestFindMoves(t *testing.T) {
	pos := LoadInitialPosition("nbqrknbr/pppppppp/8/8/8/8/PPPPPPPP/NBQRKNBR w KQkq - 0 1")
	tree := MoveTree{
		position: pos,
	}
	var nodes int
	t.Run("d=4", func(t *testing.T) {
		nodes = tree.FindMoves(4)
		if nodes != 166056 {
			t.Errorf("FindMoves(4) = %d; want 166056", nodes)
		}
	})
	t.Run("d=5", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping d=5 in short mode.")
		}
		nodes = tree.FindMoves(5)
		if nodes != 3995752 {
			t.Errorf("FindMoves(5) = %d; want 3995752", nodes)
		}

	})
}

func BenchmarkFindMoves(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pos := LoadInitialPosition("nbqrknbr/pppppppp/8/8/8/8/PPPPPPPP/NBQRKNBR w KQkq - 0 1")
		tree := MoveTree{
			position: pos,
		}
		tree.FindMoves(4)
	}
}
