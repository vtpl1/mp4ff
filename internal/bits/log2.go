package bits

// CeilLog2 returns nr bits needed to represent numbers 0 - n-1.
func CeilLog2(n uint) int {
	for i := range 32 {
		maxNr := uint(1 << i)
		if maxNr >= n {
			return i
		}
	}
	return 32
}
