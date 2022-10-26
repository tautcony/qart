package utils

import (
	"fmt"
	"reflect"
	"rsc.io/qr/coding"
	"testing"
)

func TestRotate(t *testing.T) {
	m01 := [][]coding.Pixel{{1, 2, 3, 4}, {12, 13, 14, 5}, {11, 16, 15, 6}, {10, 9, 8, 7}}
	m02 := [][]coding.Pixel{{1, 2, 3, 4}, {12, 13, 14, 5}, {11, 16, 15, 6}, {10, 9, 8, 7}}
	m03 := [][]coding.Pixel{{1, 2, 3, 4}, {12, 13, 14, 5}, {11, 16, 15, 6}, {10, 9, 8, 7}}
	m1 := [][]coding.Pixel{{10, 11, 12, 1}, {9, 16, 13, 2}, {8, 15, 14, 3}, {7, 6, 5, 4}}
	m2 := [][]coding.Pixel{{7, 8, 9, 10}, {6, 15, 16, 11}, {5, 14, 13, 12}, {4, 3, 2, 1}}
	m3 := [][]coding.Pixel{{4, 5, 6, 7}, {3, 14, 15, 8}, {2, 13, 16, 9}, {1, 12, 11, 10}}

	testCase := func(origin [][]coding.Pixel, expected [][]coding.Pixel, turn int) {
		t.Run(fmt.Sprintf("Rotate for %v degree", 90*turn), func(t *testing.T) {
			RotatePixel(origin, turn)
			if !reflect.DeepEqual(origin, expected) {
				t.Errorf("got %v but given %v, ", origin, expected)
			}
		})
	}
	testCase(m01, m1, 1)
	testCase(m02, m2, 2)
	testCase(m03, m3, 3)
}
