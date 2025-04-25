package scraper

import (
	"fmt"
	"regexp"

	"github.com/go-rod/rod"
)

// fetch the first published date (or expected date of publication) for the book
func fetchDate(page *rod.Page) (string, error) {

	published, err := getTextFromSelector(page, `p[data-testid="publicationInfo"]`)

	if err != nil {
		return published, err
	}

	re := regexp.MustCompile(`([A-Za-z]+\s+\d{1,2}\,\s+\d{1,4})`)

	match := re.FindString(published)

	if match == "" {
		return published, fmt.Errorf("failed to parse date")
	}

	return match, nil
}
