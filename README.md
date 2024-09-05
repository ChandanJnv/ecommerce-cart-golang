# E-commerce Cart API in Golang
A backend API for managing an e-commerce cart system, built using Golang, gin framework, jwt, mongoDB. It includes functionalities like adding, updating, and removing items from the cart, managing user addresses, Product, and user authentication.

## Features
- User authentication (signup, login)
- Product management (add-admin, view, search)
- Cart management (add, remove, view, checkout, instantBuy)
- Address management (add, update, delete)

## Installation
1. **Clone the repository:**
    ```bash
    git clone https://github.com/ChandanJnv/ecommerce-cart-golang.git
    cd ecommerce-cart-golang
    ```

2. **Install dependencies:**
    ```bash
    go mod download
    ```

3. **Install and Setup MongoDB:**
    [install mongodb](https://www.mongodb.com/docs/manual/administration/install-community/)

3. **Set SECRET_KEY:**
    ```bash
    export SECRET_KEY="your-secret-key"
    ```

3. **Run the application:**
    ```bash
    go run main.go
    ```

## API Endpoints

### User Endpoints

#### **Register User**
- **URL**: `/users/signup`
- **Method**: `POST`
- **Body**:
    ```json
    {
        "first_name": "alpha",
        "last_name": "beta",
        "password": "alpha@123",
        "email": "alpha@beta.com",
        "phone": "9876543210"
    }
    ```
- **Response**:
    ```json
    {
      "message": "Successfully signed in."
    }
    ```

#### **Login User**
- **URL**: `/users/login`
- **Method**: `POST`
- **Body**:
    ```json
	{
	    "email":"alpha@beta.com",
	    "password":"alpha@123"
	}
    ```
- **Response**:
    ```json
    {
      "token": "your-authentication-token"
    }
    ```

### Product Endpoints

#### **Admin Add Product**
- **URL**: `/admin/addproduct`
- **Method**: `POST`
- **Body**:
    ```json
	{
	    "product_name": "laptop",
	    "price": 200,
	    "Rating": 4,
	    "Image": "/img/path/dotjpg"
    }
    ```
- **Response**:
    ```
    successfully added
    ```

#### **View All Products**
- **URL**: `/users/productview`
- **Method**: `GET`
- **Response**:
    ```json
    {
        "Product_ID": "66d4321250820c57cfb26557",
        "product_name": "laptop",
        "price": 200,
        "rating": 4,
        "image": "/img/path/dotjpg"
    },
    {
        "Product_ID": "66d4330450820c57cfb26558",
        "product_name": "mobile",
        "price": 20000,
        "rating": 4,
        "image": "/img/path/dotjpg"
    },
    {
        "Product_ID": "66d4331d50820c57cfb26559",
        "product_name": "table",
        "price": 500,
        "rating": 4,
        "image": "/img/path/dotjpg"
    }

    ```

#### **Search Product**
- **URL**: `/users/search?name={product_name}`
- **Method**: `GET`
- **Headers**: 
    - `token`: `<token>`
- **Response**:
    ```json
    {
        "Product_ID": "66d4321250820c57cfb26557",
        "product_name": "laptop",
        "price": 200,
        "rating": 4,
        "image": "/img/path/dotjpg"
    }
    ```

### Cart Endpoints

#### **Add Item to Cart**
- **URL**: `/addtocart?id={product_id}&userID={user_id}`
- **Method**: `POST`
- **Headers**: 
    - `token`: `<token>`
- **Body**:
    ```
    Successfully added to the cart
    ```

#### **Remove Item from Cart**
- **URL**: `/removeitem?userID={user_id}&id={product_id}`
- **Method**: `DELETE`
- **Headers**: 
    - `token`: `<token>`
- **Body**:
    ```
    Successfully removed item from the cart
    ```

#### **Get Cart Details**
- **URL**: `/cart?id={user_id}`
- **Method**: `GET`
- **Headers**: 
    - `token`: `<token>`
- **Response**:
    ```json
    {
        20000[
        {
            "Product_ID": "66d4330450820c57cfb26558",
            "product_name": "mobile",
            "price": 20000,
            "rating": 4,
            "image": "/img/path/dotjpg"
        }
        ]
    }
    ```


#### **Buy From Cart**
- **URL**: `/cartcheckout?userID={user_id}`
- **Method**: `POST`
- **Headers**: 
    - `token`: `<token>`
- **Response**:
    ```
    Successfully placed the order
    ```


#### **Buy Now**
- **URL**: `/instantbuy?id={product_id}&userID={user_id}s`
- **Method**: `POST`
- **Headers**: 
    - `token`: `<token>`
- **Response**:
    ```
    Successfully placed the order
    ```

### Address Endpoints

#### **Add New Address**
- **URL**: `/address/addaddress?id={user_id}`
- **Method**: `POST`
- **Headers**: 
    - `token`: `<token>`
- **Body**:
    ```json
	{
	    "house_name":"my address",
	    "street_name":"my street",
	    "city_name":"my city",
	    "pin_code":"654321"
	}
    ```
- **Response**:
    ```
    Address added successfully
    ```

#### **Edit Home Address**
- **URL**: `/address/edithomeaddress?id={user_id}`
- **Method**: `POST`
- **Headers**: 
    - `token`: `<token>`
- **Body**:
    ```json
	{
	    "house_name":"home address",
	    "street_name":"home street",
	    "city_name":"home city",
	    "pin_code":"654321"
	}
    ```
- **Response**:
    ```
    successfully updated the home address
    ```

#### **Edit Work Address**
- **URL**: `/address/editworkaddress?id={user_id}`
- **Method**: `POST`
- **Headers**: 
    - `token`: `<token>`
- **Body**:
    ```json
    {
		"house_name": "work address",
		"street_name": "work street",
		"city_name": "work city",
		"pin_code": "123456"
    }
    ```
- **Response**:
    ```
    Successfully updated address
    ```

#### **Delete Address**
- **URL**: `/address/deleteaddress?id={user_id}`
- **Method**: `DELETE`
- **Headers**: 
    - `token`: `<token>`
- **Response**:
    ```
    Successfully Deleted
    ```

## License

This project is licensed under the MIT License.
