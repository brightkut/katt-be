package handler

import (
	"fmt"
	"katt-be/internal/domain"
	"katt-be/internal/service"

	"github.com/gofiber/fiber/v2"
)

type walletHandler struct {
	walletService service.WalletService
}

func NewWalletHandler(walletService service.WalletService) walletHandler {
	return walletHandler{walletService: walletService}
}

func (w *walletHandler) Create(c *fiber.Ctx) error {
	dto := new(domain.CreateWalletDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := w.walletService.Create(dto)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Create wallet success with email: %v", dto.Email),
	})
}

func (w *walletHandler) GetByEmail(c *fiber.Ctx) error {
	dto := new(domain.GetWalletByEmailDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"data": w.walletService.GetByEmail(dto),
	})
}
