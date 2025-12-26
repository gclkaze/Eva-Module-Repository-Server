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
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/magiconair/properties/assert"
)

func TestModuleReleaseAcceptThatDoesNotExist(t *testing.T) {
	res := AdminLogin(t)
	releaseID := utils.GetRandomUintNumber(700)

	rec := ModuleReleaseAccept(releaseID, &res.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusInternalServerError)
	ErrorResultsContains(fmt.Sprintf("couldn't find release with id %d", releaseID), rec, t)
}

func TestModuleReleaseAccept(t *testing.T) {
	respModSuggestCreation := ModuleSuggestedBySimpleUser(t)

	res := AdminLogin(t)
	releaseID := respModSuggestCreation.Value

	rec := ModuleReleaseAccept(releaseID, &res.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease.ID == respModSuggestCreation.Value, true)
	assert.Equal(t, theRelease.Status.Label == repositories.Accepted.String(), true)

	ReleaseFolderShouldExist(true, theRelease, t)
}

func TestModuleReleaseAcceptBySimpleUser(t *testing.T) {
	respModSuggestCreation := ModuleSuggestedBySimpleUser(t)

	res := UserLogin(t)
	releaseID := respModSuggestCreation.Value

	rec := ModuleReleaseAccept(releaseID, &res.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusUnauthorized)

	ErrorResultsContains("", rec, t)
	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease.ID == respModSuggestCreation.Value, true)
	assert.Equal(t, theRelease.Status.Label == repositories.Pending.String(), true)
}

func TestModuleSuggestAndRejection(t *testing.T) {
	respModSuggestCreation := ModuleSuggestedBySimpleUser(t)

	res := AdminLogin(t)
	releaseID := respModSuggestCreation.Value

	rec := ModuleReleaseReject(releaseID, &res.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease.ID == respModSuggestCreation.Value, true)
	assert.Equal(t, theRelease.Status.Label == repositories.Rejected.String(), true)

	ReleaseFolderShouldExist(false, theRelease, t)
}

func TestModuleSuggestAndCancel(t *testing.T) {
	respModSuggestCreation := ModuleSuggestedBySimpleUser(t)

	res := AdminLogin(t)
	releaseID := respModSuggestCreation.Value

	rec := ModuleReleaseCancel(releaseID, &res.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusInternalServerError)

	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease.ID == respModSuggestCreation.Value, true)
	assert.Equal(t, theRelease.Status.Label == repositories.Pending.String(), true)

	ReleaseFolderShouldExist(false, theRelease, t)
}

func TestModuleSuggestAcceptReleaseAndCancelRelease(t *testing.T) {
	respModSuggestCreation := ModuleSuggestedBySimpleUser(t)

	admin := AdminLogin(t)
	releaseID := respModSuggestCreation.Value

	rec := ModuleReleaseAccept(releaseID, &admin.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease.ID == respModSuggestCreation.Value, true)
	assert.Equal(t, theRelease.Status.Label == repositories.Accepted.String(), true)

	ReleaseFolderShouldExist(true, theRelease, t)

	rec = ModuleReleaseCancel(releaseID, &admin.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)
	err = json.Unmarshal(rec.Body.Bytes(), &respModSuggestCreation)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, respModSuggestCreation.Value, releaseID)

	theRelease, err = TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease.ID == respModSuggestCreation.Value, true)
	assert.Equal(t, theRelease.Status.Label == repositories.Canceled.String(), true)

	ReleaseFolderShouldExist(false, theRelease, t)
}

func TestSuggestedModuleReleaseDeletion(t *testing.T) {
	user := UserLogin(t)
	modID, respModSuggestCreation := ModuleSuggestedBySimpleUserGetAllInfo(t)

	rec := ModuleSuggestedReleaseDelete(modID, respModSuggestCreation.Value, &user.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var response models.RequestResult[bool]
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, response.Value && response.Result, true)

	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease == nil, true)
}

func TestAcceptedReleaseDeletion(t *testing.T) {
	modID, respModSuggestCreation := ModuleSuggestedBySimpleUserGetAllInfo(t)
	admin := AdminLogin(t)

	//Accepted from Admin
	rec := ModuleReleaseAccept(respModSuggestCreation.Value, &admin.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)
	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	ReleaseFolderShouldExist(true, theRelease, t)

	releasePath, _ := TheTestServer.GetBackend().GetReleaseService().GetReleaseFolder(theRelease)
	distFile := TheTestServer.GetBackend().GetReleaseService().GetDefaultDistFilename()

	//Delete by the Admin
	rec = ModuleSuggestedReleaseDelete(modID, respModSuggestCreation.Value, &admin.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var response models.RequestResult[bool]
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, response.Value && response.Result, true)

	theRelease, err = TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease == nil, true)

	assert.Equal(t, utils.FolderExists(releasePath), false)
	//	assert.Equal(t, utils.FileExists(distFile), false)
	assert.Equal(t, utils.FileExists(path.Join(releasePath, distFile)), false)

}
func TestAcceptedReleaseDeletionBySimpleUser(t *testing.T) {
	user := UserLogin(t)
	modID, respModSuggestCreation := ModuleSuggestedBySimpleUserGetAllInfo(t)
	admin := AdminLogin(t)

	//Accepted from Admin
	rec := ModuleReleaseAccept(respModSuggestCreation.Value, &admin.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)
	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	ReleaseFolderShouldExist(true, theRelease, t)

	releasePath, _ := TheTestServer.GetBackend().GetReleaseService().GetReleaseFolder(theRelease)
	distFile := TheTestServer.GetBackend().GetReleaseService().GetDefaultDistFilename()

	//Delete by the Creator
	rec = ModuleSuggestedReleaseDelete(modID, respModSuggestCreation.Value, &user.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var response models.RequestResult[bool]
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, response.Value && response.Result, true)

	theRelease, err = TheTestServer.GetBackend().GetReleaseService().GetRelease(respModSuggestCreation.Value)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease == nil, true)

	assert.Equal(t, utils.FolderExists(releasePath), false)
	//	assert.Equal(t, utils.FileExists(distFile), false)
	assert.Equal(t, utils.FileExists(path.Join(releasePath, distFile)), false)
}

func TestModuleReleaseAcceptMultipleVersions(t *testing.T) {
	versions := []string{
		"1.0.0",
		"2.0.0",
		"3.0.0",
	}
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
}

func TestModuleReleaseAcceptSameVersions(t *testing.T) {
	versions := []string{
		"1.0.0",
		"1.0.0",
	}
	t.Helper()
	res := AdminLogin(t)
	modID, resp := ModuleCreated(t)

	var releaseIDS []uint
	for i := range versions {

		version := versions[i]
		_, respModSuggestCreation := ModuleSuggestedBySimpleUserGetAllInfoWithSpecificVersion(modID, resp, version, t)

		releaseID := respModSuggestCreation.Value
		rec := ModuleReleaseAccept(releaseID, &res.Value, t, TheTestServer.GetRouter())
		releaseIDS = append(releaseIDS, releaseID)
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

	assert.Equal(t, sum, 2)
}

// module with multiple releases -> DELETE the Module :) By Simple User
func TestModuleDeletionOnMultipleVersions(t *testing.T) {
	versions := []string{
		"1.0.0",
		"2.0.0",
		"3.3.3",
	}
	modID := ModuleVersionsCreated(versions, t)
	theModule, err := TheTestServer.GetBackend().GetModuleService().GetModule(modID)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithName(theModule.Repr)), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[0])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[1])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[2])), true)

	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := UserLogsIn(u, pass, t, TheTestServer.GetRouter())

	rec := ModuleDelete(modID, &resp.Value, t, TheTestServer.GetRouter())
	var respModDeletion models.RequestResult[bool]
	err = json.Unmarshal(rec.Body.Bytes(), &respModDeletion)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, respModDeletion.Value, true)
	assert.Equal(t, respModDeletion.Result, true)

	assert.Equal(t, err == nil, true)
	moduleFolder := TheTestServer.GetDeveloperModuleFolder(u.ID, modID)
	assert.Equal(t, utils.FolderExists(moduleFolder), false)
	assert.Equal(t, utils.FileExists(TheTestServer.GetReleaseBasePathWithName(theModule.Repr)), false)
}

// module with multiple releases -> DELETE the Module :) By Admin
func TestModuleDeletionOnMultipleVersionsByAdmin(t *testing.T) {
	versions := []string{
		"1.0.0",
		"2.0.0",
		"3.3.3",
	}
	modID := ModuleVersionsCreated(versions, t)
	theModule, err := TheTestServer.GetBackend().GetModuleService().GetModule(modID)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithName(theModule.Repr)), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[0])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[1])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[2])), true)

	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := AdminLogin(t)

	rec := ModuleDelete(modID, &resp.Value, t, TheTestServer.GetRouter())
	var respModDeletion models.RequestResult[bool]
	err = json.Unmarshal(rec.Body.Bytes(), &respModDeletion)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, respModDeletion.Value, true)
	assert.Equal(t, respModDeletion.Result, true)

	assert.Equal(t, err == nil, true)
	moduleFolder := TheTestServer.GetDeveloperModuleFolder(u.ID, modID)
	assert.Equal(t, utils.FolderExists(moduleFolder), false)
	assert.Equal(t, utils.FileExists(TheTestServer.GetReleaseBasePathWithName(theModule.Repr)), false)
}

// module with multiple releases -> DELETE a single release
func TestModuleSingleReleaseDeletion(t *testing.T) {
	versions := []string{
		"1.0.0",
		"2.0.0",
		"3.3.3",
	}
	modID := ModuleVersionsCreated(versions, t)
	theModule, err := TheTestServer.GetBackend().GetModuleService().GetModule(modID)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithName(theModule.Repr)), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[0])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[1])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[2])), true)

	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := UserLogsIn(u, pass, t, TheTestServer.GetRouter())

	theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetModuleReleaseByVersion(modID, versions[1])
	assert.Equal(t, err == nil, true)

	releasePath, _ := TheTestServer.GetBackend().GetReleaseService().GetReleaseFolder(theRelease)
	distFile := TheTestServer.GetBackend().GetReleaseService().GetDefaultDistFilename()

	assert.Equal(t, utils.FolderExists(releasePath), true)
	assert.Equal(t, utils.FileExists(path.Join(releasePath, distFile)), true)

	rec := ModuleSuggestedReleaseDelete(modID, theRelease.ID, &resp.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusOK)

	var response models.RequestResult[bool]
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, response.Value && response.Result, true)

	theRelease, err = TheTestServer.GetBackend().GetReleaseService().GetRelease(theRelease.ID)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, theRelease == nil, true)

	assert.Equal(t, utils.FolderExists(releasePath), false)
	assert.Equal(t, utils.FileExists(path.Join(releasePath, distFile)), false)
}

// module with multiple releases -> DELETE a single release that does not exist :)
func TestModuleSingleReleaseDeletionWithUnknownVersion(t *testing.T) {
	versions := []string{
		"1.0.0",
		"2.0.0",
		"3.3.3",
	}
	modID := ModuleVersionsCreated(versions, t)
	theModule, err := TheTestServer.GetBackend().GetModuleService().GetModule(modID)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithName(theModule.Repr)), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[0])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[1])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theModule.Repr, versions[2])), true)

	pass := GetDefaultUserPassword()
	u, err := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.User.String())
	assert.Equal(t, err == nil, true)
	resp := UserLogsIn(u, pass, t, TheTestServer.GetRouter())

	assert.Equal(t, err == nil, true)
	distFile := TheTestServer.GetBackend().GetReleaseService().GetDefaultDistFilename()
	rec := ModuleSuggestedReleaseDelete(modID, 10001, &resp.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusInternalServerError)

	var response models.ErrorResult
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, response.Details == "unknown release 10001", true)

	for i := range versions {
		theRelease, err := TheTestServer.GetBackend().GetReleaseService().GetModuleReleaseByVersion(modID, versions[i])
		releasePath, _ := TheTestServer.GetBackend().GetReleaseService().GetReleaseFolder(theRelease)
		assert.Equal(t, err == nil, true)
		assert.Equal(t, theRelease != nil, true)

		assert.Equal(t, utils.FolderExists(releasePath), true)
		assert.Equal(t, utils.FileExists(path.Join(releasePath, distFile)), true)

	}
}

// suggest a module with a given version -> accept it -> suggest another one with the given version
func TestMultipleModulesSameVersions(t *testing.T) {
	versions := []string{
		"1.0.0",
		"2.0.0",
		"3.3.3",
	}
	firstModID := NamedModuleVersionsCreated("first", versions, t)
	theFirstModule, err := TheTestServer.GetBackend().GetModuleService().GetModule(firstModID)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithName(theFirstModule.Repr)), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theFirstModule.Repr, versions[0])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theFirstModule.Repr, versions[1])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theFirstModule.Repr, versions[2])), true)

	secondModID := NamedModuleVersionsCreated("second", versions, t)
	theSecondModule, err := TheTestServer.GetBackend().GetModuleService().GetModule(secondModID)
	assert.Equal(t, err == nil, true)

	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithName(theFirstModule.Repr)), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theSecondModule.Repr, versions[0])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theSecondModule.Repr, versions[1])), true)
	assert.Equal(t, utils.FolderExists(TheTestServer.GetReleaseBasePathWithNameAndVersion(theSecondModule.Repr, versions[2])), true)
}
