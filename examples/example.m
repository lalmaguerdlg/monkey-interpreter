


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
let x = 25
let result = fib(x, counter)

puts("fib(" + string(x) + ") = " + string(result));

puts("iterations: " + string(counter() - 1));

