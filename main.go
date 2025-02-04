package main

import (
	"fmt"
	"math/big"
)

type Point struct {
	X big.Int
	Y big.Int
}

func (p1 *Point) addPoint(p2 *Point) Point {
	var res Point

	// Coords
	x1 := p1.X
	y1 := p1.Y
	x2 := p2.X
	y2 := p2.Y

	var topTmp, lowTmp big.Int

	topTmp.Sub(&y2, &y1)
	lowTmp.Sub(&x2, &x1)

	lowTmp.ModInverse(&lowTmp, &mod)

	// tmp = (y2 - y1)/(x2 - x1)
	var tmp big.Int
	tmp.Mul(&topTmp, &lowTmp)

	//tmp.Mod(&tmp, &mod)

	// tmp^2
	var tmpSquare big.Int
	tmpSquare.Mul(&tmp, &tmp)

	// tmp^2 - x2 - x1
	res.X.Sub(&tmpSquare, &x2)
	res.X.Sub(&res.X, &x1)

	res.X.Mod(&res.X, &mod)

	// y2 + m * (res.X - x2)
	res.Y.Sub(&res.X, &x2)
	res.Y.Mul(&tmp, &res.Y)
	res.Y.Add(&y2, &res.Y)

	res.Y.Mod(&res.Y, &mod)

	return res
}

func (p *Point) doublePoint() Point {
	// ((3x^2 + a)/(2y)) and ^2
	var gamma big.Int
	gamma.Mul(&p.X, &p.X)
	gamma.Mul(&gamma, big.NewInt(3))
	var y2 big.Int
	y2.Mul(&p.Y, big.NewInt(2))
	y2.ModInverse(&y2, &mod)
	gamma.Mul(&gamma, &y2)
	var gammaPow big.Int
	gammaPow.Mul(&gamma, &gamma)

	// Coords
	var res Point
	var x2 big.Int
	x2.Mul(&p.X, big.NewInt(2))
	res.X.Sub(&gammaPow, &x2)

	res.X.Mod(&res.X, &mod)

	var div big.Int
	div.Sub(&p.X, &res.X)
	res.Y.Mul(&gamma, &div)
	res.Y.Sub(&res.Y, &p.Y)

	res.Y.Mod(&res.Y, &mod)

	return res
}

func getPubKey(privKey big.Int) Point {
	// Public key
	var pubKey Point
	// Acumulate sum
	var res Point

	res.X.SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	res.Y.SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)

	// num >> 1 same with num / 2 (we can lost 1, if num is odd!)
	for ; privKey.Cmp(big.NewInt(0)) != 0; privKey.Rsh(&privKey, 1) {

		if privKey.Bit(0) == 1 {
			// Never used before
			if pubKey.X.Cmp(big.NewInt(0)) == 0 && pubKey.Y.Cmp(big.NewInt(0)) == 0 {
				pubKey.X = res.X
				pubKey.Y = res.Y
			} else {
				pubKey = pubKey.addPoint(&res)
			}
		}

		res = res.doublePoint()
	}

	return pubKey
}

// Global var
var mod big.Int

func main() {
	mod.SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	var privKey big.Int
	// Example key
	privKey.SetString("733115b84f8151f8e1f15e2c80fa938ed9c4da3b052ce79ae702db33e022fd91", 16)

	pubKey := getPubKey(privKey)

	fmt.Println("X: " + pubKey.X.Text(16))
	fmt.Println("Y: " + pubKey.Y.Text(16))
}