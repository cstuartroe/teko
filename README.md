# Teko

Teko is a statically typed scripting language that aims to combine the safety and scalability of Java with the quick start and powerful data structures of Python. The plan is for Teko to be both interpreted and compiled, à la Python; its interpretation-oriented syntax gives it the ease of use of a scripting language like JS, Python, or Ruby, while its static typing will hopefully allow compiled code to run at speeds comparable to other compiled languages like Java.

## Language features

### General design philosophy

Like other scripting languages, a Teko file can consist of just one line of code - no need to wrap in a class or a `main()` method.

Teko is not whitespace sensitive, using braces and semicolons to delimit syntax.

Other than being a statically typed scripting language, other elements of its design philosophy were the avoidance of implicit type coercion and of label overloading.

### Primitive Types

Teko's primitive types include arbitrary-precision integers, floating-point numbers, booleans, chars, and strings. Variables are declared with a Java/C-like syntax, supporting multiple declaration:

```
int n = 3;
real x;
bool b = true;
char c = '?';
str s1 = "Hello, World!", s2;
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
numset; <2,4,3,1> (or some other order - sets aren't ordered)
3 in numset; // true

// Maps
{str:int} ages = {"Bob":42,"Alice":35};
ages{"Bob"}; // 42
ages{"Carol"} = 16;

// Tuples
(str,int) alice = ("Alice",35);
alice(0); // "Alice"

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

Control structures should look familiar to anyone who writes Java, JS, or C:

```
if ( cond1 ) {
    ...
} else if ( cond2 ) {
    ...
} else { 
    ...
}

while (cond) {
    ...
}

for (int i = 1; i < 10; i++) { 
    ...
}

for (int i in numlist) {
    ...
}
```

I'd also like to implement an asynchronous loop:

```
each (int i in numlist) {
    ...
}
```

And maybe an easy syntax for threading:

```
begin {
    ...
}
```

### Functions

Defining functions in Teko looks similar to any other variable assignment:

```
int add(int n1, int n2) = { return n1 + n2; };
```

However, there are a few spooky things about Teko functions. Their return type and parameter set are immutable, and the parameter set is actually a struct. The definition is mutable though, and can be assigned on a different line than declaration, like any other variable.

```
add(2,2); // 4
add.args; // (int n1, int n2)
add.args = (str s, bool b); // bad!
add = { return n1 - n2; }; // fine
add(2,2); // 0
```

Structs are actually pretty powerful data types in Teko. They support default values, and by extension so do functions. Parameter defaults work essentially identically to Python, allowing for positional arguments followed by keyword arguments.

```
struct body = (int nose ? 1, int eyes ? 2, int fingers ? 10);
alice = body();
alice; // (nose = 1, eyes = 2, fingers = 10)
carol = body(fingers = 12); // supports keyword arguments
carol; // (nose = 1, eyes = 2, fingers = 12)

void do_something(int n ? 0, str s, bool b ? true, real{} xs) { ... };
do_something(2,"Hi!",xs = {1.1,2.2});
```

### Classes and Objects

Yeah, I'd like to build classes, class hierarchies, and generics into Teko at some point.

## Open Design Issues

### Function Ordering

As a type-checked interpreted language, Teko is open to a type of issue that both typical compiled languages and typical scripting languages avoid, in different ways. Consider the following equivalent snippets in Teko, Python, Java, and Haskell:

```
[ordering.to]

int f1(int n) = {
    if (n < 0) {
        return n;
    } else {
        return f2(n);
    }
};

int f2(int n) = {
    return f1(n-1);
};

f1(5);
```

```python
[ordering.py]

def f1(n):
    if n < 0:
        return n
    else:
        return f2(n)

def f2(n):
    return f1(n-1)

f1(5)
```

```java
[Ordering.java]

public class Ordering {
    static int f1(int n) {
        if (n < 0) {
            return n;
        } else {
            return f2(n);
        }
    }
    
    static int f2(int n) {
        return f1(n-1);
    }
    
    public static void main(String[] args) {
        f1(5);
    }
}
```

```haskell
[Ordering.hs]

module Ordering where

f1 :: Int -> Int
f1 n | n < 0 = n 
     | otherwise = f2 n 

f2 :: Int -> Int
f2 n = f1 (n-1)
```

The issue is how `f2` gets type-checked in the body of `f1`, before it has been defined. In Java, all method declarations are put into a symbol table before their bodies are type-checked, so `f2`'s type is known by the time it is type-checked in the body of `f1`. In Python, the existence and type of `f2` isn't considered until `f1` is executed. When the definition of `f1` is interpreted, `f2` is not yet in scope, but by the time `f1` is actually executed `f2` has been brought into scope, and the type-checker doesn't look it up until then. The Haskell is a rather interesting case - when compiled, Haskell exhibits a similar behavior to Java. Haskell also has an interpreter, GHCi, which actually rejects the same code when I attempt to define `f1` before `f2`.

For Teko, though, there is a problem. It will attempt to type-check `f2` when it is interpreting the definition of `f1`, but can't yet. Note that this issue is only unavoidable with cyclically referential functions or type definitions. (Are there type definitions in Teko which are also open to this problem?) If `f2` did not mention `f1`, then its declaration could simply be brought above `f1`. Class definitions, for instance, are not subject to this issue, since subtyping cannot be cyclical.

A potential solution to this is Teko's unusual style of function definition, which allows the separation of declaration and body. This gives way to the following solution:

```
[ordering.to]

int f1(int n);
int f2(int n);

f1 = {
    if (n < 0) {
        return n;
    } else {
        return f2(n);
    }
};

f2 = {
    return f1(n-1);
};

f1(5);
```

More inquiry is needed to decide whether this is a sufficient solution.

### Function parameters as named structs

In Teko, function parameter sets are structs.

```
int f1(int n);
f1.args; // (int n)
```

Structs are objects which can be named with a variable, like anything else. This raises the possibility of declaring a function's parameter set as being a previously defined struct.

```
struct a = (int n);
int f1(a) = { return n; };
```

This further raises an even stranger possibility!

```
struct a;
int f1(a);
a = (int n);
f1 = { return n; };
```

Even I admit that's probably a step too far. A struct passed in a function declaration would presumably be passed by value; in fact, allowing anything like this opens the door to mutating a function's parameter set, which is definitely outlawed. It was just too interesting an idea not to write down.

### Complex number syntax

I'd really like for complex numbers to be a primitive type in Teko, but expressing them raises all manner of hairy issues for my lexical and syntactic structure. Here was my original prototype for complex numbers:

```
comp z = 3 + 2i;
```

Firstly, this looks an awful lot like `2i` is a complex number, and `3` is being implicitly coerced into adding. Secondly, is `i` a quantity (entailing an implicit multiplication in `2i`)? Is it a conversion shorthand, like `2$`? Does it affect the validity of `i` as a label elsewhere?

On a related note, is `2.` two tokens, an integer `2` followed by the existing conversion shorthand `.`, or is it a single token denoting a real? Either way results in the same type for this particular expression, but is there a situation in which it matters?

## Acknowledgements

I want to thank [@iafisher](https://github.com/iafisher) and [@goldfirere](https://github.com/goldfirere/) for their assistance so far in thinking through the design and implementation of Teko.
