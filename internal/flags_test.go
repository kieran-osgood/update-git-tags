package internal

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCheckFlag(t *testing.T) {
	v := reflect.ValueOf("non-zero value")
	err := CheckFlag("property_path", v)
	assert.Nil(t, err, "it returns nil if non-zero value")

	v = reflect.ValueOf("")
	err = CheckFlag("property_path", v)
	assert.NotNil(t, err, "it returns error if zero value")
}

func TestParseFlags3(t *testing.T) {
	var tests = []struct {
		name      string
		args      []string
		conf      AllFlags
		wantError bool
	}{
		{"it errors when no flags passed in", []string{}, AllFlags{}, true},
		{"it parses flags into AllFlags", []string{"--RepositoryUrl=test_url", "--Branch=test_main", "--PropertyPath=version", "--sshKey=123", "--PreviousHash=abc", "--FilePath=app.json"},
			AllFlags{
				RepositoryUrl: "test_url",
				Branch:        "test_main",
				SshKey:        "123",
				PreviousHash:  "abc",
				PropertyPath:  "version",
				FilePath:      "app.json",
				args:          []string{},
			},
			false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, output, err := ParseFlags("test program", tt.args)
			if tt.wantError {
				assert.NotNil(t, err)
				return
			}
			if err != nil {
				t.Errorf("err got %v, want nil", err)
			}
			if output != "" {
				t.Errorf("output got %q, want empty", output)
			}
			if !reflect.DeepEqual(*conf, tt.conf) {
				t.Errorf("conf got %+v, want %+v", *conf, tt.conf)
			}
		})
	}
}
