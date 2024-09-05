package routes

import (
	"github.com/ChandanJnv/ecommerce-cart-golang/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
}

func AddAddressRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/address/addaddress", controllers.AddAddress())
	incomingRoutes.POST("/address/edithomeaddress", controllers.EditHomeAddress())
	incomingRoutes.POST("/address/editworkaddress", controllers.EditWorkAddress())
	incomingRoutes.DELETE("/address/deleteaddress", controllers.DeleteAddress())
}
