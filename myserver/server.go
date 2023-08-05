// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hello is a simple hello, world demonstration web server.
//
// It serves version information on /version and answers
// any other request like /name by saying "Hello, name!".
//
// See golang.org/x/example/outyet for a more sophisticated server.
package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.Static("/static", "static")
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/login", getLogin)
	e.GET("/comments", getComments)
	e.POST("/comments", postComments)

	e.Logger.Fatal(e.Start(":1323"))
}

func getLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", "test")
}

func getComments(c echo.Context) error {
	return c.Render(http.StatusOK, "comments", "test")
}

func postComments(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{})
}
