package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateusrlopez/go-market/models"
	"github.com/mateusrlopez/go-market/repositories"
	"github.com/mateusrlopez/go-market/shared/types"
	"github.com/mateusrlopez/go-market/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewHandler struct {
	ReviewRepository repositories.ReviewRepository
}

func (h *ReviewHandler) Index(w http.ResponseWriter, r *http.Request) {
	query := types.ReviewIndexQuery{}

	if err := utils.DecodeQuery(&query, r.URL.Query()); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	reviews, err := h.ReviewRepository.RetrieveAll(&query)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	review := models.Review{}

	if err = json.Unmarshal(body, &review); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = review.Validate(); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	result, err := h.ReviewRepository.Create(&review)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{"id": result.InsertedID.(primitive.ObjectID).Hex()})
}

func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	review := models.Review{}

	if err = json.Unmarshal(body, &review); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = review.Validate(); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = h.ReviewRepository.Update(id, &review)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusNoContent, map[string]interface{}{})
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = h.ReviewRepository.Delete(id)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusNoContent, map[string]interface{}{})
}
