output "list_of_maps" {
  value = [
    {
      one   = "one"
      two   = "two"
      three = "three"
    },
    {
      one   = "uno"
      two   = "dos"
      three = "tres"
    }
  ]
}

output "not_list_of_maps" {
  value = "Just a string"
}
