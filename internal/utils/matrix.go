package utils

import (
	"math"
	"rsc.io/qr/coding"
)

func RotatePixel(matrix [][]coding.Pixel, turn int) {
	if turn == 0 {
		return
	}
	n := len(matrix)
	x := int(math.Floor(float64(n) / 2))
	y := n - 1
	for i := 0; i < x; i++ {
		for j := i; j < y-i; j++ {
			switch turn {
			case 0: // pass
			case 1:
				matrix[i][j], matrix[j][y-i], matrix[y-j][i], matrix[y-i][y-j] = matrix[y-j][i], matrix[i][j], matrix[y-i][y-j], matrix[j][y-i]
			case 2:
				matrix[i][j], matrix[y-i][y-j] = matrix[y-i][y-j], matrix[i][j]
				matrix[j][y-i], matrix[y-j][i] = matrix[y-j][i], matrix[j][y-i]
			case 3:
				matrix[i][j], matrix[j][y-i], matrix[y-j][i], matrix[y-i][y-j] = matrix[j][y-i], matrix[y-i][y-j], matrix[i][j], matrix[y-j][i]
			}
		}
	}
}
