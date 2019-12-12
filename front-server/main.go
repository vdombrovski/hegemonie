// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/pongo2"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LoginForm struct {
	UserId string `form:"userid" binding:"Required"`
	Passwd string `form:"passwd" binding:"Required"`
}

func routeForms(m *macaron.Macaron) {
	m.Post("/action/login", binding.Bind(LoginForm{}),
		func(ctx *macaron.Context, s session.Store, info LoginForm) {
			ctx.SetSecureCookie("session", "NYI")
			s.Set("userid", info.UserId)
			ctx.Redirect("/game/land")
		})
	m.Post("/action/logout",
		func(ctx *macaron.Context, s session.Store) {
			ctx.SetSecureCookie("session", "")
			s.Flush()
			ctx.Redirect("/index.html")
		})
}

func routePages(m *macaron.Macaron) {
	m.Get("/",
		func(ctx *macaron.Context, f *session.Flash) {
			if f.Get("userid") == "" {
				ctx.Redirect("/index.html")
			} else {
				ctx.Redirect("/game/land")
			}
		})
	m.Get("/index.html",
		func(ctx *macaron.Context) {
			ctx.HTML(200, "index")
		})
	m.Get("/game/user",
		func(ctx *macaron.Context, s session.Store) {
			ctx.Data["userid"] = s.Get("userid")
			ctx.HTML(200, "user")
		})
	m.Get("/game/character",
		func(ctx *macaron.Context, s session.Store) {
			ctx.Data["userid"] = s.Get("userid")
			ctx.HTML(200, "character")
		})
	m.Get("/game/land",
		func(ctx *macaron.Context, s session.Store) {
			ctx.Data["userid"] = s.Get("userid")
			ctx.HTML(200, "land")
		})
}

func routeMiddlewares(m *macaron.Macaron, dirTemplates, dirStatic string) {
	// TODO(jfs): The secret has to be shared among all the running instances
	m.SetDefaultCookieSecret(randomSecret())
	m.Use(macaron.Static(dirStatic, macaron.StaticOptions{
		Prefix: "static",
	}))
	m.Use(pongo2.Pongoer(pongo2.Options{
		Directory:       dirTemplates,
		Extensions:      []string{".tpl", ".html", ".tmpl"},
		HTMLContentType: "text/html",
		Charset:         "UTF-8",
		IndentJSON:      true,
		IndentXML:       true,
	}))
	m.Use(session.Sessioner())
	m.Use(func(ctx *macaron.Context, s session.Store) {
		auth := func() {
			uid := s.Get("userid")
			if uid == "" {
				ctx.Redirect("/index.html")
			}
		}
		// Pages under the /game/* prefix require an established authentication
		switch {
		case strings.HasPrefix(ctx.Req.URL.Path, "/game/"):
			auth()
		case strings.HasPrefix(ctx.Req.URL.Path, "/action/"):
			auth()
		}
	})
}

func randomSecret() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(time.Now().UnixNano(), 16))
	sb.WriteRune('-')
	sb.WriteString(strconv.FormatUint(uint64(rand.Uint32()), 16))
	sb.WriteRune('-')
	sb.WriteString(strconv.FormatUint(uint64(rand.Uint32()), 16))
	return sb.String()
}

func main() {
	var err error

	var north, world, templates, static string
	flag.StringVar(&north, "north", "127.0.0.1:8080", "TCP/IP North endpoint")
	flag.StringVar(&world, "world", "127.0.0.1:8081", "World Server to be contacted")
	flag.StringVar(&templates, "templates", "/var/lib/hegemonie/templates", "Directory with the HTML tmeplates")
	flag.StringVar(&static, "static", "/var/lib/hegemonie/static", "Directory with the static files")
	flag.Parse()

	m := macaron.Classic()
	routeMiddlewares(m, templates, static)
	routeForms(m)
	routePages(m)

	err = http.ListenAndServe(north, m)
	if err != nil {
		log.Printf("Server error: %s", err.Error())
	}
}
