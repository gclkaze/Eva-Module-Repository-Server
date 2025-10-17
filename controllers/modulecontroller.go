package controllers

import (
	"net/http"

	"github.com/gclkaze/evamodulerepositoryserver/models"

	"github.com/gin-gonic/gin"
)

var albums = []models.Module{
	{ID: "1", Title: "The Dark Side of the Moon", Artist: "Pink Floyd", Price: 9.99},
	{ID: "2", Title: "Back in Black", Artist: "AC/DC", Price: 8.99},
}

func SearchModules(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func CreateAlbum(c *gin.Context) {
	var newAlbum models.Module
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
