# minimal-example-with-error

This terraform template should fail because the EC2 KeyPair `key-that-does-not-exist` does not exist.  It is used to 
test that:
 
- terratest will retry with an expected error message
- terratest will fail when the terraform apply fails