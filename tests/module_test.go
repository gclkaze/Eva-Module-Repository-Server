package tests

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gclkaze/evamodulerepositoryserver/tests/testmodels"
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

func assertTags(ourTags []string, keywords []models.Keyword, t *testing.T) {
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
	assertTags(ourTags, theMod.Keywords, t)

	return modID, &resp.Value
}

func TestModuleCreation(t *testing.T) {
	title := "My Super Module"
	repr := "my-super-module"
	tags := "super,module,tagged"
	description := "This is a description!"
	theFile := "assignvalue.eva"
	filePath := "test_resources/" + theFile

	testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
}

func TestModuleSuggestion(t *testing.T) {
	title := "My Super Module To Be Suggested"
	repr := "my-super-module-to-be-suggested"
	tags := "super,module,tagged,suggested"
	description := "This is a description of the Suggested Module!"
	theFile := "expr_length.eva"
	filePath := "test_resources/" + theFile
	modID, resp := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)

	myVersion := "1.0.1"
	ModuleSuggest(modID, myVersion, resp, t, TheTestServer.GetRouter())

	//this should create a new release
}

func TestModuleSuggestionWithUnknownModule(t *testing.T) {

}

func TestMultipleModuleSuggestion(t *testing.T) {
	title := "My Super Module To Be Suggested"
	repr := "my-super-module-to-be-suggested"
	tags := "super,module,tagged,suggested"
	description := "This is a description of the Suggested Module!"
	theFile := "expr_length.eva"
	filePath := "test_resources/" + theFile
	modID, resp := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)

	myVersion := "1.0.1"

	maxReleaseID, err := TheTestServer.GetBackend().GetReleaseService().GetMaxID()
	assert.Equal(t, err == nil, true)
	rec := ModuleSuggest(modID, myVersion, resp, t, TheTestServer.GetRouter())

	//this should create a new release
	assert.Equal(t, rec.Code, http.StatusOK)

	//then check the module in DB
	var respModSuggestCreation models.RequestResult[uint]
	err = json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, respModSuggestCreation.Value, maxReleaseID+1)

	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, theRelease.Version, myVersion)
	assert.Equal(t, theRelease.Description, description)
	assert.Equal(t, theRelease.ModuleID, modID)

	//check the status
	st, err := TheTestServer.GetBackend().GetReleaseService().GetStatus(repositories.Pending)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease.StatusID, st.ID)
}

func TestUserWithSuggestionCannotMakeNewSuggestion(t *testing.T) {
	title := "My Super Module To Be Suggested"
	repr := "my-super-module-to-be-suggested"
	tags := "super,module,tagged,suggested"
	description := "This is a description of the Suggested Module!"
	theFile := "expr_length.eva"
	filePath := "test_resources/" + theFile
	modID, resp := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)

	myVersion := "1.0.1"

	maxReleaseID, err := TheTestServer.GetBackend().GetReleaseService().GetMaxID()
	assert.Equal(t, err == nil, true)
	rec := ModuleSuggest(modID, myVersion, resp, t, TheTestServer.GetRouter())

	//this should create a new release
	assert.Equal(t, rec.Code, http.StatusOK)

	//then check the module in DB
	var respModSuggestCreation models.RequestResult[uint]
	err = json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, respModSuggestCreation.Value, maxReleaseID+1)

	theOtherVersion := "2.0.1"
	rec = ModuleSuggest(modID, theOtherVersion, resp, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusInternalServerError)

	var errorResult models.ErrorResult
	err = json.Unmarshal(rec.Body.Bytes(), &errorResult)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, strings.Contains(errorResult.Details, "there are is a pending release of that module, need to cancel, reject, accept it first to create a new release"), true)
}

func TestCancelModuleSuggestion(t *testing.T) {

}

func TestCancelModuleSuggestionWithUnknownModule(t *testing.T) {

}

func TestModuleReleaseAccept(t *testing.T) {
}

func TestModuleSuggestAndRejection(t *testing.T) {
}

func TestModuleSuggestAndCancel(t *testing.T) {
}

func TestModuleSuggestAcceptReleaseAndCancelRelease(t *testing.T) {
}

func TestModuleDeletionWithoutBeingSuggested(t *testing.T) {
}
func TestSuggestedModuleDeletion(t *testing.T) {
}
