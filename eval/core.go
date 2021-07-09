package main

import (
	"fmt"
	"math"
	"math/big"
)

type Int = big.Int

type Figure struct {
	Edges [][]int `json:"edges"`
	Vertices [][]Int `json:"vertices"`
}

type Problem struct {
	Hole [][]Int `json:"hole"`
	Epsilon Int `json:"epsilon"`
	Figure Figure `json:"figure"`
}

type Point = []Int

func toFloat(i Int) float64{
	f, _ := new(big.Float).SetInt(&i).Float64()
	return f
}

func intersect(p1, p2, p3, p4 Point) bool {
	tc1 := (toFloat(p1[0]) - toFloat(p2[0])) * (toFloat(p3[1]) - toFloat(p1[1])) + (toFloat(
		p1[1]) - toFloat(p2[1])) * (toFloat(p1[0]) - toFloat(p3[0]))
	tc2 := (toFloat(p1[0]) - toFloat(p2[0])) * (toFloat(p4[1]) - toFloat(p1[1])) + (toFloat(
		p1[1]) - toFloat(p2[1])) * (toFloat(p1[0]) - toFloat(p4[0]))
	td1 := (toFloat(p3[0]) - toFloat(p4[0])) * (toFloat(p1[1]) - toFloat(p3[1])) + (toFloat(
		p3[1]) - toFloat(p4[1])) * (toFloat(p3[0]) - toFloat(p1[0]))
	td2 := (toFloat(p3[0]) - toFloat(p4[0])) * (toFloat(p2[1]) - toFloat(p3[1])) + (toFloat(
		p3[1]) - toFloat(p4[1])) * (toFloat(p3[0]) - toFloat(p2[0]))
	return tc1*tc2<0 && td1*td2<0
}

func distance(a []Int, b[]Int) *Int{
	var diffX, diffY, XX, YY Int
	diffX.Sub(&a[0], &b[0])
	diffY.Sub(&a[1], &b[1])
	XX.Mul(&diffX, &diffX)
	YY.Mul(&diffY, &diffY)
	var sum Int
	sum.Add(&XX, &YY)
	return &sum
}

type Pose struct {
	Vertices [][]Int `json:"vertices"`
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
		origD, _ := new(big.Float).SetInt(distance(origV[i], origV[j])).Float64()
		nowD, _ := new(big.Float).SetInt(distance(nowV[i], nowV[j])).Float64()
		eps, _ := new(big.Float).SetInt(&problem.Epsilon).Float64()
		result := (math.Abs(nowD / origD - 1) * 1000000 <= eps)
		if !result {
			return false, fmt.Sprintf("Edge between (%d, " +
				"%d) has an invalid length: original: %.0f pose: %.0f", i, j, origD, nowD)
		}
		for ii, _ := range problem.Hole {
			h := problem.Hole
			jj := (ii+1) % len(h)
			if intersect(h[ii], h[jj], nowV[i], nowV[j]){
				return false, fmt.Sprintf("Edge between (%d, " +
					"%d) intersects: hole(%d, %d)", i, j, ii, jj)
			}

		}
	}
	return true, "OK"
}

