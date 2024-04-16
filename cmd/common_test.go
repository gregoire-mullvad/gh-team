package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTeam(t *testing.T) {
	org, slug, err := parseTeam("mullvad/services")
	assert.Equal(t, "mullvad", org)
	assert.Equal(t, "services", slug)
	assert.NoError(t, err)
}
func Test_parseTeam_error(t *testing.T) {
	inputs := []string{"", "mullvad", "mullvad/", "/services"}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			_, _, err := parseTeam(input)
			assert.Error(t, err)
		})
	}
}
