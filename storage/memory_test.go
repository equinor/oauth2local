package storage

import (
	"testing"
)

func EmptyMemory() *MemoryStorage {
	return Memory()
}

func PrefilledMemory() *MemoryStorage {
	m := Memory()
	err := m.SetToken(AccessToken, "accesstokenstring")
	if err != nil {
		panic(err)
	}
	err = m.SetToken(IDToken, "idtokenstring")
	if err != nil {
		panic(err)
	}
	err = m.SetToken(RefreshToken, "refreshtokenstring")
	if err != nil {
		panic(err)
	}
	return m
}
func TestMemoryStorage_GetToken(t *testing.T) {
	type args struct {
		tt TokenType
	}
	tests := []struct {
		name    string
		m       Storage
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Load token", EmptyMemory(), args{AccessToken}, "", true},
		{"Load token", PrefilledMemory(), args{AccessToken}, "accesstokenstring", false},
		{"Load token", PrefilledMemory(), args{RefreshToken}, "refreshtokenstring", false},
		{"Load token", PrefilledMemory(), args{IDToken}, "idtokenstring", false},
		{"Load token", PrefilledMemory(), args{100}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetToken(tt.args.tt)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryStorage.GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemoryStorage.GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStorage_DeleteToken(t *testing.T) {
	type args struct {
		tt TokenType
	}
	tests := []struct {
		name    string
		m       Storage
		args    args
		wantErr bool
	}{
		{"Delete token", EmptyMemory(), args{AccessToken}, false},
		{"Delete access token", PrefilledMemory(), args{AccessToken}, false},
		{"Delete refresh token", PrefilledMemory(), args{RefreshToken}, false},
		{"Delete id token", PrefilledMemory(), args{IDToken}, false},
		{"Delete unknown token type", PrefilledMemory(), args{100}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.DeleteToken(tt.args.tt); (err != nil) != tt.wantErr {
				t.Errorf("MemoryStorage.DeleteToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryStorage_SetToken(t *testing.T) {

	m := EmptyMemory()
	type args struct {
		tt    TokenType
		token string
	}
	tests := []struct {
		name    string
		m       Storage
		args    args
		wantErr bool
	}{
		{"Store access token", m, args{AccessToken, "eqyeqy"}, false},
		{"Store id token", m, args{IDToken, "eqyeqy"}, false},
		{"Store id token again ", m, args{IDToken, "sssss"}, false},
		{"Store refresh token", m, args{RefreshToken, "eqyeqy"}, false},
		{"Store unknown token type", m, args{100, "eqyeqy"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetToken(tt.args.tt, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("MemoryStorage.SetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
