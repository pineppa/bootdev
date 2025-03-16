package src

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

const rssUrl = "https://www.wagslane.dev/index.xml"

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var rssFeed = RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &rssFeed, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return &rssFeed, err
	}
	defer resp.Body.Close()

	// Check if HTTP status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &rssFeed, err
	}
	if err := xml.Unmarshal(body, &rssFeed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}
	return &rssFeed, nil
}

// func printRssItem(rssItem RSSItem) {
// 	fmt.Println(" - ", html.UnescapeString(rssItem.Title))
// 	fmt.Println(" - ", html.UnescapeString(rssItem.PubDate))
// 	fmt.Println(" - ", html.UnescapeString(rssItem.Description))
// 	fmt.Println(" - ", html.UnescapeString(rssItem.Link))
// }

func printRssFeed(rssFeed *RSSFeed) {
	fmt.Println(" - ", html.UnescapeString(rssFeed.Channel.Title))
	fmt.Println(" - ", html.UnescapeString(rssFeed.Channel.Description))
	fmt.Println(" - ", html.UnescapeString(rssFeed.Channel.Link))
	// for _, item := range rssFeed.Channel.Item {
	// 	printRssItem(item)
	// }
}
