// Copyright 2022 Kyrylo Danshyn
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package evmops

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type Opcode byte

type Line struct {
	// Inst is the instruction
	Inst Instruction

	ProgramCounter uint64

	Args []string

	// Comment is a comment about the operation, that should be displayed on the same line as the operation.
	Comment string
}

func (l *Line) String() string {
	str := ""

	if l.Inst.Mnemonic != "" {
		str += l.Inst.Mnemonic
	}

	if len(l.Args) > 0 {
		str += " " + joinAsHex(l.Args, ", ")
	}

	l.Comment += fmt.Sprintf(" pc=%d", l.ProgramCounter)

	str += " ; " + l.Comment

	return str
}

func joinAsHex(elems []string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return hex.EncodeToString([]byte(elems[0]))
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(hex.EncodeToString([]byte(elems[0])))
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(hex.EncodeToString([]byte(s)))
	}
	return b.String()
}
