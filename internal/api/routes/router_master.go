package routes

import (
	"example/internal/api/handlers"
	"example/internal/api/middleware"
	"example/internal/pkg/models/ms_product"
	"example/internal/pkg/models/ms_product_category"
	"example/internal/pkg/models/ms_supplier"
	"github.com/gofiber/fiber/v2"
)

func SetupMasterRoutes(
	api fiber.Router,
	msSupplierService ms_supplier.Service,
	msProductCategoryService ms_product_category.Service,
	msProductService ms_product.Service,
) {
	supplierRoute := api.Group("/v1/master/supplier", middleware.Protected())
	supplierRoute.Get("/", handlers.HandleGetSuppliers(msSupplierService))
	supplierRoute.Get("/:id", handlers.HandleGetSupplier(msSupplierService))
	supplierRoute.Post("/", handlers.HandleAddSupplier(msSupplierService))
	supplierRoute.Put("/:id", handlers.HandleUpdateSupplier(msSupplierService))
	supplierRoute.Delete("/", handlers.HandleRemoveSupplier(msSupplierService))

	productCategoryRoute := api.Group("/v1/master/product_category", middleware.Protected())
	productCategoryRoute.Get("/", handlers.HandleGetProductCategories(msProductCategoryService))
	productCategoryRoute.Get("/:id", handlers.HandleGetProductCategory(msProductCategoryService))
	productCategoryRoute.Post("/", handlers.HandleAddProductCategory(msProductCategoryService))
	productCategoryRoute.Put("/:id", handlers.HandleUpdateProductCategory(msProductCategoryService))
	productCategoryRoute.Delete("/:id", handlers.HandleRemoveProductCategory(msProductCategoryService))

	productRoute := api.Group("/v1/master/product", middleware.Protected())
	productRoute.Get("/", handlers.HandleGetProducts(msProductService))
	productRoute.Get("/:id/:based_on", handlers.HandleGetProduct(msProductService))
	productRoute.Post("/", handlers.HandleAddProduct(msProductService))
	productRoute.Put("/:id", handlers.HandleUpdateProduct(msProductService))
	productRoute.Delete("/", handlers.HandleRemoveProduct(msProductService))
}
