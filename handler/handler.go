package handler

import (
	"fmt"
	"katt-be/category"
	"katt-be/transaction"
	"katt-be/wallet"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (h *Handler) CreateWallet(c *fiber.Ctx) error {
	dto := new(wallet.CreateWalletDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := wallet.CreateWallet(h.DB, dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("create wallet success %v", dto.Email),
	})
}

func (h *Handler) GetWalletByEmail(c *fiber.Ctx) error {
	dto := new(wallet.GetWalletByEmailDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"data": wallet.GetWalletByEmail(h.DB, dto),
	})
}

func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	dto := new(category.CreateCategoryDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := category.CreateCategory(h.DB, dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("create category success %v", dto.CategoryName),
	})
}

func (h *Handler) GetAllCategoryByWalletId(c *fiber.Ctx) error {
	dto := new(category.GetAllCategoryByWalletIdDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"data": category.GetAllCategoryByWalletId(h.DB, dto),
	})
}

func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	categoryId := c.Params("id")

	err := category.DeleteCategory(h.DB, categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("delete category success %v", categoryId),
	})
}

func (h *Handler) CreateTransaction(c *fiber.Ctx) error {
	dto := new(transaction.CreateTransactionDto)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := transaction.CreateTransaction(h.DB, dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("create transaction success %v", dto.TransactionName),
	})
}

func (h *Handler) GetAllTransactionByWalletId(c *fiber.Ctx) error {
	walletId := c.Query("walletId")
	pageNumber, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	dto := &transaction.GetAllTransactionByWalletIdDto{
		WalletId: walletId,
		Page:     int32(pageNumber),
		PageSize: int32(pageSize),
	}

	transactions, err := transaction.GetAllTransactionByWalletId(h.DB, dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"data": transactions,
	})
}

func (h *Handler) DeleteTransaction(c *fiber.Ctx) error {
	transactionId := c.Params("id")

	err := transaction.DeleteTransaction(h.DB, transactionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("delete transaction success %v", transactionId),
	})
}
