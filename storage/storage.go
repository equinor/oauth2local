package storage

type Storage interface {
	SetToken(TokenType, string) error
	GetToken(TokenType) (string, error)
	DeleteToken(TokenType) error
}

type TokenType int

const (
	Empty = iota
	RefreshToken
	AccessToken
	IDToken
)

func (tt TokenType) String() string {
	tts := []string{
		"Unkown",
		"Refresh Token",
		"Access Token",
		"ID Token",
	}

	if tt < Empty || tt > IDToken {
		return tts[0]
	}

	return tts[tt]
}
