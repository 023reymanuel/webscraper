package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ScrapeData holds the scraped information from a webpage.
type ScrapeData struct {
	Links  []string // URLs from <a> tags
	Texts  []string // Text from <p> tags
	Images []string // Src from <img> tags
}

// scrapePage fetches and scrapes a webpage, returning collected data.
func scrapePage(url string) (ScrapeData, error) {
	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return ScrapeData{}, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return ScrapeData{}, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	// Load HTML into goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ScrapeData{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	// Collect data
	data := ScrapeData{}

	// Extract links from <a> tags
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists && strings.HasPrefix(href, "http") {
			data.Links = append(data.Links, href)
		}
	})

	// Extract text from <p> tags
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			data.Texts = append(data.Texts, text)
		}
	})

	// Extract image sources from <img> tags
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			data.Images = append(data.Images, src)
		}
	})

	return data, nil
}

// saveToFile writes the scraped data to a file.
func saveToFile(data ScrapeData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintln(writer, "Scraped Links:")
	for i, link := range data.Links {
		fmt.Fprintf(writer, "%d. %s\n", i+1, link)
	}

	fmt.Fprintln(writer, "\nScraped Text (Paragraphs):")
	for i, text := range data.Texts {
		fmt.Fprintf(writer, "%d. %s\n", i+1, text)
	}

	fmt.Fprintln(writer, "\nScraped Images:")
	for i, src := range data.Images {
		fmt.Fprintf(writer, "%d. %s\n", i+1, src)
	}

	return writer.Flush()
}

func main() {
	// Parse URL flag
	url := flag.String("url", "", "URL to scrape (e.g., https://example.com)")
	flag.Parse()

	if *url == "" {
		log.Fatal("Please provide a URL using the -url flag")
	}

	// Scrape the page
	data, err := scrapePage(*url)
	if err != nil {
		log.Fatalf("Failed to scrape: %v", err)
	}

	// Print results
	fmt.Println("Scraped Links:")
	for i, link := range data.Links {
		fmt.Printf("%d. %s\n", i+1, link)
	}

	fmt.Println("\nScraped Text (Paragraphs):")
	for i, text := range data.Texts {
		fmt.Printf("%d. %s\n", i+1, text)
	}

	fmt.Println("\nScraped Images:")
	for i, src := range data.Images {
		fmt.Printf("%d. %s\n", i+1, src)
	}

	// Ask user if they want to save the data
	fmt.Print("\nWould you like to save the scraped data to a file? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" {
		filename := "output.txt"
		if err := saveToFile(data, filename); err != nil {
			log.Printf("Error saving to file: %v", err)
		} else {
			fmt.Printf("Data saved to %s\n", filename)
		}
	}
}
