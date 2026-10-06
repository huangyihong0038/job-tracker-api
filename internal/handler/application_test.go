package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"job-tracker-api/internal/model"
	"job-tracker-api/internal/repository"
	"job-tracker-api/internal/service"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeApplicationService struct {
	createFunc  func(model.CreateApplicationRequest) (model.Application, error)
	getByIDFunc func(int64) (model.Application, error)
	deleteFunc  func(int64) error
}

func (service *fakeApplicationService) Create(
	_ context.Context,
	request model.CreateApplicationRequest,
) (model.Application, error) {
	return service.createFunc(request)
}

func (service *fakeApplicationService) GetByID(
	_ context.Context,
	id int64,
) (model.Application, error) {
	return service.getByIDFunc(id)
}

func (service *fakeApplicationService) Delete(
	_ context.Context,
	id int64) error {
	if service.deleteFunc != nil {
		return service.deleteFunc(id)
	}
	return nil
}

func newTestRouter(
	applicationService ApplicationService,
) *gin.Engine {
	gin.SetMode(gin.TestMode)

	applicationHandler := NewApplicationHandler(applicationService)

	router := gin.New()

	router.POST(
		"/api/v1/applications",
		applicationHandler.Create,
	)

	router.GET(
		"/api/v1/applications/:id",
		applicationHandler.GetByID,
	)

	router.DELETE(
		"/api/v1/applications/:id",
		applicationHandler.Delete,
	)
	return router
}

func performCreateRequest(
	router *gin.Engine,
	body string,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/applications",
		bytes.NewBufferString(body),
	)

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	return responseRecorder
}

func TestApplicationHandlerCreateSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	applicationService := &fakeApplicationService{
		createFunc: func(
			request model.CreateApplicationRequest,
		) (model.Application, error) {
			if request.Company != "ByteDance" {
				t.Fatalf(
					"expected company ByteDance, got %q",
					request.Company,
				)
			}

			if request.Position != "Backend Intern" {
				t.Fatalf(
					"expected position Backend Intern, got %q",
					request.Position,
				)
			}

			if request.Status != "planned" {
				t.Fatalf(
					"expected status planned, got %q",
					request.Status,
				)
			}

			return model.Application{
				ID:       1,
				Company:  request.Company,
				Position: request.Position,
				Status:   request.Status,
			}, nil
		},
	}

	applicationHandler := NewApplicationHandler(applicationService)

	router := gin.New()

	router.POST(
		"/api/v1/applications",
		applicationHandler.Create,
	)

	requestBody := []byte(`{
		"company": "ByteDance",
		"position": "Backend Intern",
		"status": "planned"
	}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/applications",
		bytes.NewReader(requestBody),
	)

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(
		responseRecorder,
		request,
	)

	if responseRecorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusCreated,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	var response model.Application

	if err := json.Unmarshal(
		responseRecorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"expected valid JSON response, got error: %v",
			err,
		)
	}

	if response.ID != 1 {
		t.Fatalf("expected ID 1, got %d", response.ID)
	}

	if response.Company != "ByteDance" {
		t.Fatalf(
			"expected company ByteDance, got %q",
			response.Company,
		)
	}

}

func TestApplicationHandlerCreateValidationError(t *testing.T) {
	applicationService := &fakeApplicationService{
		createFunc: func(
			request model.CreateApplicationRequest,
		) (model.Application, error) {
			return model.Application{}, service.ErrInvalidStatus
		},
	}

	router := newTestRouter(applicationService)

	responseRecorder := performCreateRequest(
		router,
		`{
			"company": "ByteDance",
			"position": "Backend Intern",
			"status": "unknown"
		}`,
	)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusBadRequest,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	body := responseRecorder.Body.String()

	if !strings.Contains(body, `"code":"INVALID_ARGUMENT"`) {
		t.Fatalf(
			"expected INVALID_ARGUMENT, got: %s",
			body,
		)
	}

	if !strings.Contains(body, `"message":"invalid status"`) {
		t.Fatalf(
			"expected invalid status message, got: %s",
			body,
		)
	}
}

func TestApplicationHandlerCreateInternalError(t *testing.T) {
	applicationService := &fakeApplicationService{
		createFunc: func(
			request model.CreateApplicationRequest,
		) (model.Application, error) {
			return model.Application{},
				errors.New("database connection failed")
		},
	}

	router := newTestRouter(applicationService)

	responseRecorder := performCreateRequest(
		router,
		`{
			"company": "ByteDance",
			"position": "Backend Intern",
			"status": "planned"
		}`,
	)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusInternalServerError,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	body := responseRecorder.Body.String()

	if !strings.Contains(body, `"code":"INTERNAL_ERROR"`) {
		t.Fatalf(
			"expected INTERNAL_ERROR, got: %s",
			body,
		)
	}

	if strings.Contains(body, "database connection failed") {
		t.Fatalf(
			"response leaked internal error: %s",
			body,
		)
	}
}

func TestApplicationHandlerGetByIDSuccess(t *testing.T) {
	applicationService := &fakeApplicationService{
		getByIDFunc: func(id int64) (model.Application, error) {
			if id != 1 {
				t.Fatalf("expected ID 1, got %d", id)
			}

			return model.Application{
				ID:       1,
				Company:  "ByteDance",
				Position: "Backend Intern",
				Status:   "planned",
			}, nil
		},
	}

	router := newTestRouter(applicationService)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/applications/1",
		nil,
	)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s",
			recorder.Code, recorder.Body.String())
	}

	var response model.Application
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	if response.ID != 1 || response.Company != "ByteDance" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestApplicationHandlerGetByIDErrors(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		serviceErr error
		wantStatus int
		wantCode   string
		wantCalled bool
	}{
		{
			name:       "invalid ID",
			path:       "/api/v1/applications/abc",
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_ARGUMENT",
		},
		{
			name:       "not found",
			path:       "/api/v1/applications/999",
			serviceErr: service.ErrApplicationNotFound,
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
			wantCalled: true,
		},
		{
			name:       "internal error",
			path:       "/api/v1/applications/1",
			serviceErr: errors.New("database connection failed"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_ERROR",
			wantCalled: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceCalled := false
			fake := &fakeApplicationService{
				getByIDFunc: func(id int64) (model.Application, error) {
					serviceCalled = true
					return model.Application{}, test.serviceErr
				},
			}

			router := newTestRouter(fake)
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("expected HTTP %d, got %d; body: %s",
					test.wantStatus, recorder.Code, recorder.Body.String())
			}
			if serviceCalled != test.wantCalled {
				t.Fatalf("expected service called=%t, got %t",
					test.wantCalled, serviceCalled)
			}
			if !strings.Contains(recorder.Body.String(), test.wantCode) {
				t.Fatalf("expected error code %s; body: %s",
					test.wantCode, recorder.Body.String())
			}
			if strings.Contains(recorder.Body.String(), "database connection failed") {
				t.Fatalf("internal error leaked: %s", recorder.Body.String())
			}

		})
	}
}

func TestApplicationHandlerCreateInvalidJSON(t *testing.T) {
	serviceCalled := false
	fake := &fakeApplicationService{
		createFunc: func(
			request model.CreateApplicationRequest,
		) (model.Application, error) {
			serviceCalled = true
			return model.Application{}, nil
		},
	}

	router := newTestRouter(fake)
	recorder := performCreateRequest(router, `{"company":`)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d; body: %s",
			recorder.Code, recorder.Body.String())
	}
	if serviceCalled {
		t.Fatal("service was called for invalid JSON")
	}
	if !strings.Contains(recorder.Body.String(), `"code":"INVALID_ARGUMENT"`) {
		t.Fatalf("unexpected error response: %s", recorder.Body.String())
	}
}

func TestApplicationHandlerDeleteSuccess(t *testing.T) {
	called := false
	applicationService := &fakeApplicationService{
		deleteFunc: func(id int64) error {
			if id != 1 {
				t.Fatalf("expected id 1, got %d", id)
			}
			called = true
			return nil
		},
	}

	router := newTestRouter(applicationService)
	request := httptest.NewRequest(
		http.MethodDelete, "/api/v1/applications/1", nil,
	)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if !called {
		t.Fatal("service.Delete was not called")
	}
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d; body: %s",
			response.Code, response.Body.String())
	}
	if response.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", response.Body.String())
	}
}

func TestApplicationHandlerDeleteInvalidID(t *testing.T) {
	called := false
	applicationService := &fakeApplicationService{
		deleteFunc: func(id int64) error {
			called = true
			return nil
		},
	}

	router := newTestRouter(applicationService)
	request := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/applications/abc",
		nil,
	)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d; body: %s",
			response.Code, response.Body.String())
	}
	if called {
		t.Fatal("service.Delete should not be called for an invalid ID")
	}
}

func TestApplicationHandlerDeleteNotFound(t *testing.T) {
	applicationService := &fakeApplicationService{
		deleteFunc: func(id int64) error {
			if id != 99 {
				t.Fatalf("expected id 99, got %d", id)
			}
			return service.ErrApplicationNotFound
		},
	}

	router := newTestRouter(applicationService)
	request := httptest.NewRequest(
		http.MethodDelete, "/api/v1/applications/99", nil,
	)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d; body: %s",
			response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"code":"NOT_FOUND"`) {
		t.Fatalf("unexpected error response: %s", response.Body.String())
	}
}

func TestApplicationLifecycle(t *testing.T) {

	repository := repository.NewMemoryApplicationRepository()
	service := service.NewApplicationService(repository)

	router := newTestRouter(service)

	//创建
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/applications",
		strings.NewReader(`{
		"company": "ByteDance",
		"position":"Backend Intern",
		"status":"planned"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("预期状态码201，实际得到 %d，响应内容：%s", response.Code, response.Body.String())
	}

	var application model.Application

	if err := json.Unmarshal(response.Body.Bytes(), &application); err != nil {
		t.Fatalf("解析返回json失败: %v", err)
	}
	if application.ID <= 0 {
		t.Fatalf("expected a positive ID, got %d", application.ID)
	}

	path := "/api/v1/applications/" + strconv.FormatInt(application.ID, 10)

	//查询
	request = httptest.NewRequest(
		http.MethodGet,
		path,
		nil,
	)

	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("get want 200, got %d", response.Code)
	}

	//删除
	request = httptest.NewRequest(
		http.MethodDelete,
		path,
		nil,
	)

	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("delete want 204, got %d", response.Code)
	}

	//再次查询

	request = httptest.NewRequest(
		http.MethodGet,
		path,
		nil,
	)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("after delete, want 404, got %d", response.Code)
	}

}
