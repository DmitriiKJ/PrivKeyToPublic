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

	if x1.Cmp(&x2) == 0 && y1.Cmp(&y2) == 0 {
        return p1.doublePoint()
    }

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

	// y2 + tmp * (res.X - x2)
	res.Y.Sub(&x1, &res.X)
	res.Y.Mul(&tmp, &res.Y)
	res.Y.Sub(&res.Y, &y1)

	res.Y.Mod(&res.Y, &mod)

	return res
}

func (p *Point) doublePoint() Point {
	// ((3x^2)/(2y)) and ^2
	var lambda big.Int
	lambda.Mul(&p.X, &p.X)
	lambda.Mul(&lambda, big.NewInt(3))
	var y2 big.Int
	y2.Mul(&p.Y, big.NewInt(2))
	y2.ModInverse(&y2, &mod)
	lambda.Mul(&lambda, &y2)

	var lambdaPow big.Int
	lambdaPow.Mul(&lambda, &lambda)

	// Coords
	var res Point
	var x2 big.Int
	x2.Mul(&p.X, big.NewInt(2))
	res.X.Sub(&lambdaPow, &x2)

	res.X.Mod(&res.X, &mod)

	res.Y.Sub(&p.X, &res.X)
	res.Y.Mul(&lambda, &res.Y)
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
	
	for ; privKey.Cmp(big.NewInt(1)) != 0; privKey.Rsh(&privKey, 1){

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

	pubKey = pubKey.addPoint(&res)

	return pubKey
}

// Global var
var mod big.Int

func main() {
	mod.SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	var privKey big.Int
	// Example key
	privKey.SetString("fbdfa5e4a198c9b24003200452b410a9000c0ea236e2ca9657a15ed376dc416d", 16)

	pubKey := getPubKey(privKey)

	fmt.Println("Public key: 04" + pubKey.X.Text(16) + pubKey.Y.Text(16))
}