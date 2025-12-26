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

package utils

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gin-gonic/gin"
)

func Ok[T any](value T) models.RequestResult[T] {
	return models.RequestResult[T]{
		Result: true,
		Value:  value,
	}
}

func OkWithMessage[T any](value T, msg string) models.RequestResult[T] {
	return models.RequestResult[T]{
		Result:  true,
		Value:   value,
		Message: msg,
	}
}

func SimpleOkMessage(msg string) models.RequestResult[models.EmptyRequestResult] {
	return models.RequestResult[models.EmptyRequestResult]{
		Result:  true,
		Message: msg,
	}
}

func Err(err error, msg string) models.ErrorResult {
	resp := models.ErrorResult{
		Result: false,
		Error:  msg,
	}

	if gin.Mode() == gin.DebugMode && err != nil {
		resp.Details = err.Error()
	}

	return resp
}
