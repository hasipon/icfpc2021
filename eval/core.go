package main

import "math/big"

type Int = big.Int

type Figure struct {
	Edges [][]Int `json:"edges"`
	Vertices [][]Int `json:"vertices"`
}

type Problem struct {
	Hole [][]Int `json:"hole"`
	Epsilon Int `json:"epsilon"`
	Figure Figure `json:"figure"`
}

func distance(a []Int, b[]Int) Int{
	var diffX, diffY, XX, YY Int
	diffX.Sub(&a[0], &b[0])
	diffY.Sub(&a[1], &b[1])
	XX.Mul(&diffX, &diffX)
	YY.Mul(&diffY, &diffY)
	var sum Int
	sum.Add(&XX, &YY)
	return sum
}

type Pose struct {
	Vertices [][]Int `json:"vertices"`
}

func dislike(problem *Problem, pose *Pose) Int {
	var sum Int
	for _, h := range problem.Hole {
		first := true
		var min Int
		for _, v := range pose.Vertices {
			tmp := distance(h, v)
			if first {
				min = tmp
			} else {
				if tmp.Cmp(&min) < 0 {
					min = tmp
				}
			}
			first = false
		}
		sum.Add(&sum, &min)
	}
	return sum
}

