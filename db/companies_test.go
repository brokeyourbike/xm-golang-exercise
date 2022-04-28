package db

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CompaniesSuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	db         *gorm.DB
	repository *CompaniesRepo
}

func (s *CompaniesSuite) SetupDatabase() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{SkipDefaultTransaction: false})
	require.NoError(s.T(), err)
}

func (s *CompaniesSuite) SetupTest() {
	s.SetupDatabase()
	s.repository = NewCompaniesRepo(s.db)
}

func (s *CompaniesSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(CompaniesSuite))
}

func (s *CompaniesSuite) TestItCanCreateCompany() {
	user := models.Company{
		Name:    "test",
		Code:    "super-hash",
		Country: "US",
		Website: "example.com",
		Phone:   "+12345",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO `companies`").
		WithArgs(user.Name, user.Code, user.Country, user.Website, user.Phone).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	u, err := s.repository.Create(user)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), uint64(1), u.ID)
}
