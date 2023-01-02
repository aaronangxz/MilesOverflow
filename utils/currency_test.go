package utils

import (
	"testing"
)

func TestCConvertFCYToSGD(t *testing.T) {
	type args struct {
		original float64
		currency string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "MYR-SGD",
			args:    args{original: 100, currency: "MYR"},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertFCYToSGD(tt.args.original, tt.args.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("Conversion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Conversion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
