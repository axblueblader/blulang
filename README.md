# BluLang

- Programming language based on Go and supports Vietnamese language syntax
- Example:

```
; this is a single line comment

; 'if' is an expression so it can stand alone or be assigned
; example of conditional statement
let a = if 10 != 10 {
    11
} else {
    12
}
; the last statement is always returned as a value

let b = 0

; 'while' is also an expression so it can stand alone or be assigned
; 'c' will be 10, 'b' will be 10
; example of looping statement
let c = while b != 10 {
    b = b + 1
}

; example of normal function declaration
fn minus(num1, num2, num3) {
    num1 - num2 - num3
}

; example of functions as a variable
let sum = fn (num1, num2, num3) {
    num1 + num2 + num3
}

let d = sum(a,b,c)
print("Hello world")
print(d) ; final result is 32
d ; returned result for testing
```

- More examples in [sample folder](/sample)

## Syntax

| Construct             | Syntax                            |
|-----------------------|-----------------------------------|
| Variable declaration  | let a = 10                        |
| Conditional statement | if 1 == 1 {} else {}              |
| Looping statement     | while 1 == 1 {}                   |
| Function declaration  | fn function(arg1,arg2) {}         |
| Printing              | print("ok")                       |
| Reading user input    | let a = input()                   |
| Return statement      | last statement is always returned |
| Array declaration     | let arr = [1,2,3]                 |
| Array usage           | arr[2] = 3, arr = arr + [4]       |
| Array element count   | count(arr)                        |

## Run

- Requirements: Go >= 1.18
- Checkout this repository then run:

```shell
go build
blulang ./sample/hello.blu
```