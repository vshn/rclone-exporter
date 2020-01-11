package main

import "testing"

func TestGetFriendlyUrlString(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"should leave URL alone",
			args{url: "http://rclone"},
			"http://rclone",
		},
		{
			"should remove User from URL",
			args{url: "http://user:pass@rclone:5572"},
			"http://rclone:5572",
		},
		{
			"should not parse",
			args{url: "user:pass:rclone:5572"},
			"user:pass:rclone:5572",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFriendlyUrlString(tt.args.url); got != tt.want {
				t.Errorf("GetFriendlyUrlString() = %v, want %v", got, tt.want)
			}
		})
	}
}
