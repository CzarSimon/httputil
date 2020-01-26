package jwt_test

import (
	"testing"
	"time"

	"github.com/CzarSimon/httputil/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWTIssueAndVerify(t *testing.T) {
	creds := jwt.Credentials{
		Issuer: "issuer-name",
		Secret: "super-secret-token",
	}

	issuer := jwt.NewIssuer(creds)
	verifier := jwt.NewVerifier(creds, time.Minute)

	wrongIssuerCreds := jwt.Credentials{
		Issuer: "wrong-issuer-name",
		Secret: "super-secret-token",
	}
	wrongIssuerVerfifier := jwt.NewVerifier(wrongIssuerCreds, time.Minute)

	wrongSecretCreds := jwt.Credentials{
		Issuer: "issuer-name",
		Secret: "super-secret-token-but-wrong",
	}
	wrongSecretVerfifier := jwt.NewVerifier(wrongSecretCreds, time.Minute)

	tests := []struct {
		name                string
		issuer              jwt.Issuer
		verifier            jwt.Verifier
		user                jwt.User
		tokenLifetime       time.Duration
		wantErr             error
		wantVerificationErr error
	}{
		{
			name:     "happy-path-0",
			issuer:   issuer,
			verifier: verifier,
			user: jwt.User{
				ID:    "user-id-t0",
				Roles: []string{"USER"},
			},
			tokenLifetime:       time.Hour,
			wantErr:             nil,
			wantVerificationErr: nil,
		},
		{
			name:     "happy-path-1",
			issuer:   issuer,
			verifier: verifier,
			user: jwt.User{
				ID:    "user-id-t1",
				Roles: []string{"USER"},
			},
			tokenLifetime:       time.Hour,
			wantErr:             nil,
			wantVerificationErr: nil,
		},
		{
			name:   "sad-path-missing-id",
			issuer: issuer,
			user: jwt.User{
				Roles: []string{"ADMIN"},
			},
			tokenLifetime: time.Hour,
			wantErr:       jwt.ErrInvalidTokenContent,
		},
		{
			name:   "sad-path-missing-role",
			issuer: issuer,
			user: jwt.User{
				ID: "user-id-t3",
			},
			tokenLifetime: time.Hour,
			wantErr:       jwt.ErrInvalidTokenContent,
		},
		{
			name:     "sad-path-expired-token",
			issuer:   issuer,
			verifier: verifier,
			user: jwt.User{
				ID:    "user-id-t4",
				Roles: []string{"SYSTEM"},
			},
			tokenLifetime:       -2 * time.Minute,
			wantErr:             nil,
			wantVerificationErr: jwt.ErrExpiredToken,
		},
		{
			name:     "sad-path-wrong-issuer-name-in-verifier",
			issuer:   issuer,
			verifier: wrongIssuerVerfifier,
			user: jwt.User{
				ID:    "user-id-t5",
				Roles: []string{"SYSTEM"},
			},
			tokenLifetime:       time.Hour,
			wantErr:             nil,
			wantVerificationErr: jwt.ErrInvalidToken,
		},
		{
			name:     "sad-path-wrong-secret-name-in-verifier",
			issuer:   issuer,
			verifier: wrongSecretVerfifier,
			user: jwt.User{
				ID:    "user-id-t6",
				Roles: []string{"ANONYMOUS"},
			},
			tokenLifetime:       time.Hour,
			wantErr:             nil,
			wantVerificationErr: jwt.ErrInvalidToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawToken, err := tt.issuer.Issue(tt.user, tt.tokenLifetime)
			if err != tt.wantErr {
				t.Errorf("jwt.Issuer.Issue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != nil {
				return
			}

			if rawToken == "" {
				t.Errorf("jwt.Issuer.Issue() token is empty, should not be")
				return
			}

			user, err := tt.verifier.Verify(rawToken)
			if err != tt.wantVerificationErr {
				t.Errorf("jwt.Verifier.Verify() error = %v, wantVerificationErr %v", err, tt.wantVerificationErr)
				return
			}

			if tt.wantVerificationErr != nil {
				return
			}

			assertJWTUser(t, tt.user, user)
		})
	}
}

func assertJWTUser(t *testing.T, expected, found jwt.User) {
	assert := assert.New(t)

	assert.Equal(expected.ID, found.ID)
	assert.Equal(expected.Roles, found.Roles)
}
