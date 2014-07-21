package avg

import (
	. "testing"
)

func BenchmarkRolling(b *B) {
	a := MovingAverage{Size: 16}

	for i := 0; i < b.N; i++ {
		a.Update(float64(i))
		a.Average()
	}
}

func BenchmarkAlpha64(b *B) {
	a := AlphaAverage{Alpha: 1 / 16}

	for i := 0; i < b.N; i++ {
		a.Update(float64(i))
		a.Average()
	}
}
