package services

import (
	"context"
	"lesson-proj/internal/database"
	"lesson-proj/internal/models"
	productUtils "lesson-proj/internal/services/products/utils"
)

type ProductService struct {
	repository *database.ProductRepository
}

func NewProductService(repository *database.ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (productService *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	product, err := productService.repository.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (productService *ProductService) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	product, err := productService.repository.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (productService *ProductService) CreateProduct(ctx context.Context, inputProduct models.CreateProduct) (*models.Product, error) {
	if err := productUtils.ValidateCreateProductInput(inputProduct.Title, inputProduct.Description, inputProduct.Price); err != nil {
		return nil, err
	}
	createdProduct, err := productService.repository.CreateProduct(ctx, inputProduct)
	if err != nil {
		return nil, err
	}
	return createdProduct, nil
}

func (productService *ProductService) UpdateProduct(ctx context.Context, id int, inputProduct models.UpdateProduct) (*models.Product, error) {
	if err := productUtils.ValidateUpdateProductInput(
		inputProduct.Title, 
		inputProduct.Description, 
		inputProduct.Price,
	); err != nil {
		return nil, err
	}
	updatedProduct, err := productService.repository.UpdateProduct(ctx, id, inputProduct)
	if err != nil {
		return nil, err
	}
	return updatedProduct, nil
}

func (productService *ProductService) DeleteProduct(ctx context.Context, id int) error {
	err := productService.repository.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}
	return nil
}