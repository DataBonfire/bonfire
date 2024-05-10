package utils

import "testing"

// @author: Haxqer
// @email: haxqer666@gmail.com
// @since: 10/19/23
// @desc: TODO

func TestHashPassword(t *testing.T) {
	type args struct {
		s    string
		salt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "testcase-01", args: args{s: "123456", salt: "kolplanet"}, want: "59ce4af11891b33b6bbb9cd7206b44e3ffad1360"},
		{name: "testcase-02", args: args{s: "5#cX4O#f1Z", salt: "kolplanet"}, want: "013807b8aa936110341e6818986806ee42a67b83"},
		{name: "testcase-03", args: args{s: "81Kad1d345Ze1", salt: "kolplanet"}, want: "94a2dde362cc13104497cdf640e802337339d903"},
		{name: "testcase-04", args: args{s: "a1adKad2d345Za0", salt: "kolplanet"}, want: "76f1656d95ed47e556d85c6c20fd0f377a300539"},
		{name: "testcase-05", args: args{s: "229effc05", salt: "kolplanet"}, want: "0f9befa59b62f9ef82e21e85d8a2562a9c882f17"},
		{name: "testcase-06", args: args{s: "2kshn4ksd5", salt: "kolplanet"}, want: "4cb07e62f6112f8ba7dc17e3dc68662ed221afa9"},
		{name: "testcase-07", args: args{s: "sjsnk1n4ksd5", salt: "kolplanet"}, want: "63ae37347ad40b7c00eef1f2e426bdf815ba510d"},
		{name: "testcase-08", args: args{s: "12345678", salt: "kolplanet"}, want: "b35b5bf5465f69c72a193df44d8536a9eb0ed3f5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashPassword(tt.args.s, tt.args.salt); got != tt.want {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
