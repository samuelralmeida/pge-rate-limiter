package tokens

var tokens = map[string]int{}

type tokenFetch struct{}

func NewTokenFetch() *tokenFetch {
	return &tokenFetch{}
}

func (t *tokenFetch) GetLimitByToken(token string) int {
	return tokens[token]
}
