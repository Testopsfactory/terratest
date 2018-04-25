# var-test

This test makes sure that terratest can successfully pass complicated variables (e.g. lists and maps) on the command
line to Terraform using the -var option. The templates in this example take in a number of different types of variables
and pass them through to outputs.