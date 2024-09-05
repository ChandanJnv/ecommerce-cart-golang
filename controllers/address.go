package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ChandanJnv/ecommerce-cart-golang/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// localhost:8000/address/addaddress?id={user_id}
//
//	{
//	    "house_name":"home address",
//	    "street_name":"home street",
//	    "city_name":"home city",
//	    "pin_code":"654321"
//	}
func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "user id is empty"})
			c.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal server error")
			return
		}
		var addressess models.Address
		addressess.Address_ID = primitive.NewObjectID()
		if err := c.BindJSON(&addressess); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var addressinfo []bson.M

		if err := pointCursor.All(ctx, &addressinfo); err != nil {
			panic(err)
		}

		var size int32

		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addressess}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				log.Println(err)
			}
		} else {
			c.IndentedJSON(400, "Not Allowed")
			return
		}
		c.IndentedJSON(200, "Address added successfully")

	}
}

// localhost:8000/address/edithomeaddress?id={user_id}
//
//	{
//	    "house_name":"home address",
//	    "street_name":"home street",
//	    "city_name":"home city",
//	    "pin_code":"654321"
//	}
func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal server error")
		}

		var editaddress models.Address
		if c.BindJSON(&editaddress); err != nil {
			log.Println(err)
			c.IndentedJSON(500, err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code", Value: editaddress.Pincode}}}}

		if _, err := UserCollection.UpdateOne(ctx, filter, update); err != nil {
			log.Println(err)
			c.IndentedJSON(500, "something went wrong")
			return
		}
		ctx.Done()

		c.IndentedJSON(200, "successfully updated the home address")
	}
}

// localhost:8000/address/editworkaddress?id={user_id}
//
//	{
//		"house_name": "work address",
//		"street_name": "work street",
//		"city_name": "work city",
//		"pin_code": "123456"
//	  }
func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal server error")
		}

		var editaddress models.Address
		if c.BindJSON(&editaddress); err != nil {
			log.Println(err)
			c.IndentedJSON(500, err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House}, {Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: "address.1.pin_code", Value: editaddress.Pincode}}}}
		if _, err := UserCollection.UpdateOne(ctx, filter, update); err != nil {
			log.Println(err)
			c.IndentedJSON(500, "something went wrong")
			return
		}
		ctx.Done()

		c.IndentedJSON(200, "Successfully updated address")

	}
}

// localhost:8000/address/deleteaddress?id={user_id}
func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Eroor": "Invalid Search Index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error: "+err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "wrong command: "+err.Error())
			return
		}

		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Successfully Deleted")
	}
}
