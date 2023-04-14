package main

import "sort"

type code struct {
	Index int
	Freq  int
	Code  string
}

func FanoEncoding(freq []int) []code {
	var codes []code
	i := 0
	for index, f := range freq {
		if f != 0 {
			codes = append(codes, code{
				Index: index,
				Freq:  f,
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
	sort.Slice(codes, func(i, j int) bool { return codes[i].Freq > codes[j].Freq })

	buildCodes(0, len(codes)-1, codes)
	return codes
}

func buildCodes(left int, right int, codes []code) {
	if left == right {
		return
	}
	if left+1 == right {
		codes[left].Code += "0"
		codes[right].Code += "1"
		return
	}
	mid := left
	sum := codes[mid].Freq
	for i := left + 1; i < right; i++ {
		sum += codes[i].Freq
		if sum > codes[len(codes)/2].Freq {
			mid = i
			break
		}
	}
	for i := left; i <= mid; i++ {
		codes[i].Code += "0"
	}
	for i := mid + 1; i <= right; i++ {
		codes[i].Code += "1"
	}
	buildCodes(left, mid, codes)
	buildCodes(mid+1, right, codes)
}
