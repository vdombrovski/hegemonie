// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Resources struct {
	Gold  uint64
	Wood  uint64
	Stone uint64
	Wheat uint64
}

type ResourcesMultiplier struct {
	Gold  float32
	Wood  float32
	Stone float32
	Wheat float32
}

type UnitType struct {
	// Unique Id of the Unit Type
	Id uint64

	// The display name of the Unit Type
	Name string

	// Instantiation cost of the current UnitType
	Cost Resources
}

// Both Cell and City must not be 0, and have a non-0 value.
type Unit struct {
	// Unique Id of the Unit
	Id uint64

	// The unique Id of the unit type. It must not be 0.
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

type City struct {
	// The unique ID of the current City
	Id uint64

	// The unique ID of the main Character in charge of the City.
	// The Manager may name a Deputy manager in the City.
	Owner uint64

	// The unique ID of a second Character in charge of the City.
	Admin []uint64

	// The unique ID of the Cell the current City is built on.
	// This is redundant with the City field in the Cell structure.
	// Both information must match.
	Cell uint64

	// The display name of the current City
	Name string

	// Tells if the City structure is usable. A deleted City generates "Not Found"
	// errors.
	Deleted bool

	// Resources stock owned by the current City
	Stock Resources

	// Resources produced each round by the City, before the enforcing of
	// Production Boosts ans Production Multipliers
	Production Resources

	// An array of Units guarding the current City.
	// This is redundant with the City field of the Unit type.
	// Consider it as an index.
	Units []uint64

	Buildings []Building
}

type Politics struct {
	Cities        []City
	Units         []Unit
	UnitTypes     []UnitType
	BuildingTypes []BuildingType

	NextId uint64
	rw     sync.RWMutex
}

func (p *Politics) Init() {
	p.rw.Lock()
	defer p.rw.Unlock()

	p.NextId = 1
	p.Cities = make([]City, 0)
	p.Units = make([]Unit, 0)
}

func (p *Politics) ReadLocker() sync.Locker {
	return p.rw.RLocker()
}

func (p *Politics) getNextId() uint64 {
	return atomic.AddUint64(&p.NextId, 1)
}

func (p *Politics) GetCity(id uint64) *City {
	for _, c := range p.Cities {
		if c.Id == id {
			return &c
		}
	}
	return nil
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

func (p *Politics) HasCity(id uint64) bool {
	return p.GetCity(id) != nil
}

func (p *Politics) CreateCity(id, loc uint64) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	c0 := p.GetCity(id)
	if c0 != nil {
		if c0.Deleted {
			c0.Deleted = false
			return nil
		} else {
			return errors.New("City already exists")
		}
	}

	c := City{Id: id, Cell: loc, Units: make([]uint64, 0)}
	p.Cities = append(p.Cities, c)
	return nil
}

func (p *Politics) SpawnUnit(idCity, idType uint64) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	c := p.GetCity(idCity)
	if c == nil {
		return errors.New("City not found")
	}

	t := p.GetUnitType(idType)
	if t == nil {
		return errors.New("Unit type not found")
	}

	unit := Unit{Id: p.getNextId(), Type: t.Id, City: c.Id, Cell: 0}
	p.Units = append(p.Units, unit)

	c.Units = append(c.Units, unit.Id)
	return nil
}

func (c *City) GetBuilding(id uint64) *Building {
	for _, b := range c.Buildings {
		if id == b.Id {
			return &b
		}
	}
	return nil
}

func (p *Politics) SpawnBuilding(idCity, idType uint64) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	c := p.GetCity(idCity)
	if c == nil {
		return errors.New("City not found")
	}

	t := p.GetBuildingType(idType)
	if t == nil {
		return errors.New("Building tye not found")
	}

	// TODO(jfs): consume the resources

	b := Building{Id: p.getNextId(), Type: idType}
	c.Buildings = append(c.Buildings, b)
	return nil
}

func (m *Politics) Check(w *World) error {
	return nil
}
