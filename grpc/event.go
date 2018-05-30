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

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/getamis/sirius/util"
)

func ParseEvents(abiEvents map[string]abi.Event) (methods Methods, msgs []Message) {
	for _, ev := range abiEvents {
		method, msg := ParseEvent(ev)
		methods = append(methods, method)
		msgs = append(msgs, msg...)
	}

	return
}

func ParseEvent(ev abi.Event) (Method, []Message) {
	method := Method{}

	if ev.Anonymous {
		method.Name = "onEvent" + util.ToCamelCase(ev.Id().Hex())
	} else {
		method.Name = "on" + ev.Name
	}

	method.Inputs = append(method.Inputs, parseArgs(ev.Inputs)...)

	var requiredMessages []Message
	if len(method.Inputs) > 0 {
		requiredMessages = append(requiredMessages, ToMessage(method.RequestName(), method.Inputs))
	}

	return method, requiredMessages
}
