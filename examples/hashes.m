let newUser = fn(username, type) {
  { "username": username, "type": type }
}

let people = [
  newUser("mike", 1),
  newUser("toko", 2),
  newUser("code_snippet", 3),
  newUser("rav", 3),
]

puts(people[0].username)
puts(people[1]["username"])


let typeAsString = fn (type) {
  if (type == 1) {
    "admin"
  } else {
    "user(" + string(type) + ")"
  }
}

let printUsers = fn(users) {
  let iter = fn(users) {
    if (len(users) == 0) {
      return ;
    }
    let current = first(users)
    puts("User: " + current.username)
    puts(" Type - " + typeAsString(current["type"]))
    iter(tail(users))
  }
  iter(users)
}

printUsers(people)
