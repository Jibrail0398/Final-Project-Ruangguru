package main_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"golang.org/x/crypto/bcrypt"
	"a21hc3NpZ25tZW50/repository"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReportRepository", func() {
	var (
		mockDB  *sql.DB
		mockSql sqlmock.Sqlmock
		reportRepo repository.ReportRepository
	)

	BeforeEach(func() {
		// Create a new mock database connection
		var err error
		mockDB, mockSql, err = sqlmock.New()
		Expect(err).To(BeNil(), "Failed to create mock database")

		// Create repository with mock database
		reportRepo = repository.NewReportRepository(mockDB)
	})

	AfterEach(func() {
		// Close the mock database connection
		mockDB.Close()
	})

	

	Describe("GetReportByUser", func() {

		Context("When database query fails", func() {
			It("Should return an error on query failure", func() {
				// Expect the SELECT query and simulate a database error
				mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM report WHERE id = $1`)).
					WithArgs(1).
					WillReturnError(errors.New("database query error"))

				// Attempt to get reports
				reports, err := reportRepo.GetReportByUser(1)

				// Assert error occurred and no reports returned
				Expect(err).To(HaveOccurred())
				Expect(reports).To(BeNil())
				Expect(err.Error()).To(Equal("error while querying report"))
			})
		})
	})

	Describe("Delete", func() {
		Context("When deleting a report successfully", func() {
			It("Should delete report without error", func() {
				// Expect the DELETE query
				mockSql.ExpectExec(regexp.QuoteMeta(`DELETE FROM report WHERE id = $1`)).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Attempt to delete report
				err := reportRepo.Delete(1)

				// Assert no error occurred
				Expect(err).To(BeNil())
			})
		})

		Context("When deletion fails", func() {
			It("Should return an error if database deletion fails", func() {
				// Expect the DELETE query and simulate a database error
				mockSql.ExpectExec(regexp.QuoteMeta(`DELETE FROM report WHERE id = $1`)).
					WithArgs(1).
					WillReturnError(errors.New("deletion failed"))

				// Attempt to delete report
				err := reportRepo.Delete(1)

				// Assert error occurred
				Expect(err).To(HaveOccurred())
			})
		})
	})
})


// func TestReportRepository(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "ReportRepository Suite")
// 	RunSpecs(t, "UserRepository Suite")
// 	RunSpecs(t, "ChatRepository Suite")
// }

var _ = Describe("UserRepository", func() {
	var (
		mockDB  *sql.DB
		mockSql sqlmock.Sqlmock
		userRepo repository.UserRepository
	)

	BeforeEach(func() {
		// Create a new mock database connection
		var err error
		mockDB, mockSql, err = sqlmock.New()
		Expect(err).To(BeNil(), "Failed to create mock database")

		// Create repository with mock database
		userRepo = repository.NewUserRepo(mockDB)
	})

	AfterEach(func() {
		// Close the mock database connection
		mockDB.Close()
	})

	Describe("Register", func() {
		Context("When registering a new user", func() {
			It("Should successfully register a user", func() {
				// Prepare test data
				username := "testuser"
				email := "test@example.com"
				password := "testpassword"

				// Expect the INSERT query with specific parameters
				// We'll use a matcher for the hashed password
				mockSql.ExpectExec(regexp.QuoteMeta(`INSERT INTO users(username,email,password) VALUES ($1,$2,$3)`)).
					WithArgs(username, email, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Attempt to register
				err := userRepo.Register(username, email, password)

				// Assert no error occurred
				Expect(err).To(BeNil())
			})

			It("Should return an error if database insertion fails", func() {
				// Prepare test data
				username := "testuser"
				email := "test@example.com"
				password := "testpassword"

				// Expect the INSERT query and simulate a database error
				mockSql.ExpectExec(regexp.QuoteMeta(`INSERT INTO users(username,email,password) VALUES ($1,$2,$3)`)).
					WithArgs(username, email, sqlmock.AnyArg()).
					WillReturnError(errors.New("database insertion error"))

				// Attempt to register
				err := userRepo.Register(username, email, password)

				// Assert error occurred
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("registration failed"))
			})
		})
	})

	Describe("GetUserByEmail", func() {
		Context("When user exists", func() {
			

			It("Should return error when user not found", func() {
				// Prepare test data
				email := "nonexistent@example.com"

				// Expect the SELECT query with no rows
				mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id,username,email,password FROM users WHERE email = $1`)).
					WithArgs(email).
					WillReturnError(sql.ErrNoRows)

				// Attempt to get user
				user, err := userRepo.GetUserByEmail(email)

				// Assert error occurred and no user returned
				Expect(err).To(HaveOccurred())
				Expect(user).To(BeZero())
				Expect(err.Error()).To(ContainSubstring("user with email nonexistent@example.com not found"))
			})
		})
	})

	Describe("Login", func() {
		Context("When login credentials are valid", func() {
			It("Should successfully login with correct credentials", func() {
				// Prepare test data
				email := "test@example.com"
				password := "testpassword"
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

				// Prepare mock row for GetUserByEmail
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password"}).
					AddRow(1, "testuser", email, string(hashedPassword))

				// Expect the SELECT query for GetUserByEmail
				mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id,username,email,password FROM users WHERE email = $1`)).
					WithArgs(email).
					WillReturnRows(rows)

				// Attempt to login
				user, err := userRepo.Login(email, password)

				// Assert no error and correct user returned
				Expect(err).To(BeNil())
				Expect(user.Email).To(Equal(email))
			})

			It("Should return error with incorrect password", func() {
				// Prepare test data
				email := "test@example.com"
				correctPassword := "testpassword"
				incorrectPassword := "wrongpassword"
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

				// Prepare mock row for GetUserByEmail
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password"}).
					AddRow(1, "testuser", email, string(hashedPassword))

				// Expect the SELECT query for GetUserByEmail
				mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT id,username,email,password FROM users WHERE email = $1`)).
					WithArgs(email).
					WillReturnRows(rows)

				// Attempt to login
				user, err := userRepo.Login(email, incorrectPassword)

				// Assert error occurred and no user returned
				Expect(err).To(HaveOccurred())
				Expect(user).To(BeZero())
				Expect(err.Error()).To(Equal("invalid password"))
			})
		})
	})
})





var _ = Describe("ChatRepository", func() {
	var (
		mockDB  *sql.DB
		mock    sqlmock.Sqlmock
		chatRepo repository.ChatRepository
		err     error
	)

	BeforeEach(func() {
		mockDB, mock, err = sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())

		chatRepo = repository.NewChatRepository(mockDB)
	})

	AfterEach(func() {
		mockDB.Close()
	})

	Describe("SaveChat", func() {
		It("should save chat successfully", func() {
			// Prepare mock expectations
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO chat(date,question,response,fk_id_user,fk_report_id) VALUES($1,$2,$3,$4,$5)")).
				WithArgs("2023-06-15", "Test Question", "Test Response", 1, 100).
				WillReturnResult(sqlmock.NewResult(1, 1))

			// Execute the method
			err := chatRepo.SaveChat("2023-06-15", "Test Question", "Test Response", 1, 100)

			// Assert expectations
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
		})

		It("should return error when database insertion fails", func() {
			// Prepare mock expectations
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO chat(date,question,response,fk_id_user,fk_report_id) VALUES($1,$2,$3,$4,$5)")).
				WithArgs("2023-06-15", "Test Question", "Test Response", 1, 100).
				WillReturnError(fmt.Errorf("database error"))

			// Execute the method
			err := chatRepo.SaveChat("2023-06-15", "Test Question", "Test Response", 1, 100)

			// Assert expectations
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("error while save chat"))
			Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
		})
	})

	Describe("GetChatByReport", func() {
		It("should retrieve chats for a specific report", func() {
			// Prepare mock rows
			rows := sqlmock.NewRows([]string{"id", "date", "question", "response"}).
				AddRow(1, "2023-06-15", "Q1", "R1").
				AddRow(2, "2023-06-16", "Q2", "R2")

			// Prepare mock expectations
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id,date,question,response FROM chat WHERE fk_report_id = $1")).
				WithArgs(100).
				WillReturnRows(rows)

			// Execute the method
			chats, err := chatRepo.GetChatByReport(100)

			// Assert expectations
			Expect(err).ShouldNot(HaveOccurred())
			Expect(chats).Should(HaveLen(2))
			Expect(chats[0].Id).Should(Equal(1))
			Expect(chats[0].Date).Should(Equal("2023-06-15"))
			Expect(chats[0].Question).Should(Equal("Q1"))
			Expect(chats[0].Response).Should(Equal("R1"))
			Expect(chats[1].Id).Should(Equal(2))
			Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
		})

		It("should return empty slice when no rows found", func() {
			// Prepare mock expectations
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id,date,question,response FROM chat WHERE fk_report_id = $1")).
				WithArgs(100).
				WillReturnRows(sqlmock.NewRows([]string{"id", "date", "question", "response"}))

			// Execute the method
			chats, err := chatRepo.GetChatByReport(100)

			// Assert expectations
			Expect(err).ShouldNot(HaveOccurred())
			Expect(chats).Should(BeEmpty())
			Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
		})

		It("should return error when query fails", func() {
			// Prepare mock expectations
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id,date,question,response FROM chat WHERE fk_report_id = $1")).
				WithArgs(100).
				WillReturnError(fmt.Errorf("query error"))

			// Execute the method
			chats, err := chatRepo.GetChatByReport(100)

			// Assert expectations
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("error while querying chat data"))
			Expect(chats).Should(BeNil())
			Expect(mock.ExpectationsWereMet()).ShouldNot(HaveOccurred())
		})
	})
})

func TestReportRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
	// RunSpecs(t, "UserRepository Suite")
	// RunSpecs(t, "ChatRepository Suite")
}

