package tokens

var tokens = map[string]int{
	"token-example": 100,
}

type tokenFetch struct{}

func NewTokenFetch() *tokenFetch {
	return &tokenFetch{}
}

func (t *tokenFetch) GetLimitByToken(token string) int {
	return tokens[token]
}
