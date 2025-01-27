package tokens

var tokens = map[string]int{}

func GetLimitByToken(token string) int {
	return tokens[token]
}
