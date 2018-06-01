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
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/getamis/sirius/util"
)

type Argument struct {
	Name    string
	Type    string
	IsSlice bool
}

func (arg Argument) String() string {
	if arg.IsSlice {
		arg.Type = "repeated " + arg.Type
	}
	if arg.Name == "" {
		arg.Name = "arg"
	}
	return strings.TrimSpace(strings.Join([]string{
		arg.Type,
		util.ToUnderScore(arg.Name),
	}, " "))
}

type Method struct {
	Const   bool
	Name    string
	Inputs  []Argument
	Outputs []Argument
}

var methodTemplate = `rpc {{ .Name }}({{ ToInputMsg }}) returns ({{ ToOutputMsg }}) {}`

func (m Method) String() string {
	tmpl, err := template.New("method").
		Funcs(template.FuncMap(
			map[string]interface{}{
				"ToInputMsg": func() string {
					if len(m.Inputs) > 0 {
						return m.RequestName()
					} else if m.Const {
						return Empty.Name
					}
					return TransactionReq.Name
				},
				"ToOutputMsg": func() string {
					// if it's not a const method, we return
					// the transaction hash
					if m.Const {
						if len(m.Outputs) > 0 {
							return m.ResponseName()
						}
						return Empty.Name
					}
					return TransactionResp.Name
				},
			})).Parse(methodTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, m)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}

	return result.String()
}

func (m Method) RequestName() string {
	return util.ToCamelCase(m.Name) + "Req"
}

func (m Method) ResponseName() string {
	return util.ToCamelCase(m.Name) + "Resp"
}

type Methods []Method

// Len is part of sort.Interface.
func (m Methods) Len() int {
	return len(m)
}

// Swap is part of sort.Interface.
func (m Methods) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (m Methods) Less(i, j int) bool {
	return strings.Compare(m[i].Name, m[j].Name) < 0
}

var TransactionReq = Message{
	Name: "TransactionReq",
	Args: []Argument{
		{
			Name:    "opts",
			Type:    TransactOptsReq.Name,
			IsSlice: false,
		},
	},
}

var TransactionResp = Message{
	Name: "TransactionResp",
	Args: []Argument{
		{
			Name:    "hash",
			IsSlice: false,
			Type:    "string",
		},
	},
}

var Empty = Message{
	Name: "Empty",
}

var TransactOptsReq = Message{
	Name: "TransactOpts",
	Args: []Argument{
		{
			Name:    "private_key",
			IsSlice: false,
			Type:    "string",
		},
		{
			Name:    "nonce",
			IsSlice: false,
			Type:    "int64",
		},
		{
			Name:    "value",
			IsSlice: false,
			Type:    "int64",
		},
		{
			Name:    "gas_price",
			IsSlice: false,
			Type:    "int64",
		},
		{
			Name:    "gas_limit",
			IsSlice: false,
			Type:    "int64",
		},
	},
}

type Message struct {
	Name string
	Args []Argument
}

var messageTemplate = `message {{ .Name }} {
{{ PrintArgs .Args -}}
}`

func (m Message) String() string {
	tmpl, err := template.New("message").
		Funcs(template.FuncMap(
			map[string]interface{}{
				"PrintArgs": func(args []Argument) (result string) {
					for index, arg := range args {
						result = result + "    " + arg.String() + " = " + fmt.Sprintf("%d", index+1) + ";\n"
					}
					return result
				},
			})).Parse(messageTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, m)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}

	return result.String()
}

type Messages []Message

// Len is part of sort.Interface.
func (m Messages) Len() int {
	return len(m)
}

// Swap is part of sort.Interface.
func (m Messages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (m Messages) Less(i, j int) bool {
	return strings.Compare(m[i].Name, m[j].Name) < 0
}

type Sources []string

// Len is part of sort.Interface.
func (s Sources) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s Sources) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface.
func (s Sources) Less(i, j int) bool {
	return strings.Compare(s[i], s[j]) < 0
}
