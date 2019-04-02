package storage
import (
	"fmt"
	"sync"
)

type MemoryStorage struct {
	rw *sync.RWMutex

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
	case AccessToken:
		return m.accessToken, nil
	case IDToken:
		return m.idToken, nil
	case RefreshToken:
		return m.refreshToken, nil
	}
	return "", fmt.Errorf("No %v in store", tt)
}

func (m *MemoryStorage) SetToken(tt TokenType, token string) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	switch tt {
	case AccessToken:
		m.accessToken = token
	case RefreshToken:
		m.refreshToken = token
	case IDToken:
		m.idToken = token
	default:
		return fmt.Errorf("No store for %v ", tt)
	}
	return nil
}

func (m *MemoryStorage) DeleteToken(tt TokenType) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	switch tt {
	case AccessToken:
		m.accessToken = ""
	case RefreshToken:
		m.refreshToken = ""
	case IDToken:
		m.idToken = ""
	default:
		return fmt.Errorf("No %v in store", tt)
	}
	return nil
}
