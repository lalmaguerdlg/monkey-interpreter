let firstname = "Luis"
let lastname = "Almaguer"

let concat = fn(a, b, separator) { a + separator + b }

let fullname = concat(firstname, lastname, " ")


puts(fullname)

# These are comments
puts("Casting to string") # can be place after statements
puts("int: " + string(1))
puts("bool: " + string(true))
puts("function: " + string(fn(x) {x}))

