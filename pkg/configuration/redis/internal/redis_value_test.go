package internal

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import "testing"

func TestGetRedisValueAndVersion(t *testing.T) {
	type args struct {
		redisValue string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "empty value",
			args: args{
				redisValue: "",
			},
			want:  "",
			want1: "",
		},
		{
			name: "value without version",
			args: args{
				redisValue: "mockValue",
			},
			want:  "mockValue",
			want1: "",
		},
		{
			name: "value without version",
			args: args{
				redisValue: "mockValue||",
			},
			want:  "mockValue",
			want1: "",
		},
		{
			name: "value with version",
			args: args{
				redisValue: "mockValue||v1.0.0",
			},
			want:  "mockValue",
			want1: "v1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetRedisValueAndVersion(tt.args.redisValue)
			if got != tt.want {
				t.Errorf("GetRedisValueAndVersion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetRedisValueAndVersion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParseRedisKeyFromEvent(t *testing.T) {
	type args struct {
		eventChannel string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "invalid channel name",
			args: args{
				eventChannel: "invalie channel name",
			},
			want:    "",
			wantErr: true,
		}, {
			name: "valid channel name",
			args: args{
				eventChannel: channelPrefix + "key",
			},
			want:    "key",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRedisKeyFromEvent(tt.args.eventChannel)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRedisKeyFromEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseRedisKeyFromEvent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
