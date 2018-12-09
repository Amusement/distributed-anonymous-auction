package common

import (
	"math/big"
)

type lagrangePoint struct {
	X *big.Int
	Y *big.Int
}

type lagrangePoints []lagrangePoint

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
}

func ComputeLagrange(compressedPoints []CompressedPoints) map[Price]BigInt {
	lagrangeMap := make(map[Price]lagrangePoints)
	for _, cp := range compressedPoints {
		for k, v := range cp.Points {
			lagrangeMap[k] = append(lagrangeMap[k], lagrangePoint{
				X: big.NewInt(int64(v.X)),
				Y: big.NewInt(0).SetBytes(v.Y.Val.Bytes())}) // copy over the value
		}
	}

	interpolationMap := make(map[Price]BigInt)
	for k, v := range lagrangeMap {
		interpolationMap[k] = BigInt{v.lagrange()}
	}
	return interpolationMap
}
