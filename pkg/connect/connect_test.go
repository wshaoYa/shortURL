package connect

import "testing"

func TestPing(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "正常", args: args{target: "https://www.baidu.com/"}, want: true},
		{name: "不正常", args: args{target: "https://www.baidu22222222124142342342.com"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ping(tt.args.target); got != tt.want {
				t.Errorf("Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}
