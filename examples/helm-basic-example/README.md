# Helm Basic Example

This folder contains a minimal helm chart to demonstrate how you can use Terratest to test your helm charts.

There are two kinds of tests you can perform on a helm chart:

- Helm Template tests are tests designed to test the logic of the templates. These tests should run `helm template` with
  various input values and parse the yaml to validate any logic embedded in the templates (e.g by reading them in using
  client-go). Since templates are not statically typed, the goal of these tests is to promote fast cycle time
- Helm Integration tests are tests that are designed to deploy the infrastructure and validate the resource
  configurations. If you consider the templates to be syntactic tests, these are semantic tests that validate the
  behavior of the deployed resources.

The helm chart deploys a single replica Deployment resource given the container image spec. This chart requires the
`containerImageRepo` and `containerImageTag` input values.

See the corresponding terratest code for an example of how to test this chart:

- [helm_basic_example_template_test.go](/test/helm_basic_example_template_test.go): the template tests for this chart.
<!-- TODO: Append the example with integration tests and deployment instructions once terratest has deploy test functions -->

## Running automated tests against this Helm Chart

1. Install and setup [helm](https://docs.helm.sh/using_helm/#installing-helm)
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -tags kubernetes -run TestHelmBasicExampleTemplate`
