package internal

// ValidLuhn check number is valid or not based on Luhn algorithm
func ValidLuhn(number int) bool {
	return (number%10+checksum(number/10))%10 == 0 //nolint:gomnd
}

//nolint:gomnd
func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur *= 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number /= 10
	}
	return luhn % 10
}
