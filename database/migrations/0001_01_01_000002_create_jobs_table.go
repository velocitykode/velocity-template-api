package migrations

import "github.com/velocitykode/velocity/orm/migrate"

func init() {
	migrate.Register(&migrate.Migration{
		Version:     "20010101000002",
		Description: "create jobs table",
		Up: func(m *migrate.Migrator) error {
			// The database queue driver reads this table into queue.JobRecord,
			// where scheduled_at is a time.Time and reserved_at / failed_at are
			// *time.Time, so those three columns must be real timestamps - a
			// varchar column fails the row scan and no job is ever popped.
			// Timestamp rather than TimestampTz keeps one convention with the
			// created_at / updated_at pair below; the driver writes UTC and reads
			// UTC back, so no zone travels in the value. reserved_at and failed_at
			// stay nullable because they map to pointers and the driver's own
			// predicates test them with IS NULL.
			if err := m.CreateTable("jobs", func(t *migrate.TableBuilder) {
				t.ID()
				t.String("queue", 255)
				t.JSON("payload")
				t.Integer("attempts").Default("0")
				t.Timestamp("scheduled_at")
				t.Timestamp("reserved_at").Nullable()
				t.String("reserved_by", 255).Nullable()
				t.Timestamp("failed_at").Nullable()
				t.String("failed_reason", 5000).Nullable()
				t.Timestamps()
			}); err != nil {
				return err
			}

			// Every worker poll reserves the next job with `WHERE queue = ?`
			// ordered by `scheduled_at ASC, id ASC LIMIT 1`, on the table that
			// churns hardest in the schema. Leading with the equality column and
			// carrying both sort columns lets that query walk the index and stop
			// at the first candidate instead of scanning and sorting. The same
			// index serves the per-queue size and delayed-job counts.
			//
			// Not partial on unreserved rows: the pop predicate also matches rows
			// whose lease expired (reserved_at IS NOT NULL), so a partial index
			// could not answer it and would only add write cost.
			if err := m.Index("jobs", "queue", "scheduled_at", "id"); err != nil {
				return err
			}

			return m.CreateTable("failed_jobs", func(t *migrate.TableBuilder) {
				t.ID()
				t.String("queue", 255)
				t.JSON("payload")
				t.String("exception", 10000)
				t.Timestamps()
			})
		},
		Down: func(m *migrate.Migrator) error {
			if err := m.DropTable("failed_jobs"); err != nil {
				return err
			}
			return m.DropTable("jobs")
		},
	})
}
