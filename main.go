package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)

	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	if n == 1 {
		fmt.Println(1)
		return
	}

	rowNumber := 0
	colNumber := 0
	count := 1

	// countは1スタートなので、<= にした。
	for count <= n*n {
		for {
			if colNumber == n {
				rowNumber++
				colNumber--
				break
			}
			if matrix[rowNumber][colNumber] != 0 {
				rowNumber++
				colNumber--
				break
			}
			matrix[rowNumber][colNumber] = count
			count++
			colNumber++
		}
		for {
			if rowNumber == n {
				rowNumber--
				colNumber--
				break
			}
			if matrix[rowNumber][colNumber] != 0 {
				rowNumber--
				colNumber--
				break
			}
			matrix[rowNumber][colNumber] = count
			count++
			rowNumber++
		}
		for {
			if colNumber < 0 {
				rowNumber--
				colNumber++
				break
			}
			if matrix[rowNumber][colNumber] != 0 {
				rowNumber--
				colNumber++
				break
			}
			matrix[rowNumber][colNumber] = count
			count++
			colNumber--
		}
		for {
			if matrix[rowNumber][colNumber] != 0 {
				rowNumber++
				colNumber++
				break
			}
			matrix[rowNumber][colNumber] = count
			count++
			rowNumber--
		}
	}

	for _, row := range matrix {
		for _, col := range row {
			fmt.Printf("%d ", col)
		}
		fmt.Println()
	}
}
