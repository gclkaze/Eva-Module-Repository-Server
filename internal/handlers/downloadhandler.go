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

package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	service *services.DownloadService
}

func NewDownloadHandler(service *services.DownloadService) *DownloadHandler {
	return &DownloadHandler{service: service}
}

func (h *DownloadHandler) DownloadRelease(c *gin.Context) {
	//the release needs to be ACCEPTED
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	dest, filename, err := h.service.DownloadRelease(releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, err.Error()))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(dest)
}

func (h DownloadHandler) DownloadAnyRelease(c *gin.Context) {
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	dest, filename, err := h.service.DownloadAnyRelease(releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, err.Error()))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(dest)
}

func (h *DownloadHandler) AuthUserDownloadSpecificRelease(c *gin.Context) {
	release := c.Param("release")
	userID := c.GetUint("userId")

	isItTmp, dest, filename, err := h.service.AuthUserDownloadSpecificRelease(userID, release)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, err.Error()))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	if isItTmp {
		c.File(filepath.Join(dest, filename))
		// Run cleanup after response is completed
		c.Writer.Flush()
		if notifier, ok := c.Writer.(http.CloseNotifier); ok {
			// CloseNotifier is deprecated but still works in many setups;
			// better approach below without it.
			_ = notifier
		}
		go func(dir string) {
			// Best-effort delay to avoid removing while still streaming
			// (better: move tmp cleanup into service with deferred close)
			time.Sleep(2 * time.Second)
			_ = os.RemoveAll(dir)
		}(dest)
	} else {
		c.File(dest)
	}

}

func (h *DownloadHandler) DownloadPublicRelease(c *gin.Context) {
	release := c.Param("release")
	//no token attached, thus the release needs to be on ACCEPT STATUS
	dest, filename, err := h.service.DownloadPublicRelease(release)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, err.Error()))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(dest)
}
