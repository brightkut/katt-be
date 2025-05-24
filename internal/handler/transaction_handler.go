package handler

import (
	"fmt"
	"katt-be/internal/domain"
	"katt-be/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type transactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) transactionHandler {
	return transactionHandler{transactionService: transactionService}
}

func (t *transactionHandler) Create(c *fiber.Ctx) error {
	dto := new(domain.CreateTransactionDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := t.transactionService.Create(dto)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("create transaction success %v", dto.TransactionName),
	})
}

func (t *transactionHandler) FindAllByWalletId(c *fiber.Ctx) error {
	walletId := c.Query("walletId")
	pageNumber, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	dto := &domain.GetAllTransactionByWalletIdDto{
		WalletId: walletId,
		Page:     int32(pageNumber),
		PageSize: int32(pageSize),
	}

	transactions, err := t.transactionService.FindAllByWalletId(dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"data": transactions,
	})
}

func (t *transactionHandler) Delete(c *fiber.Ctx) error {
	transactionId := c.Params("id")

	err := t.transactionService.Delete(transactionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("delete transaction success %v", transactionId),
	})
}
