
# Oro Language

Oro is an interpreted, expressive, it has a noiseless syntax, free of useless semi colons, braces or parentheses, and treats everything as an expression.

It features mutable and immutable values, if and match conditionals, functions, type hinting, repeat loops, modules, the pipe operator, use and many more. All of that while retaining it's expressiveness, clean syntax and easy of use.


## Table of Contents

* [Usage](#usage)
    * [Run a Source File](#run-a-source-file)
    * [REPL](#repl)
* [Variables](#variables)
    * [Constants](#constants)
    * [Type Lock](#type-lock)
* [Data Types](#data-types)
    * [String](#string)
    * [Symbol](#symbol)
    * [Int](#int)
    * [Float](#float)
    * [Boolean](#boolean)
    * [Array](#array)
    * [Dictionary](#dictionary)
    * [Nil](#nil)
    * [Type Conversion](#type-conversion)
    * [Type Checking](#type-checking)
* [Operators](#operators)
    * [Shorthand Assignment](#shorthand-assignment)
* [Functions](#functions)
    * [Type Hinting](#type-hinting)
    * [Default Parameters](#default-parameters)
    * [Return Statement](#return-statement)
    * [Variadic](#variadic)
    * [Arrow Functions](#arrow-functions)
    * [Closures](#closures)
    * [Recursion](#recursion)
    * [Tricks](#tricks)
* [Conditionals](#conditionals)
    * [If](#if)
    * [Ternary Operator](#ternary-operator)
    * [Match](#match)
    * [Pattern Matching](#pattern-matching)
* [Repeat Loop](#repeat-loop)
* [Range Operator](#range-operator)
* [Pipe Operator](#pipe-operator)
* [Immutability](#immutability)
* [Modules](#modules)
* [Uses](#uses)
* [Comments](#comments)
* [Standard Library](#standard-library)

## Usage

If you want to play with the language, but have no interest in toying with its code, you can download a built binary for your operating system. Just head to the [latest release](https://github.com/luiscm/oro/releases/latest) and download one of the archives.

The other option, where you get to play with the code and run your changes, is to `go get github.com/luiscm/oro` and install it as a local binary with `go install`. Obviously, you'll need `GOROOT` in your path, but I guess you already know what you're doing.

### Run a source file

To run an Oro source file, give it a path relative to the current directory.

```
oro run path/to/file.oro
```

### REPL

As any serious language, Oro provides a REPL too:

```
oro repl
```

## Variables

Variables in Oro start with the keyword `var`. Accessing an undeclared variable, in contrast with some languages, will not create it, but instead throw a runtime error.

```swift
var name = "Luis"
var married = true

var age = 50
age = 50
```

Names have to start with an alphabetic character and continue either with alphanumeric, underscores, questions marks or exclamation marks. When you see a question mark, don't confuse them with optionals like in some other languages. In here they have no special lexical meaning except that they allow for some nice variable names like `is_empty?` or `do_it!`.

### Constants

Constants have the same traits as variables, except that they start with `val` and are immutable. Once declared, reassigning a constant will produce a runtime error. Even data structures are locked into immutability. Elements of an Array or Dictionary can't be added, updated or removed.

```swift
val name = "Luis"
name = "Carlos" // runtime error
```

### Type Lock

Type lock is a safety feature of mutable variables. Once they're declared with a certain data type, they can only be assigned to that same type. This makes for more predictable results, as an integer variable can't be assigned to a string or array. In this regard, Oro works as a strong typed language.

This will work:

```swift
var nr = 10
nr = 15
```

This won't:

```swift
var nr = 10
nr = "ten" // runtime error
```

## Data Types

Oro supports 7 data types: `Boolean`, `String`, `Int`, `Float`, `Array`, `Dictionary`, `Symbol` and `Nil`.

### String

Strings are UTF-8 encoded, meaning that you can stuff in there anything, even emojis.

```swift
val weather = "Hot"
val price = "円900"
```

String concatenation is handled with the `+` operator. Concatenating between a string and another data type will result in a runtime error.

```swift
val name = "Luis" + " " + "Carlos" 
```

Additionally, strings are treated as enumerables. They support subscripting and iteration in `repeat in` loops.

```swift
"world"[2] // "r" 
```

Escape sequences are there too if you need them: `\"`, `\n`, `\t`, `\r`, `\a`, `\b`, `\f` and `\v`. Nothing changes from other languages, so I'm sure you can figure out by yourself what every one of them does.

```swift
val code = "if name == \"luis\"\nthen\n\tprint(10)\nend"
```

### Symbol

Symbols, or atoms as some languages refer to them, are constants where the name is their value. Although they behave a lot like strings and can generally be interchanged, internally they are treated as their own type. As the language progresses, Symbols will be put to better use.

```swift
val eq = :dog == :cat
val arr = ["dog", :cat, :mouse]
val dict = [:name => "Luis", :age => 50]
val concat = "hello" + :world
```

They're interesting to use as control conditions, emulating enums as a fixed, already-known value:

```swift
val os = "linux"
match os
when :linux
  println("FREE")
when :windows
  println("NO FREE")
end
```

### Int

Integers are whole numbers that support most of the arithmetic and bitwise operators, as you'll see later. They can be represented also as: binary with the 0b prefix, hexadecimal with the 0x prefix and octal with the 0o prefix.

```swift
val dec = 27
val oct = 0o33
val hex = 0x1B
val bin = 0b11011
val arch = 2 ** 32
```

A sugar feature both in Integer and Float is the underscore:
 
```swift
val big = 27_000_000
```

It has no special meaning, as it will be ignored in the lexing phase. Writing `1_000` and `1000` is the same thing to the interpreter.

### Float

Floating point numbers are used in a very similar way to Integers. In fact, they can be mixed and matched, like `3 + 0.2` or `5.0 + 2`, where the result will always be a Float.

```swift
val pi = 3.14_159_265
val e = 2.71828182
```

Scientific notation is also supported via the `e` modifier:

```swift
val sci = 0.1e3
val negsci = 25e-5
```

### Boolean

It would be strange if this data type included anything else except `true` and `false`.

```swift
val found = true
val log = false
```

Expressions like the `if/else`, as you'll see later, will check for values that aren't necessarily boolean. Integers and Floats will be checked if they're equal to 0, and Strings, Arrays and Dictionaries if they're empty. These are called `truthy` expressions and internally, will be evaluated to boolean.

### Array

Arrays are ordered collections of any data type. You can mix and match strings with integers, or floats with other arrays.
 
```swift
val multi = [5, "Hi", ["Hello", "World"]]
val names = ["Luis", "Carlos", 1337]

val luis = names[0]
val number = names[-1]
```
 
Individual array elements can be accessed via subscripting with a 0-based index:

```swift
val names = ["Luis", "Alberto", "Carlos"]
val first = names[0] // "Luis"
val last = names[-1] // "Carlos"
```

In the same style, an index can be used to check if it exists. It will return `nil` if it doesn't:

```swift
if names[10]
  // handle it
end
```

Individual elements can be reassigned on mutable arrays:

```swift
var numbers = [5, 8, 10, 15]
numbers[1] = 7
```

Appended with an empty or placeholder index:

```swift
numbers[] = 100
numbers[_] = 200 // Same.
```

Arrays can be compared with the `==` and `!=` operators, which will check the position and value of every element of both arrays. Equal arrays should have the same exact values in the same position.

They can also be combined with the `+` operator, which adds the element of the right side to the array on the left side.

```swift
val concat = ["an", "array"] + ["and", "another"]
// ["an", "array", "and", "another"]
```

Oh and if you're that lazy, you can ommit commas too:

```swift
val nocomma = [5 7 9 "Hi"]
```
 
### Dictionary
 
Dictionaries are hashes with a key and a value of any data type. They're good to hold unordered, structured data:

```swift
val user = ["name" => "Dr. Luis Carlos", "profession" => "Developer", "age" => 50]
```

I'd argue that using Symbols for keys would make them look cleaner:

```swift
val user2 = [:name => "Dr. Luis Carlos", :profession => "Developer", :age => 50]
```

Unlike arrays, internally their order is irrelevant, so you can't rely on index-based subscripting. They only support key-based subscripting:
 
```swift
user["name"] // "Dr. Luis Carlos"
user2[:name] // "Dr. Luis Carlos"
```

Values can be reassigned or inserted by key on mutable dictionaries:

```swift
var numbers = ["one" => 1, "two" => 2]
numbers["one"] = 5
numbers["three"] = 3 // new key:value

var numbers2 = [:one => 10, :two => 20]
numbers[:one] = 50
numbers[:three] = 30 // new key:value
```

To check for a key's existence, you can access it as normal and check if it's `nil` or truthy:

```swift
if user["location"] == nil
  // do ...
end

if user2[:location] == nil
  // do ...
end
```

### Nil

Oro has a Nil type and yes, I'm totally aware of its problems. This was a choice for simplicity, at least for the time being. In the future, I plan to experiment with optionals and hopefully integrate them into the language.

```swift
val empty = nil
```

### Type Conversion

Converting between types is handled in a few ways that produce exactly the same results. The `as` operator is probably the more convenient and more expressive of the bunch. Like all type conversion methods, it can convert to `String`, `Int`, `Float` and `Array`:

```swift
val nr = 10
nr as String
nr as Int
nr as Float
nr as Array
```

Provided by the runtime are the appropriately named functions: `String()`, `Int()`, `Float()` and `Array()`.

```swift
val str = String(10)
val int = Int("10")
val fl = Float(10)
val arr = Array(10)
```

The `Type` module of the Standard Library provides interfaces to those same functions and even adds some more, like `Type.of()` and `Type.isNumber?()`.

```swift
val str = Type.toString(10)
val int = Type.toInt("10")
val fl = Type.toFloat(10)
val arr = Type.toArray(10)
```

Which method you choose to use is strictly preferential and depends on your background.

### Type Checking

There will be more than one occassion where you'll need to type check a variable. Oro provides a few ways to achieve that.

The `is` operator is specialized in checking types and should be the one you'll want to use practically everywhere.

```swift
val nr = 10
if nr is Int
  println("Yes, an integer")
end
```

There's also the `typeof()` runtime function and `Type.of()` from the Standard Library. They essentially do the same thing, but not only they're longer to write, but return strings. The above would be equivalent to:

```swift
if Type.of(nr) == "Int"
  println("Yes, an integer")
end
```

## Operators

You can't expect to run some calculations without a good batch of operators, right? Well, Oro has a range of arithmetic, boolean and bitwise operators to match your needs.

By order of precedence:

```swift
Boolean: && || (AND, OR)
Bitwise: & | ~ (Bitwise AND, OR, NOT)
Equality: == != (Equal, Not Equal)
Comparison: < <= > >=
Bitshift: << >> (Bitshift Left and Right)
Arithmetic: + - * / % ** (Addition, Subtraction, Multiplication, Division, Modulus, Exponential)
```

Arithmetic expressions can be safely used for Integers and Floats:

```swift
1 + 2 * 3 / 4.2
2 ** 8
3 % 2 * (5 - 3)
```

Addition can be used to concatenate Strings or combine Arrays and Dictionaries:

```swift
"hello" + " " + "world"
[1, 2] + [3, 4]
["a" => 1, "b" => 2] + ["c" => 3]
[:a => 10, :b => 20] + [:c => 30]
```

Comparison operators can compare Integers and Float by exact value, Strings, Arrays and Dictionaries by length:

```swift
5 > 2
3.2 <= 4.5
"one" < "three"
"two" < "four"
[1, 2] > [5]
["a" => 1] < ["b" => 2, "c" => 3]
[:a => 10] < [:b => 20, :c => 30]
```

Equality and inequality can be used for most data types. Integers, Floats and Booleans will be compared by exact value, Strings by length, Arrays by the value and position of the elements, and Dictionaries by the the combination of key and value.

```swift
1 != 4
1.0 != 2.5
true == true
"one" == "three"
"two" == "four"
[1, 2, 3] != [1, 2]
["a" => 1, "b" => 2] != ["a" => 5, "b" => 6]
[:a => 10, :b => 20] != [:a => 50, :b => 60]
```

Boolean operators can only be used with Boolean values, namely `true` or `false`. Other data types will not be converted to truthy values.

```swift
true == true
false != true
```

Bitwise and bitshift operator apply only to Integers. Float values can't be used, even those that "look" as Integers, like `1.0` or `5.0`.

```swift
10 >> 1
12 & 5 | 3
5 ~ 2
```

### Shorthand Assignment

Operators like `+`, `-`, `*` and `/` support shorthand assignment to variables. Basically, statements like this:

```swift
count = count + 1
```

Can be expressed as:

```swift
count += 1
```

## Functions

Oro treats functions as first class, like any sane language should. It checks all the boxes: they can be passed to variables, as arguments to other functions, and as elements to data structures. They also support recursion, closures, currying, variadic parameters, you name it.

```swift
val add = fn x, y
  x + y
end
```

Parentheses are optional and for simple functions like the above, I'd omit them. Calling the function needs the parentheses though:

```swift
val sum = add(1335, 2)
```

### Type Hinting

Like in strong typed languages, type hinting can be a very useful feature to validate function arguments and its return type. It's extra useful for library functions that have no assurance of the data types they're going to get.

This function call will produce output:

```swift
val add = fn (x: Int, y: Int) -> Int
  x + y
end

println(add(5, 2))
```

This however, will cause a type mismatch runtime error:

```swift
println(add(5, "two"))
```

Oro is not a strong typed language, so type hinting is completely optional. Generally, it's a good idea to use it as a validation measure. Once you enforce a certain type, you'll be sure of how the function executes.

### Default Parameters

Function parameters can have default values, used when the parameters are omitted from function calls.

```swift
val architecture = fn bits = 6
  2 ** bits
end

puts(architecture()) // 64
writeln(architecture(4)) // 16 
```

They can be combined with type hinting and, obviously, need to be of the same declared type.

```swift
val architecture = fn bits: Int = 6
  2 ** bits
end

echo(architecture())
```

### Return Statement

Until now we haven't seen a single `return` statement. Functions are expressions, so the last line is considered its return value. In most cases, especially with small functions, you don't have to bother. However, there are scenarios with multiple return points that need to explicitly tell the interpreter.

```swift
val even = fn n
  if n % 2 == 0
    return true
  end
  false
end

puts(even(-1))
println(even(2))
``` 

The last statement doesn't need a `return`, as it's the last line and will be automatically inferred. With the `if` on the other hand, the interpreter can't understand the intention, as it's just another expression. It needs the explicit `return` to stop the other statements from being interpreted.

In the case of multiple return points, I'd advise to always use `return`, no matter if it's the first or last statement. It will make for clearer intentions. 

### Variadic

Variadic functions take an indefinite number of parameters and merge them all into a single, Array argument. Their first use would be as a sugar:

```swift
val add = fn ...nums
  var count = 0
  repeat n in nums
    count = count + n
  end
  count
end

echo(add(1, 2, 3, 4, 5)) // 15
```

Even better, they can be used for functions that respond differently based on the number of arguments:

```swift
val structure = fn ...args
  if Enum.size(args) == 2
    val key = args[0]
    val value = args[1]
    return [key: value]
  end
  if Enum.size(args) > 2
    return args
  end
  args[0]
end

echo(structure("name", "Luis")) // dictionary
puts(structure(1, 2, 3)) // array
writeln(structure(5)) // integer
```

Functions may have as many parameters as needed, as long the variadic argument is the last parameter:

```swift
val calc = fn mult, ...nums
  mult * Enum.reduce(nums, 0, fn x, acc do x + acc end)
end

println(calc(10, 1, 2, 3, 4)) // 100
```

Variadic arguments can even have default values:

```swift
val joins = fn (glue: String, ...words = ["hello", "there"])
  String.join(words, glue)
end

puts(joins(" ")) // "hello there"
```

### Arrow Functions

Very useful when passing short functions as arguments, arrow functions provide a very clean syntax. They're handled internally exactly like normal functions. The only difference is that they're meant as a single line of code, while normal functions can handle blocks.

This normal function:

```swift
val sub = fn x
  x - 5
end
```

Is equivalent to:

```swift
val sub = (x) -> x - 5
```

They're not that useful to just spare a couple lines of code. They shine when passed as arguments:

```swift
Enum.map([1, 2, 3, 4], (x) -> x * 2)
Enum.reduce(1..10, 0, (x, accum) -> x + accum)
```

### Closures

Closures are functions inside functions that hold on to values from the parent and "close" them when executed. This allows for some interesting side effects, like currying:

```swift
val add = fn x
  fn y
    x + y
  end
end

echo(add(5)(7)) // 12
```

Some would prefer a more explicit way of calling:

```swift
val add_5 = add(5) // returns a function
val add_5_7 = add_5(7) // 12
```

You could nest a virtually unlimited amount of functions inside other functions, and all of them will have the scope of the parents.

### Recursion

Recursive functions calculate results by calling themselves. Although loops are probably easier to mentally visualize, recursion provides for some highly expressive and clean code. Technically, they build an intermediate stack and rewind it with the correct values in place when a finishing, non-recursive result is met. It's easier to understand them if you think of how they're executed. Val's see the classic factorial example:

```swift
val factorial = fn n
  if n == 0
    return 1
  end
  n * factorial(n - 1)
end
``` 

Keep in mind that Oro doesn't provide tail call optimization, as Go still doesn't support it. That would allow for more memory efficient recursion, especially when creating large stacks.

### Tricks

As first class, functions have their share of tricks. First, they can self-execute and return their result immediately:

```swift
val pow_2 = fn x
  x ** 2
end(2)

echo(pow_2) // 4
```

Not sure how useful, but they can be passed as elements to data structures, like arrays and dictionaries:

```swift
val add = fn x, y do x + y end
val list = [1, 2, add]
list[2](5, 7) 

puts(list[2](5,7)) // 12
```

Finally, like you may have guessed from previous examples, they can be passed as parameters to other functions:

```swift
val add = fn x, factor
  x + factor(x)
end
add(5, (x) -> x * 2)

echo(add(5, (x) -> x * 2)) // 15
```

## Conditionals

Oro provides two types of conditional statements. The `if/else` is limited to just an `if` and/or `else` block, without support for multiple `else if` blocks. That's because it advocates the use of the much better looking and flexible `match` statement.

### If

An `if/else` block looks pretty familiar:

```swift
if 1 == 2
  println("1 equal to 2.")
else
  println("1 isn't equal to 2.")
end
```

Sometimes it's useful to inline it for simple checks:

```swift
val married = true
val free_time = if married then 0 else 100_000_000 end
```

### Ternary Operator

The ternary operator `?:` is a short-hand `if/else`, mostly useful when declaring variables based on a condition or when passing function parameters. It's behaviour is exactly as that of an `if/else`.

```swift
val price = 100
val offer = 120
val status = offer > price ? "sold" : "bidding"
```

Although multiple ternary operators can be nested, I wouldn't say that would be the most readable code. Actually, except for simple checks, it generally makes for unreadable code.

### Match

`Match` expressions on the other hand are way more interesting. They can have multiple whens with multiple conditions that break automatically on each successful when, act as generic if/else, and match array elements.

```swift
val a = 5
match a
when 2, 3
  println("Is it 2 or 3?")
when 5
  println("It is 5. Magic!")
else
  println("No idea, sorry.")
end
```

Not only that, but a `match` can behave as a typical if/else when no control condition is provided. It basically becomes a `match true`.

```swift
val a = "Luis"
match true
when a == "Luis"
  println("Luis")
when a == "Carlos"
  println("Carlos") 
else
  println("Nobody")
end
```

### Pattern Matching

When fed arrays as the control condition, the `match` can pattern match its elements. Every argument to the match when is compared to the respective element of the array. Off course, for a match, the number of arguments should match the size of the array.

```swift
match ["game", "of", "thrones"]
when "game", "thrones"
  println("no match")
when "game", "of", "thrones"
  println("yes!")
end
```

That's probably useful from time to time, but it's totally achievable with array whens. The `match` can do much better than that.

```swift
match ["Luis", "Carlos", 2]
when "Luis", _, _
  println("Luis Something")
when _, _ 2
  println("Something 2")
else
  println("Lame movie pun not found")
end
```

The `_` is a placeholder that will match any type and value. That makes it powerful to compare arrays where you don't need to know every element. You can mix and match values with placeholders in any position, as long as they match the size of the array.

## Repeat Loop

Oro takes a modern approach to the `repeat` loop, evading from the traditional, 3-parts `repeat` we've been using repeat decades. Instead, it focuses on a flexible `repeat in` loop that iterates arrays, dictionaries, and as you'll see later, ranges.

```swift
repeat v in [1, 2, 3, 4]
  println(v)
end
```

Obviously, the result of the loop can be passed to a variable, and that's what makes them interesting to manipulate enumerables.

```swift
val plus_one = repeat v in [1, 2, 3, 4]
  v + 1
end

println(plus_one) // [2, 3, 4, 5]
```

Passing two arguments for arrays or strings will return the current index and value. For dictionaries, the first argument will be the key.

```swift
repeat i, v in "abcd"
  println(i + "=>" + v)
end
```

```swift
repeat k, v in ["name" => "Luis", "age" => 50]
  println(k)
  println(v)
end
```

With that power, you could build a function like `map` in no time:

```swift
val map = fn x, f
  repeat v in x
    f(v)
  end
end

val plus_one = map([1, 2, 3, 4], (x) -> x + 1)

println(plus_one) // [2, 3, 4, 5]
```

Without arguments, the `repeat` loop can behave as an infinite loop, much like a traditional `while`. Although there's not too many use cases, it does its job when needed. An example would be prompting the user for input and only breaking the infinite loop on a specific text.

```swift
repeat do
  val pass = prompt("Enter the password: ")
  if pass == "123"
    println("Good, strong password!")
    break
  end
end
```

The `break` and `continue` keywords, well break or skip the iteration. They function exactly like you're used to.

```swift
var i = 0
repeat do
  if i == 10
    break
  end
  i += 1
end

repeat i in 1..10
  if i == 5
    break
  end
end

repeat i in 1..10
  if i == 5
    continue
  end
end
```

*The `repeat` loop is currently naively parsed. It works for most cases, but still, it's not robust enough. I'm working to find a better solution.*

## Range Operator

The range operator is a special type of sugar to quickly generate an array of integers or strings. 

```swift
val numbers = 0..9
val huge = 999..100
val alphabet = "a".."z"
```

As it creates an enumerable, it can be put into a `repeat in` loop or any other function that expects an array.

```swift
repeat v in 10..20
  println(v)
end
```

Although its bounds are inclusive, meaning that the left and right expressions are included in the generated array, nothing stops you from doing calculations. This is completely valid:

```swift
val numbers = [1, 2, 3, 4]
repeat i in 0..Enum.size(numbers) - 1
  println(i)
end
```

## Pipe Operator

The pipe operator, is a very expressive way of chaining functions calls. Instead of ugly code like the one below, where the order of operations is from the inner function to the outers ones:

```swift
subtract(pow(add(2, 1)))
```

You'll be writing beauties like this one:

```swift
add(2, 1) |> pow() |> subtract()
```

The pipe starts from left to right, evaluating each left expression and passing it automatically as the first parameter to the function on the right side. Basically, the result of `add` is passed to `pow`, and finally the result of `pow` to `subtract`.

It gets even more interesting when combined with standard library functions:

```swift
["hello", "world"] |> String.join(" ") |> String.capitalize()
```

```swift
var name = "oro language !!!"
val expressive? = fn (x: String) -> String
  if x != ""
    return "Hello " + x
  end
end

val pipe = name |> expressive?() |> String.capitalize()
println(pipe) // "Hello Oro Language"
```

Enumerable functions too:

```swift
Enum.map([1, 2, 3], (x) -> x + 1) |> Enum.filter((x) -> x % 2 == 1)

// or even nicer

[1, 2, 3] |> Enum.map((x) -> x + 1) |> Enum.filter((x) -> x % 2 == 1)
```

Such a simple operator hides so much power and flexibility into making more readable code. Almost always, if you have a chain of functions, think that they could be put into a pipe.

## Immutability

Now that you've seen most of the language constructs, it's time to fight the dragon. Immutability is something you may not agree with immediately, but it makes a lot of sense the more you think about it. What you'll earn is increased clarity and programs that are easier to reason about.

Iterators are typical examples where mutability is seeked for. The dreaded `i` variable shows itself in almost every language's `repeat` loop. Oro keeps it simple with the `repeat in` loop that tracks the index and value. Even if it looks like it, the index and value aren't mutable, but instead arguments to each iteration of the loop.

```swift
val numbers = [10, 5, 9]
repeat k, v in numbers
  println(v) 
  println(numbers[k]) // same thing
end
```

But there may be more complicated scenarios, like wanting to modify an array's values. Sure, you can do it with the `repeat in` loop as we've seen earlier, but higher order functions play even better:

```swift
val plus_one = Enum.map([1, 2, 3], (x) -> x + 1)
println(plus_one) // [2, 3, 4]
```

What about accumulators? Val's say you want the product of all the integer elements of an array (factorial) and obviously, you'll need a mutable variable to hold it. Fortunately we have `reduce`:

```swift
val product = Enum.reduce(1..5, 1, (x, accum) -> x * accum)
println(product)
```

Think first of how you would write the problem with immutable values and only move to mutable ones when it's impossible, hard or counter-intuitive. In most cases, immutability is the better choice.

## Modules

Modules are very simple containers of data and nothing more. They're not an imitation of classes, as they can't be initialized, don't have any type of access control, inheritance or whatever. If you need to think in Object Oriented terms, they're like a class with only static properties and methods. They're good to give some structure to a program, but not to represent cars, trees and cats.

```swift
module Color
  val white = "#fff"
  val grey = "#666"
  val hexToRGB = fn hex
    // some calculations
  end
end

val background = Color.white
val font_color = Color.hexToRGB(Color.grey)
```

Because modules are interpreted and cached before-hand, properties and functions have access to each other. In contrast to modules, everything else in Oro is single pass and as such, it will only recognize calls to a module that has already been declared.

## Uses

Source file uses are a good way of breaking down projects into smaller, easily digestible files. There's no special syntax or rules to used files. They're included in the caller's scope and treated as if they were originally there. Uses are cached, so in multiple uses, only the first one is actually interpreted.

```swift
// other.oro
val name = "Luis"
val fr = fn x
  "friend " + x
end
```

```swift
// main.oro
use "other"

val phrase = name + " " + fr("Alberto")
println(phrase) // "Luis friend Alberto"
```

The file is relatively referenced from the caller and in this case, both `main.oro` and `other.oro` reside in the same folder. As the long as the extension is `.oro`, there's no need to write it in the use statement. Even the quotes can be omitted and the file written as an identifier, as long as it doesn't include a dot (as in `other.oro`) and isn't a reserved keyword.

A more useful pattern would be to wrap used files into a module. That would make for a more intuitive system and prevent scope leakage. The cat case above could be written simply into:

```swift
// Other.oro
module Other
  val name = "Luis"
  val fr = fn x
    "friend " + x
  end
end
```

```swift
// main.oro
use other

val phrase = Other.name + " " + Other.fr("Alberto")
```

Uses are expressions too! Technically, they can be used anywhere else an Integer or String can, even though it probably wouldn't make for the classiest code ever.

```swift
// exp.oro
val x = 10
val y = 15
x + y
```

```swift
// main.oro
val value = use exp
println(value) // 25

if use exp == 25
  println("Ok")
end
```

## Comments

Nothing ground breaking in here. You can write either single line or multi line comments:

```
# an inline comment

// an inline comment

/*
  I'm spanning multiple
  lines.
*/
```

### Standard Library

The Standard Library is fully written in Oro with the help of a few essential functions provided by the runtime. That is currently the best source to check out some "production" Oro code and see what it's capable of. [Read the documentation](https://github.com/luiscm/oro/wiki/Standard-Library). 

### Future Plans

In the near future, hopefully, I plan to:

- Improve the Standard Library with more functions.
- Support optional values for null returns.
- Write more tests!
- Write some useful benchmarks with non-trivial programs.

### License

This project's source code is released under the [MIT License](http://opensource.org/licenses/MIT).

### Buy Me a Coffee

Obviously, it's all licensed under the MIT license, so use it as you wish; but if you'd like to buy me a coffee, I won't complain.

- Bitcoin - `0794b1ce-67b9-48b1-9fa0-6b8cec498b04`
