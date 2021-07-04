package main

import "testing"

func TestFindAllMoves(t *testing.T) {
	pos := LoadInitialPosition("nbqrknbr/pppppppp/8/8/8/8/PPPPPPPP/NBQRKNBR w KQkq - 0 1")
	tree := MoveTree{
		position: pos,
	}
	var nodes int
	t.Run("d=4", func(t *testing.T) {
		nodes = tree.FindAllMoves(4)
		if nodes != 166056 {
			t.Errorf("FindAllMoves(4) = %d; want 166056", nodes)
		}
	})
	t.Run("d=5", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping d=5 in short mode.")
		}
		nodes = tree.FindAllMoves(5)
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
