# Teko

Teko is a programming language of an unusual sort - a statically typed scripting language. The plan is for Teko to be both interpreted and compiled, à la Python; its interpretation-oriented organization and built-in types and data structures give it the ease of use of a scripting language like JS, Python, or Ruby, while its static typing and "functional-by-default" semantics will provide better safety and scalability.

**Major design features include:**

 * Static typing
 * The freedom of organization of a scripting language - put variable declarations, control blocks, or whatever you want right in the body of a file, no need to wrap it in a class or `main()`!
 * Zero type coercion - made up for with a variety of lightweight conversion syntax (`5$` means `5.to_str()`)
 * Consistency of semantics:
   * First-order functions
   * Everything is a method: e.g., `+` is *always* an alias for `.add()`; `3 + 4` *is* `(3).add(4)`
   * Everything is an expression - assignments and control blocks evaluate to a value
 * Out of the box data structures, including amortized arrays, labels, maps, sets, structs, and enums
 * Elements of functional or functional-by-default style programming:
   * First-order functions and closures
   * Fastidiousness about mutability - mutable variables must be declared with the `var` keyword
   * `sees` and `modifies` annotation for all variables a function may read or write, `IO` annotation to permit functions to perform input/output actions
   * Type inference - type may be omitted from declarations
   * Lightweight currying with ellipsis `...` syntax
 * Purely structural types
   * Objects can be created without a constructor, by simply listing all members
   * "Class" definitions are just syntactic sugar - they're actually functions!

## Overview

### Syntax

The syntax is most similar to Javascript:

```
var str s;

@IO
for (int i in 1..15) {
    s = "";
    (i % 3 == 0) ? s += "Fizz ";
    (i % 5 == 0) ? s += "Buzz ";
    (s == "")    ? s = i$;
    println(s);
};

int square(int n) -> n*n;
int add3(int other) = (3).add;
let five = add3(2);

enum Location = {
  Brussels, Beijing, Cairo, Buenos_Aires
};

constructor coords(int lat, int long);

[Location: coords] city_coords = [
  Brussels:     coords(50,   4),
  Beijing:      coords(39,   116),
  Cairo:        coords(30,   31),
  Buenos_Aires: coords(-34, -58)
];



static enum Mode = {Land, Sea, Air};

type Vehicle {
  int topspeed,
  void travel(Location destination),
  Mode mode
}

class Car(int topspeed, string owner) {
  mode = Mode.Land;
  private var Location currentLocation;

  @modifies currentLocation
  void travel(Location destination) {
    println(owner + " is driving!");
    currentLocation = destination;
  }
};

void impound({string owner} &c) -> {
    c.owner = The_Government; // darn!
};
```

Teko is not whitespace sensitive.

### Declarations and Primitive Types

Teko's primitive types include arbitrary-precision integers, floating-point numbers, booleans, chars, and strings. There is also a type `label`, which has no useful attributes and exists only for its identity. Variables are declared with a Java/C-like syntax, supporting multiple declaration:

```
int n = 3;
float x;
bool q = true;
char c = '?';
str s1 = "Hello, World!", s2;
bits b = 0b10010110101;
bytes B = 0x1ac093e746f6;
```

A last, interesting primitive type is `unit`. These represent units of measure, or multiples of units. `unit`s support multiplication and division by reals and other units, as well as exponentiation.

```
unit meter, USD, year;
unit olympic_pool_volume = 2500*(meter^3);
unit us_median_salary = 56000*USD/year;
```

Variables can also be declared using the keyword `let` - expressions in Teko are not ambiguous for type, so type can be immediately inferred if the variable is set on the same line. However, variables in Teko must always know their type, so variables declared with `let` *must* be simultaneously set.

```
let n; // bad!

let n = 3; // good
```

### Composite Types

Teko also boasts a several powerful composite data types:

```
// Amortized Arrays
var int[] numarray1 = [1,2,3];
numarray1[2]; // 3
numarray1.size; // 3 and O(1)
3 in numarray; // true and O(n)
sort<int>(numarray);

numarray1[0..2]; // [1, 2] - ranges are exclusive at the end
numarray1 = numarray1[0..2]; // resizing is supported

int[4] numarray2 = [1,2,3,4]; // arrays can be declared a specific size
numarray2 = numarray2[0..2]; // error - cannot be resized

// Sets
int{} numset = {1,2,3,4};
numset; // {2,4,3,1} (or some other order - sets aren't ordered)
3 in numset; // true and O(1)

// Enums
enum<int> numbers = {ONE = 1, TWO = 2, ANOTHERONE = 1};
ONE == ANOTHERONE; // true
ONE is ANOTHERONE; // false

enum Location = {Brussels, Beijing, Cairo, Buenos_Aires};

// Maps
int[str] ages = ["Bob":42,"Alice":35];
ages["Bob"]; // 42
ages["Carol"] = 16;

// Structs
class person(str name, int age);
person bob = person("Bob", 42);
bob; // (name = "Bob", age = 42)
bob(0); // "Bob"
bob.name; // "Bob"
```

A major design philosophy of Teko is that it contains (almost) no implicit type coercion, but enjoys a variety of lightweight typecasting features:

```
n + 1.0; // bad!
n. + 1.0; // 4.0
n + " times"; // bad!
n$ + " times"; // "3 times"
numarray + numlist; // bad!
```

Exceptions to this include setting values of `struct` types in contexts when expected type is known, in which case writing out the type again is not necessary:

```
person bob = ("Bob", 42);
person[2] my_friends = [("Amanda", 31), ("Carol", 26)];

// but still bad:

let bob = ("Bob", 42); // no no
```

### Control Structures

Basic control structures should mostly look familiar to anyone who writes Java, JS, or C, although parentheses are optional and I introduce some abbreviations:

```
if ( cond1 ) {
    ...
} else if ( cond2 ) {
    ...
} else {
    ...
};

// final else and parens are optional
if cond {
    ...
} : {
    ...        
};

// even more abbreviated:
cond ? { ... } : { ... }; // if-else
cond ? { ... }; // a simple if


while cond {
    ...
};
```

I only have Python-style `for` loops, rather than C-style:

```
for (int i in numlist) {
    ...
};

// parens and type of the iterator can be omitted:

for i in numlist {
    ...
};
```

Note the semicolons at the end - control blocks actually evaluate to values, even if the programmer chooses not to use those values.

```
int x = if (n > 0) {
    n1;
} else {
    0;
};

int x = n1 > 0 ? { n1; } : { 0; }; // return can be implied
int x = n1 > 0 ? n1 : 0; // braces can be omitted for any codeblock with one line

int[] x = while ( cond ) {
    do_something(); // only last value is stored
};

int[] x = while cond do_something();

int[] xs = for (int i in 1:10) {
    i^2;
};

int[] xs = for (i in 1:10) i^2; // yield, rather than return, is always implied for fors
```

There are also asynchronous control blocks:

```
// Begins execution of the block for each element in a separate thread, and waits for the last one to finish

parallel (int i in numlist) {
    ...
};
```

And maybe an easy syntax for threading:

```
// Begins the block in a new thread

begin {
    ...
};

// Like parallel, but proceeds without waiting to finish the loops
// Essentially just abbreviates a pair of braces

begin parallel (int i in numlist) {
    ...
};

async int[] xs = begin parallel ( ... ) {
    ... // like with fors, yield is the implied keyword
};
```

### "Exploded" Data Types

Teko supports a unique class of data types called exploded types.
Exploded types are essentially iterable data structures that pass themselves to an outer context element by element.
Perhaps it's easiest to understand by example:

```
str[] names = ["Alice","Bob","Carol","Eve"];
names[0,1,2]; // ["Alice","Bob","Carol"]
names[0 to 3]; // ["Alice","Bob","Carol"]
typeof(0 to 3); // ~int - arrays slices are exploded!

int[] nums = [0,1,2];
typeof(nums); // int[]
typeof(~nums); // ~int - "~" is the explosion operator
names[~nums]; // ["Alice","Bob","Carol"]
names[~nums, 3]; // ["Alice","Bob","Carol","Eve"]

nums + 5; // type error - int[] and int cannot add
~nums + 5; // [5,6,7] - but ~int and int can!
~int <: int; // true! - an exploded int can be passed anywhere a normal int can
~nums + ~nums; // [0,1,2,1,2,3,2,3,4];

bool even(int n) -> n % 2 == 0;

even(~nums); // true, false, true - no need for a function map()!
typeof(even(~nums)); // ~bool
[even(~nums)]; // [true, false, true] - exploded types can be collected back into an array

struct pair = (int x, int y);

// without explosion, nested loops produce nested lists
let pairs1 = for (int i in nums) {
  for (int j in nums) {
    pair(i ,j);
  };
};
pairs1; // [ [(0,0),(0,1),(0,2)],
             [(1,0),(1,1),(1,2)],
             [(2,0),(2,1),(2,2)] ]

// but explosion flattens them right out!
let pairs2 = for (int i in nums) {
  ~for (int j in nums) {
    pair(i ,j);
  };
};
pairs2; // [(0,0),(0,1),(0,2),
            (1,0),(1,1),(1,2),
            (2,0),(2,1),(2,2)]
```

### Functions

Defining functions in Teko with a codeblock uses the special `->` setter:

```
int add(int n1, int n2) -> { return n1 + n2; };
```

Incidentally, I'd like to use `->` to allow lazy evaluation of normal variables, too:

```
let x -> some_horribly_time_consuming_function();

printf(x); // takes a long time to evaluate x
printf(x); // done in a flash - unlike a function call, lazy evaluation is performed once
```

The `return` keyword for functions is optional - it will force a function to terminate and may assist readability, but can be omitted:

```
int foo(var int n1, int n2, bool b) -> {
    b ? n1 += 5;
    n1 * n2;
};
```

In fact, as in control blocks, a single-line block can go without braces:

```
int add(int n1, int n2) ->  n1 + n2;
```

An alternative is the `yield` keyword - functions which use `yield` return linked lists.

```
int{} foo() -> {
  var i = 1;
  while (i < 1000) {
    i = i * (i+1);
    yield i;
  };
};
```

There are a few spooky things about Teko functions. Their return type and parameter set comprise the type and are thus immutable, and the parameter set is actually a struct. The definition can be made mutable though, and can be reassigned, like the value of any other mutable variable, though I'm not sure I endorse this behavior.

```
var int add(int n1, int n2) -> { return n1 + n2 };
add(2,2); // 4
add.args; // (int n1, int n2)
add.args = (str s, bool b); // bad!
add -> { return n1 - n2; }; // fine
add(2,2); // 0
```

Structs are actually pretty powerful data types in Teko. They support default values and passing by reference, and by extension so do functions. Parameter defaults work essentially identically to Python, allowing for positional arguments followed by keyword arguments.

```
struct body = (int nose ? 1, int eyes ? 2, var int fingers ? 10);

alice = body();
alice; // (nose = 1, eyes = 2, fingers = 10)

carol = body(fingers = 12); // supports keyword arguments
carol; // (nose = 1, eyes = 2, fingers = 12)

void do_something(int n ? 0, str s, bool b ? true, real{} xs) -> { ... };
do_something(2, "Hi!", xs = {1.1, 2.2});

void remove_finger(body) -> {finger -= 1};
typeof(remove_finger); // void(int nose ? 1, int eyes ? 2, var int fingers ? 10)

void remove_finger_real(var body b) -> body.finger -= 1;
typeof(remove_finger_real); // void(var body b)

remove_finger_real(carol);
carol.fingers; // 12 - oops!

void no_wait_really_remove_finger_this_time(var body &b) -> body.finger -= 1;
typeof(no_wait_really_remove_finger_this_time); // void(var body &b);

no_wait_really_remove_finger_this_time(carol);
carol.fingers; // 11 - ouch!
```

Teko supports anonymous functions:

```
let words = {"apple", "ball", "crayon"};

let by_length = sort(words, key = (str s) -> s.length)
```

Teko also supports currying with ellipsis `...`:

```
@IO
int foo(int n, str s, bool b) -> {
  println(n);
  println(s);
  println(b);
};

@IO
int bar(str s) = foo(5, b = true, ...);

bar("test");
// output:
// 5
// test
// true
```

and basic parameter unpacking:

```
struct foo_args = (int n, str s, bool b);
let fa = (3, "noodles", false);

foo(fa);
```

### Mutability

Teko variables and struct arguments are immutable by default. To make them mutable, the `var` keyword must be used in declaration.

```
int n = 0;
n++; // bad!

var int n = 0;
n++;

let n = 0;
n++; // also bad!

var n = 0; // "let var" is abbreviated to simply "var"
n++;

struct person = (str name, var int age);

void age_one_year(var person &p) {
  person.age++; // fine
};

void award_doctorate(var person &p) {
  person.name += ", Ph.D."; // no! name is not mutable!
};

let x1 = [1, 2, 3];
x[0] = 6; // nope

var x2 = [1, 2, 3];
x[0] = 6; // yep
```

### Asynchrony

As covered above, Teko supports asynchronous `begin` and `parallel` loops. It also supports `async` variables, which are allowed to be modified by asynchronous routines which they are not inside. `async` functions automatically act as though wrapped in `begin`s, starting in their own thread.

```
async var int n1;
var int n2;

begin {
    n1 = 3; // fine
    n2 = 4; // no!
};

@IO
async void foo() -> {
    sleep(1000);
    println("I'm awake!");
};

foo();
println("Somebody's sleeping...");
```
outputs

```
Somebody's sleeping...
I'm awake!
```

Variables set to the output of `async` functions can only be `async`, naturally. Thus, `async` is also the first example of Teko's *effect typing* - annotating functions for their permitted behaviors.

Asynchronous variables have an attribute `state`, which can be `UNINITIALIZED`, `HANGING`, `UPDATING`, or `TERMINATED`. State is set by any process or thread involved in setting the value of the variable, and provides insight into that process - an `async` variable is `UNINITIALIZED` if it has not yet been provided with a value, or to `TERMINATED` by a process as it finishes executing (more on the other two).

```
ASYNC_STATE; // <UNINITIALIZED, HANGING, UPDATING, TERMINATED>

async var int n1;
n1.state; // UNINITIALIZED

async void foo(var int &n) -> {
    sleep(500);
    n = 5;
};

foo(n1);
n1.state; // UNINITIALIZED
sleep(1000);
n1.state; // TERMINATED
```

`HANGING` and `UPDATING` are related to two other types of asynchronous behavior in Teko. Hanging functions, annotated with `@hangs` can stall and need to be resumed. Updating functions, annotated with `@updates` can output multiple values at different times, using the `yield` keyword instead of `return`. The `collect` function compiles outputs of updating functions into lists.

The `await` function sleeps until the state of an asynchronous variable is `HANGING` or `TERMINATED`, or until it is updated.

```
@hangs
async Book find_book(str title) -> {
    Book b;

    b = search_shelf(title);
    if (b is null) { hang; }
    else { return b; }

    b = search_floor(title);
    if (b is null) { hang; }
    else { return b; }

    b = search_entire_library(title);
    return b;
};

Book b;
b; // null

b = find_book("Great Expectations");
await(b); // also waits for hanging, or new update
b.state; // HANGING - guess it wasn't on the shelf....
b; // null

b.resume();
await(b);
b.state; // TERMINATED - oh good! it was on the floor!
b; // Book(title = "Great Expectations", author = "Charles Dickens", ....)

@updates
int position() -> {
    pos = 0;
    while pos <= 100 {
        sleep(50);
        pos += 1;
        yield pos;
    };
};

Sprite link = (image = "link.png");
int x = position();
while (x.state != TERMINATED) {
    await(x); // waits for an update
    screen.blit(sprite = link, x = x, y = 0); // Draws Link on the screen!
};
x; // 100

async int{} positions = collect(position);
positions.state; // UPDATING
positions; // {}

await(positions); // waits for first update
positions; // {0}

await(positions, <TERMINATED>); // waits for positions to be terminated
positions; // {0, 1, 2, 3, ... 100}

typeof(await); // <A> void(async A a, ASYNC_STATE<> ? <UPDATING, HANGING, TERMINATED>)
```

So that the variable's `state` is a useful reflection of an associated process, it is best practice that `async` variables not be made mutable, and be set a single time to some asynchronous routine.

```
@updates
async int foo() -> {
  for (i in 0:5) {
    sleep(100);
    yield i;
  };
};

@updates
async int bar() -> {
  for (i in 0:10) {
    sleep(100);
    yield i;
  };
};

async var int x;
x = foo();
x = bar();

@IO
while (x.state != TERMINATED) {
  await(x);
  print(x$ + ", "); // print does not add newline
};
// prints "0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5"
// oh no! when this point is reached, x.state has been set to TERMINATED by foo, but bar is continuing to update it!
```

### Other effect typing

So what are some other type annotations? Anything you've been seeing beginning with `@` is some type of annotation.

Functions need to be annotated with `sees` or `modifies` for any external variables they may access or modify, respectively. `IO` is required for any function that performs input/output actions like accessing files or accepting user input.

```
var int n;

@IO
@sees n
void foo() -> {
    println(n);
};

@modifies n
void bar() -> {
    n = 5;
};
```

A last form of effect typing that should be familiar to any Java coder is the `throws` syntax, required for functions that may throw particular error types:

```
void foo() -> {
  if ( weather() == RAINY ) {
    throw ClothingError("Oh no! My socks are all wet!"); // bad
  }
}

@throws ClothingError
void bar() -> {
  if ( weather() == RAINY ) {
    throw ClothingError("My socks are wet but at least my side effects are safely annotated"); // bad
  }
}

@throws ClothingError
void another_function() -> {
  bar();
}
```

### Namespaces

Teko supports namespaces within files:

```
namespace physics {
  unit c = 300000*kilometer/second;
  unit G = 6.8*(10^-11)*(meter^3)/(kilogram*(second^2));
  unit h = 6.6*(10^-34)*joule*second;
}

physics.c; // 300000*kilometer/second
```

### Classes and Objects

Everything in Teko is an object with attributes:

```
isa(3, int); // true
isa(int, int); // false
isa(int, type); // true
isa(3, obj); // true
isa(int, obj); // true
isa(type, obj); // true
isa(type, type); // true
isa(obj, type); // true!
```

Teko has a hierarchical type system very similar to Java, with classes in a hierarchy:

```
class Vehicle {
    static enum() modes = <Land, Sea, Air>;
    private int topspeed;
    public final virtual void travel(Location destination);
}

class Car extends Vehicle {
    ...
}
```

As in Java or C++, the `static` keyword makes an attribute an attribute of the class rather than of instances, and the `virtual` keyword allows final attributes to be overriden by child classes.

I need to think more about multiple inheritance.

## Acknowledgements

I want to thank [@iafisher](https://github.com/iafisher), [@Pierre-vh](https://github.com/Pierre-vh), and [@goldfirere](https://github.com/goldfirere/) for their assistance so far in thinking through the design and implementation of Teko.
