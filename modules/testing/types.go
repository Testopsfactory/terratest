package testing

// TestingT is an interface that describes the implementation of the testing object
// that the majority of Terratest functions accept as first argument.
// Using an interface that describes *testing.T instead of the actual implementation
// makes terratest usable in a wider variety of contexts (e.g. use with ginkgo : https://godoc.org/github.com/onsi/ginkgo#GinkgoT)
type TestingT interface {
	Fail()
	Fatalf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Name() string
	Logf(format string, args ...interface{})
}
