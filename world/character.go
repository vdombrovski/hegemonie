// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"errors"
)

func (p *Politics) CharacterGet(cid uint64) *Character {
	// TODO(jfs): lookup in the sorted array
	for _, c := range p.Characters {
		if c.Id == cid {
			return &c
		}
	}
	return nil
}

func (p *Politics) CharacterShow(uid, cid uint64) (Character, error) {
	if cid <= 0 || uid <= 0 {
		return Character{}, errors.New("EINVAL")
	}

	p.rw.RLock()
	defer p.rw.RUnlock()

	if pChar := p.CharacterGet(cid); pChar == nil {
		return Character{}, errors.New("Not Found")
	} else if pChar.User != uid {
		return Character{}, errors.New("Forbidden")
	} else {
		return *pChar, nil
	}
}

// Notify the caller of the cities managed by the given Character.
func (p *Politics) CharacterGetCities(id uint64, owner func(*City), deputy func(*City)) {
	if id <= 0 {
		return
	}

	p.rw.RLock()
	defer p.rw.RUnlock()

	for _, c := range p.Cities {
		if c.Owner == id {
			owner(&c)
		} else if c.Deputy == id {
			deputy(&c)
		}
	}
}
