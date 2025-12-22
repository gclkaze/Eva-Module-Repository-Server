package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gclkaze/evamodulerepositoryserver/tests/testmodels"
	"github.com/magiconair/properties/assert"
)

const TestResourcesFolder = "test_resources"

func Setup() {
	StartServer()
}

func Teardown() {
	TeardownServer()
}

func GetModuleInfo(theFile string, t *testing.T) (string, string, string, string, string) {
	title := "My Super Module To Be Suggested- " + t.Name()
	repr := "my-super-module-to-be-suggested-" + t.Name()
	tags := "super,module,tagged,suggested"
	description := "This is a description of the Suggested Module!" + t.Name()
	filePath := path.Join(TestResourcesFolder, theFile)

	return title, repr, tags, description, filePath
}
func AssertTags(ourTags []string, keywords []models.Keyword, t *testing.T) {
	assert.Equal(t, len(ourTags), len(keywords))
	sum := 0
	for i := range ourTags {
		for j := range keywords {
			if keywords[j].Label == ourTags[i] {
				sum += 1
			}
		}
	}
	assert.Equal(t, sum, len(ourTags))
}

func ErrorResultsContains(str string, rec *httptest.ResponseRecorder, t *testing.T) {
	var errorResult models.ErrorResult
	err := json.Unmarshal(rec.Body.Bytes(), &errorResult)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, strings.Contains(errorResult.Details, str), true)
}

func BooleanResultAssert(res bool, msg string, rec *httptest.ResponseRecorder, t *testing.T) {
	var reqResult models.RequestResult[bool]
	err := json.Unmarshal(rec.Body.Bytes(), &reqResult)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, reqResult.Result == res, true)
	assert.Equal(t, reqResult.Message == msg, true)
}

func AdminLogin(t *testing.T) *models.RequestResult[models.LoginResponse] {
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.Admin.String())
	assert.Equal(t, err == nil, true)
	pass := GetDefaultUserPassword()
	return UserLogsIn(u, pass, t, TheTestServer.GetRouter())
}

func UserLogin(t *testing.T) *models.RequestResult[models.LoginResponse] {
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	pass := GetDefaultUserPassword()
	return UserLogsIn(u, pass, t, TheTestServer.GetRouter())
}

func testModuleCreation(mr *testmodels.ModuleRequest, t *testing.T) (uint, *models.LoginResponse) {
	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := UserLogsIn(u, pass, t, TheTestServer.GetRouter())

	maxID, err := TheTestServer.GetBackend().GetModuleService().GetMaxID()
	assert.Equal(t, err == nil, true)

	//upload a module,
	title := mr.Title
	repr := mr.Repr
	tags := mr.Tags
	description := mr.Description
	theFile := mr.TheFile
	filePath := mr.FilePath

	rec := ModuleCreate(title, repr, tags, description, filePath, resp, t, TheTestServer.GetRouter())
	//check response
	assert.Equal(t, rec.Code, http.StatusOK)

	//then check the module in DB
	var respModCreation models.RequestResult[uint]
	err = json.Unmarshal(rec.Body.Bytes(), &respModCreation)
	assert.Equal(t, err == nil, true)

	modID := respModCreation.Value
	assert.Equal(t, modID, maxID+1)

	//check the module data
	theMod, err := TheTestServer.GetBackend().GetModuleService().GetModule(modID)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theMod.Title, title)
	assert.Equal(t, theMod.Repr, repr)
	assert.Equal(t, theMod.ID, modID)
	assert.Equal(t, theMod.Description, description)

	//check the folder
	//tests\modules\developers\2\1
	moduleFolder := TheTestServer.GetDeveloperModuleFolder(u.ID, modID)
	assert.Equal(t, utils.FolderExists(moduleFolder), true)
	assert.Equal(t, utils.FileExists(path.Join(moduleFolder, theFile)), true)

	//assert the tags
	ourTags := strings.Split(mr.Tags, ",")
	AssertTags(ourTags, theMod.Keywords, t)

	return modID, &resp.Value
}
