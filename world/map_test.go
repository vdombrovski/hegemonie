// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import "testing"

func TestMapInit(t *testing.T) {
	var m Map
	m.Init()

	if m.HasLocation(1) {
		t.Fatal()
	}
	if m.getNextId() != 1 {
		t.Fatal()
	}

	loc, err := m.CreateLocation()
	if err != nil {
		t.Fatal()
	}

	if loc != 2 {
		t.Fatal()
	}
	if m.getNextId() != 3 {
		t.Fatal()
	}
	if !m.HasLocation(loc) {
		t.Fatal()
	}
}

func TestMapEinval(t *testing.T) {
	var m Map
	m.Init()

	// Test that identical, zero or inexistant locations return an error
	for _, src := range []uint64{0, 1, 2} {
		for _, dst := range []uint64{0, 1, 2} {
			if err := m.Connect(src, dst); err == nil {
				t.Fatal()
			}
			if err := m.Disconnect(src, dst); err == nil {
				t.Fatal()
			}
		}
	}
}

func TestMapPathOneWay(t *testing.T) {
	var m Map
	m.Init()

	l0, err := m.CreateLocation()
	if err != nil {
		t.Fatal()
	}

	l1, err := m.CreateLocation()
	if err != nil {
		t.Fatal()
	}

	err = m.Connect(l0, l1)
	if err != nil {
		t.Fatal()
	}

	if step, err := m.NextStep(l0, l1); err != nil {
		t.Fatal()
	} else if step != l1 {
		t.Fatal()
	}

	if step, err := m.NextStep(l1, l0); err == nil {
		t.Fatal()
	} else if step != 0 {
		t.Fatal()
	}
}
