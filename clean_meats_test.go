package main

import (
	"testing"

	helper "github.com/penkk55/grpc-go/helper"
)

func TestCleanMeats(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "Fatback t-bone t-bone, pastrami  ..   t-bone.  pork, meatloaf jowl enim.  Bresaola t-bone.",
			expected: []string{"Fatback", "t-bone", "t-bone", "pastrami", "t-bone", "pork", "meatloaf", "jowl", "enim", "Bresaola", "t-bone"},
		},
		{
			input:    "Bresaola-t-bone, pork-burger, meatloaf.",
			expected: []string{"Bresaola-t-bone", "pork-burger", "meatloaf"},
		},
		// {
		// 	input:    "pork&cheese-burger, meatloaf!",
		// 	expected: []string{"pork-cheese-burger", "meatloaf"},
		// },
		{
			input:    "t-bone...",
			expected: []string{"t-bone"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			// Call the CleanMeats function
			result := helper.CleanMeats(test.input)

			// Compare the result with the expected output
			if len(result) != len(test.expected) {
				t.Errorf("expected %v, but got %v", test.expected, result)
			}

			for i := range result {
				if result[i] != test.expected[i] {
					t.Errorf("at index %d: expected %s, but got %s", i, test.expected[i], result[i])
				}
			}
		})
	}
}
