package handler

import (
	"fmt"
	"katt-be/internal/domain"
	"katt-be/internal/service"

	"github.com/gofiber/fiber/v2"
)

type categoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) categoryHandler {
	return categoryHandler{categoryService: categoryService}
}

func (ch *categoryHandler) Create(c *fiber.Ctx) error {
	dto := new(domain.CreateCategoryDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := ch.categoryService.Create(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("create category success %v", dto.CategoryName),
	})
}

func (ch *categoryHandler) FindAllByWalletId(c *fiber.Ctx) error {
	dto := new(domain.GetAllCategoryByWalletIdDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	categories, error := ch.categoryService.FindAllByWalletId(*dto)

	if error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(error.Error())
	}

	return c.JSON(fiber.Map{
		"data": categories,
	})
}

func (ch *categoryHandler) Delete(c *fiber.Ctx) error {
	categoryId := c.Params("id")

	err := ch.categoryService.Delete(categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("delete category success %v", categoryId),
	})
}
