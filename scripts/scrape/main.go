package main

import (
	"bookkeeper/internal/db"
	"bookkeeper/internal/model"
	"bookkeeper/internal/scraper"
	"context"
	"log"
	"os"
	"time"
)

// function to extract book data given a URL that contains a list of books (Goodreads)
func main() {
	pageURL := os.Args[1]

	logFile, err := os.OpenFile("data/extraction.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	infoLog := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errLog := log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	db.Connect()

	defer db.Client.Disconnect(context.TODO())

	if err != nil {
		errLog.Println(err)
		return
	}

	links, err := scraper.FetchBookLinks(pageURL)

	if err != nil {
		errLog.Printf("error fetching links from %s", pageURL)
		return
	} else {
		infoLog.Printf("fetched links from %s", pageURL)
	}

	for _, link := range links {

		book := scraper.FetchBookData(link, errLog)
		book.URL = link

		var doc model.Document

		doc.Source = pageURL
		doc.Work = book

		err = db.InsertBooks(doc)

		if err != nil {
			errLog.Printf("error pushing book %v, url: %s", err, link)
		} else {
			infoLog.Printf("successfully pushed %s to db", book.Title)
		}
		time.Sleep(5 * time.Second)

	}

	infoLog.Printf("extraction complete")

}
