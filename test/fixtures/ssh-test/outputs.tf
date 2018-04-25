output "example_public_ip" {
  value = "${aws_instance.example_public.public_ip}"
}

output "example_private_ip" {
  value = "${aws_instance.example_private.private_ip}"
}