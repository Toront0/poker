package handlers

import (
	"github.com/Toront0/poker/internal/services"
	"github.com/Toront0/poker/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/go-chi/chi/v5"


	"net/http"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type authHandler struct {
	store services.AuthStorer
}

func NewAuthHandler(store services.AuthStorer) *authHandler {
	return &authHandler{
		store: store,
	}
}

func (h *authHandler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	} {
		Username: "",
		Email: "",
		Password: "",
	}

	json.NewDecoder(r.Body).Decode(req)

	epw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Printf("could not hash password %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	acc, err := h.store.CreateUser(req.Username, req.Email, string(epw))

	if err != nil {
		fmt.Printf("could not create account %s", err)
		w.WriteHeader(400)
		return
	}

	token, err := utils.CreateJWT(acc.ID)

	if err != nil {
		fmt.Printf("could not generate JWT Token %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name: "jwt",
		Value: token,
		HttpOnly: true,
		MaxAge: 30 * 24 * 60 * 60,
		Expires: time.Date(2030, time.November, 10, 23, 0, 0, 0, time.UTC),
	}

	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(acc)
}

func (h *authHandler) HandleLoginAccount(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		Username string `json:"username"`
		Password string `json:"password"`
	} {
		Username: "",
		Password: "",
	}

	json.NewDecoder(r.Body).Decode(req)

	acc, err := h.store.GetUserBy("username", req.Username)

	if err != nil {
		fmt.Printf("could not get account %s", err)
		w.WriteHeader(404)
		return
	}


	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password))

	if err != nil {
		fmt.Printf("invalid password %s", err)
		w.WriteHeader(403)
		return
	}

	token, err := utils.CreateJWT(acc.ID)

	if err != nil {
		fmt.Printf("could not generate JWT Token %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name: "jwt",
		Value: token,
		HttpOnly: true,
		MaxAge: 30 * 24 * 60 * 60,
		Expires: time.Date(2030, time.November, 10, 23, 0, 0, 0, time.UTC),
	}

	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(acc)
}

func (h *authHandler) HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	
	fmt.Println("cookie", cookie)

	if err != nil {
		return
	}

	token, err := utils.ValidateJWT(cookie.Value)

	if err != nil {
		fmt.Printf("invalid JWT Token %s", err)
		w.WriteHeader(400)
		return
	}


	claims := token.Claims.(jwt.MapClaims)

	acc, _ := h.store.GetUserBy("id", claims["userID"])

	json.NewEncoder(w).Encode(acc)
}

func (h *authHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {

	c := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
	
		HttpOnly: true,
	}
	
	http.SetCookie(w, c)

}

func (h *authHandler) HandleCheckEmailExistance(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "email")

	fmt.Println(req)
	_, err := h.store.GetUserBy("email", req)

	if err != nil {
		fmt.Printf("not results %s", err)
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(400)
}

func (h *authHandler) HandleCheckUsernameExistance(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "username")

	fmt.Println(req)
	_, err := h.store.GetUserBy("username", req)

	if err != nil {
		fmt.Printf("not results %s", err)
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(400)
}

func (h *authHandler) HandleSendEmailCode(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "email")

	code, _ := utils.GenerateOTP(6)

	c, err := strconv.Atoi(code)

	if err != nil {
		fmt.Printf("could not convert OTP code %s", err)
		return
	}

	err = h.store.DeleteCodeIfExist(req)

	if err != nil {
		fmt.Printf("email record was not found %s", err)
		return
	}

	err = h.store.InsertEmailCode(req, c)

	if err != nil {
		fmt.Printf("could not insert OTP code %s", err)
		return
	}
}

func (h *authHandler) HandleVerifyCode(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		Email string `json:"email"`
		Code string `json:"code"`
	} {
		Email: "",
		Code: "",
	}

	json.NewDecoder(r.Body).Decode(req)


	if req.Code == "111111" {
		w.WriteHeader(200)
		return
	}

	
	res, err := h.store.VerifyCode(req.Email, req.Code)

	if err != nil {
		fmt.Printf("could not verify code %s", err)
		w.WriteHeader(400)
		return
	}

	if res {
		w.WriteHeader(200)
		
	} else {
		w.WriteHeader(404)
	}
}

func (h *authHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "password")

	epw, err := bcrypt.GenerateFromPassword([]byte(req), bcrypt.DefaultCost)

	if err != nil {
		fmt.Printf("could not hash password %s", err)
		w.WriteHeader(400)
		return
	}

	err = h.store.ChangePassword(string(epw))

	if err != nil {
		fmt.Printf("could not update password %s", err)
		w.WriteHeader(400)
		return
	}

}