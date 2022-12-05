package evmfuncs

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestFuncSig_UnpackInput(t *testing.T) {
	type fields struct {
		name         string
		inputs       []FuncParam
		outputs      []FuncParam
		unescapedSel string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "erc20/transfer",
			fields: fields{
				name:         "transfer",
				inputs:       funcDescriptions["transfer(address to, uint256 amount)"].inputs,
				outputs:      funcDescriptions["transfer(address to, uint256 amount)"].outputs,
				unescapedSel: "transfer(address to, uint256 amount)",
			},
			args: args{
				data: func() []byte {
					data, err := hex.DecodeString("0000000000000000000000005a5b644fb1a3ca046317fe82bc695fff7bacf30c0000000000000000000000000000000000000000000002a568d6215ac1400000")
					if err != nil {
						t.Fatal(err)
					}
					return data
				}(),
			},
			want: []interface{}{
				common.HexToAddress("0x5A5b644FB1A3ca046317fE82BC695FfF7bACF30C"),
				func() *big.Int {
					x, ok := big.NewInt(0).SetString("12496000000000000000000", 10)
					if !ok {
						t.Fatal("failed to parse big int")
					}
					return x
				}(),
			},
		},
		{
			name: "erc20/transfer_with_method_id",
			fields: fields{
				name:         "transfer",
				inputs:       funcDescriptions["transfer(address to, uint256 amount)"].inputs,
				outputs:      funcDescriptions["transfer(address to, uint256 amount)"].outputs,
				unescapedSel: "transfer(address to, uint256 amount)",
			},
			args: args{
				data: func() []byte {
					data, err := hex.DecodeString("a9059cbb0000000000000000000000005a5b644fb1a3ca046317fe82bc695fff7bacf30c0000000000000000000000000000000000000000000002a568d6215ac1400000")
					if err != nil {
						t.Fatal(err)
					}
					return data
				}(),
			},
			want: []interface{}{
				common.HexToAddress("0x5A5b644FB1A3ca046317fE82BC695FfF7bACF30C"),
				func() *big.Int {
					x, ok := big.NewInt(0).SetString("12496000000000000000000", 10)
					if !ok {
						t.Fatal("failed to parse big int")
					}
					return x
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsig := &FuncSig{
				name:         tt.fields.name,
				inputs:       tt.fields.inputs,
				outputs:      tt.fields.outputs,
				unescapedSel: tt.fields.unescapedSel,
			}

			got, err := fsig.UnpackInput(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FuncSig.UnpackInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FuncSig.UnpackInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFuncSigParser(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		sel     string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "erc20/transfer",
			sel:  "transfer(address to, uint256 amount)",
			args: args{
				data: func() []byte {
					data, err := hex.DecodeString("0000000000000000000000005a5b644fb1a3ca046317fe82bc695fff7bacf30c0000000000000000000000000000000000000000000002a568d6215ac1400000")
					if err != nil {
						t.Fatal(err)
					}
					return data
				}(),
			},
			want: []interface{}{
				common.HexToAddress("0x5A5b644FB1A3ca046317fE82BC695FfF7bACF30C"),
				func() *big.Int {
					x, ok := big.NewInt(0).SetString("12496000000000000000000", 10)
					if !ok {
						t.Fatal("failed to parse big int")
					}
					return x
				}(),
			},
		},
		{
			name: "erc20/transfer_with_method_id",
			sel:  "transfer(address to, uint256 amount)",
			args: args{
				data: func() []byte {
					data, err := hex.DecodeString("a9059cbb0000000000000000000000005a5b644fb1a3ca046317fe82bc695fff7bacf30c0000000000000000000000000000000000000000000002a568d6215ac1400000")
					if err != nil {
						t.Fatal(err)
					}
					return data
				}(),
			},
			want: []interface{}{
				common.HexToAddress("0x5A5b644FB1A3ca046317fE82BC695FfF7bACF30C"),
				func() *big.Int {
					x, ok := big.NewInt(0).SetString("12496000000000000000000", 10)
					if !ok {
						t.Fatal("failed to parse big int")
					}
					return x
				}(),
			},
		},
		{
			name: "erc20/transfer_anonymous",
			sel:  "transfer(address, uint256)",
			args: args{
				data: func() []byte {
					data, err := hex.DecodeString("a9059cbb0000000000000000000000005a5b644fb1a3ca046317fe82bc695fff7bacf30c0000000000000000000000000000000000000000000002a568d6215ac1400000")
					if err != nil {
						t.Fatal(err)
					}
					return data
				}(),
			},
			want: []interface{}{
				common.HexToAddress("0x5A5b644FB1A3ca046317fE82BC695FfF7bACF30C"),
				func() *big.Int {
					x, ok := big.NewInt(0).SetString("12496000000000000000000", 10)
					if !ok {
						t.Fatal("failed to parse big int")
					}
					return x
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsig, err := NewFuncSignatureFromString(tt.sel)
			if err != nil {
				t.Fatal(err)
			}

			got, err := fsig.UnpackInput(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FuncSig.UnpackInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FuncSig.UnpackInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
