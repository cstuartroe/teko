# The Teko Type System

## Introduction to structural typing

The most fundamental fact about the Teko type system is that it is structural, which means that
subtype relationships need not be explicitly asserted, but that the Teko type checker verifies
subtype relationships as needed according to a set of rules based on the content of the types.

Teko has multiple different *kinds* of types, each with their own internal structure, and the
exact ruleset followed when determining whether one type is a subtype of another depends
on the kind of type that the proposed subtype and supertype each are (types of different kinds
may still be in a subtype relationship - more on that later).

Before diving into a comprehensive overview of the subtyping rules in Teko, I'll give a simple example
of two types which are in a subtype relationship, to illustrate the idea of structural typing.

### Structural subtyping: an example

The two types in this example are both *object types*, because object types are the most common kind
of type in Teko, and because the subtype rules for object types are perhaps most intuitive.

Object types consist of a set of key-value pairs, with keys being field names and values being types.
The first object type that features in our example is the following:

```
type Person = {
  name: str,
  age: int,
}
```

This type describes objects which have a `name` attribute of type `str` (the Teko string type) and
an `age` attribute of type `int` (the Teko integer type). For example, the following is a correct
declaration of a variable of type `Person`:

```
alice : Person = {
  name: "Alice",
  age: 25,
}
```

It should be apparent than `alice` has two attributes: `alice.name` evaluates to the string `"Alice"`,
and `alice.age` evaluates to the integer `25`. In Teko, these two facts comprise the entire value
of `alice`, thus the type `Person` with its two fields is an appropriate description of `alice`.

The second object type for this example is a supertype of `Person`:

```
type Named = {
  name: str,
}
```

Any object with a `name` attribute of type `str` can be considered an instance of this type. This
implies that `alice`, and in fact all instances of `Person`, are also instances of `Named`, and thus
that `Person` is a subtype of `Named`!

To further illustrate the point, suppose that a function requires an argument of type `Named`:

```
fn printName(n: Named) {
  println(n.name)
}
```

The Teko type checker ought to allow any object known to have a `name` attribute of type `str` to be
passed to `printName`, even if it was annotated with another type in its declaration and there was
no explicitly written relationship between the two types. Even though `Alice` was declared
as an instance of `Person`, we want to be able to pass `alice` as an argument to `printName` 
with no issue.

This is exactly the nature of structural typing; upon checking the function call

```
printName(alice)
```

the type checker verifies that the type of `alice` - that is, `Person` - is a subtype of the asked-for
type - that is, `Named`.

What exactly is the natural ruleset to follow when making this verification? For object types,
the rule is roughly that all fields present in the proposed supertype must be present in the
proposed subtype as well (more depth later). Because the fields of `Named` are a subset of the
fields of `Person`, the type checker successfully verifies that `Person` is a subtype of `Named`,
and thus that `printName(alice)` is a valid expression.

### Why structural typing?

As seen in the example above, structural typing allows for the greater ease of use that comes from
implicit relationships between types, without sacrificing runtime safety. In fact, much of the aim
of Teko is to mimic the easy and flexible feel of writing in a dynamically typed language, while
maintaining the ability to statically analyze code and make a much stronger set of pre-runtime
guarantees than in dynamic languages.

To that end, the motivating intention behind all subtyping rules defined for Teko is that they should
be as permissive as possible while maintaining safety. Ideally, if a programmer could read a block of
code and infer that a type error is impossible, then that block of code should pass the type checker,
and the Teko type checker is willing to go to fairly great lengths to try to attain this goal.

## Kinds of types

Teko sports at least five kinds of types: object types, function types, tuple types, union types,
and generic types.

In covering each of the kinds, I will also present the criteria for subtyping between two types of
that kind, as it tends to fit naturally into discussion of the nature of the kind.
I'll reserve another section about ([subtyping between kinds](#subtyping-between-kinds)),
as there are many more possible combinations but the rules tend to either be quite simple or reduce
to the subtyping rules given in this section.

### Object types

Object types are defined by a set of key-value pairs, with keys being field names and values being types,
as in the earlier example:

```
type Person = {
  name: str,
  age: int,
}
```

Internally, object types are represented by the following Go struct:

```go
type BasicType struct {
	name   string
	fields map[string]TekoType
}
```

The `name` field is only used for debugging and not used in determining subtyping. Thus, the only
information characterizing object types as far as the type checker is concerned is the mapping
of field names to types.

Field names generally match the regular expression `[a-zA-Z_][a-zA-Z0-9_*]` (the same as Teko
variable names), and the Teko parser will not parse type definitions with fields that violate it.
However, this guideline is intentionally broken by some special object types defined by the Teko
interpreter, specifically to avoid being accidentally or maliciously subtyped by user-defined types.

In fact, many special interpreter-defined types are in fact just object types (cf. 
[Special Object Types](#special-object-types)).

### Function types

Teko function types are defined by an ordered sequence of arguments (each of which has a type, and
optionally a name and a default value) as well as a single return type.

Function types are defined with the following two structs in Go:

```go
type FunctionArgDef struct {
	Name    string
	ttype   TekoType
	Default bool
}

type FunctionType struct {
	rtype   TekoType
	argdefs []FunctionArgDef
}
```

`Name` and `Default` are both nullable fields of `FunctionArgDef`, but `ttype` is always present.
The actual AST value of `Default` is not taken into account by the typechecker; however, its presence
or absence does affect subtyping and the resolution of function calls. `Name` is also used by the
typechecker for both determining subtype relationships and resolving function calls.

#### Argument names, positional vs. keyword arguments

A simple example of a function type in Teko is

```
type StrToBool = fn(str): bool
```

This function type has a single nameless argument of type `str` with no default value, and a return
type `bool`.

An example of a function of this type is

```
fn hasEvenLength(s: str): bool {
  s.size % 2 == 0
}
```

Actually, the type checker would evaluate this function definition and record `hasEvenLength` as
having type `fn(s: str): bool`. However, in contexts where we need a function which accepts a single
string and returns a boolean value, `hasEvenLength` ought to suffice:

```
fn judgeName(criterion: StrToBool, p: Person) {
  if criterion(p.name) {
    println("I like your name!")
  } else {
    println("I don't like your name...")
  }
}

judgeName(hasEvenLength, alice)
```

This implies something about function subtyping rules! That is, if a function type has a nameless
argument, a function may be a subtype if it has the same argument with any name present.

Note that user-defined functions cannot actually have nameless arguments - it would be impossible
to use a nameless argument in the function body! Unsurprisingly, there is no syntax for defining
a function with a nameless argument. Nonetheless, function *types* with nameless arguments can be
quite useful, as above; we don't care what name a function assigns to an argument if we are only
going to pass that argument positionally, so we'd like to be able to define a function type
which is agnostic as to argument name.

Argument names are considered by the Teko type checker because function arguments aren't always
passed positionally in Teko. Function arguments may also be passed as keyword arguments, in which
case their name obviously becomes salient. For instance, the last line of the following ought not
to pass type checking:

```
type StrCalledNameToBool = fn(name: str): bool

fn judgeName(criterion: StrCalledNameToBool, p: Person) {
  if criterion(name: p.name) {
    println("I like your name!")
  } else {
    println("I don't like your name...")
  }
}

judgeName(hasEvenLength, alice)
```

This is because this version of `judgeName` passes a keyword argument to `criterion` with the
key `name`, and `hasEventLength` doesn't have an argument called `name`.

This implies something else about function subtyping - a function type with a named argument
cannot be subtyped by another function type which has either a different or null name for the same
argument!

It is a syntactic constraint in Teko that all positional arguments must be passed prior to
all keyword arguments. Consider a function type like

```
type TakesBoolAndStr = fn(b: bool, str): null
```

Because the second argument lacks a name in this type definition, it must be passed positionally;
therefore, the first argument must be passed positionally as well:

```
fn giveBoolAndStr(f: TakesBoolAndStr) {
  f(true, "Hello")
}
```

As such, the name `b` of the first argument in the definition of `TakesBoolAndStr` needlessly
reduces the set of functions which could be considered instances of `TakesBoolAndStr`. For instance,
the following function

```
fn stateEmotion(isHappy: bool, name: str) {
  if isHappy {
    println(name + " is happy")
  } else {
    println(name + " is unhappy")
  }
}
```

would not be an instance of `TakesBoolAndStr` according to the rule determined previously, even though
it ought to be usable by `giveBoolAndStr` or any other context demanding a `TakesBoolAndStr`.

Because function types with named arguments followed by nameless arguments would introduce this
type of unhelpful constraint on usage, such function types are not allowed; all nameless arguments
must precede all named arguments.

#### Number of arguments, and default values


