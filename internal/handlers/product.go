package handlers

import (
	"encoding/json"
	"lesson-proj/internal/models"
	services "lesson-proj/internal/services/products"
	"net/http"
)

type ProductHandler struct {
	service *services.ProductService
}

// NewProductHandler — factory function (constructor).
// It creates a new ProductHandler object.
func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (handler *ProductHandler) GetAllProducts(response http.ResponseWriter, request *http.Request) {
	products, err := handler.service.GetAllProducts(request.Context())
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}
	respondWithJSON(response, http.StatusOK, products)
}

func (handler *ProductHandler) GetProductByID(response http.ResponseWriter, request *http.Request) {
	// Extract the product ID from the URL path with helper function
	id, err := getIDFromPath(request)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid product ID")
		return
	}
	product, err := handler.service.GetProductByID(request.Context(), id)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Failed to retrieve product")
		return
	}
	if product == nil {
		respondWithError(response, http.StatusNotFound, "Product not found")
		return
	}
	respondWithJSON(response, http.StatusOK, product)
}

func (handler *ProductHandler) CreateProduct(response http.ResponseWriter, request *http.Request) {
	var productInput models.CreateProduct
	// Decode the JSON request body into the productInput struct
	// json.NewDecoder reads from request.Body and decodes JSON into the struct
	err := json.NewDecoder(request.Body).Decode(&productInput)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	createdProduct, err := handler.service.CreateProduct(request.Context(), productInput)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Failed to create product")
		return
	}
	respondWithJSON(response, http.StatusCreated, createdProduct)

}

func (handler *ProductHandler) UpdateProduct(response http.ResponseWriter, request *http.Request) {
	id, err := getIDFromPath(request)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid product ID")
		return
	}
	var productInput models.UpdateProduct
	err = json.NewDecoder(request.Body).Decode(&productInput)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	updatedProduct, err := handler.service.UpdateProduct(request.Context(), id, productInput)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Failed to update product")
		return
	}
	respondWithJSON(response, http.StatusOK, updatedProduct)
}

func (handler *ProductHandler) DeleteProduct(response http.ResponseWriter, request *http.Request) {
	id, err := getIDFromPath(request)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid product ID")
		return
	}
	err = handler.service.DeleteProduct(request.Context(), id)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	respondWithJSON(response, http.StatusNoContent, nil)
}
