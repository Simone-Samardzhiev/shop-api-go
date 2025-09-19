package validation

import "testing"

func Test_isValidPassword(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"Password1!", true},
		{"short1!", false},
		{"noupper1!", false},
		{"NOLOWER1!", false},
		{"NoNumber!", false},
		{"NoSpecial1", false},
		{"White space1!", false},
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			if isValidPassword(test.password) != test.valid {
				t.Errorf("Expected %v, got %v", test.valid, isValidPassword(test.password))
			}
		})
	}
}

func Fuzz_isValidPassword(f *testing.F) {
	f.Add("")
	f.Add("password")
	f.Add("Password!")

	f.Fuzz(func(t *testing.T, password string) {
		_ = isValidPassword(password)
	})
}
