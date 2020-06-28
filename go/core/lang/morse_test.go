/*
Copyright (C) 2019 Graeme Sutherland, Nodestone Limited


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package lang

import (
	"gotest.tools/assert"
	"testing"
)

func TestSearchMorseTreeAllSymbols(t *testing.T) {
	tree := MorseTree()
	// go through all characters and check map search works for all
	for k, v := range MorseMap {
		result := SearchMorseTree(v, tree)
		assert.Equal(t, k, result)
	}
}

func TestSearchMorseTreeBadSymbol(t *testing.T) {
	tree := MorseTree()
	assert.Equal(t, "", SearchMorseTree("-..--", tree))
}

func TestSearchMorseTreeNoMorse(t *testing.T) {
	tree := MorseTree()
	assert.Equal(t, "", SearchMorseTree("", tree))
}
