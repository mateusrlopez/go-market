package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mateusrlopez/go-market/models"
	"github.com/mateusrlopez/go-market/repositories"
	"github.com/mateusrlopez/go-market/shared/constants"
	"github.com/mateusrlopez/go-market/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	UserRepository  repositories.UserRepository
	TokenRepository repositories.TokenRepository
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}

	if err = json.Unmarshal(body, &user); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = user.ValidateRegister(); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	result, err := h.UserRepository.Create(&user)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	tr, err := h.TokenRepository.GenerateTokens(result.InsertedID.(primitive.ObjectID).Hex())

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, tr)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}

	if err = json.Unmarshal(body, &user); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = user.ValidateLogin(); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	retrievedUser := models.User{}

	if err = h.UserRepository.RetrieveByEmail(user.Email, &retrievedUser); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = retrievedUser.ComparePassword(user.Password); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	tr, err := h.TokenRepository.GenerateTokens(retrievedUser.ID.Hex())

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, tr)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(constants.ContextKey).(*models.User)

	utils.JSONResponse(w, http.StatusOK, user)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(constants.ContextKey).(*models.User)

	tr, err := h.TokenRepository.GenerateTokens(user.ID.Hex())

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, tr)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(constants.ContextKey).(*models.User)

	if err := h.TokenRepository.DeleteTokenMetadata(user.ID.Hex()); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONResponse(w, http.StatusNoContent, map[string]interface{}{})
}
