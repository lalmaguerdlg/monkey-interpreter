let firstname = "Luis"
let lastname = "Almaguer"

let concat = fn(a, b, separator) { a + separator + b }

let fullname = concat(firstname, lastname, " ")


puts(fullname)

puts("Casting to string")
puts("int: " + string(1))
puts("bool: " + string(true))
puts("function: " + string(fn(x) {x}))

