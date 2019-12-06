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

// A Road is a Directed Vertex in the transport graph
type Road struct {
	// Unique identifier of the source Cell
	SrcCell uint64
	// Unique identifier of the destination Cell
	DstCell uint64
	// May the road be used by Units
	Open bool
}

// A MapCell is a Node is the directed graph of the transport network.
type MapCell struct {
	// The unique identifier of the current cell.
	Id uint64

	// The unique ID of the city present at this location.
	City uint64

	// An array of the the Id of the Unit that are present on the cell.
	Units []uint64
}

// A Map is a directed graph destined to be used as a transport network,
// organised as an adjacency list.
type Map struct {
	Locations []MapCell
	Roads     []Road

	NextId uint64
	rw     sync.RWMutex
}

func (m *Map) Init() {
	m.Locations = make([]MapCell, 0)
	m.Roads = make([]Road, 0)
}

func (m *Map) ReadLocker() sync.Locker {
	return m.rw.RLocker()
}

func (m *Map) getNextId() uint64 {
	return atomic.AddUint64(&m.NextId, 1)
}

func (m *Map) HasLocation(loc uint64) bool {
	if loc == 0 {
		return false
	}
	for _, l := range m.Locations {
		if l.Id == loc {
			return true
		}
	}
	return false
}

func (m *Map) CreateLocation() (uint64, error) {
	m.rw.Lock()
	defer m.rw.Unlock()

	loc := m.getNextId()
	m.Locations = append(m.Locations, MapCell{Id: loc})
	return loc, nil
}

func (m *Map) Connect(src, dst uint64) error {
	if src == dst || src == 0 || dst == 0 {
		return errors.New("EINVAL")
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	if !m.HasLocation(src) {
		return errors.New("Source not found")
	}
	if !m.HasLocation(dst) {
		return errors.New("Destination not found")
	}
	for _, r := range m.Roads {
		if r.SrcCell == src && r.DstCell == dst {
			if r.Open {
				return errors.New("Road exists")
			} else {
				r.Open = true
				return nil
			}
		}
	}
	m.Roads = append(m.Roads, Road{src, dst, true})
	return nil
}

func (m *Map) Disconnect(src, dst uint64) error {
	if src == dst || src == 0 || dst == 0 {
		return errors.New("EINVAL")
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	if !m.HasLocation(src) {
		return errors.New("Source not found")
	}
	if !m.HasLocation(dst) {
		return errors.New("Destination not found")
	}
	for _, r := range m.Roads {
		if r.SrcCell == src && r.DstCell == dst {
			if r.Open {
				r.Open = false
				return nil
			} else {
				return errors.New("Road exists")
			}
		}
	}
	return errors.New("Road not found")
}

func (m *Map) NextStep(src, dst uint64) (uint64, error) {
	if src == dst || src == 0 || dst == 0 {
		return 0, errors.New("EINVAL")
	}

	return 0, errors.New("NYI")
}

func (m *Map) Check(w *World) error {
	return nil
}
