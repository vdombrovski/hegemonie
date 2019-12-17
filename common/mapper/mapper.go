// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapper

import (
    "time"
	"encoding/json"
    "math/rand"
)

type slot struct {
	Terrain int64
}

type biome struct {
	Type int64
	Decay int64
}

// Terrain generation primitive
func Generate() (string, string, error) {
    gameMap := genPhase1(64, 64, false)
    overlay := genPhase1(64, 64, true)
    gameMap, overlay = genPhase2(gameMap, overlay, 100, 64, 64)

    gameMapB, err := json.Marshal(gameMap)
    if err != nil {
        return "", "", err
    }

    overlayB, err := json.Marshal(overlay)
    if err != nil {
        return "", "", err
    }

    return string(gameMapB), string(overlayB), nil
}

// Init map slice, using a filler (nothing or grass)
func genPhase1(resX int64, resY int64, empty bool) []*slot {
	filler := int64(1)
	if empty {
		filler = int64(0)
	}

	res := []*slot{}
	for y := int64(0); y < resY; y++ {
		for x := int64(0); x < resX; x++ {
			res = append(res, &slot{Terrain: filler})
		}
	}
    return res
}

// Add biomes using decay-based generation
func genPhase2(gameMap, overlay []*slot, nSpawns, resX, resY int64) ([]*slot, []*slot) {
	rand.Seed(time.Now().UnixNano())
	var biomes = []biome{
		biome{Type: 2, Decay: 100},
		biome{Type: 2, Decay: 500},
		biome{Type: 3, Decay: 100},
		biome{Type: 4, Decay: 250},
		biome{Type: 5, Decay: 250},
	}
	var biomeOverlay = map[int64][]int64{
		2: []int64{81,},
	}
	dirs := []int64{-1 * resX, 1, resX, -1}

	biomeLen := int64(len(biomes))
	dirLen := int64(len(dirs))

	for b := int64(0); b < nSpawns; b++ {
		sources := []int64{makeSource(resX, resY)}

		newBiome := biomes[rand.Int63n(biomeLen)]
		for newBiome.Decay > int64(0) {
			sourceLen := int64(len(sources))
			newSource := sources[rand.Int63n(sourceLen)] + dirs[rand.Int63n(dirLen)]
			if newSource < int64(0) {
				newSource = int64(0)
			} else if newSource > (resX * resY - int64(1)) {
				newSource = (resX * resY - int64(1))
			}
			if !itemInSlice(sources, newSource) {
				sources = append(sources, newSource)
				gameMap[newSource] = &slot{Terrain: newBiome.Type}
			}
			newBiome.Decay--
		}
	}

	for i := range gameMap {
		if v, ok := biomeOverlay[gameMap[i].Terrain]; ok {
			vLen := int64(len(v))
			overlay[i] = &slot{Terrain: v[rand.Int63n(vLen)]}
		}
	}
	return gameMap, overlay
}

func itemInSlice(sl []int64, item int64) bool {
	for i := range sl {
		if sl[i] == item {
			return true
		}
	}
	return false
}

// Creates a new origin for a biome to be generated from
func makeSource(resX, resY int64) int64 {
	return rand.Int63n(resX) * (resY - rand.Int63n(resY))
}
