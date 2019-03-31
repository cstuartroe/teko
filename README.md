# Teko

Teko is a programming language of an unusual sort - a statically typed scripting language. The plan is for Teko to be both interpreted and compiled, à la Python; its interpretation-oriented syntax gives it the ease of use of a scripting language like JS, Python, or Ruby, while its static typing will hopefully allow compiled code to run at speeds comparable to other compiled languages like Java. This close coupling of compiled and interpreted, scripting-style code will make Teko projects both quick to get started and scalable.

**Major design features include:**

 * A robust type and visibility system reminiscient of Java, which provide safety and scalability and allow compiled code to run fast
 * The freedom of organization of a scripting language - put variable declarations, control blocks, or whatever you want right in the body of a file, no need to wrap it in a class or `main()`!
 * Zero type coercion - made up for with a variety of lightweight conversion syntax (`5$` means `5._tostr()`)
 * Everything is an object, even types: `typeof(int)` yields `type`, `int.to_str()` yields `"int"`
 * Everything is a method: e.g., `+` is *always* an alias for `._add()`; `3 + 4` *is* `(3)._add(4)`
 * Powerful data structures, including primitive arrays, linked lists, labels, maps, sets, structs, and enums
 * Fastidiousness about reference vs. value - functions support a `&` argument operator to pass an argument by reference
 * Easy support for asynchrony, including `begin { }` syntax for threading, `parallel` loops, asynchronous variables and functions, and even updating and halting function types
 * Support for effect typing of functions, marking functions as `async` (asynchronous, and allowed to alter asynchronous variables), `pure` (lacking side effects), `IO` (performs input/output actions), or with `throws` annotation (to denote which errors a method may throw)

## Overview

### Syntax

The syntax is broadly reminiscient of C++ or Java:

```
str s;

for (int i in 1:15) {
    s = "";
    if (i % 3 == 0) {
        s += "Fizz ";
    }
    if (i % 5 == 0) {
        s += "Buzz ";
    }
    if (s == "") {
        s = i$;
    }
    println(s);
}

int square(int n) -> {
    return n*n;
}
int add3(int other) = (3)._add;

final enum() Location = <Brussels, Beijing, Cairo, Buenos_Aires>;
struct coords = (int lat, int long);

final {Location: coords} city_coords = {Brussels:     (50,   4),
                                        Beijing:      (39,   116),
                                        Cairo:        (30,   31),
                                        Buenos_Aires: (-34, -58)}

class Vehicle {
    static enum() modes = <Land, Sea, Air>;
    private int topspeed;
    public final virtual void travel(Location destination);
}

class Car extends Vehicle {
    ...
}

void impound(Car &c) {
    c.owner = The_Government; // darn!
}
```

### Primitive Types

Teko's primitive types include arbitrary-precision integers, floating-point numbers, booleans, chars, and strings. There is also a type `label`, which has no useful attributes and exists only for its identity. Variables are declared with a Java/C-like syntax, supporting multiple declaration:

```
int n = 3;
real x;
bool b = true;
char c = '?';
str s1 = "Hello, World!", s2;
label THIS, THAT;
```

### Composite Types

Teko also boasts a several powerful composite data types:

```
// Arrays
int[4] numarray = [1,2,3,4];
numarray[2]; // 3

// Linked Lists
int{} numlist = {1,2,3,4};
numlist:; // 1
:numlist; // {2,3,4}
numlist = 0:numlist;
numlist; // {0,1,2,3,4}

// Sets
int<> numset = <1,2,3,4>;
numset; // <2,4,3,1> (or some other order - sets aren't ordered)
3 in numset; // true

// Enums
enum(int) numbers = <ONE = 1, TWO = 2, ANOTHERONE = 1>;
typeof(numbers.ONE); // int
enum() Location = <Brussels, Beijing, Cairo, Buenos_Aires>;
typeof(Location.Beijing); // label

// Maps
{str:int} ages = {"Bob":42,"Alice":35};
ages{"Bob"}; // 42
ages{"Carol"} = 16;

// Structs
struct person = (str name, int age);
person bob = person("Bob",42);
bob; // (name = "Bob", age = 42)
bob(0); // "Bob"
bob.name; // "Bob"
```

A major design philosophy of Teko is that it contains no implicit type coercion, but enjoys a variety of lightweight typecasting features:

```
n + 1.0; // bad!
n. + 1.0; // 4.0
n + " times"; // bad!
n$ + " times"; // "3 times"
numarray + numlist; // bad!
numarray{} + numlist; // {1,2,3,4,0,1,2,3,4}
```

### Control Structures

Basic control structures should mostly  look familiar to anyone who writes Java, JS, or C, although parentheses are optional and I introduce some abbreviations of ifs:

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
} {
    ...        
};

// even more abbreviated:
cond ? { ... } { ... }; // if-else
cond ? { ... }; // a simple if


while (cond) {
    ...
};
```

I only have Python-style `for` loops, rather than C-style:

```
for (int i in numlist) {
    ...
};
```

Note the semicolons at the end - control blocks actually evaluate to values, even if the programmer chooses not to use those values.

```
int x = if (n > 0) {
    return n1;
} else {
    return 0;
}

int x = n1 > 0 ? { n1; } { 0; }; // return can be implied
int x = n1 > 0 ? n1 0; // braces can be omitted for any codeblock with one line

int x = while ( cond ) {
    return do_something(); // only last value is stored
};

int x = while ( cond ) do_something();

int{} xs = for (int i in 1:10) { 
    yield i^2; 
};

int{} xs = for (int i in 1:10) { i^2; }; // yield, rather than return, is always implied for fors
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

async int{} xs = begin parallel ( ... ) {
    ... // like with fors, yield is the implied keyword
};
```

### Functions

Defining functions in Teko with a codeblock uses the special `->` setter:

```
int add(int n1, int n2) -> { return n1 + n2; };
```

The `return` keyword is optional - it will force a function to terminate and may assist readability, but can be omitted:

```
int foo(int n1, int n2, bool b) -> { 
    b ? n1 += 5;
    n1 * n2;
};
```

In fact, as in control blocks, a single-line block can go without braces:

```
int add(int n1, int n2) ->  n1 + n2;
```

There are a few spooky things about Teko functions. Their return type and parameter set comprise the type and are thus immutable, and the parameter set is actually a struct. The definition is mutable though, and can be assigned on a different line than declaration, like any other variable (of course, using `final` in their declaration prevents this).

```
int add(int n1, int n2) -> { return n1 + n2 };
add(2,2); // 4
add._args; // (int n1, int n2)
add._args = (str s, bool b); // bad!
add -> { return n1 - n2; }; // fine
add(2,2); // 0
```

Structs are actually pretty powerful data types in Teko. They support default values and passing by reference, and by extension so do functions. Parameter defaults work essentially identically to Python, allowing for positional arguments followed by keyword arguments.

```
struct body = (int nose ? 1, int eyes ? 2, int fingers ? 10);

alice = body();
alice; // (nose = 1, eyes = 2, fingers = 10)

carol = body(fingers = 12); // supports keyword arguments
carol; // (nose = 1, eyes = 2, fingers = 12)

void do_something(int n ? 0, str s, bool b ? true, real{} xs) = { ... };
do_something(2,"Hi!",xs = {1.1,2.2});

void remove_finger(body) -> {finger -= 1};
typeof(remove_finger); // void(int nose ? 1, int eyes ? 2, int fingers ? 10)

void remove_finger_real(body b) -> {body.finger -= 1};
typeof(remove_finger_real); // void(body b)

remove_finger_real(carol);
carol.fingers; // 12 - oops!

void no_wait_really_remove_finger_this_time(body &b) -> {body.finger -= 1};
typeof(no_wait_really_remove_finger_this_time); // void(body &b);

no_wait_really_remove_finger_this_time(**carol);
carol.fingers; // 11 - ouch!
```

### Classes and Objects

Everything in Teko is an object with attributes:

```
isinstance(3, int); // true
isinstance(int, int); // false
isinstance(int, type); // true
isinstance(3, obj); // true
isinstance(int, obj); // true
isinstance(type, obj); // true
isinstance(type, type); // true
isinstance(obj, type); // true!
```

Teko has a heirarchical type system very similar to Java, with classes in a heirarchy:

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

### Asynchrony and effect typing

As covered above, Teko supports asynchronous `begin` and `parallel` loops. It also supports `async` variables, which are allowed to be modified by asynchronous routines which they are not inside. `async` functions automatically act as though wrapped in `begin`s, starting in their own thread.

```
async int n1;
int n2;

begin {
    n1 = 3; // fine
    n2 = 4; // no!
}

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
 
Variables set to the output of `async` functions can only be `async`, naturally. Thus, `async` is also the first example of Teko's *effect typing* - annotating functions for their permitted behaviors. Asynchronous variables have an attribute denoting their state. The `await` function sleeps until the state of an asynchronous variable is `TERMINATED`.

```
ASYNC_STATE; // <UNINITIALIZED, HANGING, UPDATING, TERMINATED>

async int n1;
n1._state; // UNINITIALIZED

async void foo(int &n) -> {
    sleep(500);
    n = 5;
}

foo(n1);
n1.state; // UNINITIALIZED
sleep(1000);
n1.state; // TERMINATED

int n2;
foo(n2);
await(n2);
n2._state; // TERMINATED
```

You might have noticed the two other states not mention - `HANGING` and `UPDATING`. These are related to two other types of asynchronous behavior in Teko. Hanging functions can stall and need to be resumed. Updating functions can output multiple values at different times, using the `yield` keyword instead of `return`. The `collect` function compiles outputs of updating functions into lists.

```
async Book find_book(str title) hangs -> {
    Book b;
    
    b = search_shelf(title);
    if (b is null) { hang; }
    else { return b; }
    
    b = search_floor(title);
    if (b is null) { hang; }
    else { return b; }
    
    b = search_entire_library(title);
    return b;
}

Book b;
b; // null

b = find_book("Great Expectations");
await(b); // also waits for hanging, or new update
b._state; // HANGING - guess it wasn't on the shelf....
b; // null

b.resume();
await(b);
b._state; // TERMINATED - oh good! it was on the floor!
b; // Book(title = "Great Expectations", author = "Charles Dickens", ....)

int position() updates -> {
    pos = 0;
    while pos <= 100 {
        sleep(50);
        pos += 1;
        yield pos;
    }
}

Sprite link = Sprite(image = "link.png");
int x = position();
while (x._state != TERMINATED) {
    await(x); // waits for an update
    screen.blit(sprite = link, x = x, y = 0); // Draws Link on the screen!
}
x; // 100

async int{} positions = collect(position);
positions._state; // UPDATING
positions; // {}

await(positions); // waits for first update
positions; // {0}

await(positions, <TERMINATED>); // waits for positions to be terminated
positions; // {0, 1, 2, 3, ... 100}

typeof(await); // <A> void(async A var, ASYNC_STATE<> ? <UPDATING, HANGING, TERMINATED>)
```
So what are some other type annotations? By default, functions may access variables from their outer context, but `pure` prohibits this:

```
int n;

void foo() -> {
    n = 5;
}

void foo() pure -> {
    n = 5; // bad!
}
```

`IO` is required for any function that performs input/output actions like accessing files or accepting user input:

```
str s;

void foo() -> {
    s = input(); // no!
}

void foo() IO -> {
    s = input();
}
```

A last form of effect typing that should be familiar to any Java coder is the `throws` syntax, required for functions that may throw particular error types:

```
void foo() -> {
    if ( weather() == RAINY ) {
        throw ClothingError("Oh no! My socks are all wet!"); // bad
    }
}

void bar() throws ClothingError -> {
    if ( weather() == RAINY ) {
        throw ClothingError("My socks are wet but at least my side effects are safely annotated"); // bad
    }
}
```

## Acknowledgements

I want to thank [@iafisher](https://github.com/iafisher), [@Pierre-vh](https://github.com/Pierre-vh), and [@goldfirere](https://github.com/goldfirere/) for their assistance so far in thinking through the design and implementation of Teko.
