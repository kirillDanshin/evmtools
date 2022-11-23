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
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/vm"
)

func TestListCompleteness(t *testing.T) {
	for i := 0; i < 256; i++ {

		op := vm.OpCode(i)
		instr := InstructionSet[i]
		instrUndefined := strings.Contains(instr.Mnemonic, "__UNDEFINED_INSTRUCTION")

		if strings.Contains(op.String(), "not defined") {
			if !instrUndefined {
				t.Errorf("opcode %d should not be defined but defined as %q", i, instr.Mnemonic)
			}
			continue
		}

		if instrUndefined {
			t.Errorf("opcode %d should be defined but not defined, expected %q", i, op.String())
		}
	}
}
