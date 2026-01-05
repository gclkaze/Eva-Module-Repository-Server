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

package dto

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
)

type ModuleDTO struct {
	ID          uint     `json:"id" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Repr        string   `json:"repr" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
}

func NewModuleDTO(mod models.Module) *ModuleDTO {
	var tags []string
	for i := range mod.Keywords {
		tags = append(tags, mod.Keywords[i].Label)
	}
	return &ModuleDTO{ID: mod.ID, Title: mod.Title, Repr: mod.Repr, Description: mod.Description, Tags: tags}
}
