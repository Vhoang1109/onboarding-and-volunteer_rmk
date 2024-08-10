package storage

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cesc1802/onboarding-and-volunteer-service/feature/authentication/dto"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.DB()

	repo := NewAuthenticationRepository(db)

	t.Run("successful retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "status"}).
			AddRow(1, "test@example.com", "password123", 1)
		mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT 1$").
			WithArgs("test@example.com").
			WillReturnRows(rows)

		user, errMsg := repo.GetUserByEmail("test@example.com", "password123")
		assert.Empty(t, errMsg)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT 1$").
			WithArgs("unknown@example.com").
			WillReturnError(gorm.ErrRecordNotFound)

		user, errMsg := repo.GetUserByEmail("unknown@example.com", "password123")
		assert.Equal(t, gorm.ErrRecordNotFound.Error(), errMsg)
		assert.Nil(t, user)
	})

	t.Run("inactive user", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "status"}).
			AddRow(1, "inactive@example.com", "password123", 0)
		mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT 1$").
			WithArgs("inactive@example.com").
			WillReturnRows(rows)

		user, errMsg := repo.GetUserByEmail("inactive@example.com", "password123")
		assert.Equal(t, "User is inactive", errMsg)
		assert.Nil(t, user)
	})

	t.Run("incorrect password", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "status"}).
			AddRow(1, "test@example.com", "password123", 1)
		mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT 1$").
			WithArgs("test@example.com").
			WillReturnRows(rows)

		user, errMsg := repo.GetUserByEmail("test@example.com", "wrongpassword")
		assert.Equal(t, "Password is incorrect", errMsg)
		assert.Nil(t, user)
	})
}

func TestRegisterUser(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer db.DB()

	repo := NewAuthenticationRepository(db)

	t.Run("successful registration", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO `users` (.+) VALUES (.+)$").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		request := &dto.RegisterUserRequest{
			Email:    "newuser@example.com",
			Name:     "New User",
			Password: "newpassword",
		}

		response, err := repo.RegisterUser(request)
		assert.NoError(t, err)
		assert.Equal(t, "User registered successfully", response.Message)
	})

	t.Run("registration error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO `users` (.+) VALUES (.+)$").
			WillReturnError(errors.New("some error"))
		mock.ExpectRollback()

		request := &dto.RegisterUserRequest{
			Email:    "newuser@example.com",
			Name:     "New User",
			Password: "newpassword",
		}

		_, err := repo.RegisterUser(request)
		assert.Error(t, err)
	})
}
