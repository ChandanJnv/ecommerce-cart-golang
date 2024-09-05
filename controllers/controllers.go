package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ChandanJnv/ecommerce-cart-golang/database"
	"github.com/ChandanJnv/ecommerce-cart-golang/models"
	generate "github.com/ChandanJnv/ecommerce-cart-golang/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserCollection    *mongo.Collection = database.UserDatabase(database.Client, "Users")
	ProductCollection *mongo.Collection = database.UserDatabase(database.Client, "Products")
	Validate                            = validator.New()
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	valid := true
	msg := ""
	if err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword)); err != nil {
		msg = "Login or password is incorrect"
		valid = false
	}
	return valid, msg
}

// localhost:8000/users/signup
//
//	{
//	    "first_name": "alpha",
//	    "last_name": "beta",
//	    "password": "alpha@123",
//	    "email": "alpha@beta.com",
//	    "phone": "9876543210"
//	}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			log.Println("Failed to bind user", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			log.Println("Failed to validate user", validationErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already registerd"})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone no. is already used"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshToken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, insetErr := UserCollection.InsertOne(ctx, user)
		if insetErr != nil {
			log.Println(insetErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}

		c.JSON(http.StatusCreated, "Successfully signed in")
	}
}

// localhost:8000/users/login
//
//	{
//	    "email":"alpha@beta.com",
//	    "password":"alpha@123"
//	}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		var foundUser models.User
		if err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refreshToken, err := generate.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_ID)
		if err != nil {
			fmt.Println("failed to generate token.", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		generate.UpdateAllTokens(token, refreshToken, foundUser.User_ID)

		c.JSON(http.StatusFound, fmt.Sprintf("token: %s", token))
	}

}

// localhost:8000/admin/addproduct
//
//	{
//	    "product_name": "laptop",
//	    "price": 200,
//	    "Rating": 4,
//	    "Image": "/img/path/dotjpg"
//	}
func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var products models.Product

		if err := c.BindJSON(&products); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		products.Product_ID = primitive.NewObjectID()
		if _, err := ProductCollection.InsertOne(ctx, products); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not inserted"})
			return
		}

		c.JSON(http.StatusOK, "successfully added")
	}
}

// localhost:8000/users/productview
func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something wend wrong: "+err.Error())
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &productList); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid: "+err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, productList)

	}
}

// localhost:8000/users/search?name=laptop
func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search index"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchQueryDB, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			log.Println(err)
			c.IndentedJSON(404, "something wend wrong: "+err.Error())
			return
		}
		defer searchQueryDB.Close(ctx)

		if err := searchQueryDB.All(ctx, &searchProducts); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid: "+err.Error())
			return
		}

		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request: "+err.Error())
			return
		}

		c.IndentedJSON(200, searchProducts)
	}
}
