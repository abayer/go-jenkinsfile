package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"

	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {
	var testFiles []string
	testBase := filepath.Join("testdata", "json")
	err := filepath.Walk(testBase, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		testFiles = append(testFiles, path)
		return nil
	})
	assert.NoError(t, err)

	for _, tc := range testFiles {
		testFile, err := filepath.Rel(testBase, tc)
		assert.NoError(t, err)
		t.Run(strings.TrimSuffix(testFile, ".json"), func(t *testing.T) {
			contents, err := ioutil.ReadFile(tc)
			assert.NoError(t, err)

			got := &Root{}

			err = json.Unmarshal(contents, got)
			assert.NoError(t, err)

			assert.NotNil(t, got.Pipeline)

			// Make sure we don't have any null fields
			outYaml, err := yaml.Marshal(got)
			assert.NoError(t, err)
			// Scrub "value: null" because that's not a problem
			outStr := strings.ReplaceAll(string(outYaml), "value: null", "value: EMPTY")
			assert.NotContains(t, outStr, ": null")

			// Make sure we can marshal to JSON and reparse it.
			outJSON, err := json.Marshal(got)
			assert.NoError(t, err)

			gotTwo := &Root{}
			err = json.Unmarshal(outJSON, gotTwo)
			assert.NoError(t, err)
		})
	}

}
