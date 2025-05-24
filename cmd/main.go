package main

import (
	"context"
	"katt-be/internal/config"
	"katt-be/internal/handler"
	// "katt-be/internal/middleware"
	"katt-be/internal/repository"
	"katt-be/internal/service"
	"log"

	"github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"

	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	_ "github.com/brightkut/rest-api-go-fiber/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var (
	fiberLambda *fiberadapter.FiberLambda
)

// init the Fiber Server
func Init() {
	log.Printf("Fiber cold start")

	// TODO: update this when run on lambda
	config.LoadEnv("dev")

	// TODO: update this when run on lambda
	db := config.NewPostgres("prod")

	// auto create and update table but not for delete case
	//db.AutoMigrate(&wallet.Wallet{}, &category.Category{}, &transaction.Transaction{})

	var app *fiber.App
	app = fiber.New()

	// allow cors
	app.Use(cors.New())

	// Setup DI
	walletRepo := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(&walletRepo)
	walletHandler := handler.NewWalletHandler(&walletService)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(&categoryRepo)
	categoryHandler := handler.NewCategoryHandler(&categoryService)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(&walletRepo, &transactionRepo)
	transactionHandler := handler.NewTransactionHandler(&transactionService)

	app.Get("/hello", handler.Hello)

	// app.Use(middleware.JwtMiddleware)

	// TODO: for lambda
	// fiberLambda = fiberadapter.New(app)

	app.Post("/wallets", walletHandler.Create)
	app.Post("/wallets-by-email", walletHandler.GetByEmail)
	app.Post("/categories", categoryHandler.Create)
	app.Post("/categories-by-wallet-id", categoryHandler.FindAllByWalletId)
	app.Delete("/categories/:id", categoryHandler.Delete)
	app.Post("/transactions", transactionHandler.Create)
	app.Get("/transactions", transactionHandler.FindAllByWalletId)
	app.Delete("/transactions/:id", transactionHandler.Delete)

	//TODO: Running local
	app.Listen(":8080")
}

// // Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return fiberLambda.ProxyWithContext(ctx, req)
}

func main() {
	//NOTE: Running local
	Init()

	// TODO: uncomment for running on lambda
	// lambda.Start(Handler)
}
