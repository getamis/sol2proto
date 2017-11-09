// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package grpc

import "github.com/ethereum/go-ethereum/accounts/abi"

func ToGrpcArgument(in abi.Argument) Argument {
	arg := Argument{
		Name:    in.Name,
		IsSlice: in.Type.IsSlice,
	}

	arg.Type = toGrpcType(in.Type)
	return arg
}

func toGrpcType(t abi.Type) string {
	switch t.T {
	case abi.IntTy:
		if t.Size == 8 {
			return "byte"
		} else if t.Size == 32 {
			return "int32"
		} else if t.Size == 64 {
			return "int64"
		}
		return "bytes"
	case abi.UintTy:
		if t.Size == 8 {
			return "byte"
		} else if t.Size == 32 {
			return "uint32"
		} else if t.Size == 64 {
			return "uint64"
		}
		return "bytes"
	case abi.BoolTy:
		return "bool"
	case abi.StringTy:
		return "string"
	case abi.AddressTy:
		return "string"
	case abi.FixedBytesTy:
		return "bytes"
	case abi.BytesTy:
		return "bytes"
	case abi.HashTy:
		return "string"
	case abi.FixedPointTy:
	case abi.FunctionTy:
		fallthrough
	default:
	}

	return "bytes"
}

func ToMessage(name string, args []Argument) Message {
	return Message{
		Name: name,
		Args: args,
	}
}
