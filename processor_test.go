package yamlprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Mission struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Status      string `yaml:"status"`
	Priority    string `yaml:"priority"`
}

type testDoc struct {
	Name     string    `yaml:"name"`
	Age      int       `yaml:"age"`
	Missions []Mission `yaml:"missions"`
}

func TestProcessor_Unmarshal(t *testing.T) {
	var d testDoc
	err := DefaultProcessor.Unmarshal([]byte(`
name: "Snake Eyes"
age: 33
missions:
  ${file("testdata/missions.yaml")}
`), &d)
	require.NoError(t, err)
	require.Equal(t, testDoc{
		Name: "Snake Eyes",
		Age:  33,
		Missions: []Mission{
			{
				Name:        "Mission 1",
				Description: "This is the first mission",
				Status:      "in progress",
				Priority:    "high",
			},
			{
				Name:        "Mission 2",
				Description: "This is the second mission",
				Status:      "completed",
				Priority:    "medium",
			},
			{
				Name:        "Mission 3",
				Description: "This is the third mission",
				Status:      "pending",
				Priority:    "low",
			},
		},
	}, d)
}
