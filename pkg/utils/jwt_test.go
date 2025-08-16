package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	// 测试数据
	tests := []struct {
		name     string
		userID   uint
		username string
		role     string
		wantErr  bool
	}{
		{
			name:     "valid token",
			userID:   1,
			username: "testuser",
			role:     "user",
			wantErr:  false,
		},
		{
			name:     "empty username",
			userID:   2,
			username: "",
			role:     "user",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 生成token
			token, err := GenerateToken(tt.userID, tt.username, tt.role)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// 解析token
			claims, err := ParseToken(token)
			assert.NoError(t, err)
			assert.NotNil(t, claims)

			// 验证claims
			assert.Equal(t, tt.userID, claims.UserID)
			assert.Equal(t, tt.username, claims.Username)
			assert.Equal(t, tt.role, claims.Role)
			assert.True(t, claims.ExpiresAt > time.Now().Unix())
		})
	}
}

func TestParseInvalidToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "invalid token",
			token: "invalid.token.string",
		},
		{
			name:  "empty token",
			token: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token)
			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}
