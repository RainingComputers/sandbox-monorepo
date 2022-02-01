package storage

import (
	"context"
	"examples/bloggy/pkg/models"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getMongoStore() Storage {
	// Create Mongo store
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := CreateMongoStore(ctx, "test", "test")
	if err != nil {
		panic(err)
	}

	store.Clean(ctx)

	return store
}

func assertPostsSlice(t *testing.T, store Storage, expectedPosts []models.Post) {
	posts, err := store.All(context.Background())

	if err != nil {
		t.Error(err)
		return
	}

	if len(posts) != len(expectedPosts) {
		t.Errorf("Expected item count to be %d but got %d", len(expectedPosts), len(posts))
		return
	}

	for i := range expectedPosts {
		if !expectedPosts[i].IsEqual(posts[i]) {
			t.Errorf("Expected index %d post to be %s but got %s", i, posts[i], expectedPosts[i])
		}
	}
}

func assertNil(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assertError(t *testing.T, err error, expected error) {
	if err != expected {
		t.Errorf("Expected error to be %s but got %s", expected.Error(), err.Error())
	}
}

func TestCount(t *testing.T) {
	mongoStore := getMongoStore()
	defer mongoStore.Disconnect(context.Background())

	memStore := CreateMemoryStore()

	for _, store := range []Storage{mongoStore, memStore} {
		testPostOne := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Hello",
			Name:    "Vishnu",
			Content: "Hello world",
		}

		testPostTwo := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "World",
			Name:    "Shankar",
			Content: "This is golang!!",
		}

		testPostThree := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Golang",
			Name:    "Bob",
			Content: "Golang is awesome!!",
		}

		store.Insert(context.Background(), testPostOne)
		store.Insert(context.Background(), testPostTwo)
		store.Insert(context.Background(), testPostThree)

		expectedPosts := []models.Post{testPostOne, testPostTwo, testPostThree}

		assertPostsSlice(t, store, expectedPosts)
	}
}

func TestCreate(t *testing.T) {
	mongoStore := getMongoStore()
	defer mongoStore.Disconnect(context.Background())

	memStore := CreateMemoryStore()

	for _, store := range []Storage{mongoStore, memStore} {

		testPost := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Hello",
			Name:    "Vishnu",
			Content: "Hello world",
		}

		err := store.Insert(context.Background(), testPost)

		assertNil(t, err)

		expectedPosts := []models.Post{testPost}

		assertPostsSlice(t, store, expectedPosts)
	}
}

func TestCreateAlreadyExists(t *testing.T) {
	mongoStore := getMongoStore()
	defer mongoStore.Disconnect(context.Background())

	memStore := CreateMemoryStore()

	for _, store := range []Storage{mongoStore, memStore} {
		testPost := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Hello",
			Name:    "Vishnu",
			Content: "Hello world",
		}

		err := store.Insert(context.Background(), testPost)

		assertNil(t, err)

		err = store.Insert(context.Background(), testPost)

		assertError(t, err, ErrAlreadyExists)

		expectedPosts := []models.Post{testPost}

		assertPostsSlice(t, store, expectedPosts)
	}
}

func TestFind(t *testing.T) {
	mongoStore := getMongoStore()
	defer mongoStore.Disconnect(context.Background())

	memStore := CreateMemoryStore()

	for _, store := range []Storage{mongoStore, memStore} {
		testPostOne := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Hello",
			Name:    "Vishnu",
			Content: "Hello world",
		}

		testPostTwo := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "World",
			Name:    "Shankar",
			Content: "This is golang!!",
		}

		testPostThree := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Golang",
			Name:    "Bob",
			Content: "Golang is awesome!!",
		}

		store.Insert(context.Background(), testPostOne)
		store.Insert(context.Background(), testPostTwo)
		store.Insert(context.Background(), testPostThree)

		post, err := store.Find(context.Background(), "World")

		assertNil(t, err)

		if !post.IsEqual(testPostTwo) {
			t.Errorf("Expected %s but got %s", testPostTwo, post)
		}
	}
}

func TestFindDoesNotExist(t *testing.T) {
	mongoStore := getMongoStore()
	defer mongoStore.Disconnect(context.Background())

	memStore := CreateMemoryStore()

	for _, store := range []Storage{mongoStore, memStore} {

		_, err := store.Find(context.Background(), "Hello")

		assertError(t, err, ErrDoesNotExist)
	}
}

func TestModify(t *testing.T) {
	mongoStore := getMongoStore()
	defer mongoStore.Disconnect(context.Background())

	memStore := CreateMemoryStore()

	for _, store := range []Storage{mongoStore, memStore} {

		testPost := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Hello",
			Name:    "Vishnu",
			Content: "Hello world",
		}

		testPostModify := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Hello",
			Name:    "Vishnu",
			Content: "Golang is awesome!!",
		}

		err := store.Insert(context.Background(), testPost)

		assertNil(t, err)

		err = store.Modify(context.Background(), "Hello", testPostModify)

		assertNil(t, err)

		expectedPosts := []models.Post{testPostModify}

		assertPostsSlice(t, store, expectedPosts)
	}
}

func TestModifyDoesNotExist(t *testing.T) {
	mongoStore := getMongoStore()
	defer mongoStore.Disconnect(context.Background())

	memStore := CreateMemoryStore()

	for _, store := range []Storage{mongoStore, memStore} {
		testPost := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Hello",
			Name:    "Vishnu",
			Content: "Hello world",
		}

		testPostModify := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Hello",
			Name:    "Vishnu",
			Content: "Golang is awesome!!",
		}

		err := store.Insert(context.Background(), testPost)

		assertNil(t, err)

		err = store.Modify(context.Background(), "DoesNotExist", testPostModify)

		assertError(t, err, ErrDoesNotExist)

		expectedPosts := []models.Post{testPost}

		assertPostsSlice(t, store, expectedPosts)
	}
}

func TestRemove(t *testing.T) {
	mongoStore := getMongoStore()
	defer mongoStore.Disconnect(context.Background())

	memStore := CreateMemoryStore()

	for _, store := range []Storage{mongoStore, memStore} {
		testPostOne := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Hello",
			Name:    "Vishnu",
			Content: "Hello world",
		}

		testPostTwo := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "World",
			Name:    "Shankar",
			Content: "This is golang!!",
		}

		testPostThree := models.Post{
			ID:      primitive.NilObjectID,
			Title:   "Golang",
			Name:    "Bob",
			Content: "Golang is awesome!!",
		}

		store.Insert(context.Background(), testPostOne)
		store.Insert(context.Background(), testPostTwo)
		store.Insert(context.Background(), testPostThree)

		err := store.Remove(context.Background(), "World")

		assertNil(t, err)

		expectedPosts := []models.Post{testPostOne, testPostThree}

		assertPostsSlice(t, store, expectedPosts)
	}
}

func TestRemoveDoesNotExist(t *testing.T) {
	mongoStore := getMongoStore()
	defer mongoStore.Disconnect(context.Background())

	memStore := CreateMemoryStore()
	for _, store := range []Storage{mongoStore, memStore} {
		err := store.Remove(context.Background(), "Hello")

		assertError(t, err, ErrDoesNotExist)
	}
}
