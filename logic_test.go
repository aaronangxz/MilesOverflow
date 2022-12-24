package main

import (
	"testing"
	"time"
)

func Test_calculateLocal(t *testing.T) {
	type args struct {
		t       Transaction
		c       Card
		current float64
	}
	tests := []struct {
		name            string
		args            args
		expectedRewards float64
		expectedMiles   float64
	}{
		{
			name: "valid HRV",
			args: args{t: Transaction{
				description: "Mock Trx",
				category:    CATEGORY_ONLINE,
				paymentType: PAYMENT_TYPE_CONTACTLESS,
				amount:      2000,
				currency:    "SGD",
				time:        time.Now().Unix(),
				cardId:      1,
			}, c: mockCard()[0], current: 0},
			expectedRewards: 200,
			expectedMiles:   80,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := calculateLocal(tt.args.t, tt.args.c, tt.args.current)
			if got != tt.expectedRewards {
				t.Errorf("calculateLocal() actualRewards = %v, want %v", got, tt.expectedRewards)
			}
			if got1 != tt.expectedMiles {
				t.Errorf("calculateLocal() actualMiles = %v, want %v", got1, tt.expectedMiles)
			}
		})
	}
}
