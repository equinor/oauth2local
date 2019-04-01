package storage

import "sync"

type MemoryStorage struct {
	rw           *sync.RWMutex
	authCode     string
	refreshToken string
	accessToken  string
	idToken      string
}

func Memory() *MemoryStorage {
	return &MemoryStorage{rw: new(sync.RWMutex)}
}

func (m *MemoryStorage) GetToken(tt TokenType) (string, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	switch tt {
	case AuthorizationCode:
		return m.authCode, nil

	}
	return "", nil
}

func (m *MemoryStorage) SetToken(tt TokenType, token string) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	switch tt {
	case AuthorizationCode:
		m.authCode = token

	}
	return nil
}

func (m *MemoryStorage) DeleteToken(tt TokenType) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	switch tt {
	case AuthorizationCode:
		m.authCode = ""

	}
	return nil
}
