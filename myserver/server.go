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
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type User struct {
	Username string
	Password []byte // パスワードはハッシュ化して保存
}

// ユーザー情報を仮想的に保存するデータベース。これをredisやMySQLなどに置き換える
var users = map[string]User{
	"user1": {Username: "user1", Password: hashPassword("password1")},
	"user2": {Username: "user2", Password: hashPassword("password2")},
}

type Comment struct {
	Date    time.Time
	Message string
}

// 投稿一覧を仮想的に保存するデータベース
var comments = []Comment{
	{
		Date:    time.Date(2023, 8, 1, 10, 24, 59, 0, time.UTC),
		Message: "こんにちは",
	},
	{
		Date:    time.Date(2023, 8, 3, 13, 20, 59, 0, time.UTC),
		Message: "おはよう",
	},
	{
		Date:    time.Date(2023, 8, 5, 23, 34, 35, 0, time.UTC),
		Message: "おやすみ",
	},
}

// パスワードをハッシュ化する関数
func hashPassword(password string) []byte {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.Static("/static", "static")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(middleware.Logger())

	// Routing
	e.GET("/login", getLogin)
	e.POST("/login", postLogin)
	e.GET("/comments", getComments)
	e.POST("/comments", postComments)
	e.GET("/logout", logout)

	e.Logger.Fatal(e.Start(":1323"))
}

func getLogin(c echo.Context) error {
	session, _ := session.Get("session", c)
	username := session.Values["username"]
	if username != nil {
		return c.Render(http.StatusOK, "comments", "")
	}

	return c.Render(http.StatusOK, "login", "")
}

func postLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, ok := users[username]
	if !ok || bcrypt.CompareHashAndPassword(user.Password, []byte(password)) != nil {
		return echo.ErrUnauthorized
	}

	session, _ := session.Get("session", c)
	session.Values["username"] = username
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return c.Render(http.StatusOK, "comments", session)
}

func getComments(c echo.Context) error {
	session, _ := session.Get("session", c)
	username := session.Values["username"]
	if username == nil {
		return c.Render(http.StatusOK, "login", "")
	}
	return c.Render(http.StatusOK, "comments", comments)
}

func postComments(c echo.Context) error {
	m := c.FormValue("message")
	newComment := Comment{
		Date:    time.Now(),
		Message: m,
	}
	comments = append(comments, newComment)
	return c.Render(http.StatusOK, "comments", comments)
}

func logout(c echo.Context) error {
	session, _ := session.Get("session", c)
	session.Options.MaxAge = -1
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return c.Render(http.StatusOK, "login", "")
}
