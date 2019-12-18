// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"errors"
)

func (w *World) CharacterGet(cid uint64) *Character {
	// TODO(jfs): lookup in the sorted array
	for _, c := range w.Characters {
		if c.Id == cid {
			return &c
		}
	}
	return nil
}

func (w *World) CharacterShow(uid, cid uint64) (Character, error) {
	if cid <= 0 || uid <= 0 {
		return Character{}, errors.New("EINVAL")
	}

	w.rw.RLock()
	defer w.rw.RUnlock()

	if pChar := w.CharacterGet(cid); pChar == nil {
		return Character{}, errors.New("Not Found")
	} else if pChar.User != uid {
		return Character{}, errors.New("Forbidden")
	} else {
		return *pChar, nil
	}
}

// Notify the caller of the cities managed by the given Character.
func (w *World) CharacterGetCities(id uint64, owner func(*City), deputy func(*City)) {
	if id <= 0 {
		return
	}

	w.rw.RLock()
	defer w.rw.RUnlock()

	for _, c := range w.Cities {
		if c.Meta.Owner == id {
			owner(&c)
		} else if c.Meta.Deputy == id {
			deputy(&c)
		}
	}
}
