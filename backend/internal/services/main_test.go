package services_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/testlib"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type ServiceTestCtx struct {
	db      *db.QueriesExt
	logger  zerolog.Logger
	user    dtos.UserRegister
	dbuser  dtos.User
	admin   dtos.UserRegister
	dbadmin dtos.User
}

func SetupServiceTest(t *testing.T) ServiceTestCtx {
	t.Helper()

	var (
		logger  = testlib.Logger(t)
		queries = testlib.NewDatabase(t, logger)
	)

	svcuser := services.NewUserService(logger, queries)
	svcadmin := services.NewAdminService(logger, queries)

	user := dtos.UserRegister{
		Email:    faker.Email(),
		Username: faker.Username(),
		Password: faker.Password(),
	}

	admin := dtos.UserRegister{
		Email:    faker.Email(),
		Username: faker.Username(),
		Password: faker.Password(),
	}

	dbuser, err := svcuser.Register(context.Background(), user)
	require.NoError(t, err)

	dbadmin, err := svcadmin.Register(context.Background(), admin)
	require.NoError(t, err)

	return ServiceTestCtx{
		db:      queries,
		logger:  logger,
		user:    user,
		admin:   admin,
		dbuser:  dbuser,
		dbadmin: dbadmin,
	}
}
