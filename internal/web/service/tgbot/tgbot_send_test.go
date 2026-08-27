package tgbot

import (
	"errors"
	"strings"
	"testing"
)

func TestIsTelegramNotModifiedError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"not modified", errors.New("Bad Request: message is not modified"), true},
		{"No fields to modify", errors.New("Bad Request: No fields to modify"), true},
		{"unrelated error", errors.New("Bad Request: message to edit not found"), false},
		{"network error", errors.New("connection reset"), false},
		{"empty string", errors.New(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTelegramNotModifiedError(tt.err)
			if got != tt.want {
				t.Errorf("isTelegramNotModifiedError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestPageMessageSplitsLinkListWithoutBlankLines(t *testing.T) {
	var message strings.Builder
	message.WriteString("Individual links:\r\n")
	for range 50 {
		message.WriteString("<code>vless://" + strings.Repeat("a", 300) + "</code>\r\n")
	}

	pages := pageMessage(message.String(), telegramPageLimit)
	if len(pages) < 2 {
		t.Fatalf("pageMessage() returned %d page, want multiple", len(pages))
	}

	links := 0
	for index, page := range pages {
		if len(page) > telegramPageLimit {
			t.Errorf("page %d has %d bytes, want at most %d", index, len(page), telegramPageLimit)
		}
		openingTags := strings.Count(page, "<code>")
		closingTags := strings.Count(page, "</code>")
		if openingTags != closingTags {
			t.Errorf("page %d has %d opening tags and %d closing tags", index, openingTags, closingTags)
		}
		links += openingTags
	}
	if links != 50 {
		t.Errorf("pages contain %d links, want 50", links)
	}
}
