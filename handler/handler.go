package handler

import (
	"net/http"
	"urlShortner/mapInterface"
	"urlShortner/shortner"

	"github.com/gin-gonic/gin"
)

const shortPath = "/link/"

// const host = "/short/"

const maxRetries = 3

type PostUrl struct {
	Url string `json:"url"`
}

func Hello(m mapInterface.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello!",
		})
	}
}

func Shorten(m mapInterface.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		var url PostUrl
		if err := c.ShouldBindJSON(&url); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		id, err := shortner.CreateShortUrl(url.Url, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = m.Add(id, url.Url)
			if err == nil {
				c.JSON(200, gin.H{
					"short_link": "localhost:8080" + shortPath + id,
				})

				// fmt.Println(mapS.Save())
				return
			}

			if attempt == maxRetries {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Out of free links",
				})
				return
			}
		}
	}
}

func Link(m mapInterface.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		url, err := m.Get(id)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Redirect(http.StatusSeeOther, url)
	}
}
