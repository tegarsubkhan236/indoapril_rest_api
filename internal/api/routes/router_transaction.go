package routes

import (
	"example/internal/api/handlers"
	"example/internal/api/middleware"
	"example/internal/pkg/models/tr_purchase_order"
	"example/internal/pkg/models/tr_receiving_order"
	"example/internal/pkg/models/tr_sales_order"
	"github.com/gofiber/fiber/v2"
)

func SetupTransactionRoutes(
	api fiber.Router,
	trPurchaseOrderService tr_purchase_order.Service,
	trReceivingOrderService tr_receiving_order.Service,
	trSalesOrderService tr_sales_order.Service,
) {
	purchaseOrderRoute := api.Group("/v1/transaction/purchase_order", middleware.Protected())
	purchaseOrderRoute.Get("/", handlers.GetPurchaseOrders(trPurchaseOrderService))
	purchaseOrderRoute.Get("/detail", handlers.GetPurchaseOrder(trPurchaseOrderService))
	purchaseOrderRoute.Post("/", handlers.CreatePurchaseOrder(trPurchaseOrderService))

	receivingOrderRoute := api.Group("/v1/transaction/receiving_order", middleware.Protected())
	receivingOrderRoute.Post("/", handlers.CreateReceivingOrder(trReceivingOrderService))

	salesOrderRoute := api.Group("/v1/transaction/sales_order", middleware.Protected())
	salesOrderRoute.Get("/", handlers.GetSalesOrders(trSalesOrderService))
	salesOrderRoute.Get("/:id", handlers.GetSalesOrder(trSalesOrderService))
	salesOrderRoute.Post("/", handlers.CreateSalesOrder(trSalesOrderService))
}
