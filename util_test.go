package main

import "testing"

func Test_ParsePostIdFromArg(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "post id",
			args: args{
				arg: "123",
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "post url",
			args: args{
				arg: "https://debug.esa.io/posts/123",
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "invalid post id (1)",
			args: args{
				arg: "1a2b3c",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid post id (2)",
			args: args{
				arg: "-123",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid post id (3)",
			args: args{
				arg: "0",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePostIdFromArg(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePostIdFromArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePostIdFromArg() got = %v, want %v", got, tt.want)
			}
		})
	}
}
