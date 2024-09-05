package main

import (
	"log"
	"os"

	"github.com/ChandanJnv/ecommerce-cart-golang/controllers"
	"github.com/ChandanJnv/ecommerce-cart-golang/database"
	"github.com/ChandanJnv/ecommerce-cart-golang/middleware"
	"github.com/ChandanJnv/ecommerce-cart-golang/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("--------------------------- Starting ecommerce-cart-golang ---------------------------")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserDatabase(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	routes.AddAddressRoutes(router)
	router.Use(middleware.Authentication())

	router.POST("/addtocart", app.AddToCart())
	router.DELETE("/removeitem", app.RemoveItem())
	router.GET("/cart", app.GetItemFromCart())
	router.POST("/cartcheckout", app.BuyFromCart())
	router.POST("/instantbuy", app.InstantBuy())

	log.Println("Server starting on localhost:" + port)
	log.Fatal(router.Run(":" + port))

}
