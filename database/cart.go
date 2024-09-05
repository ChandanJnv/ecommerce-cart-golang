package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ChandanJnv/ecommerce-cart-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	errCantFindProduct    = errors.New("can't find product")
	ErrCantDecodeProduct  = errors.New("can't find product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot update user")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item form the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
	ErrCantFindProduct    = errors.New("cannot find the product")
)

// localhost:8000/addtocart?id={product_id}&userID={user_id}
func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDB, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return errCantFindProduct
	}

	var productCart []models.ProductUser
	if err := searchFromDB.All(ctx, &productCart); err != nil {
		log.Println(err)
		return ErrCantDecodeProduct
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	fmt.Println("productCart: ", productCart)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}
	if _, err := userCollection.UpdateOne(ctx, filter, update); err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil
}

// localhost:8000/removeitem?userID={user_id}&id={product_id}
func RemoveCartIterm(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}

	if _, err := userCollection.UpdateMany(ctx, filter, update); err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}

	return nil
}

// localhost:8000/cartcheckout?userID={user_id}
func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{{Key: "$sum", Value: "$usercart.price"}}}}}}
	currentResult, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getUserCart []bson.M
	if err := currentResult.All(ctx, &getUserCart); err != nil {
		panic(err)
	}
	var totalPrice int64
	for _, userItem := range getUserCart {
		price := userItem["total"]
		totalPrice = (price.(int64))
	}

	orderCart.Price = &totalPrice
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	if _, err := userCollection.UpdateMany(ctx, filter, update); err != nil {
		log.Println(err)
		return err
	}

	if err := userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems); err != nil {
		log.Println(err)
		return err
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}
	if _, err := userCollection.UpdateOne(ctx, filter2, update2); err != nil {
		log.Println(err)
		return err
	}

	getUserCartEmpty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: getUserCartEmpty}}}}

	if _, err := userCollection.UpdateOne(ctx, filter3, update3); err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var productDetails models.ProductUser
	var ordersDetails models.Order

	ordersDetails.Order_ID = primitive.NewObjectID()
	ordersDetails.Ordered_At = time.Now()
	ordersDetails.Order_Cart = (make([]models.ProductUser, 0))
	ordersDetails.Payment_method.COD = true

	if err := prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails); err != nil {
		log.Println(err)
		return err
	}
	ordersDetails.Price = productDetails.Price
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordersDetails}}}}
	if _, err := userCollection.UpdateOne(ctx, filter, update); err != nil {
		log.Println(err)
		return err
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}
	if _, err := userCollection.UpdateOne(ctx, filter2, update2); err != nil {
		log.Panicln(err)
		return err
	}

	return nil
}
