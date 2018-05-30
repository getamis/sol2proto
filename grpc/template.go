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

import "strings"

type Service struct {
	Package  string
	Name     string
	Methods  Methods
	Events   Methods
	Messages map[string]Message
	Sources  Sources
}

var ServiceTemplate string = `// Automatically generated by sol2proto. DO NOT EDIT!
// sources: {{ range .Sources }}
//     {{ . }}
{{- end }}
syntax = "proto3";

package {{ .Package }};

import "messages.proto";
service {{ .Name }} {
{{- range .Methods }}
    {{ . }}
{{- end }}

    // Not supported yet
{{- range .Events }}
    // {{ . }}
{{- end }}
}
`

var MessagesTemplate string = `// Automatically generated by sol2proto. DO NOT EDIT!
// sources: {{ range .Sources }}
//     {{ . }}
{{- end }}
syntax = "proto3";
package {{ .Package }};

import public "github.com/getamis/sol2proto/pb/messages.proto";
{{ range .Messages }}
{{ . }}
{{ end }}
`

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
