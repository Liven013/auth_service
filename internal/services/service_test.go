package services

import (
	"auth_service/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockTokenGenerator struct{}

func (m *MockTokenGenerator) GenerateAccessToken(guid string, timeCreated int64) (string, error) {
	if guid == "error" {
		return "", errors.New("accessToken error")
	}
	return "mockAccessToken", nil
}

func (m *MockTokenGenerator) GenerateRefreshToken(guid string, ip string, timeCreated int64) (string, error) {
	if guid == "error" {
		return "", errors.New("refreshToken error")
	}
	return "mockRefreshToken", nil
}

func (m *MockTokenGenerator) StoreRefreshToken(guid string, refreshToken string) error {
	if guid == "error" {
		return errors.New("store error")
	}
	return nil
}

func (m *MockTokenGenerator) TestGenerateAccessToken(t *testing.T) {
	tests := []struct {
		name    string
		user    models.User
		ip      string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "token generation complete",
			user:    models.User{GUID: "validGUID"},
			ip:      "127.0.0.1",
			wantErr: false,
		},
		{
			name:    "AccessToken generation failure",
			user:    models.User{GUID: "error"},
			ip:      "127.0.0.1",
			wantErr: true,
			errMsg:  "accessToken error",
		},
		{
			name:    "RefreshToken generation failure",
			user:    models.User{GUID: "error"},
			ip:      "127.0.0.2",
			wantErr: true,
			errMsg:  "refreshToken error",
		},
		{
			name:    "StoreRefreshToken failure",
			user:    models.User{GUID: "error"},
			ip:      "127.0.0.3",
			wantErr: true,
			errMsg:  "store error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessToken, refreshToken, err := GeneratePairToken(tt.user, tt.ip)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, accessToken)
				assert.NotEmpty(t, refreshToken)
			}
		})
	}
}
