package utils

import "testing"

func TestParseStrNumber(t *testing.T) {
	type args struct {
		src  string
		dest interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				src:  "123",
				dest: int(5),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseStrNumber(tt.args.src, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("ParseStrNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(tt.args.dest)
		})
	}
}
