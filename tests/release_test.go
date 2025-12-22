package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gclkaze/evamodulerepositoryserver/tests/testmodels"
	"github.com/magiconair/properties/assert"
)

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
}

func TestModuleReleaseAcceptBySimpleUser(t *testing.T) {
	respModSuggestCreation := ModuleSuggestedBySimpleUser(t)

	res := UserLogin(t)
	releaseID := respModSuggestCreation.Value

	rec := ModuleReleaseAccept(releaseID, &res.Value, t, TheTestServer.GetRouter())
	assert.Equal(t, rec.Code, http.StatusUnauthorized)

	ErrorResultsContains("", rec, t)
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

func TestSuggestedReleaseDeletion(t *testing.T) {
}
