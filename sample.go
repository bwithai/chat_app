package main

import (
	"errors"
	"fmt"

	"chatapp/jwt-go"

	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var SECRET = []byte("super-secret-auth-key")

func CreateJWT(userID string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = userID
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	tokenStr, err := token.SignedString(SECRET)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tokenStr, nil
}

func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Token")
		if tokenStr == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized"))
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return SECRET, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized: " + err.Error()))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
			return
		}

		reqUserID := mux.Vars(r)["userID"]
		if reqUserID != userID {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RemoveAuthorizationFromJWT(tokenString string) (string, error) {
	// Parse the token string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify that the signing method is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		// Return the secret key used to sign the token
		return []byte(SECRET), nil
	})

	// Check if there was an error parsing the token
	if err != nil {
		return "", err
	}

	// Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Remove the authorization claim from the token
	delete(claims, "auth")

	// Create a new token with the modified claims
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the new token with the secret key
	newTokenStr, err := newToken.SignedString([]byte(SECRET))
	if err != nil {
		return "", err
	}

	// Return the new token string without the authorization claim
	return newTokenStr, nil
}

func GetJwt(userID string) string {
	token, err := CreateJWT(userID)
	if err != nil {
		return ""
	}
	return token
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "super secret area")
}

func Jwt(w http.ResponseWriter, r *http.Request) {
	id := "1"
	token := GetJwt(id)

	fmt.Fprintf(w, "Token: %v", token)
}

func main() {
	r := mux.NewRouter()
	//router.Use(ValidateJWT)
	//
	//// Add your routes here
	//router.HandleFunc("/users/{userID}/profile", GetProfileHandler).Methods("GET")

	r.Handle("/api/{userID}/profile", ValidateJWT(http.HandlerFunc(Home))).Methods("GET")
	r.HandleFunc("/jwt", Jwt).Methods("GET")

	http.ListenAndServe(":3500", r)
}
