package evmtools

import (
	"reflect"
	"testing"
)

func TestMethodID(t *testing.T) {
	type args struct {
		signature string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "totalSupply()",
			args: args{
				signature: "totalSupply()",
			},
			want: []byte{0x18, 0x16, 0x0d, 0xdd},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MethodID(tt.args.signature); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MethodID() = %x, want %x", got, tt.want)
			}
		})
	}
}
