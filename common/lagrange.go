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
	result := new(big.Float)
	lenPS := len(ps)

	for i := 0; i < lenPS; i++ {
		Yterm := new(big.Float).SetInt(ps[i].Y)
		Zterm := big.NewFloat(1)
		for j := 0; j < lenPS; j++ {
			if j != i {
				numo := new(big.Float)
				numo.Sub(big.NewFloat(0), new(big.Float).SetInt(ps[j].X))
				Yterm.Mul(Yterm, numo)

				deno := new(big.Float)
				deno.Sub(new(big.Float).SetInt(ps[i].X), new(big.Float).SetInt(ps[j].X))
				Zterm.Mul(Zterm, deno)
			}
		}
		result.Add(result, new(big.Float).Quo(Yterm, Zterm))
	}
	intResult := new(big.Int)
	result.Int(intResult)
	return intResult
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
