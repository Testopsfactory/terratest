package aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRecommendedRdsInstanceType(t *testing.T) {
	type TestingScenerios struct {
		name           string
		region         string
		databaseEngine string
		instanceTypes  []string
		expected       string
	}
}

func TestGetRecommendedRdsInstanceTypeHappyPath(t *testing.T) {
	type TestingScenerios struct {
		name           string
		region         string
		databaseEngine string
		instanceTypes  []string
		expected       string
	}

	testingScenerios := []TestingScenerios{
		{
			name:           "US region, mysql, first offering available",
			region:         "us-east-2",
			databaseEngine: "mysql",
			instanceTypes:  []string{"db.t2.micro", "db.t3.micro"},
			expected:       "db.t2.micro",
		},
		{
			name:           "EU region, postgres, 2nd offering available based on region",
			region:         "eu-north-1",
			databaseEngine: "postgres",
			instanceTypes:  []string{"db.t2.micro", "db.m5.large"},
			expected:       "db.m5.large",
		},
		{
			name:           "US region, oracle-ee, 2nd offering available based on db type",
			region:         "us-west-2",
			databaseEngine: "oracle-ee",
			instanceTypes:  []string{"db.m5d.xlarge", "db.m5.large"},
			expected:       "db.m5.large",
		},
	}

	for _, scenerio := range testingScenerios {
		scenerio := scenerio

		t.Run(scenerio.name, func(t *testing.T) {
			t.Parallel()

			actual, err := GetRecommendedRdsInstanceTypeE(t, scenerio.region, scenerio.databaseEngine, scenerio.instanceTypes)
			assert.NoError(t, err)
			assert.Equal(t, scenerio.expected, actual)
		})
	}
}

func TestGetRecommendedRdsInstanceTypeErrors(t *testing.T) {
	type TestingScenerios struct {
		name           string
		region         string
		databaseEngine string
		instanceTypes  []string
	}

	testingScenerios := []TestingScenerios{
		{
			name:           "All empty",
			region:         "",
			databaseEngine: "",
			instanceTypes:  nil,
		},
		{
			name:           "No engine or instance type",
			region:         "us-east-2",
			databaseEngine: "",
			instanceTypes:  nil,
		},
		{
			name:           "No instance types",
			region:         "us-east-2",
			databaseEngine: "mysql",
			instanceTypes:  nil,
		},
		{
			name:           "Invalid instance types",
			region:         "us-east-2",
			databaseEngine: "mysql",
			instanceTypes:  []string{"garbage"},
		},
		{
			name:           "Region has no instance type available",
			region:         "eu-north-1",
			databaseEngine: "mysql",
			instanceTypes:  []string{"db.t2.micro"},
		},
		{
			name:           "No instance type available for engine",
			region:         "us-east-1",
			databaseEngine: "oracle-ee",
			instanceTypes:  []string{"db.r5d.large"},
		},
	}

	for _, scenerio := range testingScenerios {
		scenerio := scenerio

		t.Run(scenerio.name, func(t *testing.T) {
			t.Parallel()

			_, err := GetRecommendedRdsInstanceTypeE(t, scenerio.region, scenerio.databaseEngine, scenerio.instanceTypes)
			fmt.Println(err)
			assert.EqualError(t, err, NoRdsInstanceTypeError{InstanceTypeOptions: scenerio.instanceTypes, DatabaseEngine: scenerio.databaseEngine}.Error())
		})
	}
}
