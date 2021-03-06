package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateusrlopez/go-market/models"
	"github.com/mateusrlopez/go-market/repositories"
	"github.com/mateusrlopez/go-market/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	ProductRepository repositories.ProductRepository
}

func (h *ProductHandler) Index(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductRepository.RetrieveAll()

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, products)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	product := models.Product{}

	if err = json.Unmarshal(body, &product); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = product.Validate(); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	result, err := h.ProductRepository.Create(&product)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{"id": result.InsertedID.(primitive.ObjectID).Hex()})
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	product := models.Product{}

	if err = h.ProductRepository.RetrieveByID(id, &product); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	product := models.Product{}

	if err = json.Unmarshal(body, &product); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = product.Validate(); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = h.ProductRepository.Update(id, &product)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusNoContent, map[string]interface{}{})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = h.ProductRepository.Delete(id)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.JSONResponse(w, http.StatusNoContent, map[string]interface{}{})
}
