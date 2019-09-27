package testing

// TestingT is an interface that describes the implementation of the testing object
// that the majority of Terratest functions accept as first argument.
// Using an interface that describes *testing.T instead of the actual implementation
// makes terratest usable in a wider variety of contexts (e.g. use with ginkgo : https://godoc.org/github.com/onsi/ginkgo#GinkgoT)
type TestingT interface {
	//Fail marks the function as having failed but continues execution.
	Fail()
	// FailNow marks the function as having failed and stops its execution
	// by calling runtime.Goexit (which then runs all deferred calls in the
	// current goroutine).
	// Execution will continue at the next test or benchmark.
	// FailNow must be called from the goroutine running the
	// test or benchmark function, not from other goroutines
	// created during the test. Calling FailNow does not stop
	// those other goroutines.
	FailNow()
	// Fatalf is equivalent to Logf followed by FailNow.
	Fatalf(format string, args ...interface{})
	// Error is equivalent to Log followed by Fail.
	Error(args ...interface{})
	// Errorf is equivalent to Logf followed by Fail.
	Errorf(format string, args ...interface{})
	// Log formats its arguments using default formatting, analogous to Println,
	// and records the text in the error log. For tests, the text will be printed only if
	// the test fails or the -test.v flag is set. For benchmarks, the text is always
	// printed to avoid having performance depend on the value of the -test.v flag.
	Log(args ...interface{})
	// Logf formats its arguments according to the format, analogous to Printf, and
	// records the text in the error log. A final newline is added if not provided. For
	// tests, the text will be printed only if the test fails or the -test.v flag is
	// set. For benchmarks, the text is always printed to avoid having performance
	// depend on the value of the -test.v flag.
	Logf(format string, args ...interface{})
	// Name returns the name of the running test or benchmark.
	Name() string
}
