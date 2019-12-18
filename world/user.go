// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func (p *World) UserCreate(mail, pass string) (uint64, error) {
	if !validMail(mail) || !validPass(pass) {
		return 0, errors.New("EINVAL")
	}

	p.rw.Lock()
	defer p.rw.Unlock()

	h := p.hashPassword(pass)
	for _, u := range p.Users {
		if u.Email == mail && u.Password == h {
			return 0, errors.New("User exists")
		}
	}

	id := p.getNextId()
	u := User{Id: id, Name: "No-Name", Email: mail, Password: pass}
	p.Users = append(p.Users, u)
	return id, nil
}

func (p *World) UserGet(id uint64) (User, error) {
	if id <= 0 {
		return User{}, errors.New("EINVAL")
	}

	p.rw.RLock()
	defer p.rw.RUnlock()

	// TODO(jfs): lookup in the sorted array
	for _, u := range p.Users {
		if u.Id == id {
			var copy User = u
			copy.Password = ""
			return copy, nil
		}
	}

	return User{}, errors.New("User not found")
}

func (p *World) UserAuth(mail, pass string) (uint64, error) {
	if mail == "" || pass == "" {
		return 0, errors.New("EINVAL")
	}

	p.rw.RLock()
	defer p.rw.RUnlock()

	h := p.hashPassword(pass)
	for _, u := range p.Users {
		if u.Email == mail {
			if u.Password == h {
				// Hashed password matches
				return u.Id, nil
			} else if u.Password[0] == ':' && u.Password[1:] == pass {
				// Clear password matches
				return u.Id, nil
			} else {
				return 0, nil
			}
		}
	}

	return 0, errors.New("User not found")
}

func (p *World) UserGetCharacters(id uint64, hook func(*Character)) {
	for _, c := range p.Characters {
		if c.User == id {
			hook(&c)
		}
	}
}

func (p *World) hashPassword(pass string) string {
	checksum := sha256.New()
	checksum.Write([]byte(p.Salt))
	checksum.Write([]byte(pass))
	return hex.EncodeToString(checksum.Sum(nil))
}

func validMail(m string) bool {
	// TODO(jfs): Not yet implemented
	return len(m) > 0
}

func validPass(m string) bool {
	// TODO(jfs): Not yet implemented
	return len(m) > 0
}
