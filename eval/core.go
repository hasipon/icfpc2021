package main

import (
	"fmt"
	"math/big"
)

type Int = big.Int

type Figure struct {
	Edges [][]int `json:"edges"`
	Vertices [][]*Int `json:"vertices"`
}

type Problem struct {
	Hole [][]*Int `json:"hole"`
	Epsilon Int `json:"epsilon"`
	Figure Figure `json:"figure"`
}

type Point = []*Int

func intersect(p []Point) bool {
	sub := func(i, j, k, l int) *Int {
		i--
		k--
		return new(Int).Sub(p[i][j], p[k][l])
	}
	tc1_A := new(Int).Mul(sub(1, 0, 2, 0), sub(3, 1, 1, 1))
	tc1_B := new(Int).Mul(sub(1, 1, 2, 1), sub(1, 0, 3, 0))
	tc1 := new(Int).Add(tc1_A, tc1_B)


	tc2_A := new(Int).Mul(sub(1, 0, 2, 0), sub(4, 1, 1, 1))
	tc2_B := new(Int).Mul(sub(1, 1, 2, 1), sub(1, 0, 4, 0))
	tc2 := new(Int).Add(tc2_A, tc2_B)


	td1_A := new(Int).Mul(sub(3, 0, 4, 0), sub(1, 1, 3, 1))
	td1_B := new(Int).Mul(sub(3, 1, 4, 1), sub(3, 0, 1, 0))
	td1 := new(Int).Add(td1_A, td1_B)

	td2_A := new(Int).Mul(sub(3, 0, 4, 0), sub(2, 1, 3, 1))
	td2_B := new(Int).Mul(sub(3, 1, 4, 1), sub(3, 0, 2, 0))
	td2 := new(Int).Add(td2_A, td2_B)

	tc := new(Int).Mul(tc1, tc2)
	td := new(Int).Mul(td1, td2)
	if tc.Cmp(new(Int)) < 0 && td.Cmp(new(Int)) < 0 {
		return true
	}
	return false
}

func distance(a []*Int, b[]*Int) *Int{
	var diffX, diffY, XX, YY Int
	diffX.Sub(a[0], b[0])
	diffY.Sub(a[1], b[1])
	XX.Mul(&diffX, &diffX)
	YY.Mul(&diffY, &diffY)
	var sum Int
	sum.Add(&XX, &YY)
	return &sum
}

type Pose struct {
	Vertices [][]*Int `json:"vertices"`
}

func dislike(problem *Problem, pose *Pose) *Int {
	sum := new(Int)
	for _, h := range problem.Hole {
		first := true
		min := new(Int)
		for _, v := range pose.Vertices {
			tmp := distance(h, v)
			if first {
				min = tmp
			} else {
				if tmp.Cmp(min) < 0 {
					min = tmp
				}
			}
			first = false
		}
		sum.Add(sum, min)
	}
	return sum
}

func validate(problem *Problem, pose *Pose) (bool, string) {
	origV := problem.Figure.Vertices
	nowV := pose.Vertices
	if len(origV) != len(nowV) {
		return false, "mismatch length"
	}
	for _, e := range problem.Figure.Edges {
		i := e[0]
		j := e[1]
		origD := distance(origV[i], origV[j])
		nowD := distance(nowV[i], nowV[j])
		var diff *Int
		if nowD.Cmp(origD) >= 0 {
			diff = new(Int).Sub(nowD, origD)
		} else {
			diff = new(Int).Sub(origD, nowD)
		}
		diff = diff.Mul(diff, new(Int).SetInt64(1000000))
		eps := new(Int).Mul(&problem.Epsilon, origD)
		result := diff.Cmp(eps)
		if result > 0 {
			return false, fmt.Sprintf("Edge between (%d, " +
				"%d) has an invalid length: original: %d pose: %d", i, j, origD, nowD)
		}
		for ii, _ := range problem.Hole {
			h := problem.Hole
			jj := (ii+1) % len(h)
			if intersect([]Point{h[ii], h[jj], nowV[i], nowV[j]}){
				return false, fmt.Sprintf("Edge between (%d, " +
					"%d) intersects: hole(%d, %d)", i, j, ii, jj)
			}

		}
	}
	return true, "OK"
}

