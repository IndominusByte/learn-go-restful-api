package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
)

const (
	prefix = "/categories"
)

func TestValidationCreateCategory(t *testing.T) {
	s := setupEnvironment()

	var data map[string]interface{}
	// login
	req, _ := http.NewRequest(http.MethodGet, "/token/get-token", nil)
	response := executeRequest(req, s)
	result := response.Result()

	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &data)

	accessToken := data["results"].(map[string]interface{})["access_token"]

	// validation
	b := new(bytes.Buffer)
	writer := multipart.NewWriter(b)
	writer.Close()

	req, _ = http.NewRequest(http.MethodPost, prefix, b)
	req.Header.Add("Authorization", "Bearer "+accessToken.(string))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	response = executeRequest(req, s)

	body, _ = io.ReadAll(response.Result().Body)

	fmt.Println(string(body))
}

func TestCreateCategory(t *testing.T) {
	// Create a New Server Struct
	s := setupEnvironment()

	// Create a New Request
	req, _ := http.NewRequest("GET", "/token/get-token", nil)

	// Execute Request
	response := executeRequest(req, s)

	body, _ := io.ReadAll(response.Result().Body)

	fmt.Println(string(body))
	fmt.Println(response.Result().StatusCode)
}
