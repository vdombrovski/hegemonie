// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"sync"
	"sync/atomic"
)

const (
	ResourceMax = 4
)

type Resources struct {
	Amounts [ResourceMax]uint64
}

type ResourcesMultiplier struct {
	Ratios [ResourceMax]float32
}

type UnitType struct {
	// Unique Id of the Unit Type
	Id uint64

	// The number of health point for that type of unit.
	// A health equal to 0 means the death of the unit.
	Health uint

	// How afftected is that type of unit by a loss of Health.
	// Must be between 0 and 1.
	// 0 means that the capacity of the Unit isn't affected by a health reduction.
	// 1 means that the capacity of the Unit loses an equal percentage of its capacity
	// for a loss of health (in other words, a HealthFactor of 1 means that the Unit
	// will hit at 90% of its maximal power if it has 90% of its health points).
	HealthFactor float32

	// The display name of the Unit Type
	Name string

	// Instantiation cost of the current UnitType
	Build Resources

	//
	Maintenance Resources
}

// Both Cell and City must not be 0, and have a non-0 value.
type Unit struct {
	// Unique Id of the Unit
	Id uint64

	// The number of health points of the unit, Health should be less or equal to HealthMax
	Health uint

	// A copy of the definition for the current Unit.
	Type uint64

	// The unique Id of the map cell the current Unit is on.
	Cell uint64

	// The unique Id of the City the Unit is in.
	City uint64
}

type BuildingType struct {
	// Unique ID of the BuildingType
	Id uint64

	// Display name of the current BuildingType
	Name string

	// Multiplier of the City production
	Multiplier ResourcesMultiplier

	// Increment of resources produced by this building.
	Boost Resources

	// How much does the production cost
	BuildingCost Resources
}

type Building struct {
	// The unique ID of the current Building
	Id uint64

	// The unique ID of the BuildingType associated to the current Building
	Type uint64
}

type CityCore struct {
	// The unique ID of the current City
	Id uint64

	// The unique ID of the main Character in charge of the City.
	// The Manager may name a Deputy manager in the City.
	Owner uint64

	// The unique ID of a second Character in charge of the City.
	Deputy uint64

	// The unique ID of the Cell the current City is built on.
	// This is redundant with the City field in the Cell structure.
	// Both information must match.
	Cell uint64

	// The display name of the current City
	Name string

	// Resources stock owned by the current City
	Stock Resources

	// Resources produced each round by the City, before the enforcing of
	// Production Boosts ans Production Multipliers
	Production Resources
}

type City struct {
	Meta CityCore

	Deleted bool

	// An array of Units guarding the current City.
	// This is redundant with the City field of the Unit type.
	// Consider it as an index.
	Units []uint64

	Buildings []Building
}

type CityView struct {
	Core CityCore

	Units []Unit

	Buildings []Building
}

type Character struct {
	// The unique identifier of the current Character
	Id uint64

	// The unique identifier of the only User that controls the Character.
	User uint64

	// The display name of the current Character
	Name string
}

type User struct {
	// The unique identifier of the current User
	Id uint64

	// The display name of the current User
	Name string

	// The unique email that authenticates the User.
	Email string

	// The hashed password that authenticates the User
	Password string

	// Has the current User the permission to manage the service.
	Admin bool
}

type Politics struct {
	Users         []User
	Characters    []Character
	Cities        []City
	Units         []Unit
	UnitTypes     []UnitType
	BuildingTypes []BuildingType

	NextId uint64
	Salt   string
	rw     sync.RWMutex
}

func (p *Politics) Init() {
	p.rw.Lock()
	defer p.rw.Unlock()

	if p.NextId <= 0 {
		p.NextId = 1
	}
	p.Users = make([]User, 0)
	p.Characters = make([]Character, 0)
	p.Cities = make([]City, 0)
	p.Units = make([]Unit, 0)
	p.UnitTypes = make([]UnitType, 0)
	p.BuildingTypes = make([]BuildingType, 0)
}

func (p *Politics) ReadLocker() sync.Locker {
	return p.rw.RLocker()
}

func (p *Politics) getNextId() uint64 {
	return atomic.AddUint64(&p.NextId, 1)
}

func (p *Politics) GetUnit(id uint64) *Unit {
	for _, c := range p.Units {
		if c.Id == id {
			return &c
		}
	}
	return nil
}

func (p *Politics) GetUnitType(id uint64) *UnitType {
	for _, c := range p.UnitTypes {
		if c.Id == id {
			return &c
		}
	}
	return nil
}

func (p *Politics) GetBuildingType(id uint64) *BuildingType {
	for _, i := range p.BuildingTypes {
		if i.Id == id {
			return &i
		}
	}
	return nil
}

func (p *Politics) Check(w *World) error {
	p.rw.RLock()
	defer p.rw.RUnlock()

	return nil
}
