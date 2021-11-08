package main

// We can think of a board as an array of arrays (2-dimensional array), with each bool representing the cell's state.
// The board has r rows and c columns
type board [][]bool


func UpdateBoard(currentBoard board) board {
	newBoard := InitializeBoard()
	for rowIndex, row := range currentBoard { // loop through rows
		for cellIndex, _ := range row { // loop through cells
			newBoard[rowIndex][cellIndex] = UpdateCell(currentBoard, rowIndex, cellIndex)
		}
	}
	return newBoard
}

func InitializeBoard() [][]bool {
	var b board
	rows := screenHeight/gridSize
	cols := screenWidth/gridSize
	for i := 1; i <= rows; i++ {
		b = append(b, make([]bool,cols))
	}
	return b
}

func UpdateCell(currentBoard board, row int, col int) bool {
	numNeighbors := CountLiveNeighbors(currentBoard, row, col)
	// apply rules when current cell is alive
	if currentBoard[row][col] {
		if numNeighbors == 2 || numNeighbors == 3 { // Rule of propagation
			return true
		} else { // lack of mates / overpopulation the cell dies
			return false
		}
	} else { // the cell is currently dead
		if numNeighbors == 3 { // birth to new life
			return true
		} else { // remain dead
			return false
		}
	}
}

func CountLiveNeighbors(currentBoard board, row int, col int) int {
	count := 0
	for r := row-1; r <= row+1; r++ { // we loop through every eligible row
		for c := col-1; c <= col+1; c++ { // and eligible column
			if !(r == row && c == col) && InField(currentBoard, r, c) { // excluding the current cell, and the neighbor is on the board
				if currentBoard[r][c] { // neighbor is alive
					count++
				}
			}
		}
	}
	return count
}

func InField(currentBoard board, row int, col int) bool {
	rows := screenHeight/gridSize
	cols := screenWidth/gridSize
	if row < 0 || row > (rows-1) || col <0 || col >(cols-1) {
		return false
	}else {
		return true
	}
}