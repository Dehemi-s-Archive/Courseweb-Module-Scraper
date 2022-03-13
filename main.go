package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	site := flag.String("site", "", "Set the Courseweb site")
	username := flag.String("u", "", "Set the Username")
	password := flag.String("p", "", "Set the Username")
	flag.Parse()

	if len(*site) < 1 {
		log.Println("Enter a site")
		return
	}

	// Initialize the http client
	jar, err := cookiejar.New(nil)

	if err != nil {
		log.Panicln(err)
	}

	client := http.Client{
		Jar: jar,
	}
	// Login to the site
	_, pErr := client.PostForm(
		"https://courseweb.sliit.lk/login/index.php",
		url.Values{
			"username": {*username},
			"password": {*password},
		},
	)

	if pErr != nil {
		log.Panicln(pErr)
	}

	// Fetch the page
	page, err := GetGoQuery(&client, *site)

	if err != nil {
		log.Panicln("Could not fetch the page")
	}

	// for each week
	page.Find(".section.main").Each(func(i int, s *goquery.Selection) {
		// list all the links
		fmt.Printf("\n--- %s ---\n", s.Find(".sectionname").First().Text())

		s.Find(".activityinstance").Each(func(j int, activity *goquery.Selection) {
			link := activity.Find("a").AttrOr("href", "Empty Link")
			mime, err := GetMime(&client, link)

			if err != nil {
				mime = "Mime err"
			}

			fmt.Printf("%s - %s - mime : %s\n", activity.Text(), link, mime)
		})

	})
}

func GetGoQuery(client *http.Client, site string) (*goquery.Document, error) {
	resp, err := client.Get(site)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return goquery.NewDocumentFromReader(resp.Body)
}

func GetMime(client *http.Client, site string) (string, error) {
	resp, err := client.Get(site)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return resp.Header.Get("Content-Type") + " " + resp.Request.Host, nil
}
