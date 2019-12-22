// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world_client

import (
	"hegemonie/world"
)

type ErrorReply struct {
	Code int
	Msg  string
}

type AuthReply struct {
	Id  uint64
	Msg string
}

type NamedItem struct {
	Name string
	Id   uint64
}

type UserShowReply struct {
	Meta       world.User
	Characters []NamedItem
}

type CharacterShowReply struct {
	Meta     world.Character
	OwnerOf  []NamedItem
	DeputyOf []NamedItem
}

type CityShowReply struct {
	Meta      world.CityCore
	Units     []NamedItem
	Buildings []NamedItem
}
