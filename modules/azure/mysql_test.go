// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure MySQL server and database, these tests can be extended
*/

func TestGetMYSQLServerSkuNameE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := GetMYSQLServerSkuNameE(t, subscriptionID, resGroupName, serverName)
	require.Error(t, err)
}

func TestGetMYSQLServerStorageMBE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := GetMYSQLServerStorageMBE(t, subscriptionID, resGroupName, serverName)
	require.Error(t, err)
}

func TestGetMYSQLServerStateE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := GetMYSQLServerStateE(t, subscriptionID, resGroupName, serverName)
	require.Error(t, err)
}

func TestGetMYSQLDBCharsetE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""
	dbName := ""

	_, err := GetMYSQLDBCharsetE(t, subscriptionID, resGroupName, serverName, dbName)
	require.Error(t, err)
}

func TestGetMYSQLDBCollationE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""
	dbName := ""

	_, err := GetMYSQLDBCollationE(t, subscriptionID, resGroupName, serverName, dbName)
	require.Error(t, err)
}
