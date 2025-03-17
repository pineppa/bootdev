package src

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"gator/internal/database"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	Guid        string `xml:"guid"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

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

func scrapeFeeds(s *CliState) error {
	nextFeed, err := s.DbQueries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving the next feed - %v", err)
	}
	if err := s.DbQueries.MarkFeedFetched(context.Background(), nextFeed.ID); err != nil {
		return fmt.Errorf("error while marking the feed as fetched - %v", err)
	}
	fmt.Printf("Scraping the next feed: %s @ %s\n", nextFeed.Name, nextFeed.Url)
	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error retrieving the RSS feed content - %v", err)
	}
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	var pubLink string
	for _, item := range rssFeed.Channel.Item {
		pubDate, err := time.Parse(layout, item.PubDate)
		if err != nil {
			return fmt.Errorf("failed to parse the time correctly - %v", err)
		}
		if item.Link != "" {
			pubLink = html.UnescapeString(item.Link)
		} else {
			pubLink = html.UnescapeString(item.Guid)
		}
		_, err = s.DbQueries.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: html.UnescapeString(item.Title), Valid: true},
			Url:         sql.NullString{String: pubLink, Valid: true},
			Description: sql.NullString{String: html.UnescapeString(item.Description), Valid: true},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: true},
			FeedID:      uuid.NullUUID{UUID: nextFeed.ID, Valid: true},
		})
		if pqErr, ok := err.(*pq.Error); err != nil && !(ok && pqErr.Code == "23505") {
			log.Printf("error in create post - %v", err)
		}
	}
	return nil
}

func printPost(post database.Post) {
	fmt.Println("Title: ", post.Title.String)
	fmt.Println("Publish time: ", post.PublishedAt.Time)
	fmt.Println("Post description: ", post.Description.String)
}

func printRssFeed(rssFeed *RSSFeed) {
	fmt.Println(" - Title:", html.UnescapeString(rssFeed.Channel.Title))
	fmt.Println(" - Description:", html.UnescapeString(rssFeed.Channel.Description))
}
