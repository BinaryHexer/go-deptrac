package validator_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/BinaryHexer/go-deptrac/validator"
)

func init() {
	validator.Log.SetOutput(os.Stderr)
}

func TestValidator_Validate(t *testing.T) {
	testCases := []struct {
		ConfigPath string
		IsValid    bool
	}{
		{ConfigPath: "../examples/simple-mvc/depfile.yaml", IsValid: true},
		{ConfigPath: "../examples/simple-cleanarch/depfile.yaml", IsValid: true},
		{ConfigPath: "../examples/simple-invalid-mvc/depfile.yaml", IsValid: false},
	}

	for _, c := range testCases {
		t.Run(c.ConfigPath, func(t *testing.T) {
			config := validator.ParseConfig([]string{c.ConfigPath})
			v := validator.NewValidator(config)
			valid, errors, err := v.Validate(true)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println("errors: ", errors)

			if valid != c.IsValid {
				t.Errorf("path %s should be %t, but is %t", c.ConfigPath, c.IsValid, valid)
			}
			if !c.IsValid && len(errors) == 0 {
				t.Error("module is invalid, but errors are empty")
			}
		})
	}
}
