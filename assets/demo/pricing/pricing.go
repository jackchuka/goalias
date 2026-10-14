package pricing

func Total(cents ...int) int {
	sum := 0
	for _, c := range cents {
		sum += c
	}
	return sum
}
