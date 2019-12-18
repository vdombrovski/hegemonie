// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

func (p *World) GetBuildingType(id uint64) *BuildingType {
	for _, i := range p.BuildingTypes {
		if i.Id == id {
			return &i
		}
	}
	return nil
}
