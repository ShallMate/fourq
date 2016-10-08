package fourq

import (
	"testing"

	"crypto/rand"
	"math/big"
)

func TestIsOnCurve(t *testing.T) {
	x, y := g.Int()
	if !IsOnCurve(x, y) {
		t.Fatal("Generator is not on curve.")
	}

	x2, y2 := ScalarMult(x, y, unpack(Order))
	if !IsOnCurve(x2, y2) {
		t.Fatal("Identity point is not on curve.")
	}

	k := make([]byte, 32)
	rand.Read(k)

	x3, y3 := ScalarMult(x, y, k)
	if !IsOnCurve(x3, y3) {
		t.Fatal("Random multiple of generator is not on curve.")
	}
}

func TestScalarMultOrder(t *testing.T) {
	x, y := g.Int()
	x3, y3 := ScalarMult(x, y, unpack(Order))

	y3real := "100000000000000000000000000000000000000000000000000000000000000"
	if x3.Sign() != 0 || y3.Text(16) != y3real {
		t.Fatal("ScalarMult(Generator, Order) was not identity.")
	}
}

func TestScalarMult(t *testing.T) {
	// Source: https://github.com/bifurcation/fourq/blob/master/impl/curve4q.py#L549
	scalar := [4]uint64{0x3AD457AB55456230, 0x3A8B3C2C6FD86E0C, 0x7E38F7C9CFBB9166, 0x0028FD6CBDA458F0}

	pt := g
	for i := 0; i < 1000; i++ {
		scalar[1] = scalar[2]
		scalar[2] += scalar[0]
		scalar[2] &= 0xffffffffffffffff

		k := new(big.Int).SetUint64(scalar[3])
		k.Lsh(k, 64)
		k.Add(k, new(big.Int).SetUint64(scalar[2]))
		k.Lsh(k, 64)
		k.Add(k, new(big.Int).SetUint64(scalar[1]))
		k.Lsh(k, 64)
		k.Add(k, new(big.Int).SetUint64(scalar[0]))

		sum := newPoint()
		tmp := newPoint()

		for i := k.BitLen() - 1; i >= 0; i-- {
			tmp.Dbl(sum)
			if k.Bit(i) != 0 {
				sum.Add(tmp, pt)
			} else {
				sum.Set(tmp)
			}
		}

		pt = sum
	}

	x, y := pt.Int()
	realX := "ef4b49bd77b4d2df1b4ac9bf2b127c2559c4377254939576011fb1b50cf89b46"
	realY := "44336f9967501c286c930e7c81b3010945125f9129c4e84f10e2acac8e940b57"

	if x.Text(16) != realX || y.Text(16) != realY {
		t.Fatal("Point is wrong!")
	}
}

func BenchmarkScalarBaseMult(b *testing.B) {
	k := make([]byte, 32)
	rand.Read(k)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ScalarBaseMult(k)
	}
}

func BenchmarkPointAdd(b *testing.B) {
	pt1, pt2 := g, newPoint()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt2.Add(pt1, pt2)
	}
}

func BenchmarkPointDbl(b *testing.B) {
	pt := g

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Dbl(pt)
	}
}
