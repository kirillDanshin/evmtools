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

package evmdis

import (
	"bytes"
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/core/asm"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/kirillDanshin/evmtools/evmops"
)

type Disassembler struct {
}

func NewDisassembler() *Disassembler {
	return &Disassembler{}
}

type Results struct {
	CompiledLen     int
	Lines           []evmops.Line
	FoundSignatures map[string]struct{}
}

func (r *Results) String() string {
	var buf bytes.Buffer
	for _, line := range r.Lines {
		buf.WriteString(line.String())
		buf.WriteByte('\n')
	}
	return buf.String()
}

var ERC20Iface = []string{
	"totalSupply()",
	"balanceOf(address)",
	"allowance(address,address)",
	"transfer(address,uint256)",
	"approve(address,uint256)",
	"transferFrom(address,address,uint256)",
}

var ERC721Iface = []string{
	"balanceOf(address)",
	"ownerOf(uint256)",
	"safeTransferFrom(address,address,uint256)",
	"safeTransferFrom(address,address,uint256,bytes)",
	"transferFrom(address,address,uint256)",
	"approve(address,uint256)",
	"setApprovalForAll(address,bool)",
	"getApproved(uint256)",
	"isApprovedForAll(address,address)",
}

func (r *Results) Implements(iface []string) bool {
	for _, sig := range iface {
		if _, ok := r.FoundSignatures[sig]; !ok {
			return false
		}
	}

	return true
}

// Disassemble disassembles the given bytecode and returns the disassembled code as lines.
// It will return results even if the bytecode is invalid,
// for cases like ENS, where the geth's asm package fails to disassemble the bytecode,
// but it the results still are sufficient to determine that the ENS contract supports ERC721.
func (d *Disassembler) Disassemble(code string) (*Results, error) {
	script, err := hex.DecodeString(code)
	if err != nil {
		return nil, err
	}

	sigs := map[string]struct{}{}
	fourBytesToSigs := map[string][]string{}
	lines := make([]evmops.Line, 0)

	it := asm.NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			if len(it.Arg()) == 4 && it.Op() == vm.PUSH4 {
				comment := ""
				hexArg := hex.EncodeToString(it.Arg())
				if _, ok := fourBytesToSigs[hexArg]; !ok {
					signatures := getSigsFromFourBytes(it.Arg())
					for _, sig := range signatures {
						sigs[sig] = struct{}{}
						fourBytesToSigs[hexArg] = append(fourBytesToSigs[hexArg], sig)
						comment += " " + sig
					}
				} else {
					comment = strings.Join(fourBytesToSigs[hexArg], ", ")
				}

				lines = append(lines, evmops.Line{
					Inst:           evmops.InstructionSet[it.Op()],
					Args:           []string{string(it.Arg())},
					ProgramCounter: it.PC(),
					Comment:        comment,
				})
				continue
			}

			if it.Op() == vm.PUSH32 {
				ascii := ""
				pos := bytes.IndexByte(it.Arg(), 0)
				ok := true
				for i := 0; i < pos; i++ {
					ok = ok && it.Arg()[i] < 0x80
				}
				if ok {
					ascii = string(it.Arg())
					if pos >= 0 {
						ascii = ascii[:pos]
					}
				}

				hexArg := hex.EncodeToString([]byte(ascii))

				lines = append(lines, evmops.Line{
					Inst:           evmops.InstructionSet[it.Op()],
					Args:           []string{hexArg},
					ProgramCounter: it.PC(),
				})
				continue
			}
		}

		lines = append(lines, evmops.Line{
			Inst:           evmops.InstructionSet[it.Op()],
			ProgramCounter: it.PC(),
		})
	}

	if err := it.Error(); err != nil {
		return &Results{
			Lines:           lines,
			FoundSignatures: sigs,
			CompiledLen:     hex.DecodedLen(len(code)),
		}, err
	}

	return &Results{
		Lines:           lines,
		FoundSignatures: sigs,
		CompiledLen:     hex.DecodedLen(len(code)),
	}, nil
}
