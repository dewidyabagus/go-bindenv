package bindenv_test

import (
	"os"
	"testing"

	. "github.com/dewidyabagus/go-bindenv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalEnv(t *testing.T) {
	t.Run("Input value is not pointer data type", func(t *testing.T) {
		env := New()
		for _, input := range []any{nil, []string{"input"}} {
			err := env.Bind(input)
			require.NotNil(t, err)
			assert.EqualError(t, err, "data type sent must be pointer struct")
		}
	})

	t.Run("Input value is not pointer struct data type", func(t *testing.T) {
		input := map[string]string{"key": "value"}
		err := New().Bind(&input)
		require.NotNil(t, err)
		assert.EqualError(t, err, "can only process struct data type")
	})

	t.Run("Binding an unsupported data type", func(t *testing.T) {
		type simple struct {
			Slice []string
		}

		input := simple{}
		err := New().Bind(&input)
		require.NotNil(t, err)
		assert.Error(t, err, "data type is not supported")
	})

	t.Run("Binding of empty struct", func(t *testing.T) {
		input := struct{}{}
		require.Nil(t, New().Bind(&input))
		assert.Equal(t, struct{}{}, input)
	})

	t.Run("Binding default value", func(t *testing.T) {
		type simpe struct {
			Name   string
			Age    uint8
			Weight float32
		}
		env := New()
		env.Set("WEIGHT", "39.54")
		env.SetDefault("NAME", "John Wick")
		env.SetDefault("AGE", "19")

		input := simpe{}
		require.Nil(t, env.Bind(&input))
		assert.Equal(t, "39.54", env.Get("WEIGHT"))
		assert.Equal(t, simpe{Name: "John Wick", Age: 19, Weight: 39.54}, input)
	})

	t.Run("Binding a simple struct", func(t *testing.T) {
		type simple struct {
			ID      int     `env:"SIMPLE_ID"`
			Name    string  `env:"SIMPLE_NAME"`
			Student bool    `env:"SIMPLE_STUDENT"`
			Age     uint    `env:"SIMPLE_AGE"`
			Weight  float64 `env:"SIMPLE_WEIGHT"`
		}
		variables := map[string]string{
			"SIMPLE_ID":      "999999",
			"SIMPLE_NAME":    "John Wick",
			"SIMPLE_STUDENT": "TRUE",
			"SIMPLE_AGE":     "12",
			"SIMPLE_WEIGHT":  "40.499",
		}

		env := New()

		input := simple{}
		require.Nil(t, env.Sets(variables))
		require.Nil(t, env.Bind(&input))
		assert.Equal(t, simple{ID: 999999, Name: "John Wick", Student: true, Age: 12, Weight: 40.499}, input)
		for key, value := range variables {
			assert.Equal(t, value, env.Get(key))
		}
	})

	t.Run("Binding nested struct", func(t *testing.T) {
		type sequence struct {
			Version uint32 `env:"SEQUENCE_VERSION"`
		}

		type tech struct {
			Status   bool `env:"TECH_STATUS"`
			Sequence sequence
		}

		type simpe struct {
			Name string `env:"SIMPLE_NAME"`
			Tech tech
		}

		os.Clearenv()
		env := New()
		env.Sets(map[string]string{
			"SIMPLE_NAME":      "new tech name",
			"TECH_STATUS":      "true",
			"SEQUENCE_VERSION": "123",
		})

		input := simpe{}
		expected := simpe{
			Name: "new tech name",
			Tech: tech{
				Status: true,
				Sequence: sequence{
					Version: 123,
				},
			},
		}
		require.Nil(t, env.Bind(&input))
		assert.Equal(t, expected, input)
	})
}
