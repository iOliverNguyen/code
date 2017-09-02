package main

import (
	"fmt"
	"math/bits"
)

var (
	A = []int{1, 2, 3, 4, 5, 6}
	N uint16
	G []int // GCD
	M []MM

	max    int
	result []int16
)

type MM struct {
	M int    // max
	C uint32 // child
	L int16  // left
	R int16  // right
}

func output(s uint32, m MM) {
	max = m.M
	result = make([]int16, 0, N)
	routput(s, m, m.L, m.R)

	fmt.Print("\nMax:    ", m.M)
	fmt.Print("\nInput:  ")
	for _, n := range A {
		fmt.Print(n, " ")
	}
	fmt.Print("\nOutput: ")
	for _, n := range result {
		fmt.Print(A[n], " ")
	}
	fmt.Print("\nGCD:    ")
	for i := uint16(1); i < N; i++ {
		fmt.Print(GCD(result[i], result[i-1]), " ")
	}
	fmt.Println()
}

func routput(s uint32, m MM, L, R int16) {
	switch {
	case s == 0:
		return
	case m.L == L, m.R == R:
	case m.R == L, m.L == R:
		m.L, m.R = m.R, m.L
	default:
		return
	}

	c1 := m.C
	c2 := inv(s, c1)
	routput(c1, M[c1], m.L, -2)
	routput(c2, M[c2], m.L, -2)
	result = append(result, tip(s))
	routput(c1, M[c1], -2, m.R)
	routput(c2, M[c2], -2, m.R)
}

func debug(s uint32, m MM) {
	fmts, fmtb, fmtd := fmt.Sprintf("%% %ds ", N), fmt.Sprintf("%%0%db ", N), `% 6d `
	fmt.Printf(fmts+fmts+fmts+"% 6s % 6s % 6s\n", "s", "c1", "c2", ".M", ".L", ".R")
	format := fmtb + fmtb + fmtb + fmtd + fmtd + fmtd + "\n"
	for i := uint32(0); i < uint32(len(M)); i++ {
		fmt.Printf(format, i, M[i].C, inv(i, M[i].C), M[i].M, M[i].L, M[i].R)
	}
	if s > 0 {
		fmt.Println("...")
		fmt.Printf(format, s, m.C, inv(s, m.C), m.M, m.L, m.R)
	}
}

func debugGCD() {
	for i := int16(0); i < int16(N); i++ {
		for j := int16(0); j < int16(N); j++ {
			s := GCD(i, j)
			fmt.Print(s, " ")
		}
		fmt.Println()
	}
}

func gcd(a, b int) int {
	if b > a {
		a, b = b, a
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func GCD(a, b int16) int {
	if a == -1 || b == -1 {
		return 0
	}
	return G[int(a*int16(N)+b)]
}

func down(s uint32) uint32 {
	return s & (1<<uint(bits.Len32(s)-1) - 1)
}

func down2(s uint32) uint32 {
	return s & (1<<uint(bits.Len32(s)-2) - 1)
}

func tip(s uint32) int16 {
	return int16(bits.Len32(s) - 1)
}

func span(n uint16) uint32 {
	return 1<<n - 1
}

func inv(s, c uint32) uint32 {
	return ^c & down(s)
}

func calc(s uint32) MM {
	n := tip(s)
	mask := down(s)
	mx := MM{-1, 0, 0, 0}
	var m [4]MM
	for c, C := uint32(0), down2(s); c <= C; c++ {
		if s|c != s {
			continue
		}
		d := ^c & mask
		m[0] = MM{M[c].M + M[d].M + GCD(M[c].L, n) + GCD(M[d].L, n), c, M[c].R, M[d].R}
		m[1] = MM{M[c].M + M[d].M + GCD(M[c].L, n) + GCD(M[d].R, n), c, M[c].R, M[d].L}
		m[2] = MM{M[c].M + M[d].M + GCD(M[c].R, n) + GCD(M[d].L, n), c, M[c].L, M[d].R}
		m[3] = MM{M[c].M + M[d].M + GCD(M[c].R, n) + GCD(M[d].R, n), c, M[c].L, M[d].L}
		for i := 0; i < 4; i++ {
			if m[i].M > mx.M {
				mx = m[i]
			}
		}
	}
	if mx.L == -1 {
		mx.L = n
	}
	if mx.R == -1 {
		mx.R = n
	}
	return mx
}

func main() {
	N = uint16(len(A))
	if N < 2 || N > 20 {
		panic("N: out of range")
	}

	G = make([]int, N*N)
	for a := uint16(0); a < N; a++ {
		for b := uint16(0); b < a; b++ {
			s := gcd(A[a], A[b])
			G[a*N+b] = s
			G[b*N+a] = s
		}
	}
	// debugGCD()

	M = make([]MM, 1<<(N-1))
	M[0] = MM{0, 0, -1, -1}
	for i := 1; i < len(M); i++ {
		M[i].M = -1
	}
	for i := uint16(0); i < N-1; i++ {
		M[1<<i] = MM{0, 0, int16(i), int16(i)}
	}
	// debug(0, MM{})

	for i, S := uint32(3), span(N-1); i <= S; i++ {
		if M[i].M >= 0 {
			continue
		}
		M[i] = calc(i)
	}
	s := span(N)
	m := calc(s)
	// debug(s, m)
	output(s, m)
}
