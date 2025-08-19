let firstname = "Luis"
let lastname = "Almaguer"

let concat = fn(a, b, separator) { a + separator + b }

let fullname = concat(firstname, lastname, " ")


puts(fullname)

puts("foo " + 1)
puts("bar " + true)
puts("baz " + puts(1))
puts("baz " + fn(x) { x })
