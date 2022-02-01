package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"examples/bloggy/pkg/models"
	"examples/bloggy/pkg/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getRouterAndStorage() (*gin.Engine, storage.Storage) {
	// Create Mongo store
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := storage.CreateMongoStore(ctx, "test", "test")
	if err != nil {
		panic(err)
	}

	store.Clean(ctx)

	// Create routes
	router := gin.Default()
	CreateRoutes(store, router)

	return router, store
}

func assertCount(t *testing.T, store storage.Storage, expectedCount int) {
	allPosts, err := store.All(context.Background())
	count := len(allPosts)

	if err != nil {
		t.Fatal(err)
	}

	if count != expectedCount {
		t.Errorf("Expected all item count of %d but got %d", expectedCount, count)
	}
}

func assertPost(t *testing.T, store storage.Storage, idx int, expectedPost models.Post) {
	allPosts, err := store.All(context.Background())

	if err != nil {
		t.Fatal(err)
		return
	}

	post := allPosts[0]

	if !post.IsEqual(expectedPost) {
		t.Errorf("Expected %s but got %s", expectedPost, post)
	}
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	if w.Code != expected {
		t.Errorf("Expected status to be %d but got %d", expected, w.Code)
	}
}

func assertBody(t *testing.T, w *httptest.ResponseRecorder, expectedPost models.Post) {
	var post models.Post

	err := json.Unmarshal(w.Body.Bytes(), &post)

	if err != nil {
		t.Errorf("Invalid json body: %s", err.Error())
	}

	post.ID = expectedPost.ID // Ignore ID field

	if post != expectedPost {
		t.Errorf("Expected %s but got %s", expectedPost, post)
	}
}

func assertBodyErrorMessage(t *testing.T, w *httptest.ResponseRecorder, message string) {
	if w.Body.String() != message {
		t.Errorf("Expected body message to be %s but got %s", message, w.Body.String())
	}
}

func TestCreate(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	testPost := models.Post{
		ID:      primitive.NilObjectID,
		Title:   "Hello",
		Name:    "Vishnu",
		Content: "Hello world",
	}

	testBody, _ := json.Marshal(testPost)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/create", bytes.NewBuffer(testBody))
	router.ServeHTTP(w, req)

	assertStatus(t, w, 200)

	assertCount(t, store, 1)

	assertPost(t, store, 0, testPost)
}

func TestCreateBadJSONBody(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	testBody := []byte("Invalid json")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/create", bytes.NewBuffer(testBody))
	router.ServeHTTP(w, req)

	assertStatus(t, w, 400)

	assertCount(t, store, 0)
}

func TestCreateAlreadyExists(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	testPost := models.Post{
		ID:      primitive.NilObjectID,
		Title:   "Hello",
		Name:    "Vishnu",
		Content: "Hello world",
	}

	testBody, _ := json.Marshal(testPost)

	req, _ := http.NewRequest("POST", "/v1/create", bytes.NewBuffer(testBody))
	router.ServeHTTP(httptest.NewRecorder(), req)

	w := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/v1/create", bytes.NewBuffer(testBody))
	router.ServeHTTP(w, req)

	assertStatus(t, w, 400)

	assertBodyErrorMessage(t, w, "storage: post already exists")

	assertCount(t, store, 1)

	assertPost(t, store, 0, testPost)
}

func TestFind(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	testPostOne := models.Post{
		ID:      primitive.NilObjectID,
		Title:   "Hello",
		Name:    "Vishnu",
		Content: "Hello world",
	}

	testPostTwo := models.Post{
		ID:      primitive.NilObjectID,
		Title:   "World",
		Name:    "Vishnu",
		Content: "This is golang",
	}

	testBodyOne, _ := json.Marshal(testPostOne)
	testBodyTwo, _ := json.Marshal(testPostTwo)

	req, _ := http.NewRequest("POST", "/v1/create", bytes.NewBuffer(testBodyOne))
	router.ServeHTTP(httptest.NewRecorder(), req)

	req, _ = http.NewRequest("POST", "/v1/create", bytes.NewBuffer(testBodyTwo))
	router.ServeHTTP(httptest.NewRecorder(), req)

	w := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/find/World", nil)
	router.ServeHTTP(w, req)

	assertStatus(t, w, 200)

	assertBody(t, w, testPostTwo)
}

func TestFindDoesNotExist(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/find/DoesNotExist", nil)
	router.ServeHTTP(w, req)

	assertStatus(t, w, 400)

	assertBodyErrorMessage(t, w, "storage: post does not exists")
}

func TestModify(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	testPost := models.Post{
		ID:      primitive.NilObjectID,
		Title:   "Hello",
		Name:    "Vishnu",
		Content: "Hello world",
	}

	testPostModified := models.Post{
		ID:      primitive.NilObjectID,
		Title:   "Hello",
		Name:    "SomeoneElse",
		Content: "Hello world, modified!",
	}

	testBody, _ := json.Marshal(testPost)
	testModifiedBody, _ := json.Marshal(testPostModified)

	req, _ := http.NewRequest("POST", "/v1/create", bytes.NewBuffer(testBody))
	router.ServeHTTP(httptest.NewRecorder(), req)

	w := httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/v1/modify/Hello", bytes.NewBuffer(testModifiedBody))
	router.ServeHTTP(w, req)

	assertStatus(t, w, 200)

	assertCount(t, store, 1)

	assertPost(t, store, 0, testPostModified)
}

func TestModifyBadJSON(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	testBody := []byte("Invalid json")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/modify/test", bytes.NewBuffer(testBody))
	router.ServeHTTP(w, req)

	assertStatus(t, w, 400)
}

func TestModifyDoesNotExist(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	testPost := models.Post{
		ID:      primitive.NilObjectID,
		Title:   "Hello",
		Name:    "Vishnu",
		Content: "Hello world",
	}

	testBody, _ := json.Marshal(testPost)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/modify/test", bytes.NewBuffer(testBody))
	router.ServeHTTP(w, req)

	assertStatus(t, w, 400)

	assertBodyErrorMessage(t, w, "storage: post does not exists")
}

func TestRemove(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	testPost := models.Post{
		ID:      primitive.NilObjectID,
		Title:   "Hello",
		Name:    "Vishnu",
		Content: "Hello world",
	}

	testBody, _ := json.Marshal(testPost)

	req, _ := http.NewRequest("POST", "/v1/create", bytes.NewBuffer(testBody))
	router.ServeHTTP(httptest.NewRecorder(), req)

	w := httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/v1/remove/Hello", nil)
	router.ServeHTTP(httptest.NewRecorder(), req)

	assertStatus(t, w, 200)

	assertCount(t, store, 0)
}

func TestRemoveDoesNotExist(t *testing.T) {
	router, store := getRouterAndStorage()
	defer store.Disconnect(context.Background())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/remove/Hello", nil)
	router.ServeHTTP(w, req)

	assertStatus(t, w, 400)

	assertBodyErrorMessage(t, w, "storage: post does not exists")
}
