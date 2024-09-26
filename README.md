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
- Vietnamese example calculating the Fibonacci number:
```
hàm fiboViệt(trỏ) {
    nếu trỏ == 0 {
        0
    } hay nếu trỏ = 1 {
        1
    } hay {
        fiboViệt(trỏ-2) + fiboViệt(trỏ-1)
    }
}
; vị trí  0 1 2 3 4 5 6 7
; giá trị 0 1 1 2 3 5 8 13
; 8 + 13
in(fiboEnglish(6) + fiboViệt(7)) ; kết quả là 21
```
- More examples in [sample folder](/sample)

## Syntax

| Construct             | Syntax                                                                                              |
|-----------------------|-----------------------------------------------------------------------------------------------------|
| Variable declaration  | let a = 10                                                                                          |
| Conditional statement | if 1 == 1 { print("ok") } else { print("what?") }                                                   |
| Looping statement     | while 1 == 1 { print("forever") }                                                                   |
| Function declaration  | fn add(arg1,arg2) { arg1 + arg2}                                                                    |
| Printing              | print("ok")                                                                                         |
| Reading user input    | let a = input()                                                                                     |
| Return statement      | fn inc(a) { while a<100 {if a == 50 {a return}} }<br/>last value before 'return' is always returned |
| Break statement       | while 1 == 1 { if a=3 {break} }<br/>last value before 'break' is returned                           |
| Array declaration     | let arr = [1,2,3]                                                                                   |
| Array usage           | arr[2] = 3, arr = arr + [4]                                                                         |
| Array element count   | count(arr)                                                                                          |
| Comment               | from ';' to end of line e.g ```; this is a comment```                                               |
| Object                | let body = { head: { eyes: 2, nose: 1}, torso: 1}                                                   |

**Important note** is that all construct returns the last statement's value so these syntax are allowed
```
let value1 = if 1+1 == 2 {
    1
} else {
    2
}
; value1 is now 1

let start = 0
let sum = while start < 3 {
    start = start + 1
}
; result will be sum = start = 3

let multiply2 = fn (a) {
    a*2
}   
; no need explicit return statement in function
```


## Run

- Requirements: Go >= 1.18
- Checkout this repository then run:

```shell
go build
blulang ./sample/hello.blu
```