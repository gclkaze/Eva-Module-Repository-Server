package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/routes"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

func UserLogsIn(u *models.UserAccount, pass string, t *testing.T, r *gin.Engine) *models.RequestResult[models.LoginResponse] {
	assert.Equal(t, pass != "", true)
	body := map[string]interface{}{
		"email":    u.Email,
		"password": pass,
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", routes.APIGroup, routes.AuthGroup, routes.LoginEndpoint), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	s := w.Body.String()
	fmt.Printf("%s", s)

	var resp models.RequestResult[models.LoginResponse]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, err == nil, true)
	return &resp
}

func ModuleCreate(title string, repr string, tags string, description string, filePath string,
	access *models.RequestResult[models.LoginResponse], t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["title"] = title
	fields["repr"] = repr
	fields["tags"] = tags
	fields["description"] = description

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		assert.Equal(t, false, true)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		assert.Equal(t, false, true)
	}

	if _, err := io.Copy(part, file); err != nil {
		assert.Equal(t, false, true)
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", routes.APIGroup, routes.ModulesGroup, routes.ModuleUploadEndpoint), body)
	req.Header.Set("Authorization", "Bearer "+(*access).Value.AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ModuleCreateMultipleFiles(title string, repr string, tags string, description string, filePaths []string,
	access *models.RequestResult[models.LoginResponse], t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["title"] = title
	fields["repr"] = repr
	fields["tags"] = tags
	fields["description"] = description

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	for i := range filePaths {
		file, err := os.Open(filePaths[i])
		if err != nil {
			assert.Equal(t, false, true)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("file", file.Name())
		if err != nil {
			assert.Equal(t, false, true)
		}

		if _, err := io.Copy(part, file); err != nil {
			assert.Equal(t, false, true)
		}
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", routes.APIGroup, routes.ModulesGroup, routes.ModuleUploadEndpoint), body)
	req.Header.Set("Authorization", "Bearer "+(*access).Value.AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ModuleCreateNoFile(title string, repr string, tags string, description string,
	access *models.RequestResult[models.LoginResponse], t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["title"] = title
	fields["repr"] = repr
	fields["tags"] = tags
	fields["description"] = description

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", routes.APIGroup, routes.ModulesGroup, routes.ModuleUploadEndpoint), body)
	req.Header.Set("Authorization", "Bearer "+(*access).Value.AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ModuleSuggest(modID uint, version string,
	access *models.LoginResponse, t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["modId"] = utils.UintToString(modID)
	fields["version"] = version

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", routes.APIGroup, routes.ModulesGroup, routes.ModuleSuggestEndpoint), body)
	req.Header.Set("Authorization", "Bearer "+(*access).AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ModuleDelete(modID uint,
	access *models.LoginResponse, t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["id"] = utils.UintToString(modID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", routes.APIGroup, routes.ModulesGroup, routes.ModuleDeleteEndpoint), body)
	req.Header.Set("Authorization", "Bearer "+(*access).AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ModuleCancelSuggestion(modID uint, releaseID uint,
	access *models.LoginResponse, t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["id"] = utils.UintToString(modID)
	fields["releaseId"] = utils.UintToString(releaseID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s/%d%s/%d", routes.APIGroup, routes.ReleasesGroup, modID, routes.ReleaseCancelEndpoint, releaseID), body)
	req.Header.Set("Authorization", "Bearer "+(*access).AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func ModuleSuggestedReleaseDelete(modID uint, releaseID uint,
	access *models.LoginResponse, t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["id"] = utils.UintToString(modID)
	fields["releaseId"] = utils.UintToString(releaseID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s/%d%s/%d", routes.APIGroup, routes.ReleasesGroup, modID, routes.ReleaseDeleteEndpoint, releaseID), body)
	req.Header.Set("Authorization", "Bearer "+(*access).AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ModuleReleaseAccept(releaseID uint,
	access *models.LoginResponse, t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["releaseId"] = utils.UintToString(releaseID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s/%d", routes.APIGroup, routes.SuperviseGroup, routes.SuperviseAcceptReleaseEndpoint, releaseID), body)
	req.Header.Set("Authorization", "Bearer "+(*access).AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ModuleReleaseReject(releaseID uint,
	access *models.LoginResponse, t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["releaseId"] = utils.UintToString(releaseID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s/%d", routes.APIGroup, routes.SuperviseGroup, routes.SuperviseRejectReleaseEndpoint, releaseID), body)
	req.Header.Set("Authorization", "Bearer "+(*access).AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func ModuleReleaseCancel(releaseID uint,
	access *models.LoginResponse, t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {

	fields := make(map[string]string)
	fields["releaseId"] = utils.UintToString(releaseID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			assert.Equal(t, false, true)
		}
	}

	if err := writer.Close(); err != nil {
		assert.Equal(t, false, true)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s/%d", routes.APIGroup, routes.SuperviseGroup, routes.SuperviseCancelReleaseEndpoint, releaseID), body)
	req.Header.Set("Authorization", "Bearer "+(*access).AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
