package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"vyynl/gator/internal/database"

	"github.com/google/uuid"
)

/* RSS Handling */
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ERROR fetching RSS Feed: %v", err)
	}
	req.Header.Set("User-Agent", "gator") // Common practice to ID self to server

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ERROR client executing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ERROR reading resp body: %v", err)
	}

	var rss RSSFeed
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, fmt.Errorf("ERROR unmarshalling XML: %v", err)
	}

	return &rss, nil
}

func scrapeFeeds(s *state) error {
	nextFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("DATABASE ERROR fetching next feed to fetch: %v", err)
	}

	feed, err := s.db.MarkFeedFetched(
		context.Background(),
		nextFetch.ID,
	)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR marking feed fetched: %v", err)
	}

	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, item := range rss.Channel.Item {
		var nullDesc sql.NullString
		if len(item.Description) != 0 {
			nullDesc.String = item.Description
			nullDesc.Valid = true
		}

		var nullPub sql.NullTime
		if len(item.PubDate) != 0 {
			time, _ := time.Parse("2006-01-02 15:04", item.PubDate)
			nullPub.Time = time
			nullDesc.Valid = true
		}
		_, err := s.db.GetPostForURL(
			context.Background(),
			item.Link,
		)
		if err == nil {
			continue
		}

		post, err := s.db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Url:         item.Link,
				Description: nullDesc,
				PublishedAt: nullPub,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
			Error.Println(err)
		}
		fmt.Printf("Post: %s - FeedID: %s - Successfully scraped\n", post.Title, post.FeedID)
	}

	return nil
}
