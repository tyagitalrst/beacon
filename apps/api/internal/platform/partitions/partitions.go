package partitions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func tableName(day time.Time) string {
	return "events_" + day.Format("2006_01_02")
}

// EnsureDaily creates the events partition for today and the next `days` days (UTC).
func EnsureDaily(ctx context.Context, db *pgxpool.Pool, days int) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	for i := 0; i <= days; i++ {
		from := today.AddDate(0, 0, i)
		to := from.AddDate(0, 0, 1)
		name := tableName(from)
		sql := fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s PARTITION OF events FOR VALUES FROM ('%s') TO ('%s')",
			pgx.Identifier{name}.Sanitize(), from.Format(time.RFC3339), to.Format(time.RFC3339),
		)
		if _, err := db.Exec(ctx, sql); err != nil {
			return fmt.Errorf("create partition %s: %w", name, err)
		}
	}
	return nil
}

// DropOlderThan drops events partitions whose whole day is older than retention.
// It returns the names of the dropped tables.
func DropOlderThan(ctx context.Context, db *pgxpool.Pool, retention time.Duration) ([]string, error) {
	rows, err := db.Query(ctx, `
                  SELECT c.relname
                  FROM pg_inherits i
                  JOIN pg_class c ON c.oid = i.inhrelid
                  WHERE i.inhparent = 'events'::regclass`)
	if err != nil {
		return nil, fmt.Errorf("list partitions: %w", err)
	}
	names, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("list partitions: %w", err)
	}

	cutoff := time.Now().UTC().Add(-retention)
	var dropped []string
	for _, name := range names {
		day, err := time.Parse("2006_01_02", strings.TrimPrefix(name, "events_"))
		if err != nil {
			continue // not one of our daily tables: never drop it
		}
		if day.AddDate(0, 0, 1).After(cutoff) {
			continue // still has rows inside the retention window
		}
		if _, err := db.Exec(ctx, "DROP TABLE "+pgx.Identifier{name}.Sanitize()); err != nil {
			return dropped, fmt.Errorf("drop partition %s: %w", name, err)
		}
		dropped = append(dropped, name)
	}
	return dropped, nil
}
