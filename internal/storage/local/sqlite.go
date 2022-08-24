package local

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os/user"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/go-rfe/gpwd/internal/agent/migrations"
	"github.com/go-rfe/gpwd/internal/logging/log"
	pb "github.com/go-rfe/gpwd/internal/proto"
)

const (
	timestamppbDateFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

var (
	_ Secrets  = (*sqliteStorage)(nil)
	_ Accounts = (*sqliteStorage)(nil)

	ErrAccountNotExists = errors.New("account doesn't exist")
	ErrAccountExists    = errors.New("account already exists")
)

type sqliteStorage struct {
	conn *sql.DB
}

func NewSQLiteStorage(path string, masterPassword []byte) (*sqliteStorage, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	username := currentUser.Username

	databasePath := fmt.Sprintf(
		"%s/%s?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=SHA256",
		path, "secrets.db", username, masterPassword,
	)

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}

	ss := &sqliteStorage{
		conn: db,
	}

	if err := ss.migrate(); err != nil {
		return nil, err
	}

	return ss, nil
}

func (ss *sqliteStorage) CreateSecret(ctx context.Context, secret *pb.Secret) (string, error) {
	var err error
	var metadata []byte
	if secret.Labels != nil {
		metadata, err = json.Marshal(secret.Labels)
		if err != nil {
			return "", err
		}
	}

	_, err = ss.conn.ExecContext(ctx, `
		INSERT INTO secrets (id, labels, created_at, synced, data) VALUES (?, ?, ?, ?);
	`, secret.ID, metadata, secret.CreatedAt.AsTime().String(), secret.GetStatus().GetSynced(), secret.GetData())
	if err != nil {
		return "", err
	}

	return secret.ID, nil
}

func (ss *sqliteStorage) ListSecrets(ctx context.Context) ([]*pb.Secret, error) {
	var secrets []*pb.Secret
	rows, err := ss.conn.QueryContext(ctx, `
		SELECT id, labels, 
		created_at, updated_at, deleted_at,
		synced, deleted,
		data 
		FROM secrets;
	`)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

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
			&createdAtString, &updatedAtString, &deletedAtString,
			&secret.Status.Synced, &secret.Status.Deleted,
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
			createdAt, err := time.Parse(timestamppbDateFormat, createdAtString.String)
			if err != nil {
				return nil, err
			}
			secret.CreatedAt = timestamppb.New(createdAt)
		}

		if updatedAtString.String != "" {
			updatedAt, err := time.Parse(timestamppbDateFormat, updatedAtString.String)
			if err != nil {
				return nil, err
			}
			secret.UpdatedAt = timestamppb.New(updatedAt)
		}

		if deletedAtString.String != "" {
			deletedAt, err := time.Parse(timestamppbDateFormat, deletedAtString.String)
			if err != nil {
				return nil, err
			}
			secret.DeletedAt = timestamppb.New(deletedAt)
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (ss *sqliteStorage) GetSecret(ctx context.Context, id string) (*pb.Secret, error) {
	secret := &pb.Secret{
		Labels: make(map[string]string, 0),
		Data:   make([]byte, 0),
		Status: &pb.Status{},
	}

	labels := make([]byte, 0)

	var createdAtString, updatedAtString, deletedAtString sql.NullString
	row := ss.conn.QueryRowContext(ctx, `
		SELECT id, labels, created_at, updated_at, deleted_at, synced, deleted, data FROM secrets WHERE id=?;
	`, id)

	err := row.Scan(
		&secret.ID, &labels,
		&createdAtString, &updatedAtString, &deletedAtString,
		&secret.Status.Synced, &secret.Status.Deleted,
		&secret.Data,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoSecretFound
	}

	if err != nil {
		return nil, err
	}

	if len(labels) > 0 {
		if err := json.Unmarshal(labels, &secret.Labels); err != nil {
			return nil, err
		}
	}

	if createdAtString.String != "" {
		createdAt, err := time.Parse(timestamppbDateFormat, createdAtString.String)
		if err != nil {
			return nil, err
		}
		secret.CreatedAt = timestamppb.New(createdAt)
	}

	if updatedAtString.String != "" {
		updatedAt, err := time.Parse(timestamppbDateFormat, updatedAtString.String)
		if err != nil {
			return nil, err
		}
		secret.UpdatedAt = timestamppb.New(updatedAt)
	}

	if deletedAtString.String != "" {
		deletedAt, err := time.Parse(timestamppbDateFormat, deletedAtString.String)
		if err != nil {
			return nil, err
		}
		secret.DeletedAt = timestamppb.New(deletedAt)
	}

	return secret, nil
}

func (ss *sqliteStorage) UpdateSecret(ctx context.Context, secret *pb.Secret) error {
	metadata, err := json.Marshal(secret.GetLabels())
	if err != nil {
		return err
	}

	var updatedAt string
	if secret.GetUpdatedAt() != nil {
		updatedAt = secret.GetUpdatedAt().AsTime().String()
	}

	_, err = ss.conn.ExecContext(ctx, `
		UPDATE secrets SET labels=?, updated_at=?, synced=?, data=? WHERE id=?;
	`, metadata, updatedAt, secret.GetStatus().GetSynced(), secret.GetData(), secret.GetID())
	if err != nil {
		return err
	}

	return nil
}

func (ss *sqliteStorage) DeleteSecret(ctx context.Context, secret *pb.Secret) error {
	deletedAt := timestamppb.Now()
	_, err := ss.conn.ExecContext(ctx, `
		UPDATE secrets SET
		labels=NULL,
		created_at=NULL, 
		updated_at=NULL, 
		deleted_at=?, 
		data=NULL,
		synced=?, 
		deleted=true 
		WHERE id=?;
	`, deletedAt.AsTime().String(), secret.Status.GetSynced(), secret.GetID())
	if err != nil {
		return err
	}

	return nil
}

func (ss *sqliteStorage) CreateAccount(ctx context.Context, account *pb.Account) (string, error) {
	_, err := ss.conn.ExecContext(ctx, `
		INSERT INTO accounts (id, server, username, password) VALUES (?, ?, ?, ?);
	`, account.GetID(), account.GetServerAddress(), account.GetUserName(), account.GetUserPassword())
	if err != nil {
		return "", err
	}

	return account.ID, nil
}

func (ss *sqliteStorage) GetAccount(ctx context.Context) (*pb.Account, error) {
	account := &pb.Account{
		UserPassword: make([]byte, 0),
	}

	row := ss.conn.QueryRowContext(ctx, `SELECT id, server, username, password, registered FROM accounts LIMIT 1;`)

	err := row.Scan(&account.ID, &account.ServerAddress, &account.UserName, &account.UserPassword, &account.Registered)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotExists
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (ss *sqliteStorage) UpdateAccount(ctx context.Context, account *pb.Account) error {
	_, err := ss.conn.ExecContext(ctx, `
		UPDATE accounts SET server=?, username=?, password=?, registered=? WHERE id=?;
	`, account.GetServerAddress(), account.GetUserName(), account.GetUserPassword(), account.GetRegistered(), account.GetID())
	if err != nil {
		return err
	}

	return nil
}

func (ss *sqliteStorage) DeleteAccount(ctx context.Context) error {
	_, err := ss.conn.ExecContext(ctx, `DELETE FROM accounts;`)
	if err != nil {
		return err
	}

	return nil
}

func (ss *sqliteStorage) Close() error {
	return ss.conn.Close()
}

func (ss *sqliteStorage) migrate() error {
	data := bindata.Resource(migrations.AssetNames(), migrations.Asset)

	sourceDriver, err := bindata.WithInstance(data)
	if err != nil {
		return err
	}

	db, err := sqlite3.WithInstance(ss.conn, &sqlite3.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithInstance("go-bindata", sourceDriver, "sqlite3", db)
	if err != nil {
		return err
	}

	if err := migration.Up(); !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func closeRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't close rows")
	}
}
