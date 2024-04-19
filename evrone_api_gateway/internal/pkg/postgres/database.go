package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	pgxadapter "github.com/pckhoi/casbin-pgx-adapter/v2"

	errorspkg "evrone_api_gateway/internal/errors"
	"evrone_api_gateway/internal/pkg/config"
)

// PostgresDB ...
type PostgresDB struct {
	*pgxpool.Pool
	Sq *Squirrel
}

// New provides PostgresDB struct init
func New(config *config.Config) (*PostgresDB, error) {

	db := PostgresDB{Sq: NewSquirrel()}

	if err := db.connect(config); err != nil {
		return nil, err
	}

	return &db, nil
}

func GetStrConfig(config *config.Config) string {
	var conn []string

	if len(config.DB.Host) != 0 {
		conn = append(conn, "host="+config.DB.Host)
	}

	if len(config.DB.Port) != 0 {
		conn = append(conn, "port="+config.DB.Port)
	}

	if len(config.DB.User) != 0 {
		conn = append(conn, "user="+config.DB.User)
	}

	if len(config.DB.Password) != 0 {
		conn = append(conn, "password="+config.DB.Password)
	}

	if len(config.DB.Name) != 0 {
		conn = append(conn, "dbname="+config.DB.Name)
	}

	if len(config.DB.SSLMode) != 0 {
		conn = append(conn, "sslmode="+config.DB.SSLMode)
	}

	return strings.Join(conn, " ")
}

func GetPgxPoolConfig(config *config.Config) (*pgx.ConnConfig, error) {
	return pgx.ParseConfig(GetStrConfig(config))
}

func GetAdapter(config *config.Config) (*pgxadapter.Adapter, error) {
	pgxPoolConfig, err := GetPgxPoolConfig(config)
	if err != nil {
		return nil, err
	}
	return pgxadapter.NewAdapter(pgxPoolConfig, pgxadapter.WithDatabase(config.DB.Name))
}

func (p *PostgresDB) connect(config *config.Config) error {
	pgxpoolConfig, err := pgxpool.ParseConfig(GetStrConfig(config))
	if err != nil {
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), pgxpoolConfig)
	if err != nil {
		return fmt.Errorf("unable to connect database config: %w", err)
	}

	p.Pool = pool

	return nil
}

func (p *PostgresDB) Close() {
	p.Pool.Close()
}

func (p *PostgresDB) Error(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return errorspkg.ErrorConflict
		}
	}
	if err == pgx.ErrNoRows {
		return errorspkg.ErrorNotFound
	}
	return err
}

func (p *PostgresDB) ErrSQLBuild(err error, message string) error {
	return fmt.Errorf("error during sql build, %s: %w", message, err)
}
