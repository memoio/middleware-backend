package server

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Type         int  `json:"type,omitempty"`
	IsRegistered bool `json:"isRegistered,omitempty"`
	// Nonce string `json:"nonce,omitempty"`
	jwt.StandardClaims
}

func VerifyAccessToken(tokenString string) (string, error) {
	parts := strings.SplitN(tokenString, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrNullToken
	}

	claims := &Claims{}
	_, _, err := new(jwt.Parser).ParseUnverified(parts[1], claims)
	if err != nil {
		return "", ErrValidToken
	}

	// check Audience
	if claims.Audience != Domain || claims.Issuer != Domain {
		return "", ErrValidToken
	}

	// check token type
	if claims.Type != AccessToken {
		return "", ErrValidTokenType
	}

	// check signature, Expires time and Issued time
	token, err := parseToken(parts[1])
	if err != nil || !token.Valid {
		return "", ErrValidToken
	}

	return claims.Subject, nil
}

func VerifyRefreshToken(tokenString string) (string, error) {
	parts := strings.SplitN(tokenString, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrNullToken
	}

	claims := &Claims{}
	_, _, err := new(jwt.Parser).ParseUnverified(parts[1], claims)
	if err != nil {
		return "", ErrValidToken
	}

	// check Audience
	if claims.Audience != Domain || claims.Issuer != Domain {
		return "", ErrValidToken
	}

	// check token type
	if claims.Type != RefreshToken {
		return "", ErrValidTokenType
	}

	token, err := parseToken(parts[1])
	if err != nil || !token.Valid {
		return "", ErrValidToken
	}

	return genAccessTokenWithFlag(claims.Subject, claims.IsRegistered)
}

func genAccessToken(subject string) (string, error) {
	expireTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Type: AccessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Audience:  Domain,
			Issuer:    Domain,
			Subject:   subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTKey)
}

func genAccessTokenWithFlag(subject string, isRegistered bool) (string, error) {
	expireTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Type:         AccessToken,
		IsRegistered: isRegistered,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Audience:  Domain,
			Issuer:    Domain,
			Subject:   subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTKey)
}

func genRefreshToken(subject string) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Type: RefreshToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Audience:  Domain,
			Issuer:    Domain,
			Subject:   subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTKey)
}

func genRefreshTokenWithFlag(subject string, isRegistered bool) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Type:         RefreshToken,
		IsRegistered: isRegistered,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Audience:  Domain,
			Issuer:    Domain,
			Subject:   subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTKey)
}

// func ParseDidToken(tokenString string, did string) (*jwt.Token, error) {
//     return jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
//     	parts := strings.Split(did, ":")
//     	if len(parts) != 3 || parts[0] != "did" || parts[1] != "eth" {
//     		return nil, ErrValidToken
//     	}

//     	pubKeyBytes, err := hex.DecodeString(parts[2])
//     	if err != nil {
//     		return nil, err
//     	}

//         return crypto.UnmarshalPubkey(pubKeyBytes)
//     })
// }

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return JWTKey, nil
	})
}
