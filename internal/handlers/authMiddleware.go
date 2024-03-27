package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

const SecretKey = "secret key"

type userIDKey struct {}

type Claims struct {
    jwt.RegisteredClaims
    UserID uint
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		var userID uint = 0;
		cookie, err := r.Cookie("JWT")
		if err == nil {
			token := cookie.Value
			userID, err = getUserID(token)
		}

		if err != nil {
			userID = urlStorage.CreateNewUser()
			token := buildJWTString(userID)
			http.SetCookie(w, &http.Cookie{Name: "JWT", Value: token})
		}

		r = r.WithContext(context.WithValue(r.Context(), userIDKey{}, userID))
		next.ServeHTTP(w, r)
	})
}

func buildJWTString(userID uint) string {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims {
        RegisteredClaims: jwt.RegisteredClaims{},
        UserID: userID,
    })

    tokenString, _ := token.SignedString([]byte(SecretKey))
    return tokenString
}

func getUserID(tokenString string) (uint, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(SecretKey), nil
    })

	if err != nil {
        return 0, err
    }

    if !token.Valid {
        return 0, errors.New("token is not valid")
    }

    return claims.UserID, nil
}