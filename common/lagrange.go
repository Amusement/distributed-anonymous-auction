package common

import (
	"math/big"
)

type lagrangePoint struct {
	X *big.Int
	Y *big.Int
}

type lagrangePoints []lagrangePoint

// For some reason, the algorithm I found (I had to convert it to use big.Int and some minor tweaks)
//  is off by 1 all the time. Not exactly sure why, but I tested with couple of random points, but always off by 1
//  For now, we just minus 1 at the final result until I can figure out what the heck I did wrong
func (ps lagrangePoints) lagrange() *big.Int {
	result := new(big.Int)
	lenPS := len(ps)

	for i := 0; i < lenPS; i++ {
		Yterm := ps[i].Y
		for j := 0; j < lenPS; j++ {
			if j != i {
				numo := new(big.Int)
				numo.Sub(big.NewInt(0), ps[j].X)
				numo.Mul(Yterm, numo)

				deno := new(big.Int)
				deno.Sub(ps[i].X, ps[j].X)

				Yterm.Div(numo, deno)
			}
		}
		result.Add(result, Yterm)
	}
	return result
	//return result.Sub(result, big.NewInt(1)) // Hack.. I am off by 1..?
}

//func ComputeLagrange(points []Point) *big.Int {
func ComputeLagrange(compressedPoints []CompressedPoints) map[Price]*big.Int {
	lagrangeMap := make(map[Price]lagrangePoints)
	for _, cp := range compressedPoints {
		for k, v := range cp.Points {
			lagrangeMap[k] = append(lagrangeMap[k], lagrangePoint{
				X: big.NewInt(int64(v.X)),
				Y: v.Y.Val})
		}
	}
	interpolationMap := make(map[Price]*big.Int)
	for k, v := range lagrangeMap {
		interpolationMap[k] = v.lagrange()
	}
	return interpolationMap
}
