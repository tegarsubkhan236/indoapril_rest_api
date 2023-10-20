package router

import (
	"example/api/controller"
	"example/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func setupMasterRoutes(v1 fiber.Router) {
	masterRoute := v1.Group("/master", middleware.Protected())

	supplierRoute := masterRoute.Group("/supplier")
	supplierRoute.Get("/", controller.GetSuppliers)
	supplierRoute.Get("/:id", controller.GetSupplier)
	supplierRoute.Post("/", controller.CreateSupplier)
	supplierRoute.Put("/:id", controller.UpdateSupplier)
	supplierRoute.Delete("/", controller.DeleteSupplier)

	productCategoryRoute := masterRoute.Group("/product_category")
	productCategoryRoute.Get("/", controller.GetProductCategories)
	productCategoryRoute.Get("/:id", controller.GetProductCategory)
	productCategoryRoute.Post("/", controller.CreateProductCategory)
	productCategoryRoute.Put("/:id", controller.UpdateProductCategory)
	productCategoryRoute.Delete("/:id", controller.DeleteProductCategory)

	productRoute := masterRoute.Group("/product")
	productRoute.Get("/", controller.GetProducts)
	productRoute.Get("/:id/:based_on", controller.GetProduct)
	productRoute.Post("/", controller.CreateProduct)
	productRoute.Put("/:id", controller.UpdateProduct)
	productRoute.Delete("/", controller.DeleteProduct)
}
