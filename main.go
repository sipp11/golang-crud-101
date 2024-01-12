package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Customer model
type Customer struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *gorm.DB
var err error

func main() {
	// Connect to SQLite database
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Customer{})

	// Create some initial customer data
	createInitialData()

	// Initialize Gin router
	router := gin.Default()

	// Define REST API routes
	router.POST("/customers", createCustomer)
	router.PUT("/customers/:id", updateCustomer)
	router.DELETE("/customers/:id", deleteCustomer)
	router.GET("/customers/:id", getCustomerByID)

	// Run the server
	router.Run(":8080")
}

func createInitialData() {
	customer1 := Customer{Name: "John Doe", Age: 21}
	customer2 := Customer{Name: "Jane Doe", Age: 20}

	db.Create(&customer1)
	db.Create(&customer2)
}

func createCustomer(c *gin.Context) {
	var customer Customer
	c.BindJSON(&customer)

	db.Create(&customer)
	c.JSON(http.StatusOK, customer)
}

func updateCustomer(c *gin.Context) {
	ID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.AbortWithStatus(401)
		fmt.Println(err)
		return
	}
	var body Customer
	c.BindJSON(&body)
	if body.ID != 0 && ID != int(body.ID) {
		// if body contain ID, then check if both values match
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID not matched"})
		return
	}

	var customer Customer
	if err := db.Where("id = ?", ID).First(&customer).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	customer.Name = body.Name
	customer.Age = body.Age
	db.Save(&customer)
	c.JSON(http.StatusOK, customer)
}

func deleteCustomer(c *gin.Context) {
	id := c.Params.ByName("id")
	var customer Customer

	if err := db.Where("id = ?", id).First(&customer).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		db.Delete(&customer)
		c.JSON(http.StatusOK, gin.H{"id #" + id: "deleted"})
	}
}

func getCustomerByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var customer Customer

	if err := db.Where("id = ?", id).First(&customer).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, customer)
	}
}
