package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouterCRUD(t *testing.T) {
	previousUsers := mockUsers
	mockUsers = newMockUserStore()
	t.Cleanup(func() { mockUsers = previousUsers })
	router := NewRouter()

	response := call(t, router, http.MethodGet, "/users", "")
	if response.Code != 200 {
		t.Fatalf("select status = %d, want 200", response.Code)
	}
	var result struct {
		Users []User `json:"users"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil || len(result.Users) != 2 {
		t.Fatalf("select body = %s, want 2 users", response.Body)
	}

	response = call(t, router, http.MethodPost, "/users", `{"name":"Carol","email":"carol@example.com"}`)
	if response.Code != 201 {
		t.Fatalf("create status = %d, want 201", response.Code)
	}
	var created User
	if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	response = call(t, router, http.MethodPut, "/users/"+created.ID, `{"name":"Carol Updated","email":"carol.updated@example.com"}`)
	if response.Code != 200 || !strings.Contains(response.Body.String(), "Carol Updated") {
		t.Fatalf("update response = %#v", response)
	}

	response = call(t, router, http.MethodGet, "/users/"+created.ID, "")
	if response.Code != 200 || !strings.Contains(response.Body.String(), "carol.updated@example.com") {
		t.Fatalf("select one response = %#v", response)
	}

	response = call(t, router, http.MethodDelete, "/users/"+created.ID, "")
	if response.Code != 200 {
		t.Fatalf("delete status = %d, want 200", response.Code)
	}

	response = call(t, router, http.MethodGet, "/users/"+created.ID, "")
	if response.Code != 404 {
		t.Fatalf("deleted user status = %d, want 404", response.Code)
	}
}

func TestUploadFile(t *testing.T) {
	router := NewRouter()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/files", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("upload status = %d, want 200", recorder.Code)
	}
	var result struct {
		Filename    string `json:"filename"`
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Filename != "hello.txt" || result.Size != 5 || result.ContentType != "application/octet-stream" {
		t.Fatalf("upload response = %#v", result)
	}

	response := call(t, router, http.MethodPost, "/files", "")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("missing file status = %d, want 400", response.Code)
	}
}

func call(t *testing.T, router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}
