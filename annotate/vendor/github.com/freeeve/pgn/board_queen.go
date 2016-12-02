package pgn

func (b Board) findAttackingQueen(pos Position, color Color, check bool) (Position, error) {
	count := 0
	retPos := NoPosition

	// straight
	r := pos.GetRank()
	f := pos.GetFile()
	for {
		f--
		testPos := PositionFromFileRank(f, r)
		if b.checkQueenColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f++
		testPos := PositionFromFileRank(f, r)
		if b.checkQueenColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		r++
		testPos := PositionFromFileRank(f, r)
		if b.checkQueenColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		r--
		testPos := PositionFromFileRank(f, r)
		if b.checkQueenColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	// diagonals
	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f--
		r--
		testPos := PositionFromFileRank(f, r)
		if b.checkQueenColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f--
		r++
		testPos := PositionFromFileRank(f, r)
		if b.checkQueenColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f++
		r++
		testPos := PositionFromFileRank(f, r)
		if b.checkQueenColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f++
		r--
		testPos := PositionFromFileRank(f, r)
		if b.checkQueenColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if count == 0 {
		return NoPosition, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) findAttackingQueenFromFile(pos Position, color Color, file File) (Position, error) {
	count := 0
	retPos := NoPosition

	//diagonal
	r := pos.GetRank()
	f := pos.GetFile()

	for {
		f--
		r--
		testPos := PositionFromFileRank(f, r)
		if file == f && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()

	for {
		f++
		r++
		testPos := PositionFromFileRank(f, r)
		if file == f && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()

	for {
		f++
		r--
		testPos := PositionFromFileRank(f, r)
		if file == f && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()

	for {
		f--
		r++
		testPos := PositionFromFileRank(f, r)
		if file == f && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}
	//fim diagonal
	r = pos.GetRank()
	f = pos.GetFile()

	for {
		f--
		testPos := PositionFromFileRank(f, r)
		if file == f && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f++
		testPos := PositionFromFileRank(f, r)
		if file == f && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	if file == pos.GetFile() {
		r = pos.GetRank()
		f = pos.GetFile()
		for {
			r++
			testPos := PositionFromFileRank(f, r)
			if b.checkQueenColor(testPos, color) {
				retPos = testPos
				count++
				break
			} else if testPos == NoPosition || b.containsPieceAt(testPos) {
				break
			}
		}

		r = pos.GetRank()
		f = pos.GetFile()
		for {
			r--
			testPos := PositionFromFileRank(f, r)
			if b.checkQueenColor(testPos, color) {
				retPos = testPos
				count++
				break
			} else if testPos == NoPosition || b.containsPieceAt(testPos) {
				break
			}
		}
	}
	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if count == 0 {
		return NoPosition, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) findAttackingQueenFromRank(pos Position, color Color, rank Rank) (Position, error) {
	count := 0
	retPos := NoPosition

	//diagonal
	r := pos.GetRank()
	f := pos.GetFile()

	for {
		f--
		r--
		testPos := PositionFromFileRank(f, r)
		if rank == r && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()

	for {
		f++
		r++
		testPos := PositionFromFileRank(f, r)
		if rank == r && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()

	for {
		f++
		r--
		testPos := PositionFromFileRank(f, r)
		if rank == r && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()

	for {
		f--
		r++
		testPos := PositionFromFileRank(f, r)
		if rank == r && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}
	//fim diagonal

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		r--
		testPos := PositionFromFileRank(f, r)
		if rank == r && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		r++
		testPos := PositionFromFileRank(f, r)
		if rank == r && b.checkQueenColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	if rank == pos.GetRank() {
		r := pos.GetRank()
		f := pos.GetFile()
		for {
			f--
			testPos := PositionFromFileRank(f, r)
			if b.checkQueenColor(testPos, color) {
				retPos = testPos
				count++
				break
			} else if testPos == NoPosition || b.containsPieceAt(testPos) {
				break
			}
		}

		r = pos.GetRank()
		f = pos.GetFile()
		for {
			f++
			testPos := PositionFromFileRank(f, r)
			if b.checkQueenColor(testPos, color) {
				retPos = testPos
				count++
				break
			} else if testPos == NoPosition || b.containsPieceAt(testPos) {
				break
			}
		}
	}

	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if count == 0 {
		return NoPosition, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) checkQueenColor(pos Position, color Color) bool {
	return (b.GetPiece(pos) == WhiteQueen && color == White) ||
		(b.GetPiece(pos) == BlackQueen && color == Black)
}
