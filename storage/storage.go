package storage

type Storage interface {
	SetToken(TokenType, string) error
	GetToken(TokenType) (string, error)
	DeleteToken(TokenType) error
}

type TokenType int

const (
	Empty = iota
	AuthorizationCode
	RefreshToken
	AccessToken
	IDToken
)
