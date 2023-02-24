package main

import (
	"net/http"
	"web-service-gin/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	docs.SwaggerInfo.Title = "Vintage Recordings on Vinyl Store API"

	router := gin.Default()

	v1 := router.Group("api/v1")
	{
		albums := v1.Group("/albums")
		{
			albums.GET("", getAlbums)
			albums.GET(":id", getAlbumByID)
			albums.POST("", postAlbum)
			albums.DELETE(":id", deleteAlbum)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8080")
}

// getAlbums godoc
//	@Summary		List albums
//	@Description	List all albums
//	@Tags			albums
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	album
//	@Router			/api/v1/albums [get]
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// getAlbumByID godoc
//	@Summary		Show album
//	@Description	Get album by ID
//	@Tags			albums
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		200	{object}	album
//	@Failure		404	{object}	HTTPError
//	@Router			/api/v1/albums/{id} [get]
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, createAlbumNotFoundMessage())
}

// deleteAlbum godoc
//	@Summary		Delete album
//	@Description	Delete album by ID
//	@Tags			albums
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Album ID"
//	@Success		204	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Router			/api/v1/albums/{id} [delete]
func deleteAlbum(c *gin.Context) {
	id := c.Param("id")
	for i, a := range albums {
		if a.ID == id {
			albums = remove(albums, i)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, createAlbumNotFoundMessage())
}

// postAlbum godoc
//	@Summary		Add new album
//	@Description	Add new album
//	@Tags			albums
//	@Accept			json
//	@Produce		json
//	@Param			album	body		album	true	"Album to add"
//	@Success		201		{object}	album
//	@Router			/api/v1/albums [post]
func postAlbum(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func remove[T any](slice []T, idx int) []T {
	return append(slice[:idx], slice[idx+1:]...)
}

func createAlbumNotFoundMessage() map[string]any {
	return createMessage("album not found")
}

func createMessage(message string) map[string]any {
	return gin.H{"message": message}
}
