let map = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      return accumulated
    }
    return iter(tail(arr), push(accumulated, f(first(arr))))
  };
  iter(arr, [])
}

let a = [1, 2, 3]
let double = fn(x) { x * 2 }

puts(map(a, double))


let reduce = fn(arr, initial, f) {
  let iter = fn(arr, result) {
    if (len(arr) == 0) {
      return result
    }
    return iter(tail(arr), f(result, first(arr)))
  };
  return iter(arr, initial)
}

let sum = fn(arr) {
  return reduce(arr, 0, fn(initial, el) { initial + el });
}

let testSum = [1, 2, 3, 4, 5]
puts("sum of: " + string(testSum))
puts(sum(testSum))

