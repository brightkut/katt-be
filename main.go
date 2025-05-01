package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"katt-be/category"
	"katt-be/handler"
	"katt-be/middleware"
	"katt-be/transaction"
	"katt-be/wallet"
	"log"
	"os"
	"time"

	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	_ "github.com/brightkut/rest-api-go-fiber/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Running Local
// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "secret"
// 	dbname   = "katt"
// )

var (
	db          *gorm.DB
	fiberLambda *fiberadapter.FiberLambda
)

// init the Fiber Server
func init() {
	log.Printf("Fiber cold start")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,
		},
	)
	var err error

	// Running local
	// dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), 5432, os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect DB")
	}

	// auto create and update table but not for delete case
	db.AutoMigrate(&wallet.Wallet{}, &category.Category{}, &transaction.Transaction{})

	fmt.Printf("Connect DB successfully")

	var app *fiber.App
	app = fiber.New()

	// load env - make this optional for Lambda environment
	//if err := godotenv.Load(); err != nil {
	//	log.Println("No .env file found, using environment variables")
	//}
	// allow cors
	app.Use(cors.New())

	handler := handler.NewHandler(db)

	app.Get("/hello", handler.Hello)

	// check token middleware
	app.Use(middleware.LoginMiddleware)

	// TODO for lambda
	fiberLambda = fiberadapter.New(app)

	// Wallet Handler
	app.Post("/wallets", handler.CreateWallet)
	app.Post("/wallets-by-email", handler.GetWalletByEmail)
	app.Post("/categories", handler.CreateCategory)
	app.Post("/categories-by-wallet-id", handler.GetAllCategoryByWalletId)
	app.Delete("/categories/:id", handler.DeleteCategory)
	app.Post("/transactions", handler.CreateTransaction)
	app.Get("/transactions", handler.GetAllTransactionByWalletId)
	app.Delete("/transactions/:id", handler.DeleteTransaction)

	// listen port
	// app.Listen(":8080")
}

// // Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return fiberLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)
}
