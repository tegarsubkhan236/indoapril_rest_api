package main

import (
	"example/internal/api/routes"
	"example/internal/pkg"
	"example/internal/pkg/models/cr_permission"
	"example/internal/pkg/models/cr_role"
	"example/internal/pkg/models/cr_team"
	"example/internal/pkg/models/cr_user"
	"example/internal/pkg/models/ms_product"
	"example/internal/pkg/models/ms_product_category"
	"example/internal/pkg/models/ms_product_price"
	"example/internal/pkg/models/ms_stock"
	"example/internal/pkg/models/ms_supplier"
	"example/internal/pkg/models/tr_back_order"
	"example/internal/pkg/models/tr_purchase_order"
	"example/internal/pkg/models/tr_receiving_order"
	"example/internal/pkg/models/tr_sales_order"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	// Setup Database
	db, err := connectDB()
	if err != nil {
		log.Fatal("Database Connection Error $s", err)
	}

	// Setup Fiber
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:4000",
		AllowHeaders:     "Authorization, Content-Type",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Setup Service
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Welcome to the jungle!"))
	})
	api := app.Group("/api", logger.New())
	setupAuthService(db, api)
	setupCoreService(db, api)
	setupMasterService(db, api)
	setupTransactionService(db, api)

	// Start App
	log.Fatal(app.Listen(pkg.GetEnv("APP_PORT")))
}

func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		pkg.GetEnv("DB_USER"),
		pkg.GetEnv("DB_PASSWORD"),
		pkg.GetEnv("DB_HOST"),
		pkg.GetEnv("DB_PORT"),
		pkg.GetEnv("DB_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = pkg.RunMigration(db); err != nil {
		return nil, err
	}
	if err = pkg.RunSeeder(db); err != nil {
		return nil, err
	}

	return db, nil
}

func setupAuthService(db *gorm.DB, api fiber.Router) {
	crAuthRepo := cr_user.NewRepo(db)
	crAuthService := cr_user.NewService(crAuthRepo)

	routes.SetupAuthRoutes(api, crAuthService)
}

func setupCoreService(db *gorm.DB, api fiber.Router) {
	crPermissionRepo := cr_permission.NewRepo(db)
	crPermissionService := cr_permission.NewService(crPermissionRepo)

	crRoleRepo := cr_role.NewRepo(db)
	crRoleService := cr_role.NewService(crRoleRepo)

	crUserRepo := cr_user.NewRepo(db)
	crUserService := cr_user.NewService(crUserRepo)

	crTeamRepo := cr_team.NewRepo(db)
	crTeamService := cr_team.NewService(crTeamRepo)

	routes.SetupCoreRoutes(api, crPermissionService, crRoleService, crUserService, crTeamService)
}

func setupMasterService(db *gorm.DB, api fiber.Router) {
	msSupplierRepo := ms_supplier.NewRepo(db)
	msSupplierService := ms_supplier.NewService(msSupplierRepo)

	msProductCategoryRepo := ms_product_category.NewRepo(db)
	msProductCategoryService := ms_product_category.NewService(msProductCategoryRepo)

	msProductRepo := ms_product.NewRepo(db)
	msProductService := ms_product.NewService(msProductRepo)

	routes.SetupMasterRoutes(api, msSupplierService, msProductCategoryService, msProductService)
}

func setupTransactionService(db *gorm.DB, api fiber.Router) {
	msStockRepo := ms_stock.NewRepo(db)
	msProductPriceRepo := ms_product_price.NewRepo(db)
	trBORepo := tr_back_order.NewRepo(db)

	trPORepo := tr_purchase_order.NewRepo(db)
	trPOService := tr_purchase_order.NewService(trPORepo, trBORepo)

	trRORepo := tr_receiving_order.NewRepo(db)
	trROService := tr_receiving_order.NewService(msStockRepo, msProductPriceRepo, trPORepo, trBORepo, trRORepo)

	trSORepo := tr_sales_order.NewRepo(db)
	trSOService := tr_sales_order.NewService(trSORepo, msStockRepo, msProductPriceRepo)

	routes.SetupTransactionRoutes(api, trPOService, trROService, trSOService)
}
