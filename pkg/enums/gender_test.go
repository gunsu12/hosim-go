package enums_test

import (
	"testing"

	"hosim-go/pkg/enums"
)

func TestGender_IsValid(t *testing.T) {
	tests := []struct {
		gender enums.Gender
		valid  bool
	}{
		{enums.GenderMale, true},
		{enums.GenderFemale, true},
		{enums.GenderUnknown, true},
		{enums.Gender("X"), false},
		{enums.Gender(""), false},
	}

	for _, tt := range tests {
		if got := tt.gender.IsValid(); got != tt.valid {
			t.Errorf("Gender(%q).IsValid() = %v, expected %v", tt.gender, got, tt.valid)
		}
	}
}

func TestNormalizeGender(t *testing.T) {
	tests := []struct {
		input    string
		expected enums.Gender
		ok       bool
	}{
		{"L", enums.GenderMale, true},
		{"male", enums.GenderMale, true},
		{"laki-laki", enums.GenderMale, true},
		{"m", enums.GenderMale, true},
		{"P", enums.GenderFemale, true},
		{"female", enums.GenderFemale, true},
		{"perempuan", enums.GenderFemale, true},
		{"f", enums.GenderFemale, true},
		{"other", enums.GenderUnknown, true},
		{"O", enums.GenderUnknown, true},
		{"unknown", enums.GenderUnknown, true},
		{"tidak tahu", enums.GenderUnknown, true},
		{"invalid", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		got, ok := enums.NormalizeGender(tt.input)
		if ok != tt.ok || got != tt.expected {
			t.Errorf("NormalizeGender(%q) = (%q, %v), expected (%q, %v)", tt.input, got, ok, tt.expected, tt.ok)
		}
	}
}
