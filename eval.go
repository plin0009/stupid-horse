package main

func (p Piece) Value() int {
	switch p.Type() {
	case Pawn:
		return 10
	case Knight:
		return 70 // stupid-horse!
	case Bishop:
		return 30
	case Rook:
		return 50
	case Queen:
		return 90
	case King:
		return 999
	case NoPieceType:
		return 0
	}
	panic("Invalid piece")
}

var colourMultiplier = map[PieceColour]int{White: 1, Black: -1}

func Eval(p Position) int {
	//fmt.Println(p.board)
	score := 0
	for _, square := range Squares {
		piece := p.board[square]
		// TODO: make use of piece square
		if piece == NoPiece {
			continue
		}
		score += colourMultiplier[piece.Colour()] * piece.Value()
		if piece.Type() == Pawn {
			fileDiff, rankDiff := Diff(ToSquare(File(4), pawnInfo[piece.Colour()].promotionRank), square)
			score += fileDiff + rankDiff
		}
	}
	return score
}

const BotDepth = 6

func SortByCapture(moves []Move) {
	lastCapture := 0
	for i, move := range moves {
		if move.capture {
			temp := moves[lastCapture]
			moves[lastCapture] = move
			moves[i] = temp
			lastCapture++
		}
	}
}

func Think(mt *MoveTree) int {
	var minimax func(alpha, beta int) func(*MoveTree, []Move, int) int
	minimax = func(alpha, beta int) func(*MoveTree, []Move, int) int {
		return func(mt *MoveTree, moves []Move, depth int) int {
			if depth == 0 {
				mt.eval = Eval(mt.position)
				return mt.eval
			}
			mt.eval = colourMultiplier[mt.position.turn] * -9999999
			SortByCapture(moves)
			for _, move := range moves {
				child := new(MoveTree)
				child.parent = mt
				child.move = move
				child.position = mt.position.ProcessMove(move)
				child.FindMoves(depth-1, minimax(alpha, beta))
				if !child.legal {
					continue
				}
				mt.children = append(mt.children, child)
				if mt.position.turn == White {
					if child.eval > mt.eval {
						mt.eval = child.eval
						mt.follow = child
					}
					if child.eval > alpha {
						alpha = child.eval
					}
					if alpha >= beta {
						mt.discard = true
						break
					}
				} else {
					if child.eval < mt.eval {
						mt.eval = child.eval
						mt.follow = child
					}
					if child.eval < beta {
						beta = child.eval
					}
					if beta <= alpha {
						mt.discard = true
						break
					}
				}
			}
			return mt.eval
		}
	}

	return mt.FindMoves(BotDepth, minimax(-9999999, 9999999))
}
