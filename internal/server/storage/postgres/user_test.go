package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	pg "github.com/lib/pq"
	"gophkeeper/internal/server/model"
	"reflect"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	newUUID := uuid.New()

	mock.ExpectQuery(`INSERT INTO users`).WithArgs("good@example.org", "Password").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(newUUID.String()),
	)
	mock.ExpectQuery(`INSERT INTO users`).WithArgs("existing@example.org", "Password").WillReturnError(
		&pg.Error{
			Code:    pgerrcode.IntegrityConstraintViolation,
			Message: "some error",
		})
	mock.ExpectQuery(`INSERT INTO users`).WithArgs("failing@example.org", "Password").WillReturnError(
		errors.New("you shall not pass"),
	)
	defer func() {
		_ = mdb.Close()
	}()

	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "create user",
			args: args{
				context.TODO(),
				&model.User{
					Email:    "good@example.org",
					Password: "Password",
				},
			},
			want: &model.User{
				ID:       newUUID,
				Email:    "good@example.org",
				Password: "Password",
			},
			wantErr: false,
		},
		{
			name: "create existing user",
			args: args{
				context.TODO(),
				&model.User{
					Email:    "existing@example.org",
					Password: "Password",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "create failing user",
			args: args{
				context.TODO(),
				&model.User{
					Email:    "failing@example.org",
					Password: "Password",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{
				db: mdb,
			}
			got, err := r.Create(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_Read(t *testing.T) {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	goodUUID := uuid.New()
	missingUUID := uuid.New()
	failingUUID := uuid.New()

	mock.ExpectQuery(`SELECT (.+) FROM users`).WithArgs(goodUUID.String()).WillReturnRows(
		sqlmock.NewRows([]string{"id", "email"}).AddRow(goodUUID.String(), "good@example.org"),
	)
	mock.ExpectQuery(`SELECT (.+) FROM users`).WithArgs(missingUUID.String()).WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(`SELECT (.+) FROM users`).WithArgs(failingUUID.String()).WillReturnError(
		errors.New("you shall not pass"),
	)
	defer func() {
		_ = mdb.Close()
	}()

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "read good user",
			args: args{
				context.TODO(),
				goodUUID,
			},
			want: &model.User{
				ID:    goodUUID,
				Email: "good@example.org",
			},
			wantErr: false,
		},
		{
			name: "read missing user",
			args: args{
				context.TODO(),
				missingUUID,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "read failing user",
			args: args{
				context.TODO(),
				failingUUID,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{
				db: mdb,
			}
			got, err := r.Read(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_ReadByEmailAndPassword(t *testing.T) {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	goodUUID := uuid.New()

	mock.ExpectQuery(`SELECT (.+) FROM users`).WithArgs("good@example.org", "Password").WillReturnRows(
		sqlmock.NewRows([]string{"id", "email"}).AddRow(goodUUID.String(), "good@example.org"),
	)
	mock.ExpectQuery(`SELECT (.+) FROM users`).WithArgs("good@example.org", "BadPassword").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(`SELECT (.+) FROM users`).WithArgs("failing@example.org", "Password").WillReturnError(
		errors.New("you shall not pass"),
	)
	defer func() {
		_ = mdb.Close()
	}()

	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "read with ok password",
			args: args{
				context.TODO(),
				"good@example.org",
				"Password",
			},
			want: &model.User{
				ID:    goodUUID,
				Email: "good@example.org",
			},
			wantErr: false,
		},
		{
			name: "read with bad password",
			args: args{
				context.TODO(),
				"good@example.org",
				"BadPassword",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "read failing user",
			args: args{
				context.TODO(),
				"failing@example.org",
				"BadPassword",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{
				db: mdb,
			}
			got, err := r.ReadByEmailAndPassword(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadByEmailAndPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadByEmailAndPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}
