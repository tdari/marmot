package connection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/marmotdata/marmot/internal/crypto"
	"github.com/rs/zerolog/log"
)

var (
	ErrNotFound      = errors.New("connection not found")
	ErrAlreadyExists = errors.New("connection with this name already exists")
	ErrInvalidInput  = errors.New("invalid connection input")
)

func isUniqueConstraintViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

type DeleteOptions struct {
	TeardownSchedules bool // Whether to delete schedules associated with this connection
}

type Repository interface {
	Create(ctx context.Context, conn *Connection) error
	Get(ctx context.Context, id string) (*Connection, error)
	GetByName(ctx context.Context, name string) (*Connection, error)
	List(ctx context.Context, opts *ListOptions) ([]*Connection, int, error)
	ListByType(ctx context.Context, connType string) ([]*Connection, error)
	Update(ctx context.Context, conn *Connection) error
	Delete(ctx context.Context, id string) error
	DeleteWithOptions(ctx context.Context, id string, opts *DeleteOptions) error
}

type ListOptions struct {
	Type      string
	CreatedBy string
	Search    string
	Limit     int
	Offset    int
}

type PostgresRepository struct {
	db        *pgxpool.Pool
	encryptor *crypto.Encryptor
}

func NewPostgresRepository(db *pgxpool.Pool, encryptor *crypto.Encryptor) *PostgresRepository {
	return &PostgresRepository{
		db:        db,
		encryptor: encryptor,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, conn *Connection) error {
	configCopy := make(map[string]interface{})
	for k, v := range conn.Config {
		configCopy[k] = v
	}

	if err := r.encryptConfig(conn.Type, configCopy); err != nil {
		return fmt.Errorf("encrypting config: %w", err)
	}

	configJSON, err := json.Marshal(configCopy)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	query := `
		INSERT INTO connections (id, name, type, description, config, tags, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.Exec(ctx, query,
		conn.ID,
		conn.Name,
		conn.Type,
		conn.Description,
		configJSON,
		conn.Tags,
		conn.CreatedBy,
		conn.CreatedAt,
		conn.UpdatedAt,
	)

	if err != nil {
		if isUniqueConstraintViolation(err) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("inserting connection: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*Connection, error) {
	query := `
		SELECT id, name, type, description, config, tags, created_by, created_at, updated_at
		FROM connections
		WHERE id = $1
	`

	var conn Connection
	var configJSON []byte
	var tags []string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&conn.ID,
		&conn.Name,
		&conn.Type,
		&conn.Description,
		&configJSON,
		&tags,
		&conn.CreatedBy,
		&conn.CreatedAt,
		&conn.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("querying connection: %w", err)
	}

	conn.Tags = tags

	var config map[string]interface{}
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	if err := r.decryptConfig(conn.Type, config); err != nil {
		return nil, fmt.Errorf("decrypting config: %w", err)
	}
	conn.Config = config

	return &conn, nil
}

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (*Connection, error) {
	query := `
		SELECT id, name, type, description, config, tags, created_by, created_at, updated_at
		FROM connections
		WHERE name = $1
	`

	var conn Connection
	var configJSON []byte
	var tags []string

	err := r.db.QueryRow(ctx, query, name).Scan(
		&conn.ID,
		&conn.Name,
		&conn.Type,
		&conn.Description,
		&configJSON,
		&tags,
		&conn.CreatedBy,
		&conn.CreatedAt,
		&conn.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("querying connection: %w", err)
	}

	conn.Tags = tags

	var config map[string]interface{}
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	if err := r.decryptConfig(conn.Type, config); err != nil {
		return nil, fmt.Errorf("decrypting config: %w", err)
	}
	conn.Config = config

	return &conn, nil
}

func (r *PostgresRepository) List(ctx context.Context, opts *ListOptions) ([]*Connection, int, error) {
	if opts == nil {
		opts = &ListOptions{}
	}

	query := `
		SELECT id, name, type, description, config, tags, created_by, created_at, updated_at
		FROM connections
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM connections WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if opts.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argPos)
		countQuery += fmt.Sprintf(" AND type = $%d", argPos)
		args = append(args, opts.Type)
		argPos++
	}

	if opts.CreatedBy != "" {
		query += fmt.Sprintf(" AND created_by = $%d", argPos)
		countQuery += fmt.Sprintf(" AND created_by = $%d", argPos)
		args = append(args, opts.CreatedBy)
		argPos++
	}

	if opts.Search != "" {
		pattern := "%" + opts.Search + "%"
		query += fmt.Sprintf(" AND (name ILIKE $%d OR type ILIKE $%d)", argPos, argPos)
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR type ILIKE $%d)", argPos, argPos)
		args = append(args, pattern)
		argPos++
	}

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting connections: %w", err)
	}

	query += " ORDER BY created_at DESC"

	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, opts.Limit)
		argPos++
	}

	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, opts.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying connections: %w", err)
	}
	defer rows.Close()

	var connections []*Connection
	for rows.Next() {
		var conn Connection
		var configJSON []byte
		var tags []string

		err := rows.Scan(
			&conn.ID,
			&conn.Name,
			&conn.Type,
			&conn.Description,
			&configJSON,
			&tags,
			&conn.CreatedBy,
			&conn.CreatedAt,
			&conn.UpdatedAt,
		)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to scan connection row")
			continue
		}

		conn.Tags = tags

		var config map[string]interface{}
		if err := json.Unmarshal(configJSON, &config); err != nil {
			log.Warn().Err(err).Str("connection_id", conn.ID).Msg("Failed to unmarshal config")
			continue
		}

		if err := r.decryptConfig(conn.Type, config); err != nil {
			log.Warn().Err(err).Str("connection_id", conn.ID).Msg("Failed to decrypt config")
			continue
		}
		conn.Config = config

		connections = append(connections, &conn)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating connection rows: %w", err)
	}

	return connections, total, nil
}

func (r *PostgresRepository) ListByType(ctx context.Context, connType string) ([]*Connection, error) {
	opts := &ListOptions{
		Type: connType,
	}
	connections, _, err := r.List(ctx, opts)
	return connections, err
}

func (r *PostgresRepository) Update(ctx context.Context, conn *Connection) error {
	// Create a copy of config to avoid modifying the original
	configCopy := make(map[string]interface{})
	for k, v := range conn.Config {
		configCopy[k] = v
	}

	if err := r.encryptConfig(conn.Type, configCopy); err != nil {
		return fmt.Errorf("encrypting config: %w", err)
	}

	configJSON, err := json.Marshal(configCopy)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	query := `
		UPDATE connections
		SET name = $2, description = $3, config = $4, tags = $5
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		conn.ID,
		conn.Name,
		conn.Description,
		configJSON,
		conn.Tags,
	)

	if err != nil {
		if isUniqueConstraintViolation(err) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("updating connection: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	return r.DeleteWithOptions(ctx, id, &DeleteOptions{TeardownSchedules: false})
}

func (r *PostgresRepository) DeleteWithOptions(ctx context.Context, id string, opts *DeleteOptions) error {
	if opts == nil {
		opts = &DeleteOptions{TeardownSchedules: false}
	}

	if opts.TeardownSchedules {
		deleteSchedulesQuery := `DELETE FROM ingestion_schedules WHERE connection_id = $1`

		_, err := r.db.Exec(ctx, deleteSchedulesQuery, id)
		if err != nil {
			return fmt.Errorf("deleting associated schedules: %w", err)
		}
	}

	query := `DELETE FROM connections WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting connection: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) encryptConfig(connType string, config map[string]interface{}) error {
	if r.encryptor == nil || len(config) == 0 {
		return nil
	}

	typeMeta, err := GetRegistry().GetMeta(connType)
	if err != nil {
		return fmt.Errorf("getting connection type metadata: %w", err)
	}

	sensitiveFields := GetSensitiveFields(typeMeta.ConfigSpec)

	return r.encryptor.EncryptMap(config, sensitiveFields)
}

func (r *PostgresRepository) decryptConfig(connType string, config map[string]interface{}) error {
	if r.encryptor == nil || len(config) == 0 {
		return nil
	}

	typeMeta, err := GetRegistry().GetMeta(connType)
	if err != nil {
		return fmt.Errorf("getting connection type metadata: %w", err)
	}

	sensitiveFields := GetSensitiveFields(typeMeta.ConfigSpec)

	return r.encryptor.DecryptMap(config, sensitiveFields)
}
