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
)

// Parse gRPC methods and required message types from methods in an Ethereum contract ABI.
func ParseMethods(abiMethods map[string]abi.Method) (methods Methods, msgs []Message) {
	for _, f := range abiMethods {
		method, msg := ParseMethod(f)
		methods = append(methods, method)
		msgs = append(msgs, msg...)
	}

	return
}

// Parse gRPC method and required message types from an Ethereum contract method.
func ParseMethod(m abi.Method) (Method, []Message) {
	method := Method{
		Const: m.Const,
		Name:  m.Name,
	}

	method.Inputs = append(method.Inputs, parseArgs(m.Inputs)...)
	method.Outputs = append(method.Outputs, parseArgs(m.Outputs)...)

	// If it is not a const method, we need to provide
	// more transaction options to send transactions.
	if !m.Const {
		method.Inputs = append(method.Inputs, Argument{
			Name:    "opts",
			Type:    TransactOptsReq.Name,
			IsSlice: false,
		})
	}

	var requiredMessages []Message
	if len(method.Inputs) > 0 {
		requiredMessages = append(requiredMessages, ToMessage(method.RequestName(), method.Inputs))
	}
	if len(method.Outputs) > 0 {
		requiredMessages = append(requiredMessages, ToMessage(method.RequestName(), method.Outputs))
	}

	return method, requiredMessages
}

func parseArgs(args []abi.Argument) (results []Argument) {
	for _, arg := range args {
		results = append(results, ToGrpcArgument(arg))
	}

	return
}
