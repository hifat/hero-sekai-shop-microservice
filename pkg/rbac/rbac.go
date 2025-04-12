package rbac

func IntToBinary(n, length int) []int {
	binaries := make([]int, length)
	for i := 0; i < length; i++ {
		binaries[i] = n % 2
		n /= 2
	}

	return binaries
}
