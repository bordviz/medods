package query

import (
	"log/slog"
	"strings"
)

func QueryToString(q string) slog.Attr {
	query := strings.TrimSpace(
		strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " "),
	)
	return slog.Attr{
		Key:   "query",
		Value: slog.StringValue(query),
	}
}
