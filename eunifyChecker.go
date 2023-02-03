package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

var Env = godotenv.Load(`C:\Users\HP\Desktop\Code\EunifyAlert\.env`)
var Sender, Recipient, Password = os.Getenv("SENDER"), os.Getenv("RECIPIENT"), os.Getenv("APP_PASSWORD")

func sendEmail(url string) {

	auth := smtp.PlainAuth("", Sender, Password, "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it

	to := []string{Recipient}

	msg := []byte(fmt.Sprintf("To: %s\r\n", Recipient) +
		"Subject: EUnify Hoodie, mamy to!?\r\n" +
		"\r\n" +
		fmt.Sprintf("URL: %s\r\n", url))
	err := smtp.SendMail("smtp.gmail.com:587", auth, Sender, to, msg)

	if err != nil {
		log.Fatal(err)
	}
}

func available(url string) (bool, error) {

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false, err
	}

	if doc.Find("button[data-action='add-to-cart']").Length() > 0 {
		return true, nil
	}
	return false, nil
}

func main() {

	urls := []string{
		"https://www.koeniggalerie.com/products/eunify-hoodie?variant=43630082785512",
		"https://www.koeniggalerie.com/products/eunify-hoodie?variant=43630082851048",
		"https://www.koeniggalerie.com/products/eunify-hoodie?variant=43630082719976",
		"https://www.koeniggalerie.com/products/eunify-hoodie?variant=43630082654440"}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			isAvailable, err := available(url)
			if err != nil {
				log.Printf("Error checking availability for %s: %v\n", url, err)
				return
			}
			if isAvailable {
				sendEmail(url)
			}
		}(url)
	}
	wg.Wait()
}
