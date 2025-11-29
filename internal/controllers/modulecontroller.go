// Package controllers contains the controllers of the Repository Server
package controllers

import (
	"net/http"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"

	"github.com/gin-gonic/gin"
)

var modules = []models.Module{
	{Title: "The Dark Side of the Moon", Repr: "Pink Floyd"},
	{Title: "Back in Black", Repr: "AC/DC"},
}

func SearchModules(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, modules)
}

func CreateAlbum(c *gin.Context) {
	var newAlbum models.Module
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	modules = append(modules, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
