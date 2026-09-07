package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type Event struct {
	EventTime   time.Time
	Link        string
	Description string
	ItemName    string
	Default     bool
}

func ParseLink(link string, defaults bool) (Event, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(link)
	if err != nil {
		return Event{}, fmt.Errorf("не удалось загрузить RSS для ID: %s", err)
	}

	if len(feed.Items) == 0 {
		return Event{}, fmt.Errorf("RSS не содержит элементов для ID %s", link)
	}
	event := parseDescription(feed.Items[0].Description, defaults)
	return event, err
}

func parseDescription(html string, defaults bool) Event {
	var e Event
	// 1. Дата и время события
	e.EventTime = parseEventTime(html)
	// 2. Описание события
	e.Description = parseEventDescription(html)
	// 3. Наименование товара/услуг
	e.ItemName = parseItemName(html)
	// 4. Default
	e.Default = defaults
	return e
}

func parseEventTime(html string) time.Time {
	re := regexp.MustCompile(`Дата и время события:\s*</strong>\s*([^<]+)`)
	m := re.FindStringSubmatch(html)
	if len(m) < 2 {
		return time.Time{}
	}
	t, err := time.Parse("02.01.2006 15:04", strings.TrimSpace(m[1]))
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseEventDescription(html string) string {
	re := regexp.MustCompile(`Описание события:\s*</strong>\s*(.*?)<br/>`)
	m := re.FindStringSubmatch(html)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(stripTags(m[1]))
}

func parseItemName(html string) string {
	re := regexp.MustCompile(`<tbody>.*?<td>(.*?)</td>`)
	m := re.FindStringSubmatch(html)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(stripTags(m[1]))
}

func stripTags(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}

func ParseAllContracts(url string) ([]string, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить RSS: %w", err)
	}

	contractRe := regexp.MustCompile(`№\s*(\d+)`)

	var contracts []string

	for _, item := range feed.Items {
		match := contractRe.FindStringSubmatch(item.Title)
		if len(match) < 2 {
			continue
		}

		contracts = append(contracts, match[1])
	}

	return contracts, nil
}
