package db

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brokeyourbike/xm-golang-exercise/api/requests"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: false,
	})
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
	company := models.Company{
		Name:    "test",
		Code:    "super-hash",
		Country: "US",
		Website: "example.com",
		Phone:   "+12345",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `companies`")).
		WithArgs(company.Name, company.Code, company.Country, company.Website, company.Phone).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	c, err := s.repository.Create(company)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), uint64(1), c.ID)
}

func (s *CompaniesSuite) TestItCanGetCompanyById() {
	company := models.Company{
		ID:      3,
		Name:    "test",
		Code:    "c12",
		Country: "UA",
		Website: "test.com",
		Phone:   "+12345",
	}

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `companies`")).
		WithArgs(company.ID).
		WillReturnRows((sqlmock.NewRows([]string{"id", "name", "code", "country", "website", "phone"})).
			AddRow("3", "test", "c12", "UA", "test.com", "+12345"))

	res, err := s.repository.Get(company.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), company, res)
}

func (s *CompaniesSuite) TestItCanReturnErrCompanyNotFound() {
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `companies`")).
		WithArgs(10).
		WillReturnError(gorm.ErrRecordNotFound)

	res, err := s.repository.Get(10)
	assert.ErrorIs(s.T(), err, models.ErrCompanyNotFound)
	assert.Equal(s.T(), models.Company{}, res)
}

func (s *CompaniesSuite) TestItCanReturnGeneralError() {
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `companies`")).
		WithArgs(10).
		WillReturnError(gorm.ErrInvalidField)

	res, err := s.repository.Get(10)
	assert.ErrorIs(s.T(), err, gorm.ErrInvalidField)
	assert.Equal(s.T(), models.Company{}, res)
}

func (s *CompaniesSuite) TestItCanDeleteCompany() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `companies`")).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repository.Delete(10)
	assert.NoError(s.T(), err)
}

func (s *CompaniesSuite) TestItCanUpdateCompany() {
	company := models.Company{
		ID:      3,
		Name:    "test",
		Code:    "c12",
		Country: "UA",
		Website: "test.com",
		Phone:   "+12345",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `companies`")).
		WithArgs(company.Name, company.Code, company.Country, company.Website, company.Phone, company.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repository.Update(company)
	assert.NoError(s.T(), err)
}

func (s *CompaniesSuite) TestItCanSearchCompany() {
	p := requests.CompanyPayload{Website: "https://google.com", Phone: "+12345"}

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `companies`")).
		WithArgs(p.Website, p.Phone).
		WillReturnRows((sqlmock.NewRows([]string{"id", "name", "code", "country", "website", "phone"})).
			AddRow("3", "test", "c12", "UA", "test.com", "+12345"))

	_, err := s.repository.GetAll(p)
	assert.NoError(s.T(), err)
}
