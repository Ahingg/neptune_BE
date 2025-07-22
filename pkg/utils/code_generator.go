package utils

func GenerateAlphabetCode(n int) string {
	result := ""
	n += 1 // Since A starts from 1 (not 0)

	for n > 0 {
		n-- // Adjust for 0-based indexing
		letter := 'A' + (n % 26)
		result = string(letter) + result
		n /= 26
	}
	return result
}
