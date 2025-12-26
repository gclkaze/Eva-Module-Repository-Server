// Copyright (c) 2025 Michail Dorgiakis - gclkaze@gmail.com - https://github.com/gclkaze
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
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

func GetModuleInfoMultipleFiles(theFiles []string, t *testing.T) (string, string, string, string, []string) {
	title := "My Super Module To Be Suggested- " + t.Name()
	repr := "my-super-module-to-be-suggested-" + t.Name()
	tags := "super,module,tagged,suggested"
	description := "This is a description of the Suggested Module!" + t.Name()

	var thePaths []string
	for i := range theFiles {
		filePath := path.Join(TestResourcesFolder, theFiles[i])
		thePaths = append(thePaths, filePath)
	}

	return title, repr, tags, description, thePaths
}

func GetNamedModuleInfo(name string, theFile string, t *testing.T) (string, string, string, string, string) {
	title := name + " My Super Module To Be Suggested- " + t.Name()
	repr := name + "-my-super-module-to-be-suggested-" + t.Name()
	tags := "super,module,tagged,suggested"
	description := name + " This is a description of the Suggested Module!" + t.Name()
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
	t.Helper()
	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := UserLogsIn(u, pass, t, TheTestServer.GetRouter())
	_, err = TheTestServer.GetBackend().GetModuleService().GetMaxID()
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
	//assert.Equal(t, modID, maxID+1)

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

// testModuleCreationMultipleFiles
func testModuleCreationMultipleFiles(mr *testmodels.ModuleRequestMultipleFiles, t *testing.T) (uint, *models.LoginResponse) {
	t.Helper()
	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := UserLogsIn(u, pass, t, TheTestServer.GetRouter())
	_, err = TheTestServer.GetBackend().GetModuleService().GetMaxID()
	assert.Equal(t, err == nil, true)

	//upload a module,
	title := mr.Title
	repr := mr.Repr
	tags := mr.Tags
	description := mr.Description
	theFiles := mr.TheFiles
	filePath := mr.FilePath

	rec := ModuleCreateMultipleFiles(title, repr, tags, description, filePath, resp, t, TheTestServer.GetRouter())
	//check response
	assert.Equal(t, rec.Code, http.StatusOK)

	//then check the module in DB
	var respModCreation models.RequestResult[uint]
	err = json.Unmarshal(rec.Body.Bytes(), &respModCreation)
	assert.Equal(t, err == nil, true)

	modID := respModCreation.Value
	//assert.Equal(t, modID, maxID+1)

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

	for i := range filePath {
		assert.Equal(t, utils.FileExists(path.Join(moduleFolder, theFiles[i])), true)

		//assert the tags
		ourTags := strings.Split(mr.Tags, ",")
		AssertTags(ourTags, theMod.Keywords, t)

	}
	return modID, &resp.Value
}

func testModuleCreationMultipleFilesWithReturnValue(mr *testmodels.ModuleRequestMultipleFiles, t *testing.T) *httptest.ResponseRecorder {
	t.Helper()
	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := UserLogsIn(u, pass, t, TheTestServer.GetRouter())
	_, err = TheTestServer.GetBackend().GetModuleService().GetMaxID()
	assert.Equal(t, err == nil, true)

	//upload a module,
	title := mr.Title
	repr := mr.Repr
	tags := mr.Tags
	description := mr.Description
	//	theFiles := mr.TheFiles
	filePath := mr.FilePath

	rec := ModuleCreateMultipleFiles(title, repr, tags, description, filePath, resp, t, TheTestServer.GetRouter())
	return rec
}

func testModuleCreationNoFile(mr *testmodels.ModuleRequest, t *testing.T) (bool, *models.LoginResponse) {
	t.Helper()
	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := UserLogsIn(u, pass, t, TheTestServer.GetRouter())
	_, err = TheTestServer.GetBackend().GetModuleService().GetMaxID()
	assert.Equal(t, err == nil, true)

	title := mr.Title
	repr := mr.Repr
	tags := mr.Tags
	description := mr.Description

	rec := ModuleCreateNoFile(title, repr, tags, description, resp, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusBadRequest)

	var respModCreation models.ErrorResult
	err = json.Unmarshal(rec.Body.Bytes(), &respModCreation)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, respModCreation.Error, "Uploaded file is empty")
	return true, &resp.Value
}

func ReleaseFolderShouldExist(expected bool, theRelease *models.ModuleRelease, t *testing.T) {
	p, err := TheTestServer.GetBackend().GetReleaseService().GetReleaseFolder(theRelease)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, utils.FolderExists(p), expected)
	assert.Equal(t, utils.FileExists(path.Join(p, TheTestServer.GetBackend().GetReleaseService().GetDefaultDistFilename())), expected)
}

func ModuleSuggestedBySimpleUser(t *testing.T) *models.RequestResult[uint] {
	theFile := "assignvalue.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	modID, res := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
	version := "11.1.1"

	rec := ModuleSuggest(modID, version, res, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var respModSuggestCreation models.RequestResult[uint]
	err := json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)

	return &respModSuggestCreation
}

func ModuleSuggestedBySimpleUseWithMultipleVersions(versions []string, t *testing.T) []*models.RequestResult[uint] {
	theFile := "assignvalue.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	modID, res := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)

	var respModSuggestions []*models.RequestResult[uint]

	for i := range versions {
		rec := ModuleSuggest(modID, versions[i], res, t, TheTestServer.GetRouter())
		assert.Equal(t, rec.Code, http.StatusOK)

		var respModSuggestCreation models.RequestResult[uint]
		err := json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
		assert.Equal(t, err == nil, true)

		respModSuggestions = append(respModSuggestions, &respModSuggestCreation)
	}
	return respModSuggestions
}

/*func ModuleSuggestedBySimpleUseWithMultipleVersionsGetAllInfo(versions []string, t *testing.T) (uint, []*models.RequestResult[uint]) {
	theFile := "assignvalue.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	modID, res := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)

	var respModSuggestions []*models.RequestResult[uint]

	for i := range versions {
		rec := ModuleSuggest(modID, versions[i], res, t, TheTestServer.GetRouter())
		assert.Equal(t, rec.Code, http.StatusOK)

		var respModSuggestCreation models.RequestResult[uint]
		err := json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
		assert.Equal(t, err == nil, true)

		respModSuggestions = append(respModSuggestions, &respModSuggestCreation)
	}
	return modID, respModSuggestions
}*/

func ModuleSuggestedBySimpleUserGetAllInfo(t *testing.T) (uint, *models.RequestResult[uint]) {
	theFile := "assignvalue.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	modID, res := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
	version := "11.1.1"

	rec := ModuleSuggest(modID, version, res, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var respModSuggestCreation models.RequestResult[uint]
	err := json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)

	return modID, &respModSuggestCreation
}

func ModuleSuggestedBySimpleUserGetAllInfoWithVersion(version string, t *testing.T) (uint, *models.RequestResult[uint]) {
	theFile := "assignvalue.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	modID, res := testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)

	rec := ModuleSuggest(modID, version, res, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var respModSuggestCreation models.RequestResult[uint]
	err := json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)

	return modID, &respModSuggestCreation
}

func ModuleSuggestedBySimpleUserGetAllInfoWithSpecificVersion(modID uint, res *models.LoginResponse, version string, t *testing.T) (uint, *models.RequestResult[uint]) {

	rec := ModuleSuggest(modID, version, res, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var respModSuggestCreation models.RequestResult[uint]
	err := json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)

	return modID, &respModSuggestCreation
}

func ModuleCreated(t *testing.T) (uint, *models.LoginResponse) {
	theFile := "assignvalue.eva"
	title, repr, tags, description, filePath := GetModuleInfo(theFile, t)

	return testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
}

func NamedModuleCreated(name string, t *testing.T) (uint, *models.LoginResponse) {
	theFile := "assignvalue.eva"
	title, repr, tags, description, filePath := GetNamedModuleInfo(name, theFile, t)

	return testModuleCreation(&testmodels.ModuleRequest{Title: title, Repr: repr, Tags: tags, Description: description, TheFile: theFile, FilePath: filePath}, t)
}

func ModuleVersionsCreated(versions []string, t *testing.T) uint {
	res := AdminLogin(t)
	modID, resp := ModuleCreated(t)

	var releaseIDS []uint
	for i := range versions {

		version := versions[i]
		_, respModSuggestCreation := ModuleSuggestedBySimpleUserGetAllInfoWithSpecificVersion(modID, resp, version, t)

		releaseID := respModSuggestCreation.Value
		releaseIDS = append(releaseIDS, releaseID)

		rec := ModuleReleaseAccept(releaseID, &res.Value, t, TheTestServer.GetRouter())
		assert.Equal(t, rec.Code, http.StatusOK)

		theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
		assert.Equal(t, err == nil, true)
		assert.Equal(t, theRelease.ID == respModSuggestCreation.Value, true)
		assert.Equal(t, theRelease.Status.Label == repositories.Accepted.String(), true)

		ReleaseFolderShouldExist(true, theRelease, t)
	}
	theReleases, err := TheTestServer.GetBackend().GetReleaseService().GetModuleReleases(modID)
	assert.Equal(t, err == nil, true)

	sum := 0
	for i := range releaseIDS {
		for j := range theReleases {
			if theReleases[j].ID == releaseIDS[i] {
				sum += 1

				releaseVersion := theReleases[j].Version
				for k := range versions {
					if versions[k] == releaseVersion {
						sum += 1
						break
					}
				}
				break
			}
		}
	}

	assert.Equal(t, sum, len(releaseIDS)*2)
	return modID
}

func NamedModuleVersionsCreated(name string, versions []string, t *testing.T) uint {
	res := AdminLogin(t)
	modID, resp := NamedModuleCreated(name, t)

	var releaseIDS []uint
	for i := range versions {

		version := versions[i]
		_, respModSuggestCreation := ModuleSuggestedBySimpleUserGetAllInfoWithSpecificVersion(modID, resp, version, t)

		releaseID := respModSuggestCreation.Value
		releaseIDS = append(releaseIDS, releaseID)

		rec := ModuleReleaseAccept(releaseID, &res.Value, t, TheTestServer.GetRouter())
		assert.Equal(t, rec.Code, http.StatusOK)

		theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
		assert.Equal(t, err == nil, true)
		assert.Equal(t, theRelease.ID == respModSuggestCreation.Value, true)
		assert.Equal(t, theRelease.Status.Label == repositories.Accepted.String(), true)

		ReleaseFolderShouldExist(true, theRelease, t)
	}
	theReleases, err := TheTestServer.GetBackend().GetReleaseService().GetModuleReleases(modID)
	assert.Equal(t, err == nil, true)

	sum := 0
	for i := range releaseIDS {
		for j := range theReleases {
			if theReleases[j].ID == releaseIDS[i] {
				sum += 1

				releaseVersion := theReleases[j].Version
				for k := range versions {
					if versions[k] == releaseVersion {
						sum += 1
						break
					}
				}
				break
			}
		}
	}

	assert.Equal(t, sum, len(releaseIDS)*2)
	return modID
}
