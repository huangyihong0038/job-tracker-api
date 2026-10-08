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
	listFunc    func(
		model.ListApplicationsRequest,
	) (model.ListApplicationsResponse, error)

	updateFunc func(
		int64,
		model.UpdateApplicationRequest,
	) (model.Application, error)
}

func (service *fakeApplicationService) Create(
	_ context.Context,
	request model.CreateApplicationRequest,
) (model.Application, error) {
	return service.createFunc(request)
}

func (service *fakeApplicationService) List(
	_ context.Context,
	request model.ListApplicationsRequest,
) (model.ListApplicationsResponse, error) {
	if service.listFunc != nil {
		return service.listFunc(request)
	}

	return model.ListApplicationsResponse{
		Data: []model.Application{},
		Pagination: model.Pagination{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    0,
		},
	}, nil
}

func (service *fakeApplicationService) GetByID(
	_ context.Context,
	id int64,
) (model.Application, error) {
	return service.getByIDFunc(id)
}

func (service *fakeApplicationService) Update(
	_ context.Context,
	id int64,
	request model.UpdateApplicationRequest,
) (model.Application, error) {
	if service.updateFunc != nil {
		return service.updateFunc(id, request)
	}
	return model.Application{}, nil
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

	router.GET(
		"/api/v1/applications",
		applicationHandler.List,
	)

	router.PATCH(
		"/api/v1/applications/:id",
		applicationHandler.Update,
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

func performUpdateRequest(
	router *gin.Engine,
	id int64,
	body string,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/applications/"+
			strconv.FormatInt(id, 10),
		bytes.NewBufferString(body),
	)

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	return response
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

	tests := []struct {
		name       string
		serviceErr error
	}{
		{
			name:       "invalid applied_at",
			serviceErr: service.ErrInvalidAppliedAt,
		},
		{
			name:       "notes too long",
			serviceErr: service.ErrNotesTooLong,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			applicationService := &fakeApplicationService{
				createFunc: func(
					request model.CreateApplicationRequest,
				) (model.Application, error) {
					return model.Application{}, test.serviceErr
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

			if !strings.Contains(body, test.serviceErr.Error()) {
				t.Fatalf(
					"expected message %q, got: %s",
					test.serviceErr.Error(),
					body,
				)
			}

		})

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

func TestApplicationHandlerListSuccess(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantRequest model.ListApplicationsRequest
	}{
		{
			name: "default pagination",
			path: "/api/v1/applications",
			wantRequest: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 20,
			},
		},
		{
			name: "explicit query parameters",
			path: "/api/v1/applications" +
				"?page=2" +
				"&page_size=10" +
				"&status=interview" +
				"&company=Byte",
			wantRequest: model.ListApplicationsRequest{
				Page:     2,
				PageSize: 10,
				Status:   "interview",
				Company:  "Byte",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			applicationService :=
				&fakeApplicationService{
					listFunc: func(
						request model.ListApplicationsRequest,
					) (
						model.ListApplicationsResponse,
						error,
					) {
						if request != test.wantRequest {
							t.Fatalf(
								"expected request %+v, got %+v",
								test.wantRequest,
								request,
							)
						}

						return model.ListApplicationsResponse{
							Data: []model.Application{
								{
									ID:      1,
									Company: "ByteDance",
									Status:  "interview",
								},
							},
							Pagination: model.Pagination{
								Page:     request.Page,
								PageSize: request.PageSize,
								Total:    1,
							},
						}, nil
					},
				}

			router := newTestRouter(
				applicationService,
			)

			request := httptest.NewRequest(
				http.MethodGet,
				test.path,
				nil,
			)
			responseRecorder :=
				httptest.NewRecorder()

			router.ServeHTTP(
				responseRecorder,
				request,
			)

			if responseRecorder.Code != http.StatusOK {
				t.Fatalf(
					"expected status %d, got %d; body: %s",
					http.StatusOK,
					responseRecorder.Code,
					responseRecorder.Body.String(),
				)
			}

			var response model.ListApplicationsResponse

			if err := json.Unmarshal(
				responseRecorder.Body.Bytes(),
				&response,
			); err != nil {
				t.Fatalf(
					"decode response: %v",
					err,
				)
			}

			if len(response.Data) != 1 {
				t.Fatalf(
					"expected one application, got %d",
					len(response.Data),
				)
			}

			if response.Data[0].ID != 1 {
				t.Fatalf(
					"expected application ID 1, got %d",
					response.Data[0].ID,
				)
			}

			if response.Pagination.Page !=
				test.wantRequest.Page {
				t.Fatalf(
					"expected page %d, got %d",
					test.wantRequest.Page,
					response.Pagination.Page,
				)
			}

			if response.Pagination.PageSize !=
				test.wantRequest.PageSize {
				t.Fatalf(
					"expected page size %d, got %d",
					test.wantRequest.PageSize,
					response.Pagination.PageSize,
				)
			}

			if response.Pagination.Total != 1 {
				t.Fatalf(
					"expected total 1, got %d",
					response.Pagination.Total,
				)
			}
		})
	}
}

func TestApplicationHandlerListErrors(t *testing.T) {
	internalError :=
		errors.New("database connection failed")

	tests := []struct {
		name              string
		path              string
		serviceErr        error
		wantStatus        int
		wantCode          string
		wantMessage       string
		wantServiceCalled bool
		forbiddenMessage  string
	}{
		{
			name:        "query parameter is not an integer",
			path:        "/api/v1/applications?page=abc",
			wantStatus:  http.StatusBadRequest,
			wantCode:    "INVALID_ARGUMENT",
			wantMessage: "request query is invalid",
		},
		{
			name: "invalid page",
			path: "/api/v1/applications" +
				"?page=0&page_size=20",
			serviceErr:        service.ErrInvalidPage,
			wantStatus:        http.StatusBadRequest,
			wantCode:          "INVALID_ARGUMENT",
			wantMessage:       service.ErrInvalidPage.Error(),
			wantServiceCalled: true,
		},
		{
			name: "invalid page size",
			path: "/api/v1/applications" +
				"?page=1&page_size=101",
			serviceErr:        service.ErrInvalidPageSize,
			wantStatus:        http.StatusBadRequest,
			wantCode:          "INVALID_ARGUMENT",
			wantMessage:       service.ErrInvalidPageSize.Error(),
			wantServiceCalled: true,
		},
		{
			name: "invalid status",
			path: "/api/v1/applications" +
				"?page=1&page_size=20&status=unknown",
			serviceErr:        service.ErrInvalidStatus,
			wantStatus:        http.StatusBadRequest,
			wantCode:          "INVALID_ARGUMENT",
			wantMessage:       service.ErrInvalidStatus.Error(),
			wantServiceCalled: true,
		},
		{
			name:              "internal error",
			path:              "/api/v1/applications",
			serviceErr:        internalError,
			wantStatus:        http.StatusInternalServerError,
			wantCode:          "INTERNAL_ERROR",
			wantMessage:       "internal server error",
			wantServiceCalled: true,
			forbiddenMessage:  internalError.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceCalled := false

			applicationService :=
				&fakeApplicationService{
					listFunc: func(
						request model.ListApplicationsRequest,
					) (
						model.ListApplicationsResponse,
						error,
					) {
						serviceCalled = true
						return model.ListApplicationsResponse{},
							test.serviceErr
					},
				}

			router := newTestRouter(
				applicationService,
			)

			request := httptest.NewRequest(
				http.MethodGet,
				test.path,
				nil,
			)
			responseRecorder :=
				httptest.NewRecorder()

			router.ServeHTTP(
				responseRecorder,
				request,
			)

			if responseRecorder.Code != test.wantStatus {
				t.Fatalf(
					"expected status %d, got %d; body: %s",
					test.wantStatus,
					responseRecorder.Code,
					responseRecorder.Body.String(),
				)
			}

			if serviceCalled != test.wantServiceCalled {
				t.Fatalf(
					"expected service called=%t, got %t",
					test.wantServiceCalled,
					serviceCalled,
				)
			}

			body := responseRecorder.Body.String()

			if !strings.Contains(
				body,
				`"code":"`+test.wantCode+`"`,
			) {
				t.Fatalf(
					"expected error code %q, got: %s",
					test.wantCode,
					body,
				)
			}

			if !strings.Contains(
				body,
				test.wantMessage,
			) {
				t.Fatalf(
					"expected message %q, got: %s",
					test.wantMessage,
					body,
				)
			}

			if test.forbiddenMessage != "" &&
				strings.Contains(
					body,
					test.forbiddenMessage,
				) {
				t.Fatalf(
					"response leaked internal error: %s",
					body,
				)
			}
		})
	}
}

func TestApplicationHandlerUpdateSuccess(
	t *testing.T,
) {
	appliedAt := "2026-10-08"

	applicationService := &fakeApplicationService{
		updateFunc: func(
			id int64,
			request model.UpdateApplicationRequest,
		) (model.Application, error) {
			if id != 7 {
				t.Fatalf(
					"expected ID 7, got %d",
					id,
				)
			}

			if !request.Status.Present {
				t.Fatal(
					"expected status to be present",
				)
			}

			if request.Status.Value == nil {
				t.Fatal(
					"expected status value",
				)
			}

			if *request.Status.Value != "interview" {
				t.Fatalf(
					"expected status interview, got %q",
					*request.Status.Value,
				)
			}

			if !request.AppliedAt.Present {
				t.Fatal(
					"expected applied_at to be present",
				)
			}

			if request.AppliedAt.Value != nil {
				t.Fatalf(
					"expected applied_at to be nil, got %q",
					*request.AppliedAt.Value,
				)
			}

			return model.Application{
				ID:        id,
				Company:   "ByteDance",
				Position:  "Backend Intern",
				Status:    "interview",
				AppliedAt: &appliedAt,
				Notes:     "一面通过",
			}, nil
		},
	}

	router := newTestRouter(
		applicationService,
	)

	response := performUpdateRequest(
		router,
		7,
		`{
			"status": "interview",
			"applied_at": null
		}`,
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusOK,
			response.Code,
			response.Body.String(),
		)
	}

	var application model.Application

	if err := json.Unmarshal(
		response.Body.Bytes(),
		&application,
	); err != nil {
		t.Fatalf(
			"decode response: %v",
			err,
		)
	}

	if application.ID != 7 {
		t.Fatalf(
			"expected ID 7, got %d",
			application.ID,
		)
	}

	if application.Status != "interview" {
		t.Fatalf(
			"expected status interview, got %q",
			application.Status,
		)
	}

	if application.Company != "ByteDance" {
		t.Fatalf(
			"expected company ByteDance, got %q",
			application.Company,
		)
	}
}

func TestApplicationHandlerUpdateErrors(
	t *testing.T,
) {
	internalError := errors.New(
		"database connection failed",
	)

	tests := []struct {
		name              string
		path              string
		body              string
		serviceErr        error
		wantStatus        int
		wantCode          string
		wantMessage       string
		wantServiceCalled bool
		forbiddenMessage  string
	}{
		{
			name: "invalid ID",
			path: "/api/v1/applications/" +
				"abc",
			body:              `{"status":"interview"}`,
			wantStatus:        http.StatusBadRequest,
			wantCode:          "INVALID_ARGUMENT",
			wantMessage:       "id must be a positive integer",
			wantServiceCalled: false,
		},
		{
			name: "invalid request body",
			path: "/api/v1/applications/" +
				"1",
			body:              `{"status":`,
			wantStatus:        http.StatusBadRequest,
			wantCode:          "INVALID_ARGUMENT",
			wantMessage:       "request body is invalid",
			wantServiceCalled: false,
		},
		{
			name: "validation error",
			path: "/api/v1/applications/" +
				"1",
			body:              `{}`,
			serviceErr:        service.ErrNoFieldsToUpdate,
			wantStatus:        http.StatusBadRequest,
			wantCode:          "INVALID_ARGUMENT",
			wantMessage:       service.ErrNoFieldsToUpdate.Error(),
			wantServiceCalled: true,
		},
		{
			name: "application not found",
			path: "/api/v1/applications/" +
				"999",
			body:              `{"status":"interview"}`,
			serviceErr:        service.ErrApplicationNotFound,
			wantStatus:        http.StatusNotFound,
			wantCode:          "NOT_FOUND",
			wantMessage:       "application not found",
			wantServiceCalled: true,
		},
		{
			name: "internal error",
			path: "/api/v1/applications/" +
				"1",
			body:              `{"status":"interview"}`,
			serviceErr:        internalError,
			wantStatus:        http.StatusInternalServerError,
			wantCode:          "INTERNAL_ERROR",
			wantMessage:       "internal server error",
			wantServiceCalled: true,
			forbiddenMessage:  internalError.Error(),
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			serviceCalled := false

			applicationService :=
				&fakeApplicationService{
					updateFunc: func(
						id int64,
						request model.UpdateApplicationRequest,
					) (model.Application, error) {
						serviceCalled = true

						return model.Application{},
							test.serviceErr
					},
				}

			router := newTestRouter(
				applicationService,
			)

			request := httptest.NewRequest(
				http.MethodPatch,
				test.path,
				bytes.NewBufferString(
					test.body,
				),
			)

			request.Header.Set(
				"Content-Type",
				"application/json",
			)

			response := httptest.NewRecorder()

			router.ServeHTTP(
				response,
				request,
			)

			if response.Code != test.wantStatus {
				t.Fatalf(
					"expected status %d, got %d; body: %s",
					test.wantStatus,
					response.Code,
					response.Body.String(),
				)
			}

			if serviceCalled !=
				test.wantServiceCalled {
				t.Fatalf(
					"expected service called=%t, got %t",
					test.wantServiceCalled,
					serviceCalled,
				)
			}

			body := response.Body.String()

			if !strings.Contains(
				body,
				`"code":"`+test.wantCode+`"`,
			) {
				t.Fatalf(
					"expected error code %q, got: %s",
					test.wantCode,
					body,
				)
			}

			if !strings.Contains(
				body,
				test.wantMessage,
			) {
				t.Fatalf(
					"expected message %q, got: %s",
					test.wantMessage,
					body,
				)
			}
			if test.forbiddenMessage != "" &&
				strings.Contains(
					body,
					test.forbiddenMessage,
				) {
				t.Fatalf(
					"response leaked internal error: %s",
					body,
				)
			}

		})

	}

}
