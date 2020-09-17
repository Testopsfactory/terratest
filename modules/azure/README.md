# The Azure Module APIs

The following document outlines the top-level API methods available for Terratest testing.

## Availability Set test APIs

These module APIs are provisioned by the [availabilityset.go](availabiilityset.go) file.

- `AvailabilitySetExists` indicates whether the speficied Azure Availability Set  exists\
    func AvailabilitySetExists(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) bool

- `GetAvailabilitySetFaultDomainCount` gets the Fault Domain Count for the specified Azure Availability Set\
    func GetAvailabilitySetFaultDomainCount(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) int32

- `CheckAvailabilitySetContainsVM` checks if the Virtual Machine is contained in the Availability Set VMs\
    func CheckAvailabilitySetContainsVM(t testing.TestingT, vmName string, avsName string, resGroupName string, subscriptionID string) bool

- `GetAvailabilitySetVMs` gets a list of VM names in the specified Azure Availability Set\
    func GetAvailabilitySetVMs(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) []string
    
- `GetAvailabilitySetE` gets an Availability Set in the specified Azure Resource Group\
    func GetAvailabilitySetE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (*compute.AvailabilitySet, error)
    
- `GetAvailabilitySetClientE` gets a new Availability Set client in the specified Azure Subscription\
    func GetAvailabilitySetClientE(subscriptionID string) (*compute.AvailabilitySetsClient, error)
