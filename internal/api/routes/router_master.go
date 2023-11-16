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
	supplierRoute.Get("/", handlers.GetSuppliers(msSupplierService))
	supplierRoute.Get("/:id", handlers.GetSupplier(msSupplierService))
	supplierRoute.Post("/", handlers.CreateSupplier(msSupplierService))
	supplierRoute.Put("/:id", handlers.UpdateSupplier(msSupplierService))
	supplierRoute.Delete("/", handlers.DeleteSupplier(msSupplierService))

	productCategoryRoute := api.Group("/v1/master/product_category", middleware.Protected())
	productCategoryRoute.Get("/", handlers.GetProductCategories(msProductCategoryService))
	productCategoryRoute.Get("/:id", handlers.GetProductCategory(msProductCategoryService))
	productCategoryRoute.Post("/", handlers.CreateProductCategory(msProductCategoryService))
	productCategoryRoute.Put("/:id", handlers.UpdateProductCategory(msProductCategoryService))
	productCategoryRoute.Delete("/:id", handlers.DeleteProductCategory(msProductCategoryService))

	productRoute := api.Group("/v1/master/product", middleware.Protected())
	productRoute.Get("/", handlers.GetProducts(msProductService))
	productRoute.Get("/:id/:based_on", handlers.GetProduct(msProductService))
	productRoute.Post("/", handlers.CreateProduct(msProductService))
	productRoute.Put("/:id", handlers.UpdateProduct(msProductService))
	productRoute.Delete("/", handlers.DeleteProduct(msProductService))
}
