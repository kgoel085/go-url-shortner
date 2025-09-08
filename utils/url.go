package utils

import (
	"fmt"

	"example.com/url-shortner/config"
	"github.com/gosimple/slug"
	"github.com/segmentio/ksuid"
)

func GetShortUrl(code string) string {
	codeUrl := fmt.Sprintf("http://%s:%s/%s", config.Config.APP.Host, config.Config.APP.Port, code)
	return codeUrl
}

func GenerateSlug(input string, maxLength int) string {
	if (maxLength <= 0) || (maxLength > 20) {
		maxLength = 20 // Default length
	}
	// If input is provided, slugify it
	if input != "" {
		s := slug.Make(input)
		if len(s) > maxLength {
			s = s[:maxLength]
		}
		return s
	}
	// Generate KSUID and truncate
	randomSlug := ksuid.New().String()
	if len(randomSlug) > maxLength {
		randomSlug = randomSlug[:maxLength]
	}
	return randomSlug
}
