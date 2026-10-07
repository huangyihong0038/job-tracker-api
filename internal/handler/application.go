package handler

import (
	"context"
	"errors"
	"job-tracker-api/internal/model"
	"job-tracker-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApplicationService interface {
	Create(context.Context, model.CreateApplicationRequest) (model.Application, error)
	GetByID(context.Context, int64) (model.Application, error)
	Delete(context.Context, int64) error
}

type ApplicationHandler struct {
	service ApplicationService
}

func NewApplicationHandler(
	service ApplicationService,
) *ApplicationHandler {
	return &ApplicationHandler{
		service: service,
	}
}

func (handler *ApplicationHandler) Create(c *gin.Context) {
	var request model.CreateApplicationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(
			c,
			http.StatusBadRequest,
			"INVALID_ARGUMENT",
			"request body is invalid",
		)
		return
	}

	application, err := handler.service.Create(c.Request.Context(), request)
	if err != nil {
		if isValidationError(err) {
			writeError(
				c,
				http.StatusBadRequest,
				"INVALID_ARGUMENT",
				err.Error(),
			)
			return
		}
		writeError(
			c,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
		)
		return
	}
	c.JSON(http.StatusCreated, application)

}

func (handler *ApplicationHandler) GetByID(
	c *gin.Context,
) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil || id <= 0 {
		writeError(
			c,
			http.StatusBadRequest,
			"INVALID_ARGUMENT",
			"id must be a positive integer",
		)
		return
	}

	application, err := handler.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(
			err,
			service.ErrApplicationNotFound,
		) {
			writeError(
				c,
				http.StatusNotFound,
				"NOT_FOUND",
				"application not found",
			)
			return
		}

		writeError(
			c,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
		)
		return
	}

	c.JSON(http.StatusOK, application)

}

func (handler *ApplicationHandler) Delete(
	c *gin.Context,
) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil || id <= 0 {
		writeError(
			c,
			http.StatusBadRequest,
			"INVALID_ARGUMENT",
			"id must be a positive integer",
		)
		return
	}

	err = handler.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(
			err,
			service.ErrApplicationNotFound,
		) {
			writeError(
				c,
				http.StatusNotFound,
				"NOT_FOUND",
				"application not found",
			)
			return
		}

		writeError(
			c,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
		)
		return
	}

	c.Status(http.StatusNoContent)

}

func isValidationError(err error) bool {
	return errors.Is(err, service.ErrCompanyRequired) ||
		errors.Is(err, service.ErrCompanyTooLong) ||
		errors.Is(err, service.ErrPositionRequired) ||
		errors.Is(err, service.ErrPositionTooLong) ||
		errors.Is(err, service.ErrInvalidStatus) ||
		errors.Is(err, service.ErrInvalidAppliedAt) ||
		errors.Is(err, service.ErrNotesTooLong)

}

func writeError(
	c *gin.Context,
	status int,
	code string,
	message string,
) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
