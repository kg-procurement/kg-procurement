package mailer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"kg/procurement/internal/common/database"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/benbjohnson/clock"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
)

func Test_newPostgresEmailStatusAccessor(t *testing.T) {
	_ = newPostgresEmailStatusAccessor(nil, nil)
}

func Test_WriteEmailStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx         = context.Background()
			c           = setupEmailStatusAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			now         = c.cmock.Now()
			emailStatus = EmailStatus{
				ID:           "123",
				EmailTo:      "email@email.com",
				Status:       "sent",
				VendorID:     "100",
				DateSent:     now,
				ModifiedDate: now,
			}
		)

		transformedQuery, args, _ := sqlx.Named(insertEmailStatus, emailStatus)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.WriteEmailStatus(ctx, emailStatus)
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("returns error on db failure", func(t *testing.T) {
		var (
			ctx         = context.Background()
			c           = setupEmailStatusAccessorTestComponent(t)
			now         = c.cmock.Now()
			emailStatus = EmailStatus{
				ID:           "123",
				EmailTo:      "email@email.com",
				Status:       "sent",
				VendorID:     "100",
				DateSent:     now,
				ModifiedDate: now,
			}
		)

		transformedQuery, args, _ := sqlx.Named(insertEmailStatus, emailStatus)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(sql.ErrConnDone)

		err := c.accessor.WriteEmailStatus(ctx, emailStatus)
		c.g.Expect(err).Should(gomega.HaveOccurred())
	})
}

func Test_GetAll(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresEmailStatusAccessor
		mock     sqlmock.Sqlmock
	)

	setup := func(t *testing.T) (*gomega.WithT, *sql.DB) {
		g := gomega.NewWithT(t)
		realClock := clock.New()
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")

		accessor = newPostgresEmailStatusAccessor(sqlxDB, realClock)
		mock = sqlMock

		return g, db
	}

	emailStatusFields := []string{
		"id",
		"email_to",
		"status",
		"vendor_id",
		"date_sent",
		"modified_date",
	}

	dataQuery := `
		SELECT DISTINCT
			es.id,
			es.email_to,
			es.status,
			es.modified_date,
			es.vendor_id,
			es.date_sent
		FROM email_status es
		ORDER BY es.modified_date DESC
		LIMIT $1
		OFFSET $2
	`

	countQuery := "SELECT COUNT(*) from email_status es"

	fixedTime := time.Date(2024, time.September, 23, 12, 30, 0, 0, time.UTC)

	spec := GetAllEmailStatusSpec{
		PaginationSpec: database.PaginationSpec{
			Order: "DESC",
			Limit: 10,
			Page:  1,
		},
	}

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(emailStatusFields).
			AddRow(
				"1",
				"test@example.com",
				"sent",
				"100",
				fixedTime,
				fixedTime,
			)

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		emailStatusExpectation := []EmailStatus{{
			ID:           "1",
			EmailTo:      "test@example.com",
			Status:       "sent",
			VendorID:     "100",
			DateSent:     fixedTime,
			ModifiedDate: fixedTime,
		}}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).ToNot(gomega.BeNil())

		expectation := &AccessorGetEmailStatusPaginationData{
			EmailStatus: emailStatusExpectation,
			Metadata:    res.Metadata,
		}

		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("success with filter by email_to", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(emailStatusFields).
			AddRow(
				"1",
				"test@example.com",
				"sent",
				"100",
				fixedTime,
				fixedTime,
			)

		customSpec := GetAllEmailStatusSpec{
			EmailTo: "test@example.com",
			PaginationSpec: database.PaginationSpec{
				Order: "DESC",
				Limit: 10,
				Page:  1,
			},
		}

		dataQuery := `SELECT DISTINCT es.id, es.email_to, es.status, es.modified_date, es.vendor_id, es.date_sent FROM email_status es WHERE es.email_to ILIKE $1 ORDER BY es.modified_date DESC LIMIT $2 OFFSET $3`

		args := []driver.Value{
			"%" + customSpec.EmailTo + "%",
			int64(customSpec.PaginationSpec.Limit),
			int64((customSpec.PaginationSpec.Page - 1) * customSpec.PaginationSpec.Limit),
		}

		mock.ExpectQuery(dataQuery).
			WithArgs(args...).
			WillReturnRows(rows)

		countQuery := `SELECT COUNT(*) from email_status es WHERE es.email_to ILIKE $1`

		mock.ExpectQuery(countQuery).
			WithArgs("%" + customSpec.EmailTo + "%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, customSpec)

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).ToNot(gomega.BeNil())

		emailStatusExpectation := []EmailStatus{{
			ID:           "1",
			EmailTo:      "test@example.com",
			Status:       "sent",
			VendorID:     "100",
			DateSent:     fixedTime,
			ModifiedDate: fixedTime,
		}}

		expectation := &AccessorGetEmailStatusPaginationData{
			EmailStatus: emailStatusExpectation,
			Metadata:    res.Metadata,
		}

		g.Expect(res).To(gomega.Equal(expectation))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("success with ordering", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(emailStatusFields).
			AddRow(
				"1",
				"test@example.com",
				"sent",
				"100",
				fixedTime,
				fixedTime,
			)

		customQuery := `
			SELECT DISTINCT
				es.id,
				es.email_to,
				es.status,
				es.modified_date,
				es.vendor_id,
				es.date_sent
			FROM email_status es
			ORDER BY es.status DESC
			LIMIT $1
			OFFSET $2
		`

		customSpec := GetAllEmailStatusSpec{
			PaginationSpec: database.PaginationSpec{
				Order:   "DESC",
				Limit:   10,
				Page:    1,
				OrderBy: "status",
			},
		}

		args := database.BuildPaginationArgs(customSpec.PaginationSpec)

		mock.ExpectQuery(customQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, customSpec)

		emailStatusExpectation := []EmailStatus{{
			ID:           "1",
			EmailTo:      "test@example.com",
			Status:       "sent",
			VendorID:     "100",
			DateSent:     fixedTime,
			ModifiedDate: fixedTime,
		}}

		expectation := &AccessorGetEmailStatusPaginationData{
			EmailStatus: emailStatusExpectation,
			Metadata:    res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("success on empty result", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(emailStatusFields)

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery(countQuery).
			WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		expectation := &AccessorGetEmailStatusPaginationData{
			EmailStatus: []EmailStatus{},
			Metadata:    res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("error on scanning email status data rows", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(emailStatusFields).AddRow(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on executing email status data query", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		wrongQuery := `
			SELECT
				"email_to",
				"status",
				"modified_date"
			FROM email_status
			ORDER BY es.status DESC
			LIMIT $1
			OFFSET $2
		`

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(wrongQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnError(fmt.Errorf("query error"))

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error while iterating rows", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(emailStatusFields).
			AddRow(
				"1",
				"test@example.com",
				"sent",
				"100",
				fixedTime,
				fixedTime,
			).AddRow(
			"2",
			"test2@example.com",
			"failed",
			"100",
			fixedTime,
			fixedTime,
		).RowError(1, fmt.Errorf("row error"))

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(2)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on scanning total entry row", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(emailStatusFields).
			AddRow(
				"1",
				"test@example.com",
				"sent",
				"100",
				fixedTime,
				fixedTime,
			)

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		mock.ExpectQuery(countQuery).WillReturnError(fmt.Errorf("count query error"))

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})
}

func Test_Close(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		c := setupEmailStatusAccessorTestComponent(t)

		c.mock.ExpectClose()
		err := c.accessor.Close()
		c.g.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}

func Test_UpdateEmailStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx         = context.Background()
			c           = setupEmailStatusAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			now         = c.cmock.Now()
			emailStatus = EmailStatus{
				ID:           "123",
				VendorID:     "100",
				EmailTo:      "email@email.com",
				Status:       "sent",
				DateSent:     now,
				ModifiedDate: now,
			}
		)

		transformedQuery, args, _ := sqlx.Named(updateEmailStatus, emailStatus)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		rows := sqlmock.NewRows([]string{"id", "vendor_id", "email_to", "status", "date_sent", "modified_date"}).
			AddRow(emailStatus.ID, emailStatus.VendorID, emailStatus.EmailTo, emailStatus.Status, emailStatus.DateSent, emailStatus.ModifiedDate)

		c.mock.ExpectQuery(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnRows(rows)

		updatedEmailStatus, err := c.accessor.UpdateEmailStatus(ctx, emailStatus)
		c.g.Expect(err).Should(gomega.BeNil())
		c.g.Expect(updatedEmailStatus).ShouldNot(gomega.BeNil())
		c.g.Expect(updatedEmailStatus.ID).Should(gomega.Equal(emailStatus.ID))
		c.g.Expect(updatedEmailStatus.EmailTo).Should(gomega.Equal(emailStatus.EmailTo))
		c.g.Expect(updatedEmailStatus.Status).Should(gomega.Equal(emailStatus.Status))
		c.g.Expect(updatedEmailStatus.ModifiedDate).Should(gomega.Equal(emailStatus.ModifiedDate))
	})

	t.Run("error on row scan", func(t *testing.T) {
		var (
			ctx         = context.Background()
			c           = setupEmailStatusAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			now         = c.cmock.Now()
			emailStatus = EmailStatus{
				ID:           "123",
				VendorID:     "100",
				EmailTo:      "email@email.com",
				Status:       "sent",
				DateSent:     now,
				ModifiedDate: now,
			}
		)

		transformedQuery, args, _ := sqlx.Named(updateEmailStatus, emailStatus)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		rows := sqlmock.NewRows([]string{"id", "vendor_id", "email_to", "status", "date_sent", "modified_date"}).
			AddRow(nil, nil, nil, nil, nil, nil)

		c.mock.ExpectQuery(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnRows(rows)

		updatedEmailStatus, err := c.accessor.UpdateEmailStatus(ctx, emailStatus)
		c.g.Expect(err).Should(gomega.HaveOccurred())
		c.g.Expect(updatedEmailStatus).Should(gomega.BeNil())
	})

	t.Run("error on db failure", func(t *testing.T) {
		var (
			ctx         = context.Background()
			c           = setupEmailStatusAccessorTestComponent(t)
			now         = c.cmock.Now()
			emailStatus = EmailStatus{
				ID:           "123",
				VendorID:     "100",
				EmailTo:      "email@email.com",
				Status:       "sent",
				DateSent:     now,
				ModifiedDate: now,
			}
		)

		transformedQuery, args, _ := sqlx.Named(updateEmailStatus, emailStatus)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectQuery(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(sql.ErrConnDone)

		updatedEmailStatus, err := c.accessor.UpdateEmailStatus(ctx, emailStatus)
		c.g.Expect(err).Should(gomega.HaveOccurred())
		c.g.Expect(updatedEmailStatus).Should(gomega.BeNil())
	})
}

type emailStatusAccessorTestComponent struct {
	g        *gomega.WithT
	mock     sqlmock.Sqlmock
	db       *sql.DB
	accessor *postgresEmailStatusAccessor
	cmock    *clock.Mock
}

type setupOptions struct {
	queryMatcher sqlmock.QueryMatcher
}

func WithQueryMatcher(matcher sqlmock.QueryMatcher) Option {
	return func(o *setupOptions) {
		o.queryMatcher = matcher
	}
}

type Option func(*setupOptions)

func setupEmailStatusAccessorTestComponent(t *testing.T, opts ...Option) emailStatusAccessorTestComponent {
	options := setupOptions{
		queryMatcher: sqlmock.QueryMatcherEqual,
	}

	for _, opt := range opts {
		opt(&options)
	}

	g := gomega.NewWithT(t)
	db, sqlMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(options.queryMatcher))
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	clockMock := clock.NewMock()

	return emailStatusAccessorTestComponent{
		g:        g,
		mock:     sqlMock,
		db:       db,
		accessor: newPostgresEmailStatusAccessor(sqlxDB, clockMock),
		cmock:    clockMock,
	}
}
