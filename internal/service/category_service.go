package service

import (
	"katt-be/internal/domain"
	"katt-be/internal/repository"

	"github.com/google/uuid"
)

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

type CategoryService interface {
	Create(dto *domain.CreateCategoryDto) error
	Delete(categoryId string) error
	FindAllByWalletId(dto domain.GetAllCategoryByWalletIdDto) ([]domain.Category, error)
}

func NewCategoryService(categoryRepository repository.CategoryRepository) categoryService {
	return categoryService{categoryRepository: categoryRepository}
}

func (c *categoryService) Create(dto *domain.CreateCategoryDto) error {
	categoryId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	category := &domain.Category{
		CategoryId:   categoryId.String(),
		WalletId:     dto.WalletId,
		CategoryName: dto.CategoryName,
	}

	return c.categoryRepository.Save(category)
}

func (c *categoryService) Delete(categoryId string) error {

	category, error := c.categoryRepository.FindById(categoryId)

	if error != nil {
		return error
	}

	error = c.categoryRepository.Delete(category)

	if error != nil {
		return error
	}

	return nil
}

func (c *categoryService) FindAllByWalletId(dto domain.GetAllCategoryByWalletIdDto) ([]domain.Category, error) {
	categories, error := c.categoryRepository.FindAllByWalletId(dto.WalletId)

	if error != nil {
		return nil, error
	}

	return categories, nil
}
