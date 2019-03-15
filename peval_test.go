package eval

import (
	"fmt"
	"reflect"
	"testing"
)

func TestHandValue(t *testing.T) {
	tests := []struct {
		name  string
		args  Cards
		want  string
		want1 int
		want2 []uint
	}{
		{
			"twoCardsFail",
			Cards("ab"),
			"",
			0,
			nil,
		},
		{
			"tooManyCards>7",
			Cards("abcdeftn"),
			"",
			0,
			nil,
		},
		{
			"repeatedCards",
			Cards("abcan"),
			"",
			0,
			nil,
		},
		{
			"twoPairs",
			Cards("ZSNFa"),
			"Two Pairs",
			3205,
			[]uint{268442665, 2102541, 69634, 2106637, 98306},
		},
		{
			"flush6cards",
			Cards("lmacifd"),
			"Flush",
			417,
			[]uint{134253349, 268471337, 2131213, 557831, 16812055},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := HandValue(tt.args)
			if got != tt.want {
				t.Errorf("HandValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("HandValue() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("HandValue() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestEvalBetterHands(t *testing.T) {
	tests := []struct {
		name string
		args Cards
		//wantCombs    map[string][]gopo0.Cards
		//wantOutsOuts map[string][]gopo0.Cards
	}{
		{
			"KingHStFlush",
			Cards("yxWVU"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotCombs, gotOutsOuts := BetterHands(tt.args)
			fmt.Printf("Combs: %v\n %v \n OutsOuts: %v \n %v \n", gotCombs, len(gotCombs), gotOutsOuts, len(gotOutsOuts))
		})
	}
}

func TestBitsToChars(t *testing.T) {
	type args struct {
		uns []uint
	}
	tests := []struct {
		name string
		args args
		want Cards
	}{
		{
			"aceKingJack",
			args{[]uint{268442665, 134224677, 33560861}},
			Cards("ZYW"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BitsToChars(tt.args.uns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitsToChars() = %v, want %v", got, tt.want)
			}
		})
	}
}
