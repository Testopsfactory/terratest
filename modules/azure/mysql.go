package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetMYSQLServerClient is a helper function that will setup a mysql server client
func GetMYSQLServerClient(subscriptionID string) (*mysql.ServersClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a mysql server client
	mysqlClient := mysql.NewServersClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	mysqlClient.Authorizer = *authorizer

	return &mysqlClient, nil
}

// GetMYSQLServerSkuName is a helper function that gets the server SKU Name
func GetMYSQLServerSkuName(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) string {
	skuName, err := GetMYSQLServerSkuNameE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return skuName
}

// GetMYSQLServerSkuNameE is a helper function that gets the server Sku Name
func GetMYSQLServerSkuNameE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (string, error) {
	// Create a mySQl Server client
	mysqlClient, err := GetMYSQLServerClient(subscriptionID)
	if err != nil {
		return "", err
	}

	// Get the corresponding server client
	mysqlServer, err := mysqlClient.Get(context.Background(), resGroupName, serverName)
	if err != nil {
		return "", err
	}

	//Return server SKU name
	return *mysqlServer.Sku.Name, nil
}

//GetMYSQLServerStorageMB  is a helper function that gets the server storage Mb
func GetMYSQLServerStorageMB(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) int32 {
	storageMb, err := GetMYSQLServerStorageMBE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return storageMb
}

// GetMYSQLServerStorageMBE is a helper function that gets the server storage Mb
func GetMYSQLServerStorageMBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (int32, error) {
	// Create a mySQl Server client
	mysqlClient, err := GetMYSQLServerClient(subscriptionID)
	if err != nil {
		return -1, err
	}

	// Get the corresponding server client
	mysqlServer, err := mysqlClient.Get(context.Background(), resGroupName, serverName)
	if err != nil {
		return -1, err
	}

	//Return server storage MB
	return *mysqlServer.ServerProperties.StorageProfile.StorageMB, nil
}

//GetMYSQLServerState  is a helper function that gets the server state
func GetMYSQLServerState(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) mysql.ServerState {
	serverState, err := GetMYSQLServerStateE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return serverState
}

//GetMYSQLServerStateE is a helper function that gets the server state
func GetMYSQLServerStateE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (mysql.ServerState, error) {
	// Create a mySQl Server client
	mysqlClient, err := GetMYSQLServerClient(subscriptionID)
	if err != nil {
		return "", err
	}

	// Get the corresponding server client
	mysqlServer, err := mysqlClient.Get(context.Background(), resGroupName, serverName)
	if err != nil {
		return "", err
	}

	//Return server state
	return mysqlServer.ServerProperties.UserVisibleState, nil
}

// GetMYSQLDBClient is a helper function that will setup a mysql DB client
func GetMYSQLDBClient(subscriptionID string) (*mysql.DatabasesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a mysql db client
	mysqlDBClient := mysql.NewDatabasesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	mysqlDBClient.Authorizer = *authorizer

	return &mysqlDBClient, nil
}

//GetMYSQLDBStatus is a helper function that gets the database status
func GetMYSQLDBStatus(t testing.TestingT, resGroupName string, serverName string, subscriptionID string, dbName string) string {
	dbStatus, err := GetMYSQLDBStatusE(t, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return dbStatus
}

//GetMYSQLDBStatusE is a helper function that gets the database status
func GetMYSQLDBStatusE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (string, error) {
	// Create a mySQl db client
	mysqldbClient, err := GetMYSQLDBClient(subscriptionID)
	if err != nil {
		return "", err
	}

	// Get the corresponding db client
	mysqlDb, err := mysqldbClient.Get(context.Background(), resGroupName, serverName, dbName)
	if err != nil {
		return "", err
	}

	//Return DB status
	return mysqlDb.Status, nil
}
