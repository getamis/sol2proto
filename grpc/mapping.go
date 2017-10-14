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
		Name: in.Name,
	}

	arg.Type, arg.IsSlice = toGrpcType(in.Type)
	return arg
}

func toGrpcType(t abi.Type) (string, bool) {
	switch t.T {
	case abi.IntTy:
		if t.Size == 8 {
			return "byte", false
		} else if t.Size == 32 {
			return "int32", false
		} else if t.Size == 64 {
			return "int64", false
		}
		return "bytes", false
	case abi.UintTy:
		if t.Size == 8 {
			return "byte", false
		} else if t.Size == 32 {
			return "uint32", false
		} else if t.Size == 64 {
			return "uint64", false
		}
		return "bytes", false
	case abi.BoolTy:
		return "bool", false
	case abi.StringTy:
		return "string", false
	case abi.SliceTy:
		elemType, _ := toGrpcType(*t.Elem)
		return elemType, true
	case abi.AddressTy:
		return "string", false
	case abi.FixedBytesTy:
		return "bytes", false
	case abi.BytesTy:
		return "bytes", false
	case abi.HashTy:
		return "string", false
	case abi.FixedPointTy:
	case abi.FunctionTy:
		fallthrough
	default:
	}

	return "bytes", false
}

func ToMessage(name string, args []Argument) Message {
	return Message{
		Name: name,
		Args: args,
	}
}
