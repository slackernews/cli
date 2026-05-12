package formatters

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/slackernews/cli/pkg/api"
	"golang.org/x/term"
)

type jsonLink struct {
	ID            string  `json:"id"`
	URL           string  `json:"url"`
	Title         string  `json:"title"`
	Score         float64 `json:"score"`
	ReplyCount    int     `json:"replyCount"`
	IsUpvoted     bool    `json:"isUpvoted"`
	Age           string  `json:"age"`
	FirstSharedBy string  `json:"firstSharedBy"`
}

// FormatTable renders links as a human-readable terminal table.
func FormatTable(links []api.RenderableLink) {
	if len(links) == 0 {
		fmt.Println("No links found")
		return
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 {
		width = 80
	}

	rows := make([][]string, len(links))
	for i, link := range links {
		title := link.Link.Title
		if title == "" {
			title = link.Link.URL
		}

		// Reserve space for other columns
		maxTitle := width - 30
		if maxTitle < 20 {
			maxTitle = 20
		}
		if len(title) > maxTitle {
			title = title[:maxTitle-3] + "..."
		}

		age := formatAge(link.FirstShare.SharedAt)
		rows[i] = []string{
			fmt.Sprintf("%d", i+1),
			title,
			fmt.Sprintf("%.0f", link.DisplayScore),
			age,
			fmt.Sprintf("%d", link.ReplyCount),
		}
	}

	t := table.New().
		Headers("Rank", "Title", "Score", "Age", "Replies").
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
			}
			return lipgloss.NewStyle()
		}).
		Width(width)

	fmt.Println(t.Render())
}

// FormatJSON renders links as a JSON array.
func FormatJSON(links []api.RenderableLink) ([]byte, error) {
	out := make([]jsonLink, len(links))
	for i, link := range links {
		title := link.Link.Title
		if title == "" {
			title = link.Link.URL
		}
		out[i] = jsonLink{
			ID:            link.Link.URL,
			URL:           link.Link.URL,
			Title:         title,
			Score:         link.DisplayScore,
			ReplyCount:    link.ReplyCount,
			IsUpvoted:     link.IsUpvoted,
			Age:           formatAge(link.FirstShare.SharedAt),
			FirstSharedBy: link.FirstShare.FullName,
		}
	}
	return json.MarshalIndent(out, "", "  ")
}

func formatAge(timestamp int64) string {
	t := time.UnixMilli(timestamp)
	d := time.Since(t)

	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	if d < 30*24*time.Hour {
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
	return fmt.Sprintf("%dmo", int(d.Hours()/24/30))
}
