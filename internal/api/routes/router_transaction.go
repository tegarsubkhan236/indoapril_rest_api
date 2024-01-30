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
	purchaseOrderRoute.Get("/", handlers.HandleGetPurchaseOrders(trPurchaseOrderService))
	purchaseOrderRoute.Get("/detail", handlers.HandleGetPurchaseOrder(trPurchaseOrderService))
	purchaseOrderRoute.Post("/", handlers.HandleAddPurchaseOrder(trPurchaseOrderService))

	receivingOrderRoute := api.Group("/v1/transaction/receiving_order", middleware.Protected())
	receivingOrderRoute.Post("/", handlers.HandleAddReceivingOrder(trReceivingOrderService))

	salesOrderRoute := api.Group("/v1/transaction/sales_order", middleware.Protected())
	salesOrderRoute.Get("/", handlers.HandleGetSalesOrders(trSalesOrderService))
	salesOrderRoute.Get("/:id", handlers.HandleGetSalesOrder(trSalesOrderService))
	salesOrderRoute.Post("/", handlers.HandleAddSalesOrder(trSalesOrderService))
}
