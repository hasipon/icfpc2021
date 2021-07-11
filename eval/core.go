package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
)

type Int = big.Int

type Figure struct {
	Edges    [][]int  `json:"edges"`
	Vertices [][]*Int `json:"vertices"`
}

type ProblemBonus struct {
	Bonus    string `json:"bonus"`
	Problem  int    `json:"problem"`
	Position []*Int `json:"position"`
}

type Problem struct {
	Hole            [][]*Int       `json:"hole"`
	Epsilon         *Int           `json:"epsilon"`
	Figure          Figure         `json:"figure"`
	Bonuses         []ProblemBonus `json:"bonuses"`
	WallHacked      bool
	Globalist       bool
	OriginalEdgeNum int
}

type Point = []*Int

func dot(a, b Point) *Int {
	x := new(Int).Mul(a[0], b[0])
	y := new(Int).Mul(a[1], b[1])
	return new(Int).Add(x, y)
}

func det(a, b Point) *Int {
	x := new(Int).Mul(a[0], b[1])
	y := new(Int).Mul(a[1], b[0])
	return new(Int).Sub(x, y)
}

var zero = new(Int).SetInt64(0)

const (
	FRONT = 1
	RIGHT = 2
	BACK  = 4
	LEFT  = 8
	ON    = 16
)

func ccw(a, b, c Point) int {
	b_a := Point{new(Int).Sub(b[0], a[0]), new(Int).Sub(b[1], a[1])}
	c_a := Point{new(Int).Sub(c[0], a[0]), new(Int).Sub(c[1], a[1])}
	s := det(b_a, c_a)
	if s == zero {
		return ON
	}
	if s.Cmp(zero) < 0 {
		return RIGHT
	}
	return LEFT
}

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

func distance(a []*Int, b []*Int) *Int {
	var diffX, diffY, XX, YY Int
	diffX.Sub(a[0], b[0])
	diffY.Sub(a[1], b[1])
	XX.Mul(&diffX, &diffX)
	YY.Mul(&diffY, &diffY)
	var sum Int
	sum.Add(&XX, &YY)
	return &sum
}

type Bonus struct {
	Bonus   string `json:"bonus,omitempty"`
	Problem int    `json:"problem,omitempty"`
	Edge    []int  `json:"edge,omitempty"`
}

type Pose struct {
	Vertices [][]*Int `json:"vertices,omitempty"`
	Bonuses  []*Bonus `json:"bonuses,omitempty"`
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

func include(problem *Problem, p Point) bool {
	x := p[0]
	y := p[1]
	cnt := 0
	for i, _ := range problem.Hole {
		j := (i + 1) % len(problem.Hole)
		x0 := new(Int).Set(problem.Hole[i][0])
		y0 := new(Int).Set(problem.Hole[i][1])
		x1 := new(Int).Set(problem.Hole[j][0])
		y1 := new(Int).Set(problem.Hole[j][1])

		x0 = x0.Sub(x0, x)
		y0 = y0.Sub(y0, y)
		x1 = x1.Sub(x1, x)
		y1 = y1.Sub(y1, y)

		cv := new(Int).Add(new(Int).Mul(x0, x1), new(Int).Mul(y0, y1))
		sv := new(Int).Sub(new(Int).Mul(x0, y1), new(Int).Mul(x1, y0))
		if sv.Cmp(zero) == 0 && cv.Cmp(zero) <= 0 {
			return true
		}

		if y0.Cmp(y1) < 0 {
		} else {
			tmp := x0
			x0 = x1
			x1 = tmp
			tmp = y0
			y0 = y1
			y1 = tmp
		}

		if y0.Cmp(zero) <= 0 && zero.Cmp(y1) < 0 {
			a := new(Int).Mul(x0, new(Int).Sub(y1, y0))
			b := new(Int).Mul(y0, new(Int).Sub(x1, x0))
			if b.Cmp(a) < 0 {
				cnt++
			}
		}
	}
	return cnt%2 == 1
}

func applyBonus(problem *Problem, pose *Pose) *Problem {
	problem.OriginalEdgeNum = len(problem.Figure.Edges)
	for _, b := range pose.Bonuses {
		if b.Bonus == "BREAK_A_LEG" {
			target := make(map[int]bool)
			v1 := problem.Figure.Vertices[b.Edge[0]]
			v2 := problem.Figure.Vertices[b.Edge[1]]
			target[b.Edge[0]] = true
			target[b.Edge[1]] = true
			mid := Point{
				new(Int).Div(new(Int).Add(v1[0], v2[0]), new(Int).SetInt64(2)),
				new(Int).Div(new(Int).Add(v1[1], v2[1]), new(Int).SetInt64(2)),
			}
			newVertexId := len(problem.Figure.Vertices)
			problem.Figure.Vertices = append(problem.Figure.Vertices, mid)
			newEdges := [][]int{
				[]int{b.Edge[0], newVertexId},
				[]int{newVertexId, b.Edge[1]},
			}
			for _, e := range problem.Figure.Edges {
				if target[e[0]] && target[e[1]] {
					continue
				}
				newEdges = append(newEdges, e)
			}
			problem.Figure.Edges = newEdges
		} else if b.Bonus == "GLOBALIST" {
			problem.Globalist = true
		} else if b.Bonus == "WALLHACK" {
			problem.WallHacked = true
		} else {
			log.Printf("Unknown bounus: %s", b.Bonus)
		}
	}
	return problem
}

func validate(problem *Problem, pose *Pose) (bool, string) {
	if len(pose.Bonuses) > 1 {
		return false, "too many bonuses"
	}

	origV := problem.Figure.Vertices
	nowV := pose.Vertices
	if len(origV) != len(nowV) {
		return false, "mismatch length"
	}

	outsideVerticesCnt := make(map[int]bool)

	for i, p := range pose.Vertices {
		if !include(problem, p) {
			outsideVerticesCnt[i] = true
		}
	}

	numOfAllowOutVertex := 0
	if problem.WallHacked {
		numOfAllowOutVertex++
	}

	if len(outsideVerticesCnt) > numOfAllowOutVertex {
		ids := ""
		cnt := 0
		for k, _ := range outsideVerticesCnt {
			ids += fmt.Sprintf("%d,", k)
			cnt++
			if cnt > numOfAllowOutVertex {
				break
			}
		}
		return false, fmt.Sprintf("vertex(%s)(0-based) is out of the hole", ids)
	}

	holePointsInd := make(map[string]int)
	for i, h := range problem.Hole {
		holePointsInd[h[0].String()+","+h[1].String()] = i
	}

	outEdges := [][]int{}

	for _, e := range problem.Figure.Edges {
		v1 := pose.Vertices[e[0]]
		v2 := pose.Vertices[e[1]]
		two := new(Int).SetInt64(2)
		sumX := new(Int).Mul(two, new(Int).Add(v1[0], v2[0]))
		sumY := new(Int).Mul(two, new(Int).Add(v1[1], v2[1]))
		mid := Point{new(Int).Div(sumX, two), new(Int).Div(sumY, two)}
		for _, v := range problem.Hole {
			v[0] = v[0].Mul(v[0], two)
			v[1] = v[1].Mul(v[1], two)
		}
		midInc := include(problem, mid)
		for _, v := range problem.Hole {
			v[0] = v[0].Div(v[0], two)
			v[1] = v[1].Div(v[1], two)
		}
		if !midInc {
			outEdges = append(outEdges, e)
		}
	}

	var globalEpsSum float64
	for _, e := range problem.Figure.Edges {
		i := e[0]
		j := e[1]
		origD := distance(origV[i], origV[j])
		nowD := distance(nowV[i], nowV[j])
		if problem.Globalist {
			nowDf, _ := new(big.Float).SetInt(nowD).Float64()
			origDf, _ := new(big.Float).SetInt(origD).Float64()
			globalEpsSum += math.Abs(nowDf/origDf-1) * 1000000
		} else {
			var diff *Int
			if nowD.Cmp(origD) >= 0 {
				diff = new(Int).Sub(nowD, origD)
			} else {
				diff = new(Int).Sub(origD, nowD)
			}
			diff = diff.Mul(diff, new(Int).SetInt64(1000000))
			eps := new(Int).Mul(problem.Epsilon, origD)
			result := diff.Cmp(eps)
			if result > 0 {
				return false, fmt.Sprintf("Edge between (%d, "+
					"%d) has an invalid length: original: %d pose: %d", i, j, origD, nowD)
			}
		}
		for ii, _ := range problem.Hole {
			h := problem.Hole
			jj := (ii + 1) % len(h)
			if intersect([]Point{h[ii], h[jj], nowV[i], nowV[j]}) {
				outEdges = append(outEdges, e)
			}
		}
	}
	if problem.Globalist {
		epsLimit, _ := new(big.Float).SetInt(problem.Epsilon).Float64()
		epsLimit = epsLimit * float64(problem.OriginalEdgeNum)
		if globalEpsSum > epsLimit {
			return false, fmt.Sprintf("global epsilon budget exceeded. limit %f, now %f",
				epsLimit, globalEpsSum)
		}
	}
	if problem.WallHacked {
		if len(outsideVerticesCnt) == 1 {
			hackedVertex := -1
			for k, _ := range outsideVerticesCnt {
				hackedVertex = k
				for _, e := range outEdges {
					if e[0] != k && e[1] != k {
						return false, fmt.Sprintf("Invalid edge between (%d, "+
							"%d)", e[0], e[1])
					}
				}
			}
			return true, fmt.Sprintf("OK: hacked vertex is %d(0-based)", hackedVertex)
		} else {
			cnt := make(map[int]int)
			for _, e := range outEdges {
				cnt[e[0]]++
				cnt[e[1]]++
			}
			found := false
			hackedVertex := -1
			for k, v := range cnt {
				if v == len(outEdges) {
					found = true
					hackedVertex = k
				}
			}
			if !found {
				return false, fmt.Sprintf("Too many Invalid edges")
			} else {
				return true, fmt.Sprintf("OK: hacked vertex is %d(0-based)", hackedVertex)
			}
		}
	} else {
		if len(outEdges) == 0 {
			return true, "OK"
		} else {
			return false, fmt.Sprintf("invalid edge(%d, %d)(0-based)", outEdges[0][0],
				outEdges[0][1])
		}
	}
}
