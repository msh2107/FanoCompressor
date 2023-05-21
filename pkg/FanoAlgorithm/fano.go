package FanoAlgorithm

import (
	"math"
	"sort"
)

type code struct {
	Index int
	freq  int
	Code  string
}

func FanoEncoding(freq []int) []code {
	var codes []code
	i := 0
	for index, f := range freq {
		if f != 0 {
			codes = append(codes, code{
				Index: index,
				freq:  f,
			})
			i++
		}
	}
	if len(codes) == 0 {
		return nil
	}

	if len(codes) == 1 {
		codes[0].Code = "0"
		return codes
	}
	sort.Slice(codes, func(i, j int) bool { return codes[i].freq > codes[j].freq })

	buildCodes(codes)
	return codes
}

func buildCodes(codes []code) {
	if len(codes) < 2 {
		return
	}

	divider := bestDividerPosition(codes)

	for i := 0; i < len(codes); i++ {
		if i >= divider {
			codes[i].Code += "0"
		} else {
			codes[i].Code += "1"
		}
	}

	buildCodes(codes[:divider])
	buildCodes(codes[divider:])
}

func bestDividerPosition(codes []code) int {
	total := 0
	for _, code := range codes {
		total += code.freq
	}

	left := 0
	prevDiff := math.MaxInt
	bestPosition := 0

	for i := 0; i < len(codes)-1; i++ {
		left += codes[0].freq

		right := total - left

		diff := abs(right - left)
		if diff >= prevDiff {
			break
		}

		prevDiff = diff
		bestPosition = i + 1
	}

	return bestPosition
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
