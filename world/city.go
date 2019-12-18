// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import "errors"

func (p *World) CityShow(userId, characterId, cityId uint64) (CityView, error) {
	p.rw.RLock()
	defer p.rw.RUnlock()

	var err error
	var result CityView

	pCity := p.CityGet(cityId)
	pChar := p.CharacterGet(characterId)
	if pCity == nil || pChar == nil {
		err = errors.New("Not Found")
	} else if pCity.Meta.Deputy != characterId && pCity.Meta.Owner != characterId {
		err = errors.New("Forbidden")
	} else if pChar.User != userId {
		err = errors.New("Forbidden")
	} else {
		result.Core = pCity.Meta
		result.Buildings = pCity.Buildings
		result.Units = make([]Unit, len(pCity.Units), len(pCity.Units))
	}

	return result, err
}

func (p *World) CityGet(id uint64) *City {
	for _, c := range p.Cities {
		if c.Meta.Id == id {
			return &c
		}
	}
	return nil
}

func (p *World) CityCheck(id uint64) bool {
	return p.CityGet(id) != nil
}

func (p *World) CityCreate(id, loc uint64) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	c0 := p.CityGet(id)
	if c0 != nil {
		if c0.Deleted {
			c0.Deleted = false
			return nil
		} else {
			return errors.New("City already exists")
		}
	}

	c := City{Meta: CityCore{Id: id, Cell: loc}, Units: make([]uint64, 0)}
	p.Cities = append(p.Cities, c)
	return nil
}

func (p *World) CitySpawnUnit(idCity, idType uint64) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	c := p.CityGet(idCity)
	if c == nil {
		return errors.New("City not found")
	}

	t := p.GetUnitType(idType)
	if t == nil {
		return errors.New("Unit type not found")
	}

	unit := Unit{Id: p.getNextId(), Health: t.Health, Type: t.Id, City: idCity, Cell: 0}
	p.Units = append(p.Units, unit)

	c.Units = append(c.Units, unit.Id)
	return nil
}

func (c *City) CityGetBuilding(id uint64) *Building {
	for _, b := range c.Buildings {
		if id == b.Id {
			return &b
		}
	}
	return nil
}

func (p *World) CitySpawnBuilding(idCity, idType uint64) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	c := p.CityGet(idCity)
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
