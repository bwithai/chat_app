package user

import (
	"chatapp/auth_jwt"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u CreateUserReq
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.Service.CreateUser(r.Context(), &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.Service.Login(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.accessToken = auth_jwt.GetJwt(u.ID)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    u.accessToken,
		MaxAge:   60 * 60 * 24,
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: false,
		Secure:   true,
	})
	json.NewEncoder(w).Encode(u.accessToken)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Token")
	fmt.Println(tokenStr)
	jwt, err := auth_jwt.RemoveAuthorizationFromJWT(tokenStr)
	if err != nil {
		fmt.Println("new token err")
	}
	fmt.Println(jwt)
	r.Header.Set("Token", jwt)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: false,
		Secure:   true,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Logout successful")
}
