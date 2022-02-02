package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	categoriesentity "github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
	"github.com/stretchr/testify/assert"
)

/*
standar validation

- empty []
- required []
- type data []
- format regex []
- minimum []
- maximum []
- file []
*/

const (
	prefix = "/categories"
	name   = "testtestingtest"
)

func TestValidationCreateCategory(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	// login
	req, _ := http.NewRequest(http.MethodGet, "/token/get-token", nil)
	response := executeRequest(req, s)
	result := response.Result()

	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &data)

	accessToken := data["results"].(map[string]interface{})["access_token"].(string)

	/*
	   standar validation

	   - empty [x]
	   - required [x]
	   - type data []
	   - format regex []
	   - minimum [x]
	   - maximum [x]
	   - file [x]
	*/

	tests := [...]struct {
		name     string
		expected string
		form     map[string]string
	}{
		{
			name:     "empty",
			expected: "Invalid input type.",
			form:     map[string]string{"icon": "@", "name": ""},
		},
		{
			name:     "required icon",
			expected: "Image is required.",
			form:     map[string]string{},
		},
		{
			name:     "required name",
			expected: "Missing data for required field.",
			form:     map[string]string{"icon": "@/app/static/test_image/image.jpeg"},
		},
		{
			name:     "minimum",
			expected: "Shorter than minimum length 3.",
			form:     map[string]string{"icon": "@/app/static/test_image/image.jpeg", "name": "a"},
		},
		{
			name:     "maximum",
			expected: "Longer than maximum length 100.",
			form:     map[string]string{"icon": "@/app/static/test_image/image.jpeg", "name": createMaximum(200)},
		},
		{
			name:     "danger file extension",
			expected: "Image must be between jpeg, png.",
			form:     map[string]string{"icon": "@/app/static/test_image/test.txt", "name": "aaa"},
		},
		{
			name:     "not valid file extension",
			expected: "Image must be between jpeg, png.",
			form:     map[string]string{"icon": "@/app/static/test_image/test.gif", "name": "aaa"},
		},
		{
			name:     "file cannot grater than 4 Mb",
			expected: "An image cannot greater than 4 Mb.",
			form:     map[string]string{"icon": "@/app/static/test_image/size.png", "name": "aaa"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ct, b, err := createForm(test.form)
			if err != nil {
				panic(err)
			}

			req, _ = http.NewRequest(http.MethodPost, prefix, b)
			req.Header.Add("Authorization", "Bearer "+accessToken)
			req.Header.Set("Content-Type", ct)

			response = executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "empty":
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_body"].(string))
			case "required icon", "danger file extension", "not valid file extension", "file cannot grater than 4 Mb":
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["icon"].(string))
			case "required name", "minimum", "maximum":
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["name"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)

		})
	}
}

func TestCreateCategory(t *testing.T) {
	repo, s := setupEnvironment()

	var data map[string]interface{}
	var form map[string]string

	// login
	req, _ := http.NewRequest(http.MethodGet, "/token/get-token", nil)
	response := executeRequest(req, s)
	result := response.Result()

	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &data)

	accessToken := data["results"].(map[string]interface{})["access_token"].(string)

	tests := [...]struct {
		name       string
		expected   string
		statusCode int
	}{
		{
			name:       "success",
			expected:   "Successfully create category.",
			statusCode: 201,
		},
		{
			name:       "duplicate name",
			expected:   "The name has already been taken.",
			statusCode: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form = map[string]string{"icon": "@/app/static/test_image/image.jpeg", "name": name}
			ct, b, err := createForm(form)
			if err != nil {
				panic(err)
			}

			req, _ = http.NewRequest(http.MethodPost, prefix, b)
			req.Header.Add("Authorization", "Bearer "+accessToken)
			req.Header.Set("Content-Type", ct)

			response = executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_app"].(string))
			assert.Equal(t, test.statusCode, response.Result().StatusCode)

			switch test.name {
			case "success":
				// check image exists in directory
				payload := categoriesentity.FormCreateUpdateSchema{Name: name}
				r, _ := repo.categoriesRepo.GetCategoryByName(context.Background(), &payload)
				assert.True(t, fileExists("/app/static/icon-categories/"+r.Icon))
			}

		})
	}
}

func TestValidationGetAllCategories(t *testing.T) {
	/*
	   standar validation

	   - empty [x]
	   - required [x]
	   - type data [x]
	   - format regex []
	   - minimum [x]
	   - maximum []
	   - file []
	*/
	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name   string
		prefix string
	}{
		{
			name:   "empty",
			prefix: prefix,
		},
		{
			name:   "required",
			prefix: prefix + "?page=&per_page=&q=",
		},
		{
			name:   "type data",
			prefix: prefix + "?page=a&per_page=a&q=1",
		},
		{
			name:   "minimum",
			prefix: prefix + "?page=-1&per_page=-1&q=",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.prefix, nil)

			response := executeRequest(req, s)

			body, _ := io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "empty", "required":
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["page"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["per_page"].(string))
			case "type data":
				assert.Equal(t, "Invalid input type.", data["detail_message"].(map[string]interface{})["_body"].(string))
			case "minimum":
				assert.Equal(t, "Must be greater than or equal to 1.", data["detail_message"].(map[string]interface{})["page"].(string))
				assert.Equal(t, "Must be greater than or equal to 1.", data["detail_message"].(map[string]interface{})["per_page"].(string))
			}

			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}
}

func TestGetAllCategories(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	req, _ := http.NewRequest(http.MethodGet, prefix+"?page=1&per_page=1&q=t", nil)

	response := executeRequest(req, s)

	body, _ := io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	assert.NotNil(t, data["results"].(map[string]interface{})["data"])
	assert.Equal(t, 200, response.Result().StatusCode)
}
