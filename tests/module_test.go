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
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

func setup() {
	StartServer()
}

func TestMain(m *testing.M) {
	// 🔧 setup (runs once)
	setup()
	// ▶ run all tests
	code := m.Run()

	// 🧹 teardown (runs once)
	teardown()
	os.Exit(code)
}
func teardown() {
	TeardownServer()
}

func userLogsIn(u *models.UserAccount, pass string, t *testing.T, r *gin.Engine) *models.RequestResult[models.LoginResponse] {
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

func moduleCreate(title string, repr string, tags string, description string, filePath string,
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

	// IMPORTANT: close writer before request
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

func TestModuleCreation(t *testing.T) {
	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := userLogsIn(u, pass, t, TheTestServer.GetRouter())

	//upload a module,
	title := "My Super Module"
	repr := "my-super-module"
	tags := "super,module,tagged"
	description := "This is a description!"

	filePath := "test_resources/assignvalue.eva"

	rec := moduleCreate(title, repr, tags, description, filePath, resp, t, TheTestServer.GetRouter())
	//check response
	assert.Equal(t, rec.Code, http.StatusOK)

	//then check the module in DB
	var respModCreation models.RequestResult[uint]
	err = json.Unmarshal(rec.Body.Bytes(), &respModCreation)
	assert.Equal(t, err == nil, true)

	modID := respModCreation.Value
	assert.Equal(t, modID, 1)

	//check the folder

}

func TestModuleSuggestion(t *testing.T) {
}

func TestModuleReleaseAccept(t *testing.T) {
}

func TestModuleSuggestAndRejection(t *testing.T) {
}

func TestModuleSuggestAndCancel(t *testing.T) {
}

func TestModuleSuggestAcceptReleaseAndCancelRelease(t *testing.T) {
}

func TestModuleDeletion(t *testing.T) {
}
