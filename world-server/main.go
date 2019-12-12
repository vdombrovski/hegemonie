// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"flag"
	"gopkg.in/macaron.v1"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	. "../world"
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

func routes(w *World, m *macaron.Macaron) {
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
}

func runServer(w *World, north string) error {
	m := macaron.Classic()
	routes(w, m)
	m.NotFound(func(ctx *macaron.Context) (int, string) {
		return 404, ""
	})
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
