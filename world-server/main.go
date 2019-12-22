// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"flag"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	. "hegemonie/world"
	. "hegemonie/world-client"
)

var (
	pathSave string
)

func makeSaveFilename() string {
	now := time.Now().Round(1 * time.Second)
	return "save-" + now.Format("20060102_030405")
}

func save(w *World) error {
	if pathSave == "" {
		return errors.New("No save path configured")
	}
	p := pathSave + "/" + makeSaveFilename()
	p = filepath.Clean(p)
	out, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = w.DumpJSON(out)
	out.Close()
	if err != nil {
		_ = os.Remove(p)
		return err
	}

	latest := pathSave + "/latest"
	_ = os.Remove(latest)
	_ = os.Symlink(p, latest)
	return nil
}

type AuthRequest struct {
	UserMail string `form:"email" binding:"Required"`
	UserPass string `form:"password" binding:"Required"`
}

func routes(w *World, m *macaron.Macaron) {
	m.Post("/user/auth", binding.Bind(AuthRequest{}),
		func(ctx *macaron.Context, form AuthRequest) {
			id, err := w.UserAuth(form.UserMail, form.UserPass)
			if id != 0 {
				ctx.JSON(200, AuthReply{Id: id})
			} else if err == nil {
				ctx.JSON(403, AuthReply{Id: 0})
			} else {
				ctx.JSON(500, AuthReply{Id: 0, Msg: err.Error()})
			}
		})

	m.Get("/user/show",
		func(ctx *macaron.Context) {
			struid := ctx.Query("uid")
			if id, err := strconv.ParseUint(struid, 10, 63); err != nil {
				ctx.JSON(400, ErrorReply{Code: 400, Msg: "Malformed User ID"})
			} else if user, err := w.UserGet(id); err != nil {
				ctx.JSON(404, ErrorReply{Code: 400, Msg: err.Error()})
			} else {
				var payload UserShowReply
				payload.Characters = make([]NamedItem, 0)
				payload.Meta = user
				w.UserGetCharacters(id, func(c *Character) {
					payload.Characters = append(payload.Characters, NamedItem{Name: c.Name, Id: c.Id})
				})
				ctx.JSON(200, &payload)
			}
		})

	m.Get("/character/show",
		func(ctx *macaron.Context) {
			struid := ctx.Query("uid")
			strcid := ctx.Query("cid")
			if uid, err := strconv.ParseUint(struid, 10, 63); err != nil {
				ctx.JSON(400, ErrorReply{Code: 400, Msg: "Malformed User ID"})
			} else if cid, err := strconv.ParseUint(strcid, 10, 63); err != nil {
				ctx.JSON(400, ErrorReply{Code: 400, Msg: "Malformed Character ID"})
			} else if character, err := w.CharacterShow(uid, cid); err != nil {
				ctx.JSON(404, ErrorReply{Code: 400, Msg: err.Error()})
			} else {
				var payload CharacterShowReply
				payload.Meta = character
				payload.OwnerOf = make([]NamedItem, 0)
				payload.DeputyOf = make([]NamedItem, 0)
				w.CharacterGetCities(cid,
					func(c *City) {
						payload.OwnerOf = append(payload.OwnerOf, NamedItem{Name: c.Meta.Name, Id: c.Meta.Id})
					},
					func(c *City) {
						payload.DeputyOf = append(payload.DeputyOf, NamedItem{Name: c.Meta.Name, Id: c.Meta.Id})
					})
				ctx.JSON(200, &payload)
			}
		})

	m.Get("/land/show",
		func(ctx *macaron.Context) {
			struid := ctx.Query("uid")
			strcid := ctx.Query("cid")
			strlid := ctx.Query("lid")
			if uid, err := strconv.ParseUint(struid, 10, 63); err != nil {
				ctx.JSON(400, ErrorReply{Code: 400, Msg: "Malformed User ID"})
			} else if cid, err := strconv.ParseUint(strcid, 10, 63); err != nil {
				ctx.JSON(400, ErrorReply{Code: 400, Msg: "Malformed Character ID"})
			} else if lid, err := strconv.ParseUint(strlid, 10, 63); err != nil {
				ctx.JSON(400, ErrorReply{Code: 400, Msg: "Malformed Land ID"})
			} else if cityView, err := w.CityShow(uid, cid, lid); err != nil {
				ctx.JSON(404, ErrorReply{Code: 400, Msg: err.Error()})
			} else {
				var payload CityShowReply
				log.Println("cityView:", cityView)
				payload.Meta = cityView.Core
				payload.Units = make([]NamedItem, 0)
				for _, u := range cityView.Units {
					payload.Units = append(payload.Units,
						NamedItem{Id: u.Id, Name: w.GetUnitType(u.Type).Name})
				}
				payload.Buildings = make([]NamedItem, 0)
				for _, b := range cityView.Buildings {
					payload.Buildings = append(payload.Buildings,
						NamedItem{Id: b.Id, Name: w.GetBuildingType(b.Type).Name})
				}
				ctx.JSON(200, &payload)
			}
		})

	m.Get("/map/dot", func(ctx *macaron.Context) (int, string) {
		return 200, w.Places.Dot()
	})
	m.Post("/map/rehash", func(ctx *macaron.Context) (int, string) {
		w.Places.Rehash()
		return 204, ""
	})
	m.Post("/map/check", func(ctx *macaron.Context) (int, string) {
		if err := w.Places.Check(w); err == nil {
			return 204, ""
		} else {
			return 502, err.Error()
		}
	})
	m.Post("/check", func(ctx *macaron.Context) (int, string) {
		if err := w.Check(); err == nil {
			return 204, ""
		} else {
			return 502, err.Error()
		}
	})
	m.Post("/save", func(ctx *macaron.Context) (int, string) {
		if err := save(w); err == nil {
			return 204, ""
		} else {
			return 501, err.Error()
		}
	})

	// Mapping routes
	m.Get("/world/places", func(ctx *macaron.Context) {
		ctx.JSON(200, &w.Places)
	})
	m.Get("/world/cities", func(ctx *macaron.Context) {
		ctx.JSON(200, &w.Cities)
	})
}

func runServer(w *World, north string) error {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	routes(w, m)
	return http.ListenAndServe(north, m)
}

func main() {
	var err error
	var w World

	w.Init()

	var north string
	var pathLoad string
	flag.StringVar(&north, "north", "127.0.0.1:8081", "File to be loaded")
	flag.StringVar(&pathLoad, "load", "", "File to be loaded")
	flag.StringVar(&pathSave, "save", "/tmp/hegemonie/data", "Directory for persistent")
	flag.Parse()

	if pathSave != "" {
		err = os.MkdirAll(pathSave, 0755)
		if err != nil {
			log.Fatalf("Failed to create [%s]: %s", pathSave, err.Error())
		}
	}

	if pathLoad != "" {
		in, err := os.Open(pathLoad)
		if err != nil {
			log.Fatalf("Failed to load the World from [%s]: %s", pathLoad, err.Error())
		}
		err = w.LoadJSON(in)
		in.Close()
		if err != nil {
			log.Fatalf("Failed to load the World from [%s]: %s", pathLoad, err.Error())
		}
	}

	err = w.Check()
	if err != nil {
		log.Fatalf("Inconsistent World: %s", err.Error())
	}

	err = runServer(&w, north)
	if err != nil {
		log.Printf("Server error: %s", err.Error())
	}

	if pathSave != "" {
		err = save(&w)
		if err != nil {
			log.Fatalf("Failed to save the World at exit: %s", err.Error())
		}
	}
}
