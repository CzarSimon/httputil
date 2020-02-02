package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CzarSimon/httputil/logger"
	"go.uber.org/zap"
	"gopkg.in/square/go-jose.v2"
	josejwt "gopkg.in/square/go-jose.v2/jwt"
)

var log = logger.MustGetLogger("httputil/jwt", zap.InfoLevel)

const roleDelimiter = ";"

// Common roles
const (
	SystemRole    = "SYSTEM"
	AdminRole     = "ADMIN"
	AnonymousRole = "ANONYMOUS"
)

// Common errors.
var (
	ErrInvalidTokenContent = errors.New("invalid token content")
	ErrInvalidToken        = errors.New("token is invalid")
	ErrExpiredToken        = errors.New("token has expired")
)

// Credentials credentials to issue and verify JWT tokens.
type Credentials struct {
	Issuer string `json:"issuer"`
	Secret string `json:"secret"`
}

// User user authenicatated in JWT.
type User struct {
	ID    string
	Roles []string
}

func (u User) String() string {
	return fmt.Sprintf("User(id=%s, roles=%v)", u.ID, u.Roles)
}

// IsSystem checks if a user has the SYSTEM role.
func (u User) IsSystem() bool {
	return u.HasRole(SystemRole)
}

// IsAdmin checks if a user has the ADMIN role.
func (u User) IsAdmin() bool {
	return u.HasRole(AdminRole)
}

// IsAnonymous checks if a user has the ANONYMOUS role.
func (u User) IsAnonymous() bool {
	return u.HasRole(AnonymousRole)
}

// HasRole checks if a users has a given role.
func (u User) HasRole(candidate string) bool {
	for _, role := range u.Roles {
		if role == candidate {
			return true
		}
	}

	return false
}

// Issuer interface for issuing auth tokens
type Issuer interface {
	Issue(user User, lifetime time.Duration) (string, error)
}

// NewIssuer creates a new Issuer using the default implementation.
func NewIssuer(creds Credentials) Issuer {
	sigingKey := jose.SigningKey{Algorithm: jose.HS256, Key: []byte(creds.Secret)}
	signer, err := jose.NewSigner(sigingKey, nil)
	if err != nil {
		log.Panic("Failed to create jose.Signer.", zap.Error(err))
	}
	return &jwtIssuer{
		name:   creds.Issuer,
		signer: signer,
	}
}

// Verifier interface for verifying tokens.
type Verifier interface {
	Verify(token string) (User, error)
}

// NewVerifier creates a new Verifier using the default implementation.
func NewVerifier(creds Credentials, leeway time.Duration) Verifier {
	return &jwtVerifier{
		secret:         []byte(creds.Secret),
		expectedIssuer: creds.Issuer,
		leeway:         leeway,
	}
}

type customClaims struct {
	ClientID   string `json:"cid,omitempty"`
	SessionID  string `json:"sid,omitempty"`
	OriginID   string `json:"org,omitempty"`
	CareUnitID string `json:"cu,omitempty"`
	Roles      string `json:"role,omitempty"`
}

type jwtIssuer struct {
	name   string
	signer jose.Signer
}

func (i *jwtIssuer) Issue(user User, lifetime time.Duration) (string, error) {
	err := i.verifyTokenContent(user)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := josejwt.Claims{
		Subject:   user.ID,
		Issuer:    i.name,
		NotBefore: josejwt.NewNumericDate(now.Add(-1 * time.Minute)),
		IssuedAt:  josejwt.NewNumericDate(now),
		Expiry:    josejwt.NewNumericDate(now.Add(lifetime)),
	}

	custCl := customClaims{
		Roles: strings.Join(user.Roles, roleDelimiter),
	}
	return josejwt.Signed(i.signer).Claims(claims).Claims(custCl).CompactSerialize()
}

func (i *jwtIssuer) verifyTokenContent(user User) error {
	if user.ID == "" {
		return ErrInvalidTokenContent
	}
	if user.Roles == nil || len(user.Roles) == 0 {
		return ErrInvalidTokenContent
	}

	return nil
}

type jwtVerifier struct {
	secret         []byte
	expectedIssuer string
	leeway         time.Duration
}

func (v *jwtVerifier) Verify(rawToken string) (User, error) {
	token, err := josejwt.ParseSigned(rawToken)
	if err != nil {
		return User{}, ErrInvalidToken
	}

	var claims josejwt.Claims
	err = token.Claims(v.secret, &claims)
	if err != nil {
		return User{}, ErrInvalidToken
	}

	var customCl customClaims
	err = token.Claims(v.secret, &customCl)
	if err != nil {
		return User{}, ErrInvalidToken
	}

	err = v.validateClaims(claims)
	if err != nil {
		return User{}, err
	}

	return getTokenFromClaims(claims, customCl), nil
}

func (v *jwtVerifier) validateClaims(claims josejwt.Claims) error {
	err := claims.ValidateWithLeeway(josejwt.Expected{Issuer: v.expectedIssuer}, v.leeway)
	if err != nil {
		return ErrInvalidToken
	}

	err = v.checkTokenExpiry(claims)
	if err != nil {
		return err
	}

	if claims.Subject == "" {
		return ErrInvalidToken
	}

	return nil
}

func (v *jwtVerifier) checkTokenExpiry(claims josejwt.Claims) error {
	earliestDate := claims.NotBefore.Time()
	latestDate := claims.Expiry.Time().Add(v.leeway)
	now := time.Now()

	if now.Before(earliestDate) {
		return ErrInvalidToken
	}

	if now.After(latestDate) {
		return ErrExpiredToken
	}

	return nil
}

func getTokenFromClaims(claims josejwt.Claims, customCl customClaims) User {
	return User{
		ID:    claims.Subject,
		Roles: strings.Split(customCl.Roles, roleDelimiter),
	}
}
