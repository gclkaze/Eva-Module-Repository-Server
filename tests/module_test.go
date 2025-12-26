package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gclkaze/evamodulerepositoryserver/tests/testmodels"
	"github.com/magiconair/properties/assert"
)

func TestMain(m *testing.M) {
	// 🔧 setup (runs once)
	Setup()
	// ▶ run all tests
	code := m.Run()

	// 🧹 teardown (runs once)
	Teardown()
	os.Exit(code)
}

func TestModuleCreation(t *testing.T) {
	theFile := "assignvalue.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
}

func TestModuleSuggestion(t *testing.T) {
	theFile := "expr_length.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

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

func TestModuleSuggestionWithUnknownModule(t *testing.T) {
	myVersion := "1.0.1"
	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := UserLogsIn(u, pass, t, TheTestServer.GetRouter())
	rec := ModuleSuggest(666, myVersion, &resp.Value, t, TheTestServer.GetRouter())

	assert.Equal(t, rec.Code, http.StatusInternalServerError)
	ErrorResultsContains(fmt.Sprintf("couldn't find module with id %d that was suggested by user %d", 666, u.ID), rec, t)
}

func TestMultipleModuleSuggestion(t *testing.T) {
	theFile := "expr_length.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	modID, resp := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)

	myVersion := "1.0.1"

	_, err := TheTestServer.GetBackend().GetReleaseService().GetMaxID()
	assert.Equal(t, err == nil, true)
	rec := ModuleSuggest(modID, myVersion, resp, t, TheTestServer.GetRouter())

	assert.Equal(t, rec.Code, http.StatusOK)

	rec = ModuleSuggest(modID, myVersion, resp, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusInternalServerError)
	ErrorResultsContains("there are is a pending release of that module, need to cancel, reject, accept it first to create a new release", rec, t)
}

func TestUserWithSuggestionCannotMakeNewSuggestion(t *testing.T) {
	theFile := "expr_length.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

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
	ErrorResultsContains("there are is a pending release of that module, need to cancel, reject, accept it first to create a new release", rec, t)
}

func TestModuleDeletionWithoutBeingSuggested(t *testing.T) {
	theFile := "expr_length.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)
	modID, resp := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)

	rec := ModuleDelete(modID, resp, t, TheTestServer.GetRouter())
	var respModDeletion models.RequestResult[bool]
	err := json.Unmarshal(rec.Body.Bytes(), &respModDeletion)

	assert.Equal(t, err == nil, true)
	assert.Equal(t, respModDeletion.Value, true)
	assert.Equal(t, respModDeletion.Result, true)

	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	moduleFolder := TheTestServer.GetDeveloperModuleFolder(u.ID, modID)
	assert.Equal(t, utils.FolderExists(moduleFolder), false)
	assert.Equal(t, utils.FileExists(path.Join(moduleFolder, theFile)), false)
}

func TestCancelModuleSuggestion(t *testing.T) {
	theFile := "expr_length.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)
	t.Helper()

	modID, resp := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
	myVersion := "21.0.661"

	maxReleaseID, err := TheTestServer.GetBackend().GetReleaseService().GetCount()
	assert.Equal(t, err == nil, true)
	rec := ModuleSuggest(modID, myVersion, resp, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var respModSuggestCreation models.RequestResult[uint]
	err = json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, respModSuggestCreation.Value > uint(maxReleaseID), true)
	releaseID := respModSuggestCreation.Value

	rec = ModuleCancelSuggestion(modID, releaseID, resp, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	BooleanResultAssert(true, fmt.Sprintf("Release %d was cancelled successfully", releaseID), rec, t)
}

func TestCancelModuleSuggestionWithUnknownModule(t *testing.T) {
	theFile := "expr_length.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	modID, resp := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
	myVersion := "1.0.1"

	count, err := TheTestServer.GetBackend().GetReleaseService().GetCount()
	assert.Equal(t, err == nil, true)
	rec := ModuleSuggest(modID, myVersion, resp, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var respModSuggestCreation models.RequestResult[uint]
	err = json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)

	fmt.Printf("Value %d maxRelease ID + 1 %d", respModSuggestCreation.Value, count+1)
	newCount, err := TheTestServer.GetBackend().GetReleaseService().GetCount()
	assert.Equal(t, err == nil, true)
	assert.Equal(t, newCount == count+1, true)
	releaseID := respModSuggestCreation.Value

	modID = utils.GetRandomUintNumber(666)
	rec = ModuleCancelSuggestion(modID, releaseID, resp, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusInternalServerError)

	ErrorResultsContains(fmt.Sprintf("couldn't find suggested release for mod %d with id %d ", releaseID, modID), rec, t)
}

func TestCancelModuleSuggestionWithUnknownRelease(t *testing.T) {
	theFile := "expr_length.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	modID, resp := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
	myVersion := "1.0.1"

	maxReleaseID, err := TheTestServer.GetBackend().GetReleaseService().GetMaxID()
	assert.Equal(t, err == nil, true)
	rec := ModuleSuggest(modID, myVersion, resp, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var respModSuggestCreation models.RequestResult[uint]
	err = json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, respModSuggestCreation.Value, maxReleaseID+1)

	releaseID := utils.GetRandomUintNumber(666)
	rec = ModuleCancelSuggestion(modID, releaseID, resp, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusInternalServerError)

	ErrorResultsContains(fmt.Sprintf("couldn't find suggested release for mod %d with id %d ", releaseID, modID), rec, t)
}

// upload zero files
func TestModuleCreationNoFile(t *testing.T) {
	theFile := "expr_length.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	res, _ := testModuleCreationNoFile(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
	assert.Equal(t, res, true)
}

func TestModuleCreationMultipleFiles(t *testing.T) {
	theFiles := []string{
		"expr_length.eva",
		"assignvalue.eva",
	}
	title, repr, tags, description, filePath := GetModuleInfoMultipleFiles(theFiles, t)
	res, _ := testModuleCreationMultipleFiles(&testmodels.ModuleRequestMultipleFiles{Title: title, Repr: repr, Tags: tags, Description: description, TheFiles: theFiles, FilePath: filePath}, t)
	assert.Equal(t, res > 0, true)
}

// upload big files
func TestModuleCreationMultipleBigFiles(t *testing.T) {
	theFiles := []string{
		"expr_length.eva",
		"assignvalue.eva",
	}
	title, repr, tags, description, filePath := GetModuleInfoMultipleFiles(theFiles, t)

	old := ResetUploadLimit(2)
	res := testModuleCreationMultipleFilesWithReturnValue(&testmodels.ModuleRequestMultipleFiles{Title: title, Repr: repr, Tags: tags, Description: description, TheFiles: theFiles, FilePath: filePath}, t)
	assert.Equal(t, res != nil, true)
	assert.Equal(t, res.Result().StatusCode, http.StatusRequestEntityTooLarge)
	ResetUploadLimit(old)
}
