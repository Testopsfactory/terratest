output "map_of_objects" {
  value = {
    one   = 1
    two   = "two"
    three = "three"
    nest = {
      four = 4
      five = "five"
    }
    nest_list = [
      {
        six   = 6
        seven = "seven"
      },
    ]
  }
}

output "not_map_of_objects" {
  value = "Just a string"
}
