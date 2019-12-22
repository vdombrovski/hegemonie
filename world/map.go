// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"errors"
	"sort"
	"strconv"
	"strings"
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

	// Biome in which the cell is
	Biome uint64
}

type SetOfCells []MapCell

func (r *SetOfCells) Len() int {
	return len(*r)
}

func (r *SetOfCells) Swap(i, j int) {
	tmp := (*r)[i]
	(*r)[i] = (*r)[j]
	(*r)[j] = tmp
}

func (r *SetOfCells) Less(i, j int) bool {
	return (*r)[i].Id < (*r)[j].Id
}

type SetOfRoads []Road

func (r *SetOfRoads) Len() int {
	return len(*r)
}

func (r *SetOfRoads) Swap(i, j int) {
	tmp := (*r)[i]
	(*r)[i] = (*r)[j]
	(*r)[j] = tmp
}

func (r *SetOfRoads) Less(i, j int) bool {
	s := (*r)[i]
	d := (*r)[j]
	return s.SrcCell < d.SrcCell || (s.SrcCell == d.SrcCell && s.DstCell < d.DstCell)
}

// A Map is a directed graph destined to be used as a transport network,
// organised as an adjacency list.
type Map struct {
	Cells  SetOfCells
	Roads  SetOfRoads
	NextId uint64

	steps      map[vector]uint64
	dirtyRoads bool
	dirtyCells bool

	rw sync.RWMutex
}

func (m *Map) Init() {
	m.Cells = make([]MapCell, 0)
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
	if m.dirtyCells {
		for _, l := range m.Cells {
			if l.Id == loc {
				return true
			}
		}
		return false
	} else {
		i := sort.Search(len(m.Cells), func(i int) bool {
			return m.Cells[i].Id >= loc
		})
		return i < len(m.Cells) && m.Cells[i].Id == loc
	}
}

func (m *Map) CreateLocation() (uint64, error) {
	m.rw.Lock()
	defer m.rw.Unlock()

	loc := m.getNextId()
	m.Cells = append(m.Cells, MapCell{Id: loc})
	return loc, nil
}

func (m *Map) ConnectRaw(src, dst uint64) error {
	if src == dst || src == 0 || dst == 0 {
		return errors.New("EINVAL")
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	m.Roads = append(m.Roads, Road{src, dst, true})
	m.dirtyRoads = true
	return nil
}

func (m *Map) firstAdjacentIndex(src uint64) int {
	m.lazySort()
	i := sort.Search(len(m.Roads), func(i int) bool {
		r := m.Roads[i]
		return r.SrcCell >= src
	})
	return i
}

func (m *Map) Connect(src, dst uint64, check bool) error {
	if src == dst || src == 0 || dst == 0 {
		return errors.New("EINVAL")
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	if check && !m.HasLocation(src) {
		return errors.New("Source not found")
	}
	if check && !m.HasLocation(dst) {
		return errors.New("Destination not found")
	}

	for i := m.firstAdjacentIndex(src); i < len(m.Roads); i++ {
		r := m.Roads[i]
		if r.SrcCell != src {
			break
		}
		if r.DstCell == dst {
			if r.Open {
				return errors.New("Road exists")
			} else {
				r.Open = true
				return nil
			}
		}
	}

	m.Roads = append(m.Roads, Road{src, dst, true})
	m.dirtyRoads = true
	return nil
}

func (m *Map) Disconnect(src, dst uint64, check bool) error {
	if src == dst || src == 0 || dst == 0 {
		return errors.New("EINVAL")
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	if check && !m.HasLocation(src) {
		return errors.New("Source not found")
	}
	if check && !m.HasLocation(dst) {
		return errors.New("Destination not found")
	}

	for i := m.firstAdjacentIndex(src); i < len(m.Roads); i++ {
		r := m.Roads[i]
		if r.SrcCell != src {
			break
		}
		if r.DstCell == dst {
			if r.Open {
				r.Open = false
				return nil
			} else {
				return errors.New("Road closed")
			}
		}
	}

	return errors.New("Road not found")
}

func (m *Map) NextStep(src, dst uint64) (uint64, error) {
	if src == dst || src == 0 || dst == 0 {
		return 0, errors.New("EINVAL")
	}

	next, ok := m.steps[vector{src, dst}]
	if ok {
		return next, nil
	} else {
		return 0, errors.New("No route")
	}
}

func (m *Map) Adjacent(src uint64) []uint64 {
	adj := make([]uint64, 0)

	for i := m.firstAdjacentIndex(src); i < len(m.Roads); i++ {
		r := m.Roads[i]
		if r.SrcCell != src {
			break
		}
		adj = append(adj, r.DstCell)
	}

	return adj
}

func (m *Map) Check(w *World) error {
	return nil
}

func (m *Map) Dot() string {
	var sb strings.Builder
	sb.WriteString("digraph g {")
	for _, c := range m.Cells {
		sb.WriteString("n" + strconv.FormatUint(c.Id, 10))
		sb.WriteRune(';')
		sb.WriteRune('\n')
	}
	for _, r := range m.Roads {
		sb.WriteRune(' ')
		sb.WriteString("n" + strconv.FormatUint(r.SrcCell, 10))
		sb.WriteString(" -> ")
		sb.WriteString("n" + strconv.FormatUint(r.DstCell, 10))
		sb.WriteRune(';')
		sb.WriteRune('\n')
	}
	sb.WriteString("}")
	return sb.String()
}

func (m *Map) lazySort() {
	if m.dirtyCells {
		sort.Sort(&m.Cells)
		m.dirtyCells = false
	}
	if m.dirtyRoads {
		sort.Sort(&m.Roads)
		m.dirtyRoads = false
	}
}

func (m *Map) Rehash() {
	next := make(map[vector]uint64)

	m.rw.Lock()
	defer m.rw.Unlock()

	// Ensure the locations are sorted
	m.lazySort()

	// Fill with the immediate neighbors
	for _, r := range m.Roads {
		next[vector{r.SrcCell, r.DstCell}] = r.DstCell
	}

	add := func(src, dst, step uint64) {
		_, found := next[vector{src, dst}]
		if !found {
			next[vector{src, dst}] = step
		}
	}

	// Call one DFS per node and shortcut when possible
	for _, cell := range m.Cells {
		already := make(map[uint64]bool)
		q := newQueue()

		// Bootstrap the DFS with adjacent nodes
		for _, next := range m.Adjacent(cell.Id) {
			q.push(next, next)
			already[next] = true
			// No need to add this in the known routes, we already did it
			// with an iteration on the roads (much faster)
		}

		for !q.empty() {
			current, first := q.pop()
			neighbors := m.Adjacent(current)
			// TODO(jfs): shuffle the neighbors
			for _, next := range neighbors {
				if !already[next] {
					// Avoid passing again in the neighbor
					already[next] = true
					// Tell to contine at that neighbor
					q.push(next, first)
					// We already learned the shortest path to that neighbor
					add(cell.Id, next, first)
				}
			}
		}
	}

	m.steps = next
}

type vector struct {
	src uint64
	dst uint64
}

type dfsTrack struct {
	current uint64
	first   uint64
}

type queue struct {
	tab   []dfsTrack
	start int
}

func newQueue() queue {
	var q queue
	q.tab = make([]dfsTrack, 0)
	return q
}

func (q *queue) push(node, first uint64) {
	q.tab = append(q.tab, dfsTrack{node, first})
}

func (q *queue) pop() (uint64, uint64) {
	v := q.tab[q.start]
	q.start++
	return v.current, v.first
}

func (q *queue) empty() bool {
	return q.start == len(q.tab)
}
