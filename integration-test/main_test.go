package integrationtest

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	pg "github.com/andreyxaxa/Sales-Tracker/pkg/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testPool *pgxpool.Pool
)

const (
	wrongID = int64(-999999)
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// container
	container, err := postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	migrateConnStr := strings.Replace(connStr, "postgres://", "pgx5://", 1)

	// migrations
	mg, err := migrate.New("file://../migrations", migrateConnStr)
	if err != nil {
		log.Fatalf("migrations failed: %s", err)
	}

	if err := mg.Up(); err != nil {
		log.Fatalf("failed to run migrate up: %s", err)
	}

	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to pgx: %s", err)
	}

	// run tests
	code := m.Run()

	testPool.Close()
	container.Terminate(ctx)

	os.Exit(code)
}

func newTestPostgres() *pg.Postgres {
	return &pg.Postgres{
		Pool:    testPool,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func newTestContext(t *testing.T) (context.Context, func()) {
	t.Helper()

	ctx := context.Background()

	tx, err := testPool.Begin(ctx)
	require.NoError(t, err)

	ctxWithKey := context.WithValue(ctx, pg.TxKey{}, tx)

	rollback := func() {
		tx.Rollback(ctx)
	}

	return ctxWithKey, rollback
}
