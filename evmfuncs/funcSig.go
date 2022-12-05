package evmfuncs

import (
	"bytes"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/kirillDanshin/evmtools"
)

type FuncParam struct {
	Name string
	Type string
}

type FuncSig struct {
	// Name is the name of the function
	name string

	// Inputs is the list of input parameters
	inputs []FuncParam

	// Outputs is the list of output parameters
	outputs []FuncParam

	unescapedSel string
}

type FuncSignature interface {
	// Name is the name of the function
	Name() string

	// Inputs is the list of input parameters
	Inputs() []FuncParam

	// Outputs is the list of output parameters
	Outputs() []FuncParam

	// String returns the string representation of the function signature
	String() string

	// WellKnown returns true if the function signature is well known or commonly used
	// in the Ethereum ecosystem (e.g. balanceOf, transfer, etc.)
	WellKnown() bool

	// Effects returns the list of side effects of the function call,
	// if given function signature is known to have specific side effects.
	Effects() Effect

	// Describe returns a string describing the function signature, including a functional description
	// if available
	Describe() string

	// UnpackInput upacks values from given function data and returns a list of Go values.
	// The function signature must match the function data.
	UnpackInput(data []byte) ([]interface{}, error)

	// UnpackOutput upacks values from given function data and returns a list of Go values.
	// The function signature must match the function data.
	UnpackOutput(data []byte) ([]interface{}, error)
}

func (fsig *FuncSig) Name() string {
	return fsig.name
}

func (fsig *FuncSig) Inputs() []FuncParam {
	return fsig.inputs
}

func (fsig *FuncSig) Outputs() []FuncParam {
	if fsig.outputs == nil && fsig.WellKnown() {
		return funcDescriptions[fsig.lookupWellKnownKey()].outputs
	}

	return fsig.outputs
}

func (fsig *FuncSig) lookupWellKnownKey() string {
	hasFuncDesc := func(s string) bool {
		_, ok := funcDescriptions[strings.TrimSpace(s)]

		return ok
	}

	if hasFuncDesc(fsig.String()) {
		return fsig.String()
	}

	if hasFuncDesc(fsig.unescapedSel) {
		return fsig.unescapedSel
	}

	dedupSpaces := func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	}

	respaced := dedupSpaces(fsig.String())
	if hasFuncDesc(respaced) {
		return respaced
	}

	anonymized := fsig.removeParamNamesFromSelector()

	if hasFuncDesc(anonymized) {
		return anonymized
	}

	return strings.ToLower(fsig.name)
}

func (fsig *FuncSig) removeParamNamesFromSelector() string {
	out := fsig.name + "("
	for i, param := range fsig.Inputs() {
		out += param.Type
		if i < len(fsig.Inputs())-1 {
			out += ","
		}
	}
	out += ")"

	return out
}

func (fsig *FuncSig) String() string {
	sig := fsig.Name()

	sig += "("
	for i, input := range fsig.Inputs() {
		sig += input.Type
		if i < len(fsig.Inputs())-1 {
			sig += ","
		}
	}
	sig += ")"

	return sig
}

func (fsig *FuncSig) Describe() string {
	desc := fsig.Name()

	desc += "("
	for i, input := range fsig.Inputs() {
		desc += input.Type + " " + input.Name
		if i < len(fsig.Inputs())-1 {
			desc += ","
		}
	}
	desc += ")"
	outs := fsig.Outputs()
	if len(outs) > 0 {
		desc += " returns ("
		for i, output := range fsig.Outputs() {
			desc += output.Type + " " + output.Name
			if i < len(fsig.Outputs())-1 {
				desc += ","
			}
		}
		desc += ")"
	}

	if description, ok := funcDescriptions[fsig.String()]; ok {
		desc += " // " + description.description
	}

	return desc
}

func (fsig *FuncSig) Effects() Effect {
	if description, ok := funcDescriptions[fsig.String()]; ok {
		return description.effects
	}

	return EffectUnknown
}

func (fsig *FuncSig) WellKnown() bool {
	_, ok := funcDescriptions[fsig.String()]

	return ok
}

func (fsig *FuncSig) unpackArgs(args abi.Arguments, data []byte, tryStripMethodID bool) ([]interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}

	sig := fsig.removeParamNamesFromSelector()
	methodID := evmtools.MethodID(sig)

	if !tryStripMethodID || !bytes.HasPrefix(data, methodID) {
		return args.UnpackValues(data)
	}

	x, err := args.UnpackValues(data[4:])

	if err != nil {
		return nil, err
	}

	return x, nil
}

func (fsig *FuncSig) UnpackInput(data []byte) ([]interface{}, error) {
	fsigABI, err := fsig.ABI()
	if err != nil {
		return nil, err
	}

	return fsig.unpackArgs(fsigABI.Methods[fsig.Name()].Inputs, data, true)
}

func (fsig *FuncSig) UnpackOutput(data []byte) ([]interface{}, error) {
	fsigABI, err := fsig.ABI()
	if err != nil {
		return nil, err
	}

	return fsig.unpackArgs(fsigABI.Methods[fsig.Name()].Outputs, data, false)
}

// ABI creates a fake ABI object for the function signature
func (fsig *FuncSig) ABI() (*abi.ABI, error) {
	abiInputs := abi.Arguments{}

	for _, param := range fsig.inputs {
		abiType, err := abi.NewType(param.Type, "", nil)
		if err != nil {
			return nil, err
		}

		abiInputs = append(abiInputs, abi.Argument{
			Name:    param.Name,
			Type:    abiType,
			Indexed: false,
		})
	}

	abiOutputs := abi.Arguments{}

	for _, param := range fsig.outputs {
		abiType, err := abi.NewType(param.Type, "", nil)
		if err != nil {
			return nil, err
		}

		abiOutputs = append(abiOutputs, abi.Argument{
			Name:    param.Name,
			Type:    abiType,
			Indexed: false,
		})
	}

	mutability := ""
	if fsig.Effects() != EffectUnknown {
		if fsig.Effects() == EffectRead {
			mutability = "view"
		}
	}

	return &abi.ABI{
		Methods: map[string]abi.Method{
			fsig.Name(): abi.NewMethod(fsig.Name(), fsig.Name(), abi.Function, mutability, false, true, abiInputs, abiOutputs),
		},
	}, nil
}

// parseFuncSig parses a function signature string and returns a FuncSig object
func parseFuncSig(rawSig string) *FuncSig {
	name, sig := parseName(rawSig)

	inputs, sig := parseInputs(sig)

	outputs, _ := parseOutputs(sig)

	return &FuncSig{
		name:         name,
		inputs:       inputs,
		outputs:      outputs,
		unescapedSel: rawSig,
	}
}

func parseName(sig string) (string, string) {
	var name string
	if strings.Contains(sig, "(") {
		name = sig[:strings.Index(sig, "(")]
		sig = sig[strings.Index(sig, "("):]
	} else {
		name = sig
		sig = ""
	}

	return name, sig
}

func parseInputs(sig string) ([]FuncParam, string) {
	var inputs []FuncParam

	if strings.Contains(sig, "(") {
		sig = sig[strings.Index(sig, "(")+1:]
		if strings.Contains(sig, ")") {
			sig = sig[:strings.Index(sig, ")")]
		}
		if sig != "" {
			inputs = parseParams(sig)
		}
	}

	return inputs, sig
}

func parseOutputs(sig string) ([]FuncParam, string) {
	var outputs []FuncParam

	if strings.Contains(sig, "returns") {
		sig = sig[strings.Index(sig, "returns")+7:]
		if strings.Contains(sig, ")") {
			sig = sig[:strings.Index(sig, ")")]
		}
		if sig != "" {
			outputs = parseParams(sig)
		}
	}

	return outputs, sig
}

func parseParams(sig string) []FuncParam {
	var params []FuncParam

	paramStrs := strings.Split(sig, ",")
	for _, paramStr := range paramStrs {
		paramStr = strings.TrimSpace(paramStr)
		param := parseParam(paramStr)
		params = append(params, param)
	}

	return params
}

func parseParam(paramStr string) FuncParam {
	var param FuncParam

	if strings.Contains(paramStr, " ") {
		param.Type = paramStr[:strings.Index(paramStr, " ")]
		param.Name = paramStr[strings.Index(paramStr, " ")+1:]
	} else {
		param.Type = paramStr
	}

	return param
}

func NewFuncSignatureFromString(sig string) (FuncSignature, error) {
	if wellKnown, ok := GetWellKnownFuncBySig(sig); ok {
		return &FuncSig{
			name:         wellKnown.name,
			inputs:       wellKnown.inputs,
			outputs:      wellKnown.outputs,
			unescapedSel: sig,
		}, nil
	}

	fsig := parseFuncSig(sig)

	fsig.unescapedSel = sig

	return fsig, nil
}
