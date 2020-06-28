/*
Copyright (C) 2020 Graeme Sutherland, Nodestone Limited


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
package hw

import (
	"github.com/G0WCZ/cwc/config"
)

type MockGeneral struct {
	Config      *config.Config
	adapterName string
	statusState map[string]string
	name        string
}

func (G *MockGeneral) Open() error {
	return nil
}

func (G *MockGeneral) ConfigChanged() error {
	return nil
}

func (G *MockGeneral) SetStatus(name string, value string) {
	G.statusState[name] = value
}

func (G *MockGeneral) GetStatus(name string) string {
	return G.statusState[name]
}

func (G *MockGeneral) Close() error {
	return nil
}

func (G *MockGeneral) Name() string {
	return G.Name()
}

func NewMockGeneral(config *config.Config, name string, adapterName string) GeneralIO {
	return &MockGeneral{Config: config,
		adapterName: adapterName,
		statusState: make(map[string]string),
		name:        name,
	}
}
