package router

import (
	"example/api/controller"
	"example/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func setupTransactionRoutes(v1 fiber.Router) {
	transactionRoute := v1.Group("/transaction", middleware.Protected())

	purchaseOrderRoute := transactionRoute.Group("/purchase_order")
	purchaseOrderRoute.Get("/", controller.GetPurchaseOrders)
	purchaseOrderRoute.Get("/detail", controller.GetPurchaseOrder)
	purchaseOrderRoute.Post("/", controller.CreatePurchaseOrder)

	receivingOrderRoute := transactionRoute.Group("/receiving_order")
	receivingOrderRoute.Post("/", controller.CreateReceivingOrder)

	salesOrderRoute := transactionRoute.Group("/sales_order")
	salesOrderRoute.Get("/", controller.GetSalesOrders)
	salesOrderRoute.Get("/:id", controller.GetSalesOrder)
	salesOrderRoute.Post("/", controller.CreateSalesOrder)
}
