package scraper

import (
	"bookkeeper/internal/model"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
)

// fetch format of the book.
// format contains the number of pages in the book and the type of book
func fetchFormat(page *rod.Page) (model.Format, error) {
	format, err := getTextFromSelector(page, `p[data-testid="pagesFormat"]`)

	if err != nil {
		return model.Format{}, nil
	}

	re := regexp.MustCompile(`(\d+).*?,\s*(.+)`)
	matches := re.FindStringSubmatch(format)

	if len(matches) < 3 {
		return model.Format{}, fmt.Errorf("no matches found")
	}

	pageNo, err := strconv.Atoi(matches[1])
	if err != nil {
		return model.Format{}, fmt.Errorf("page no. not extracted")
	}

	bookType := strings.TrimSpace(matches[2])

	return model.Format{PageNo: pageNo, Type: bookType}, nil

}
