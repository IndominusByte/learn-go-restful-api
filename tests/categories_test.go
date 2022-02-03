package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
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
		form       map[string]string
		statusCode int
	}{
		{
			name:       "success",
			expected:   "Successfully create category.",
			form:       map[string]string{"icon": "@/app/static/test_image/image.jpeg", "name": name},
			statusCode: 201,
		},
		{
			name:       "create another",
			expected:   "Successfully create category.",
			form:       map[string]string{"icon": "@/app/static/test_image/image.jpeg", "name": name + "2"},
			statusCode: 201,
		},
		{
			name:       "duplicate name",
			expected:   "The name has already been taken.",
			form:       map[string]string{"icon": "@/app/static/test_image/image.jpeg", "name": name},
			statusCode: 400,
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
		name string
		url  string
	}{
		{
			name: "empty",
			url:  prefix,
		},
		{
			name: "required",
			url:  prefix + "?page=&per_page=&q=",
		},
		{
			name: "type data",
			url:  prefix + "?page=a&per_page=a&q=1",
		},
		{
			name: "minimum",
			url:  prefix + "?page=-1&per_page=-1&q=",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)

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

func TestValidationGetCategoryById(t *testing.T) {
	/*
	   standar validation

	   - empty []
	   - required []
	   - type data [x]
	   - format regex []
	   - minimum [x]
	   - maximum []
	   - file []
	*/

	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name string
		url  string
	}{
		{
			name: "type data",
			url:  prefix + "/abc",
		},
		{
			name: "minimum",
			url:  prefix + "/-1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)

			response := executeRequest(req, s)

			body, _ := io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			assert.Equal(t, "404 page not found", strings.TrimSuffix(string(body), "\n"))
			assert.Equal(t, 404, response.Result().StatusCode)
		})
	}
}

func TestGetCategoryById(t *testing.T) {
	repo, s := setupEnvironment()

	var data map[string]interface{}

	// get id
	payload := categoriesentity.FormCreateUpdateSchema{Name: name}
	category, _ := repo.categoriesRepo.GetCategoryByName(context.Background(), &payload)

	tests := [...]struct {
		name string
		url  string
	}{
		{
			name: "not found",
			url:  prefix + "/99999999",
		},
		{
			name: "success",
			url:  prefix + "/" + strconv.Itoa(category.Id),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)

			response := executeRequest(req, s)

			body, _ := io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "not found":
				assert.Equal(t, "Category not found.", data["detail_message"].(map[string]interface{})["_app"].(string))
				assert.Equal(t, 404, response.Result().StatusCode)
			case "success":
				assert.NotNil(t, data["results"])
				assert.Equal(t, 200, response.Result().StatusCode)
			}
		})
	}
}

func TestValidationUpdateCategoryById(t *testing.T) {
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
	   - type data [x]
	   - format regex []
	   - minimum [x]
	   - maximum [x]
	   - file [x]
	*/

	tests := [...]struct {
		name     string
		expected string
		url      string
		form     map[string]string
	}{
		{
			name:     "empty",
			expected: "Invalid input type.",
			url:      prefix + "/1",
			form:     map[string]string{"icon": "@", "name": ""},
		},
		{
			name:     "optional icon & required name",
			expected: "Missing data for required field.",
			url:      prefix + "/1",
			form:     map[string]string{},
		},
		{
			name:     "type data path",
			expected: "404 page not found",
			url:      prefix + "/abc",
			form:     map[string]string{},
		},
		{
			name:     "minimum path",
			expected: "404 page not found",
			url:      prefix + "/-1",
			form:     map[string]string{},
		},
		{
			name:     "minimum form",
			expected: "Shorter than minimum length 3.",
			url:      prefix + "/1",
			form:     map[string]string{"name": "a"},
		},
		{
			name:     "maximum form",
			expected: "Longer than maximum length 100.",
			url:      prefix + "/1",
			form:     map[string]string{"name": createMaximum(200)},
		},
		{
			name:     "danger file extension",
			expected: "Image must be between jpeg, png.",
			url:      prefix + "/1",
			form:     map[string]string{"icon": "@/app/static/test_image/test.txt", "name": "aaa"},
		},
		{
			name:     "not valid file extension",
			expected: "Image must be between jpeg, png.",
			url:      prefix + "/1",
			form:     map[string]string{"icon": "@/app/static/test_image/test.gif", "name": "aaa"},
		},
		{
			name:     "file cannot grater than 4 Mb",
			expected: "An image cannot greater than 4 Mb.",
			url:      prefix + "/1",
			form:     map[string]string{"icon": "@/app/static/test_image/size.png", "name": "aaa"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ct, b, err := createForm(test.form)
			if err != nil {
				panic(err)
			}

			req, _ = http.NewRequest(http.MethodPut, test.url, b)
			req.Header.Add("Authorization", "Bearer "+accessToken)
			req.Header.Set("Content-Type", ct)

			response = executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "empty":
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_body"].(string))
				assert.Equal(t, 422, response.Result().StatusCode)
			case "optional icon & required name":
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["name"].(string))
				assert.Nil(t, data["detail_message"].(map[string]interface{})["icon"])
				assert.Equal(t, 422, response.Result().StatusCode)
			case "type data path", "minimum path":
				assert.Equal(t, "404 page not found", strings.TrimSuffix(string(body), "\n"))
				assert.Equal(t, 404, response.Result().StatusCode)
			case "minimum form", "maximum form":
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["name"].(string))
				assert.Equal(t, 422, response.Result().StatusCode)
			case "danger file extension", "not valid file extension", "file cannot grater than 4 Mb":
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["icon"].(string))
				assert.Equal(t, 422, response.Result().StatusCode)
			}
		})
	}

}

func TestUpdateCategoryById(t *testing.T) {
	repo, s := setupEnvironment()

	var data map[string]interface{}

	// get id
	payload := categoriesentity.FormCreateUpdateSchema{Name: name}
	category, _ := repo.categoriesRepo.GetCategoryByName(context.Background(), &payload)

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
		url        string
		form       map[string]string
		statusCode int
	}{
		{
			name:       "not found",
			expected:   "Category not found.",
			url:        prefix + "/99999999",
			form:       map[string]string{"name": name},
			statusCode: 404,
		},
		{
			name:       "success",
			expected:   "Successfully update the category.",
			url:        prefix + "/" + strconv.Itoa(category.Id),
			form:       map[string]string{"icon": "@/app/static/test_image/image.jpeg", "name": name},
			statusCode: 200,
		},
		{
			name:       "duplicate name",
			expected:   "The name has already been taken.",
			url:        prefix + "/" + strconv.Itoa(category.Id),
			form:       map[string]string{"icon": "@/app/static/test_image/image.jpeg", "name": name + "2"},
			statusCode: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ct, b, err := createForm(test.form)
			if err != nil {
				panic(err)
			}

			req, _ = http.NewRequest(http.MethodPut, test.url, b)
			req.Header.Add("Authorization", "Bearer "+accessToken)
			req.Header.Set("Content-Type", ct)

			response = executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "not found", "duplicate name":
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_app"].(string))
			case "success":
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_app"].(string))
				// check the previous image doesn't exists in directory
				assert.False(t, fileExists("/app/static/icon-categories/"+category.Icon))
				// check image exists in directory
				payload := categoriesentity.FormCreateUpdateSchema{Name: name}
				r, _ := repo.categoriesRepo.GetCategoryByName(context.Background(), &payload)
				assert.True(t, fileExists("/app/static/icon-categories/"+r.Icon))
			}

			assert.Equal(t, test.statusCode, response.Result().StatusCode)
		})
	}
}

func TestValidationDeleteCategoryById(t *testing.T) {
	/*
	   standar validation

	   - empty []
	   - required []
	   - type data [x]
	   - format regex []
	   - minimum [x]
	   - maximum []
	   - file []
	*/

	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name string
		url  string
	}{
		{
			name: "type data",
			url:  prefix + "/abc",
		},
		{
			name: "minimum",
			url:  prefix + "/-1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, test.url, nil)

			response := executeRequest(req, s)

			body, _ := io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			assert.Equal(t, "404 page not found", strings.TrimSuffix(string(body), "\n"))
			assert.Equal(t, 404, response.Result().StatusCode)
		})
	}
}

func TestDeleteCategoryById(t *testing.T) {
	repo, s := setupEnvironment()

	var data map[string]interface{}

	// get id
	payload := categoriesentity.FormCreateUpdateSchema{Name: name}
	category1, _ := repo.categoriesRepo.GetCategoryByName(context.Background(), &payload)

	payload = categoriesentity.FormCreateUpdateSchema{Name: name + "2"}
	category2, _ := repo.categoriesRepo.GetCategoryByName(context.Background(), &payload)

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
		url        string
		statusCode int
	}{
		{
			name:       "not found",
			expected:   "Category not found.",
			url:        prefix + "/99999999",
			statusCode: 404,
		},
		{
			name:       "delete 1",
			expected:   "Successfully delete the category.",
			url:        prefix + "/" + strconv.Itoa(category1.Id),
			statusCode: 200,
		},
		{
			name:       "delete 2",
			expected:   "Successfully delete the category.",
			url:        prefix + "/" + strconv.Itoa(category2.Id),
			statusCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ = http.NewRequest(http.MethodDelete, test.url, nil)
			req.Header.Add("Authorization", "Bearer "+accessToken)

			response = executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_app"].(string))
			assert.Equal(t, test.statusCode, response.Result().StatusCode)

			// check the previous image doesn't exists in directory
			switch test.name {
			case "delete 1":
				assert.False(t, fileExists("/app/static/icon-categories/"+category1.Icon))
			case "delete 2":
				assert.False(t, fileExists("/app/static/icon-categories/"+category2.Icon))
			}
		})
	}

}
