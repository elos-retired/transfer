package transfer_test

import (
	"errors"
	. "github.com/elos/transfer"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("APIError", func() {

	Describe("Creating ad-hoc errors", func() {
		Context("With all fields", func() {
			error500 := APIError{
				Status:           500,
				Code:             500,
				Message:          "Server Error",
				DeveloperMessage: "We screwed up",
			}

			It("should create the 500 error message", func() {
				Expect(error500.Status).To(Equal(500))
			})
		})
	})

	Describe("Writing Errors", func() {

		Context("Error Generators", func() {
			var (
				apiError *APIError
			)

			It("Not Found (404) Error", func() {
				apiError = NewNotFoundError()
				Expect(apiError).ToNot(BeNil())
				Expect(apiError.Status).To(Equal(404))
			})

			It("Server (500) Error", func() {
				apiError = NewServerError()
				Expect(apiError).ToNot(BeNil())
				Expect(apiError.Status).To(Equal(500))

				msg := "this is an error"
				err := errors.New(msg)
				apiError = NewServerErrorWithError(err)

				Expect(apiError).ToNot(BeNil())
				Expect(apiError.Status).To(Equal(500))
				Expect(apiError.DeveloperMessage).To(Equal(msg))
			})

			It("Invalid Method (405) Error", func() {
				apiError = NewInvalidMethodError()
				Expect(apiError).ToNot(BeNil())
				Expect(apiError.Status).To(Equal(405))
			})

			It("Unauthorized (401) Error)", func() {
				apiError = NewUnauthorizedError()
				Expect(apiError).ToNot(BeNil())
				Expect(apiError.Status).To(Equal(401))
			})

			It("WebSocketFailed (400) Error)", func() {
				apiError = NewWebSocketFailedError()
				Expect(apiError).ToNot(BeNil())
				Expect(apiError.Status).To(Equal(400))
			})

			It("CustomError Error", func() {
				status := 1
				code := 2
				msg := "asdfasdfasdflksadjflf"
				dmsg := "asdkjlfj2340--6948053.,a"
				apiError = NewCustomError(status, code, msg, dmsg)

				Expect(apiError).ToNot(BeNil())
				Expect(apiError.Status).To(Equal(1))
				Expect(apiError.Code).To(Equal(2))
				Expect(apiError.Message).To(Equal(msg))
				Expect(apiError.DeveloperMessage).To(Equal(dmsg))
			})
		})

		Context("Error Writers", func() {
			var (
				w *httptest.ResponseRecorder
			)

			JustBeforeEach(func() {
				w = httptest.NewRecorder()
			})

			It("Writes NotFound", func() {
				NotFound(w)
				notFoundError := NewNotFoundError()

				Expect(w.Code).To(Equal(notFoundError.Code))
				bodyBytes, err := ToJSON(notFoundError)

				Expect(err).ToNot(HaveOccurred())
				Expect(w.Body.String()).To(Equal(string(bodyBytes)))
			})

			It("Writes ServerError", func() {
				msg := "this is an error"
				err := errors.New(msg)

				ServerError(w, err)
				serverError := NewServerErrorWithError(err)

				Expect(w.Code).To(Equal(serverError.Code))
				bodyBytes, err := ToJSON(serverError)

				Expect(err).ToNot(HaveOccurred())
				Expect(w.Body.String()).To(Equal(string(bodyBytes)))
			})

			It("Writes InvalidMethodError", func() {
				InvalidMethod(w)
				invalidMethodError := NewInvalidMethodError()

				Expect(w.Code).To(Equal(invalidMethodError.Code))
				bodyBytes, err := ToJSON(invalidMethodError)

				Expect(err).ToNot(HaveOccurred())
				Expect(w.Body.String()).To(Equal(string(bodyBytes)))
			})

			It("Writes UnauthorizedError", func() {
				Unauthorized(w)
				unauthorizedError := NewUnauthorizedError()

				Expect(w.Code).To(Equal(unauthorizedError.Code))
				bodyBytes, err := ToJSON(unauthorizedError)

				Expect(err).ToNot(HaveOccurred())
				Expect(w.Body.String()).To(Equal(string(bodyBytes)))
			})

			It("Writes WebSocketFailedError", func() {
				WebSocketFailed(w)
				webSocketFailedError := NewWebSocketFailedError()

				Expect(w.Code).To(Equal(webSocketFailedError.Code))
				bodyBytes, err := ToJSON(webSocketFailedError)

				Expect(err).ToNot(HaveOccurred())
				Expect(w.Body.String()).To(Equal(string(bodyBytes)))
			})

			It("WritesCustomError", func() {
				status := 1
				code := 2
				msg := "laksdjflkjsdaf"
				dmsg := "asdfkljadksf"
				CustomError(w, status, code, msg, dmsg)

				customError := NewCustomError(status, code, msg, dmsg)
				bodyBytes, err := ToJSON(customError)

				Expect(err).ToNot(HaveOccurred())
				Expect(w.Body.String()).To(Equal(string(bodyBytes)))
			})
		})

		Context("Generic Error Writer", func() {
			It("WriteErrorResponse", func() {
				// Simple testing here because the WriteErrorResponse just calles ResourceResponse
				w := httptest.NewRecorder()
				apiError := NewCustomError(1, 2, "1", "2")

				WriteErrorResponse(w, apiError)
				Expect(w.Code).To(Equal(1))
				bodyBytes, err := ToJSON(apiError)

				Expect(err).ToNot(HaveOccurred())
				Expect(w.Body.String()).To(Equal(string(bodyBytes)))
			})
		})

	})

})
