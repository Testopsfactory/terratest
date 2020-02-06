provider "google" {
  region = "us-east1"
}

# website::tag::1:: Deploy a cloud instance
resource "google_compute_instance" "example" {
  name         = var.instance_name
  machine_type = "f1-micro"
  zone         = "us-east1-b"

  # website::tag::2:: Run Ubuntu 18.04 on the instace
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }

  network_interface {
    network = "default"
    access_config {}
  }
}

# website::tag::3:: Allow the user to pass in a custom name for the instance
variable "instance_name" {
  description = "The Name to use for the Cloud Instance."
  default     = "gcp-hello-world-example"
}
