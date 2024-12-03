package comm

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapFromYamlFileP(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		yamlText := `key: value`
		err := afero.WriteFile(fs, "file.yaml", []byte(yamlText), 0644)
		require.NoError(t, err)

		result := MapFromYamlFileP(fs, "file.yaml", false)
		assert.Equal(t, map[string]any{"key": "value"}, result)
	})

	t.Run("negative path - file not found", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		result := MapFromYamlFileP(fs, "nonexistent.yaml", false)
		assert.Panics(t, func() { _ = result })
	})

	t.Run("negative path - invalid yaml", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "file.yaml", []byte(`invalid: yaml`), 0644)
		require.NoError(t, err)

		result := MapFromYamlFileP(fs, "file.yaml", false)
		assert.Panics(t, func() { _ = result })
	})
}

func TestMapFromYamlFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		yamlText := `key: value`
		err := afero.WriteFile(fs, "file.yaml", []byte(yamlText), 0644)
		require.NoError(t, err)

		result, err := MapFromYamlFile(fs, "file.yaml", false)
		assert.NoError(t, err)
		assert.Equal(t, map[string]any{"key": "value"}, result)
	})

	t.Run("negative path - file not found", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		result, err := MapFromYamlFile(fs, "nonexistent.yaml", false)
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, result)
	})

	t.Run("negative path - invalid yaml", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "file.yaml", []byte(`invalid: yaml`), 0644)
		require.NoError(t, err)

		result, err := MapFromYamlFile(fs, "file.yaml", false)
		assert.ErrorContains(t, err, "parse yaml")
		assert.Nil(t, result)
	})
}

func TestMapFromYamlP(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		yamlText := `key: value`
		result := MapFromYamlP(yamlText, false)
		assert.Equal(t, map[string]any{"key": "value"}, result)
	})

	t.Run("negative path - invalid yaml", func(t *testing.T) {
		yamlText := `invalid: yaml`
		result := MapFromYamlP(yamlText, false)
		assert.Panics(t, func() { _ = result })
	})
}

func TestMapFromYaml(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		yamlText := `key: value`
		result, err := MapFromYaml(yamlText, false)
		assert.NoError(t, err)
		assert.Equal(t, map[string]any{"key": "value"}, result)
	})

	t.Run("negative path - invalid yaml", func(t *testing.T) {
		yamlText := `invalid: yaml`
		result, err := MapFromYaml(yamlText, false)
		assert.ErrorContains(t, err, "parse yaml")
		assert.Nil(t, result)
	})
}

func TestToJsonP(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		input := map[string]any{"key": "value"}
		result := ToJsonP(input, false)
		expected := `{"key":"value"}`
		assert.Equal(t, expected, result)
	})

	t.Run("pretty print", func(t *testing.T) {
		input := map[string]any{
			"key1": "value1",
			"key2": map[string]any{
				"subKey1": "subValue1",
				"subKey2": "subValue2",
			},
		}
		result := ToJsonP(input, true)
		expected := `{
    "key1": "value1",
    "key2": {
        "subKey1": "subValue1",
        "subKey2": "subValue2"
    }
}`
		assert.Equal(t, expected, result)
	})

	t.Run("negative path - error marshalling", func(t *testing.T) {
		input := make(chan int)
		result := ToJsonP(input, false)
		assert.Panics(t, func() { _ = result })
	})
}

func TestToJson(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		input := map[string]any{"key": "value"}
		result, err := ToJson(input, false)
		assert.NoError(t, err)
		expected := `{"key":"value"}`
		assert.Equal(t, expected, result)
	})

	t.Run("pretty print", func(t *testing.T) {
		input := map[string]any{
			"key1": "value1",
			"key2": map[string]any{
				"subKey1": "subValue1",
				"subKey2": "subValue2",
			},
		}
		result, err := ToJson(input, true)
		assert.NoError(t, err)
		expected := `{
    "key1": "value1",
    "key2": {
        "subKey1": "subValue1",
        "subKey2": "subValue2"
    }
}`
		assert.Equal(t, expected, result)
	})

	t.Run("negative path - error marshalling", func(t *testing.T) {
		input := make(chan int)
		result, err := ToJson(input, false)
		assert.ErrorIs(t, err, errors.New("failed to marshal json"))
		assert.Empty(t, result)
	})
}
