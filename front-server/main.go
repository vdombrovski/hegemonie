// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"hegemonie/common/mapper"
	. "hegemonie/world-client"
	"encoding/json"
	"flag"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/pongo2"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"io/ioutil"
)

type LoginForm struct {
	UserMail string `form:"email" binding:"Required"`
	UserPass string `form:"password" binding:"Required"`
}

type front struct {
	endpointNorth string
	endpointWorld string
	dirTemplates  string
	dirStatic     string
}

func (f *front) routePages(m *macaron.Macaron) {
	m.Get("/",
		func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
			ctx.Data["userid"] = sess.Get("uid")
			ctx.HTML(200, "index")
		})
	m.Get("/admin",
		func(ctx *macaron.Context) {

		})
	m.Get("/game/user",
		func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
			var detailUser UserShowReply

			// Validate the input
			sessid := sess.Get("userid")
			if sessid == nil {
				flash.Error("Invalid session")
				ctx.Redirect("/")
				return
			}
			strid := sessid.(string)

			// Query the World server for the user
			resp, err := http.Get("http://" + f.endpointWorld + "/user/show?uid=" + strid)
			if err != nil {
				flash.Warning("User error: " + err.Error())
				ctx.Redirect("/")
				return
			}
			// Unpack the user
			detailUser.Characters = make([]NamedItem, 0)
			if err = json.NewDecoder(resp.Body).Decode(&detailUser); err != nil {
				flash.Error("World problem: " + err.Error())
				ctx.Redirect("/")
				return
			}

			// Display the result
			ctx.Data["userid"] = detailUser.Meta.Id
			ctx.Data["User"] = &detailUser
			ctx.HTML(200, "user")
		})
	m.Get("/game/character",
		func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
			var detailCharacter CharacterShowReply

			// Validate the input
			sessid := sess.Get("userid")
			if sessid == nil {
				flash.Error("Invalid session")
				ctx.Redirect("/")
				return
			}
			userid := sessid.(string)
			charid := ctx.Query("cid")

			// Query the World server for the Character
			resp, err := http.Get("http://" + f.endpointWorld + "/character/show?uid=" + userid + "&cid=" + charid)
			if err != nil {
				flash.Warning("Character error: " + err.Error())
				ctx.Redirect("/")
				return
			}
			// Unpack the character
			detailCharacter.OwnerOf = make([]NamedItem, 0)
			detailCharacter.DeputyOf = make([]NamedItem, 0)
			if err = json.NewDecoder(resp.Body).Decode(&detailCharacter); err != nil {
				flash.Error("World problem: " + err.Error())
				ctx.Redirect("/")
				return
			}

			ctx.Data["userid"] = userid
			ctx.Data["cid"] = charid
			ctx.Data["Character"] = &detailCharacter
			ctx.HTML(200, "character")
		})
	m.Get("/game/land",
		func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
			// Validate the input
			sessid := sess.Get("userid")
			if sessid == nil {
				flash.Error("Invalid session")
				ctx.Redirect("/")
				return
			}
			userid := sessid.(string)
			charid := ctx.Query("cid")
			landid := ctx.Query("lid")
			if userid == "" || charid == "" || landid == "" {
				ctx.Redirect("/")
				return
			}

			// Query the World server for the Character
			resp, err := http.Get("http://" + f.endpointWorld + "/land/show?uid=" + userid + "&cid=" + charid + "&lid=" + landid)
			if err != nil {
				flash.Warning("Character error: " + err.Error())
				ctx.Redirect("/")
				return
			}
			// Unpack the character
			var detailLand CityShowReply
			if err = json.NewDecoder(resp.Body).Decode(&detailLand); err != nil {
				flash.Error("World problem: " + err.Error())
				ctx.Redirect("/")
				return
			}

			ctx.Data["userid"] = userid
			ctx.Data["cid"] = charid
			ctx.Data["lid"] = landid
			ctx.Data["Land"] = detailLand
			ctx.HTML(200, "land")
		})
	m.Get("/game/map",
		func(ctx *macaron.Context, s session.Store) {
			// gameMap, overlay, err := mapper.Generate()
			// if err != nil {
			// 	ctx.Resp.WriteHeader(500)
			// 	return
			// }
			// ctx.Data["map"] = gameMap
			// ctx.Data["overlay"] = overlay
			// ctx.HTML(200, "map")

			// TODO: VDO: handle error
			resp, _ := http.Get("http://" + f.endpointWorld + "/world/places")
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				// Backend error
				ctx.Resp.WriteHeader(503)
				return
			}
		    mapBytes, _ := ioutil.ReadAll(resp.Body)

			resp2, _ := http.Get("http://" + f.endpointWorld + "/world/cities")
			defer resp2.Body.Close()
			if resp2.StatusCode != http.StatusOK {
				// Backend error
				ctx.Resp.WriteHeader(503)
				return
			}
		    mapCities, _ := ioutil.ReadAll(resp2.Body)

			ctx.Data["map"] = string(mapBytes)
			ctx.Data["cities"] = string(mapCities)

			ctx.HTML(200, "map")
		})
	// TODO: VDO: disable these routes when DEBUG=false
	m.Get("/debug/map/map",
		func(ctx *macaron.Context, s session.Store) {
			gameMap, _, err := mapper.Generate()
			if err != nil {
				ctx.Resp.WriteHeader(500)
				return
			}
			ctx.JSON(200, gameMap)
		})
	m.Get("/debug/map/overlay",
		func(ctx *macaron.Context, s session.Store) {
			_, overlay, err := mapper.Generate()
			if err != nil {
				ctx.Resp.WriteHeader(500)
				return
			}
			ctx.JSON(200, overlay)
	})
}

func (f *front) routeForms(m *macaron.Macaron) {
	doLogIn := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info LoginForm) {
		// Cleanup a previous session
		sess.Flush()

		// Authenticate the user by the world-server
		var payload AuthReply
		var form url.Values = make(map[string][]string)
		form.Set("email", info.UserMail)
		form.Set("password", info.UserPass)
		resp, err := http.PostForm("http://"+f.endpointWorld+"/user/auth", form)
		log.Println(err, resp)
		if err != nil {
			flash.Error("Authentication error: " + err.Error())
			ctx.Redirect("/")
		} else if resp.StatusCode/100 != 2 {
			flash.Warning("Authentication failed")
			ctx.Redirect("/")
		} else if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			flash.Warning("Authentication problem: " + err.Error())
			ctx.Redirect("/")
		} else {
			// Establish a session for the user
			strid := strconv.FormatUint(payload.Id, 10)
			ctx.SetSecureCookie("session", strid)
			sess.Set("userid", strid)
			ctx.Redirect("/game/user")
		}
	}
	doLogOut := func(ctx *macaron.Context, s session.Store) {
		ctx.SetSecureCookie("session", "")
		s.Flush()
		ctx.Redirect("/")
	}
	m.Post("/action/login", binding.Bind(LoginForm{}), doLogIn)
	m.Post("/action/logout", doLogOut)
	m.Get("/action/logout", doLogOut)
}

func (f *front) routeMiddlewares(m *macaron.Macaron) {
	// TODO(jfs): The secret has to be shared among all the running instances
	m.SetDefaultCookieSecret(randomSecret())
	m.Use(macaron.Static(f.dirStatic, macaron.StaticOptions{
		Prefix: "static",
	}))
	m.Use(pongo2.Pongoer(pongo2.Options{
		Directory:       f.dirTemplates,
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
		case strings.HasPrefix(ctx.Req.URL.Path, "/game/"),
			strings.HasPrefix(ctx.Req.URL.Path, "/action/"):
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
	var f front

	flag.StringVar(&f.endpointNorth, "north", "127.0.0.1:8080", "TCP/IP North endpoint")
	flag.StringVar(&f.endpointWorld, "world", "127.0.0.1:8081", "World Server to be contacted")
	flag.StringVar(&f.dirTemplates, "templates", "/var/lib/hegemonie/templates", "Directory with the HTML tmeplates")
	flag.StringVar(&f.dirStatic, "static", "/var/lib/hegemonie/static", "Directory with the static files")
	flag.Parse()

	m := macaron.Classic()
	f.routeMiddlewares(m)
	f.routeForms(m)
	f.routePages(m)

	err = http.ListenAndServe(f.endpointNorth, m)
	if err != nil {
		log.Printf("Server error: %s", err.Error())
	}
}
