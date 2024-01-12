package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateCustomer(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	payload := []byte(`{"name":"Test User","age":25}`)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Customer
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Test User", response.Name)
	assert.Equal(t, 25, response.Age)
}

func TestUpdate404Customer(t *testing.T) {
	router := setupRouter()

	// Update 404 the customer
	w := httptest.NewRecorder()
	updatePayload := []byte(`{"id": 999999, "name":"Updated User","age":30}`)
	NotFoundURL := "/customers/999999"
	updateReq, _ := http.NewRequest("PUT", NotFoundURL, bytes.NewBuffer(updatePayload))
	updateReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, updateReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateCustomer(t *testing.T) {
	router := setupRouter()

	// Create a customer first
	createCustomerRequest := httptest.NewRecorder()
	payload := []byte(`{"name":"Test User","age":25}`)
	createReq, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(payload))
	createReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(createCustomerRequest, createReq)

	assert.Equal(t, http.StatusOK, createCustomerRequest.Code)

	var createdCustomer Customer
	json.Unmarshal(createCustomerRequest.Body.Bytes(), &createdCustomer)

	// Update the customer
	w := httptest.NewRecorder()
	updatePayload := []byte(`{"name":"Updated User","age":30}`)
	updateURL := fmt.Sprintf("/customers/%d", createdCustomer.ID)
	updateReq, _ := http.NewRequest("PUT", updateURL, bytes.NewBuffer(updatePayload))
	updateReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, updateReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Customer
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Updated User", response.Name)
	assert.Equal(t, 30, response.Age)
}

func TestDeleteCustomer(t *testing.T) {
	router := setupRouter()

	// Create a customer first
	createCustomerRequest := httptest.NewRecorder()
	payload := []byte(`{"name":"Jane Doe","age":30}`)
	createReq, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(payload))
	createReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(createCustomerRequest, createReq)

	assert.Equal(t, http.StatusOK, createCustomerRequest.Code)

	var createdCustomer Customer
	json.Unmarshal(createCustomerRequest.Body.Bytes(), &createdCustomer)

	// Delete it
	w := httptest.NewRecorder()
	getURL := fmt.Sprintf("/customers/%d", createdCustomer.ID)
	getOKReq, _ := http.NewRequest("DELETE", getURL, nil)
	getOKReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, getOKReq)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCustomer(t *testing.T) {
	router := setupRouter()

	// Create a customer first
	createCustomerRequest := httptest.NewRecorder()
	payload := []byte(`{"name":"Jane Doe","age":30}`)
	createReq, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(payload))
	createReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(createCustomerRequest, createReq)

	assert.Equal(t, http.StatusOK, createCustomerRequest.Code)

	var createdCustomer Customer
	json.Unmarshal(createCustomerRequest.Body.Bytes(), &createdCustomer)

	// Get - OK
	wOK := httptest.NewRecorder()
	getURL := fmt.Sprintf("/customers/%d", createdCustomer.ID)
	getOKReq, _ := http.NewRequest("GET", getURL, nil)
	getOKReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wOK, getOKReq)
	assert.Equal(t, http.StatusOK, wOK.Code)

	var getResp Customer
	json.Unmarshal(wOK.Body.Bytes(), &getResp)
	assert.Equal(t, "Jane Doe", getResp.Name)
	assert.Equal(t, 30, getResp.Age)

	// Get - Not found
	wNA := httptest.NewRecorder()
	getNAReq, _ := http.NewRequest("GET", "/customers/999", nil)
	getNAReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wNA, getNAReq)
	assert.Equal(t, http.StatusNotFound, wNA.Code)
}

func setupRouter() *gin.Engine {
	// Connect to SQLite database
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
	}

	router := gin.Default()

	router.POST("/customers", createCustomer)
	router.PUT("/customers/:id", updateCustomer)
	router.DELETE("/customers/:id", deleteCustomer)
	router.GET("/customers/:id", getCustomerByID)

	return router
}
