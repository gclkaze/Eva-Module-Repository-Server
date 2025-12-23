package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	assert.Equal(t, utils.FileExists(distFile), false)

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
	assert.Equal(t, utils.FileExists(distFile), false)
}

//module with multiple releases -> DELETE :) By Simple User
//module with multiple releases -> DELETE :) By Admin
