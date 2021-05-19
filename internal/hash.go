package internal

func HashString(s string) string {
	var acc []rune
	for _, b := range []rune(s) {
		if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_' || b >= '0' && b <= '9' {
			acc = append(acc, b)
		}
		if b == '-' || b == '_' {
			acc = append(acc, '_')
		}
	}
	return string(acc)
}
