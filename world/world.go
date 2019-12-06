// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"encoding/json"
	"io"
)

type World struct {
	People Politics
	Places Map
}

func (w *World) Init() {
	w.People.Init()
	w.Places.Init()
}

func (w *World) Check() error {
	if err := w.People.Check(w); err != nil {
		return err
	} else {
		return w.Places.Check(w)
	}
}

func (w *World) DumpJSON(dst io.Writer) error {
	return json.NewEncoder(dst).Encode(w)
}

func (w *World) LoadJSON(src io.Reader) error {
	return json.NewDecoder(src).Decode(w)
}
