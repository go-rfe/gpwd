package cloud

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib" // init postgresql driver

	"github.com/go-rfe/gpwd/internal/logging/log"
	pb "github.com/go-rfe/gpwd/internal/proto"
)

const (
	psqlDriverName           = "pgx"
	pgErrCodeUniqueViolation = "23505"
)

var (
	_ Accounts = (*DB)(nil)
)

type DB struct {
	conn *sql.DB
}

func NewDB(databaseDSN string) (*DB, error) {
	var db DB
	conn, err := sql.Open(psqlDriverName, databaseDSN)
	if err != nil {
		return nil, err
	}
	db = DB{
		conn: conn,
	}

	return &db, nil
}

func (db *DB) CreateAccount(ctx context.Context, auth *pb.Auth) error {
	var pgErr *pgconn.PgError

	_, err := db.conn.ExecContext(ctx, `
		INSERT INTO accounts (username, password) VALUES ($1, $2);`,
		auth.GetUsername(), auth.GetPassword(),
	)

	if err != nil && errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
		return ErrAccountExists
	}

	return err
}

func (db *DB) GetByName(ctx context.Context, username string) (*pb.Auth, error) {
	var userPassword []byte
	row := db.conn.QueryRowContext(ctx,
		"SELECT password FROM accounts WHERE username = $1", username)

	err := row.Scan(&userPassword)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	}

	return &pb.Auth{
		Username: username,
		Password: userPassword,
	}, nil
}

func (db *DB) CreateSecrets(ctx context.Context, auth *pb.Auth, secrets []*pb.Secret) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollbackTx(tx)

	stmtCreateSecret, err := tx.Prepare(
		`INSERT INTO secrets
    		  (id, username, labels, created_at, data) 
			  VALUES ($1, $2, $3, $4, $5)`,
	)
	if err != nil {
		return err
	}
	defer closeObject(stmtCreateSecret)

	for _, secret := range secrets {
		var metadata []byte
		if secret.Labels != nil {
			metadata, err = json.Marshal(secret.GetLabels())
			if err != nil {
				return err
			}
		}

		if _, err := stmtCreateSecret.Exec(
			secret.GetID(), auth.GetUsername(),
			metadata, secret.GetCreatedAt().AsTime(),
			secret.GetData(),
		); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateSecrets(ctx context.Context, auth *pb.Auth, secrets []*pb.Secret) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollbackTx(tx)

	stmtUpdateSecret, err := tx.Prepare(
		`UPDATE secrets set
    		  labels=$3, updated_at=$4, data=$5 
			  WHERE id=$1 AND username=$2`,
	)
	if err != nil {
		return err
	}
	defer closeObject(stmtUpdateSecret)

	for _, secret := range secrets {
		var metadata []byte
		if secret.Labels != nil {
			metadata, err = json.Marshal(secret.GetLabels())
			if err != nil {
				return err
			}
		}

		if _, err := stmtUpdateSecret.Exec(
			secret.GetID(), auth.GetUsername(),
			metadata, secret.GetUpdatedAt().AsTime(),
			secret.GetData(),
		); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteSecrets(ctx context.Context, auth *pb.Auth, secrets []*pb.Secret) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollbackTx(tx)

	stmtDeleteSecret, err := tx.Prepare(
		`UPDATE secrets SET
			labels=NULL,
			created_at=NULL, 
			updated_at=NULL, 
			deleted_at=$3, 
			data=NULL, 
			deleted=true 
			WHERE id=$1 AND username=$2;`,
	)
	if err != nil {
		return err
	}
	defer closeObject(stmtDeleteSecret)

	for _, secret := range secrets {
		if _, err := stmtDeleteSecret.Exec(
			secret.GetID(), auth.GetUsername(), secret.GetDeletedAt().AsTime(),
		); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *DB) ListSecrets(ctx context.Context, auth *pb.Auth) ([]*pb.Secret, error) {
	var secrets []*pb.Secret
	rows, err := db.conn.QueryContext(ctx, `
		SELECT id, labels, 
		created_at, updated_at, deleted_at, deleted,
		data 
		FROM secrets
		WHERE username=$1;
	`, auth.GetUsername())
	if err != nil {
		return nil, err
	}
	defer closeObject(rows)

	for rows.Next() {
		var createdAtString, updatedAtString, deletedAtString sql.NullString
		secret := &pb.Secret{
			ID:     "",
			Labels: make(map[string]string, 0),
			Data:   make([]byte, 0),
			Status: &pb.Status{},
		}

		labels := make([]byte, 0)

		err = rows.Scan(
			&secret.ID, &labels,
			&createdAtString, &updatedAtString, &deletedAtString, &secret.Status.Deleted,
			&secret.Data)
		if err != nil {
			return nil, err
		}

		if len(labels) > 0 {
			if err := json.Unmarshal(labels, &secret.Labels); err != nil {
				return nil, err
			}
		}

		if createdAtString.String != "" {
			createdAt, err := time.Parse(time.RFC3339, createdAtString.String)
			if err != nil {
				return nil, err
			}
			secret.CreatedAt = timestamppb.New(createdAt)
		}

		if updatedAtString.String != "" {
			updatedAt, err := time.Parse(time.RFC3339, updatedAtString.String)
			if err != nil {
				return nil, err
			}
			secret.UpdatedAt = timestamppb.New(updatedAt)
		}

		if deletedAtString.String != "" {
			deletedAt, err := time.Parse(time.RFC3339, deletedAtString.String)
			if err != nil {
				return nil, err
			}
			secret.DeletedAt = timestamppb.New(deletedAt)
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func closeObject(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close object")
	}
}

func rollbackTx(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		log.Error().Err(err).Msg("Failed to rollback transaction")
	}
}
