package main

import (
	"github.com/gin-gonic/gin"
	"github.com/imaskm/url-shortner/caching"
	"github.com/imaskm/url-shortner/database"
	"github.com/imaskm/url-shortner/handlers"
)

func main() {
	r := gin.Default()

	shorter := handlers.Shorter{
		Db:    database.NewDatabase(),
		Cache: caching.NewCache(),
	}

	r.POST("/short", shorter.ShortUrl)

	r.GET(":shortURL", shorter.Redirect)

	r.Run(":8989")

}
