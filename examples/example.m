
let firstname = "Luis"
let lastname = "Almaguer"

let concat = fn(a, b, separator) { a + separator + b }

let fullname = concat(firstname, lastname, " ")


puts(fullname)


let newCounter = fn(start) {
  return fn() {
    start = start + 1;
    return start
  }
}


let fib = fn(x, counter) {
  let fibImp = fn(x) {
    counter()
    if (x == 1) {
      return 1
    }
    if (x == 2) {
      return 2
    }
    return fibImp(x - 1) + fibImp(x - 2)
  }
  return fibImp(x)
}

let counter = newCounter(0)
let result = fib(25, counter)

puts("fib(25): ");
puts(result);
puts("iterations: ");
puts(counter() - 1);
