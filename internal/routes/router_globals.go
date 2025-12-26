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

// Package routes contains the Server's endpoints declarations
package routes

const (
	APIGroup         = "/api"
	AuthGroup        = "/auth"
	LoginEndpoint    = "/login"
	RefreshEndpoint  = "/refresh"
	RegisterEndpoint = "/register"

	ModulesGroup          = "/modules"
	ModuleSearchEndpoint  = "/search"
	ModuleDeleteEndpoint  = "/delete"
	ModuleUploadEndpoint  = "/upload"
	ModuleUpdateEndpoint  = "/update"
	ModuleSuggestEndpoint = "/suggest"

	ReleasesGroup         = "/releases"
	ReleaseEndpoint       = "release"
	ReleaseSearchEndpoint = "/search"
	ReleaseDeleteEndpoint = "/delete"
	ReleaseCancelEndpoint = "/cancel"

	DownloadGroup           = "/download"
	DownloadReleaseEndpoint = "/release"

	SuperviseGroup                      = "/supervise"
	SuperviseDownloadAnyReleaseEndpoint = "/download/release"
	SuperviseRejectReleaseEndpoint      = "/reject/release"
	SuperviseAcceptReleaseEndpoint      = "/accept/release"
	SuperviseCancelReleaseEndpoint      = "/cancel/release"
	SupervisePendingReleaseEndpoint     = "/pending/release"
	SuperviseBanUserEndpoint            = "/ban"
	SuperviseUnbanUserEndpoint          = "/unban"
)
