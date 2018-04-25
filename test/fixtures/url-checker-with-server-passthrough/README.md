# url-checker-with-server-passthrough

A set of templates that take in a few input variables and pass them through, unchanged, as output variables of the same
name. These templates create no resources, so they are just useful for rapid unit testing of code that works with
Terraform inputs and outputs.

Note that these templates are very similar or identical to other templates in the text-fixtures folder (that is, all
the ones that end with `-passthrough`), but as we run tests in parallel, we need to have a separate copy of these
templates for each test, or the tests will end up clobbering each others `.tfstate` files and `.terraform` folders.