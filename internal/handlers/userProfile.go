package handlers

import (
	"github.com/Toront0/poker/internal/services"

	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"encoding/json"
	"fmt"
)

type userProfileHandler struct {
	store services.UserProfileStorer
}

func NewUserProfileHandler(store services.UserProfileStorer) *userProfileHandler {
	return &userProfileHandler{
		store: store,
	}
}

func (h *userProfileHandler) HandleGetAllImages(w http.ResponseWriter, r *http.Request) {

	imgs, err := h.store.GetAllProfileImages()

	if err != nil {
		fmt.Printf("could not get all profile images %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(imgs)
}

func (h *userProfileHandler) HandleChangeProfileImage(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		UserID int `json:"userId"`
		URL string `json:"url"`
	} {
		UserID: 0,
		URL: "",
	}

	json.NewDecoder(r.Body).Decode(req)


	err := h.store.ChangeProfileImage(req.UserID, req.URL)

	if err != nil {
		fmt.Printf("could not change profile image %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *userProfileHandler) HandleGetUserDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse user ID %s", err)
		w.WriteHeader(400)
		return
	}

	u, err := h.store.GetUserByID(uID)

	if err != nil {
		fmt.Printf("could not get user %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func (h *userProfileHandler) HandleChangeUsername(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		UserID int `json:"userId"`
		NewValue string `json:"newValue"`
	}{
		UserID: 0,
		NewValue: "",
	}

	json.NewDecoder(r.Body).Decode(req)


	err := h.store.ChangeUsername(req.UserID, req.NewValue)

	if err != nil {
		fmt.Printf("could not change username %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *userProfileHandler) HandleFindUsers(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("search")

	res, err := h.store.FindUsers(s)

	if err != nil {
		fmt.Printf("could not find %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *userProfileHandler) HandleGetUserGames(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p := r.URL.Query().Get("page")
	l := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(p)
	limit, err := strconv.Atoi(l)

	fmt.Printf("page is %d, limit is %d, id is %d", page, limit, id)

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse user ID %s", err)
		w.WriteHeader(400)
		return
	}


	res, err := h.store.GetUserGames(uID, limit, page)

	if err != nil {
		fmt.Printf("could not get user games preview %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *userProfileHandler) HandleCheckPossibilityToGetMoney(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse user ID %s", err)
		w.WriteHeader(400)
		return
	}


	res, err := h.store.GetLastMoneyTransactionStatus(uID)

	if err != nil {
		fmt.Printf("could not get money status %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *userProfileHandler) HandleGetFreeMoney(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse user ID %s", err)
		w.WriteHeader(400)
		return
	}

	err = h.store.GetFreeMoney(uID)

	if err != nil {
		fmt.Printf("could not get free money %s", err)
		w.WriteHeader(400)
		return
	}
}