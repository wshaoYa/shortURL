package base62

import "testing"

func TestUint64ToBase62(t *testing.T) {
	type args struct {
		x uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		//{name: "normal", args: args{x: 6347}, want: "1En"},
		//{name: "normal", args: args{x: 2}, want: "2"},
		//{name: "abnormal", args: args{x: 1}, want: "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint64ToBase62(tt.args.x); got != tt.want {
				t.Errorf("Uint64ToBase62() = %v, want %v", got, tt.want)
			}
		})
	}
}
