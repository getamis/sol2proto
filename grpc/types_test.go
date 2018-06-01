// Copyright 2018 AMIS Technologies
// This file is part of the sol2proto
//
// The sol2proto is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The sol2proto is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the sol2proto. If not, see <http://www.gnu.org/licenses/>.
package grpc

import "testing"

func TestMethod(t *testing.T) {
	m := Method{
		Name: "Test",
		Inputs: []Argument{
			Argument{
				Name: "input1",
				Type: "uint",
			},
			Argument{
				Name: "input2",
				Type: "*big.Int",
			},
		},
		Outputs: []Argument{
			Argument{
				Name: "output1",
				Type: "Output",
			},
			Argument{
				Name:    "output2",
				Type:    "float64",
				IsSlice: true,
			},
		},
	}

	t.Log(m.String())
}
