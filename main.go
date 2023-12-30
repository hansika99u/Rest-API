package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "mydb"
)

var db *sql.DB

type Item struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	UnitPrice    float64 `json:"unit_price"`
	ItemCategory string  `json:"item_category"`
}
type Invoice struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	MobileNo    string `json:"mobileNo"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	BillingType string `json:"billingType"`
}

func init() {
	// Initialize the database connection
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to the database")
}

// func getItems(context *gin.Context) {
// 	context.IndentedJSON(http.StatusOK, items)
// }

func getItems(context *gin.Context) {
	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve items from the database"})
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.UnitPrice, &item.ItemCategory)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to scan rows"})
			return
		}
		items = append(items, item)
	}

	context.IndentedJSON(http.StatusOK, items)
}

// func addItems(context *gin.Context) {
// 	var newItems Item
// 	if err := context.BindJSON(&newItems); err != nil {
// 		return
// 	}

// 	items = append(items, newItems)

//		context.IndentedJSON(http.StatusCreated, newItems)
//	}
func addItems(context *gin.Context) {
	var newItem Item
	if err := context.BindJSON(&newItem); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	_, err := db.Exec("INSERT INTO items (name, unit_price, item_category) VALUES ($1, $2, $3)",
		newItem.Name, newItem.UnitPrice, newItem.ItemCategory)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to insert item into the database"})
		return
	}

	context.IndentedJSON(http.StatusCreated, newItem)
}

//	func getItemsById(id int) (*Item, error) {
//		for i, m := range items {
//			if m.ID == id {
//				return &items[i], nil
//			}
//		}
//		return nil, errors.New("item not found")
//	}
func getItemsById(id int) (*Item, error) {
	var item Item
	err := db.QueryRow("SELECT * FROM items WHERE id = $1", id).
		Scan(&item.ID, &item.Name, &item.UnitPrice, &item.ItemCategory)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("item not found")
		}
		return nil, errors.New("error retrieving item from the database")
	}

	return &item, nil
}

// func getItem(context *gin.Context) {
// 	id := context.Param("id")
// 	itemID, err := strconv.Atoi(id)
// 	if err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid item ID"})
// 		return
// 	}

// 	item, err := getItemsById(itemID)
// 	if err != nil {
// 		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
// 		return
// 	}

//		context.IndentedJSON(http.StatusOK, item)
//	}
func getItem(context *gin.Context) {
	id := context.Param("id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid item ID"})
		return
	}

	// Retrieve the item from the database
	item, err := getItemsById(itemID)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, item)
}

// func updateItem(context *gin.Context) {
// 	id := context.Param("id")
// 	itemID, err := strconv.Atoi(id)
// 	if err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid item ID"})
// 		return
// 	}

// 	item, err := getItemsById(itemID)
// 	if err != nil {
// 		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
// 		return
// 	}

// 	// Bind the new data from the request body
// 	if err := context.BindJSON(&item); err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
// 		return
// 	}

//		context.IndentedJSON(http.StatusOK, item)
//	}

func updateItem(context *gin.Context) {
	id := context.Param("id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid item ID"})
		return
	}

	// Retrieve the existing item from the database
	existingItem, err := getItemsById(itemID)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
		return
	}

	// Bind the new data from the request body
	if err := context.BindJSON(&existingItem); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	// Update the item in the database
	_, err = db.Exec("UPDATE items SET name = $1, unit_price = $2, item_category = $3 WHERE id = $4",
		existingItem.Name, existingItem.UnitPrice, existingItem.ItemCategory, itemID)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to update item in the database"})
		return
	}

	context.IndentedJSON(http.StatusOK, existingItem)
}

// func deleteItem(context *gin.Context) {
// 	id := context.Param("id")
// 	itemID, err := strconv.Atoi(id)
// 	if err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid item ID"})
// 		return
// 	}

// 	for i, item := range items {
// 		if item.ID == itemID {
// 			// Delete the item from the slice
// 			items = append(items[:i], items[i+1:]...)
// 			context.IndentedJSON(http.StatusNoContent, nil)
// 			return
// 		}
// 	}

// If the loop completes without finding the item
// 	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
// }

func deleteItem(context *gin.Context) {
	id := context.Param("id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid item ID"})
		return
	}

	// Retrieve the item from the database
	_, err = getItemsById(itemID)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
		return
	}

	// Delete the item from the database
	_, err = db.Exec("DELETE FROM items WHERE id = $1", itemID)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to delete item from the database"})
		return
	}

	context.IndentedJSON(http.StatusNoContent, nil)
}
func addInvoice(context *gin.Context) {
	var newInvoice Invoice
	if err := context.BindJSON(&newInvoice); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON"})
		return
	}

	_, err := db.Exec("INSERT INTO invoices (name, mobileNo, email, address, billingType) VALUES ($1, $2, $3, $4, $5)",
		newInvoice.Name, newInvoice.MobileNo, newInvoice.Email, newInvoice.Address, newInvoice.BillingType)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to insert invoice into the database"})
		return
	}

	context.IndentedJSON(http.StatusCreated, newInvoice)
}
func getInvoices(context *gin.Context) {
	rows, err := db.Query("SELECT * FROM invoices")
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve invoices from the database"})
		return
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
		err := rows.Scan(&invoice.ID, &invoice.Name, &invoice.MobileNo, &invoice.Email, &invoice.Address, &invoice.BillingType)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to scan rows"})
			return
		}
		invoices = append(invoices, invoice)
	}

	context.IndentedJSON(http.StatusOK, invoices)
}

func main() {
	router := gin.Default()
	router.GET("/items", getItems)
	router.GET("/items/:id", getItem)
	router.POST("/items", addItems)
	router.PUT("/items/:id", updateItem)
	router.DELETE("/items/:id", deleteItem)
	router.GET("/invoices", getInvoices)
	router.POST("/invoices", addInvoice)
	router.Run("localhost:9090")
}
