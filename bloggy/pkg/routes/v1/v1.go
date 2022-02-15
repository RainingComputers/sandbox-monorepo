package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"examples/bloggy/pkg/models"
	"examples/bloggy/pkg/storage"
)

func CreateHandler(s storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newPost models.Post

		err := c.BindJSON(&newPost)

		if err != nil {
			c.String(http.StatusBadRequest, "%s", err.Error())
			return
		}

		err = s.Insert(c.Request.Context(), newPost)

		if err == storage.ErrAlreadyExists {
			c.String(http.StatusBadRequest, "%s", err.Error())
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Printf("ERROR Unable to insert into store: %s\n", err)
			return
		}
	}
}

func FindHandler(s storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")

		foundPost, err := s.Find(c.Request.Context(), title)

		if err == storage.ErrDoesNotExist {
			c.String(http.StatusBadRequest, "%s", err.Error())
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Printf("ERROR Unable to find in store: %s\n", err)
			return
		}

		h
	}
}

func RemoveHandler(s storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")

		err := s.Remove(c.Request.Context(), title)

		if err == storage.ErrDoesNotExist {
			c.String(http.StatusBadRequest, "%s", err.Error())
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Printf("ERROR Unable to delete from store: %s\n", err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func ModifyHandler(s storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")

		var newPost models.Post

		err := c.BindJSON(&newPost)

		if err != nil {
			c.String(http.StatusBadRequest, "%s", err.Error())
			return
		}

		err = s.Modify(c.Request.Context(), title, newPost)

		if err == storage.ErrDoesNotExist {
			c.String(http.StatusBadRequest, "%s", err.Error())
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Printf("ERROR Unable to modify in store: %s\n", err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func CreateRoutes(s storage.Storage, router *gin.Engine) {
	v1 := router.Group("/v1")

	v1.POST("/create", CreateHandler(s))
	v1.GET("/find/:title", FindHandler(s))
	v1.DELETE("/remove/:title", RemoveHandler(s))
	v1.PATCH("/modify/:title", ModifyHandler(s))
}
