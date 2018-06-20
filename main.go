package main

import (
	"context"
	"golang.org/x/oauth2"
	"fmt"
	"log"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"google.golang.org/api/photoslibrary/v1"
	"net/http"
	"os"
	"encoding/json"
)

func main() {

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, photoslibrary.PhotoslibraryScope)

	client := getClient(config)
	srv, err := photoslibrary.New(client)

	var dates []*photoslibrary.Date
	date := &photoslibrary.Date{
		Day:   11,
		Month: 06,
		Year:  2018,}

	// create filter
	dates = append(dates, date)
	datefilter := &photoslibrary.DateFilter{Dates: dates}

	categoryList := []string{"Food"}
	contentFilter := &photoslibrary.ContentFilter{IncludedContentCategories: categoryList}

	filter := &photoslibrary.Filters{
		DateFilter: datefilter,
		ContentFilter: contentFilter}

	// create request
	req := &photoslibrary.SearchMediaItemsRequest{Filters: filter}
	caller := srv.MediaItems.Search(req)
	resp, err := caller.Do()

	if err != nil {
		log.Fatalf("Unable to call media items search: %v", err)
	}
	for key, value := range resp.MediaItems {
		log.Println(key, value);
	}

}

func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}
