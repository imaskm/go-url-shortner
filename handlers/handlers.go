package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/imaskm/url-shortner/caching"
	"github.com/imaskm/url-shortner/database"
	"github.com/imaskm/url-shortner/utility"
)

type Shorter struct {
	Db    *database.MongoDB
	Cache *caching.Cache
}

type longURL struct {
	Url       string `json:"url"`
	CustomUrl string `json:"custom_url"`
}

func (s Shorter) ShortUrl(gctx *gin.Context) {

	var body longURL

	if err := gctx.BindJSON(&body); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, urlErr := url.ParseRequestURI(body.Url)
	if urlErr != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": urlErr.Error()})
		return
	}

	if body.CustomUrl != "" {
		if len(body.CustomUrl) != 7 {
			gctx.JSON(http.StatusBadRequest, gin.H{"error": "custom url should be of 7 character"})
			return
		}

		v, _ := s.Db.GetLongURLForShortURL(body.CustomUrl)

		if v != "" {
			gctx.JSON(http.StatusBadRequest, gin.H{"error": "custom url already exist, provide new one"})
			return
		}

		err := s.Db.SaveShortURL(body.Url, body.CustomUrl)
		if err != nil {

			gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		} else {

			gctx.JSON(http.StatusCreated, body.CustomUrl)
			return

		}
	}

	// short, err := s.Db.GetShortURLForLongURL(body.Url)
	// if err != nil {
	// 	gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// if short != "" {
	// 	gctx.JSON(http.StatusCreated, short)
	// 	return
	// }

	for {
		shortUrl, err := utility.GetRandomBase58StringOfLength(7)
		if err != nil {
			gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// check if shortURL already exists
		value, err := s.Db.GetLongURLForShortURL(shortUrl)
		if err != nil {
			err := s.Db.SaveShortURL(body.Url, shortUrl)
			if err != nil {
				gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				gctx.JSON(http.StatusCreated, shortUrl)
			}
			return
		} else if value != "" {
			// creating a new short url if there is any collision
			continue

		} else {
			err := s.Db.SaveShortURL(body.Url, shortUrl)
			if err != nil {
				gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				gctx.JSON(http.StatusCreated, shortUrl)
			}
			return
		}
	}

}

func (s Shorter) Redirect(gctx *gin.Context) {
	shortURL, ok := gctx.Params.Get("shortURL")
	if !ok {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	// skipping error if issue or not found in cache
	value, _ := s.Cache.Read(shortURL)
	if value != "" {
		gctx.Redirect(http.StatusPermanentRedirect, value)
		return
	}

	longURL, err := s.Db.GetLongURLForShortURL(shortURL)
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": "short url doesn't exist"})
		return
	}

	go gctx.Redirect(http.StatusTemporaryRedirect, longURL)

	s.Cache.Write(shortURL, longURL)

}
