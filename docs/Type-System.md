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

It should be apparent that `alice` has two attributes: `alice.name` evaluates to the string `"Alice"`,
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

## A note on type equality

Throughout this document, when I refer to a type being a *subtype* of another, this includes the
possibility that they are equal types. This is because, within the logic of Teko's structural typing,
one type is a subtype of another iff all instances of the first type may be used wherever an instance
of the second type is required.

Similarly, one type being a *supertype* of another includes the possibility that they are equal.

Indeed, all Teko types are both a supertype and a subtype of themselves, and the Teko type checker's
test of type equality is whether a type is both a supertype and a subtype of another.

## Kinds of types

Teko sports at least four kinds of types: object types, function types, tuple types, and union types.

In covering each of the kinds, I will also present the criteria for subtyping between two types of
that kind, as it tends to fit naturally into discussion of the nature of the kind.
I've reserved another section about [subtyping between kinds](#subtyping-between-kinds),
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

In fact, many special interpreter-defined types are actually just object types (cf. 
[Special Object Types](#special-object-types)).

#### Object subtyping

For an object type `S` to be a subtype of `T`, two criteria must be met:

* All field names present in `T` must be present in `S`.
* For a given field name `n`, the type `S.n` must be a subtype of `T.n`.

These criteria ensure that any time an attribute of a variable of type `T` is used, an object
of type `S` has an attribute of the same name which can be used as expected:

```
fn printName(n: Named) {
  println(n.name)
}

alice : Person = {
  name: "Alice",
  age: 25,
}

printName(alice)
```

These criteria seem simple enough to assess, but become complex when a type is self-referential, e.g.:

```
type IntegerLinkedListNode = {
  n: int,
  next: null | LinkedListNode,
}
```

To determine whether another type `TwoValueLinkedListNode`

```
type TwoValueLinkedListNode = {
  m: int,
  n: int,
  next: null | TwoValueLinkedListNode,
}
```

is a subtype of `IntegerLinkedListNode`, then the type checker must determine whether
`null | TwoValueLinkedListNode` is a subtype of `null | LinkedListNode`, which requires knowing
whether `TwoValueLinkedListNode` is a subtype of `LinkedListNode`, which is exactly what the
type checker is in the middle of figuring out!

I've dedicated an entire section to this - [Self-Referential Types](#self-referential-types)

### Function types

Teko function types are defined by an ordered sequence of arguments (each of which has a type, and
optionally a name and a default value) as well as a single return type.

The Teko type checker uses the following two Go structs to represent function types:

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

`Name` is a nullable field of `FunctionArgDef`, but `ttype` and `Default` are always present.

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
case their name becomes salient. For instance, the last line of the following ought not
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
key `name`, and `hasEvenLength` doesn't have an argument called `name`.

This implies something else about function subtyping - a function type with a named argument
cannot be subtyped by another function type which has either a different or null name for the same
argument!

It is a syntactic constraint in Teko that all positional arguments must be passed prior to
all keyword arguments. Consider a hypothetical function type

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
reduces the set of functions which could be considered instances of `TakesBoolAndStr` - `b`
can never be passed as a keyword argument. For instance, the following function

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

Teko allows function definitions to state default values for arguments, which allow for the argument
to optionally not be passed as part of the function call.

```
fn describePerson(p: Person, hideAge: bool ? false) {
  println("Their name is " + p.name)
  if !hideAge {
    println("Their age is " + p.age$)
  }
}

describePerson(alice, true)
describePerson(alice)
```

The value of the default argument is not considered to comprise part of the function type, but
the presence or absence of a default argument is, and affects subtyping.

When denoting function types with default arguments, a trailing question mark is used to indicate
that an argument has a default value. For example, the type of `describePerson` above is

```
fn(p: Person, hideAge: bool?)
```

Because arguments with default values may be omitted, a function type with a default argument
is a subtype of a function type with all of the same arguments, but lacking that argument. For
example, the above function type is a subtype of

```
fn(p: Person)
```

The justification for this can be seen in the fact that function of type
`fn(p: Person, hideAge: bool?)` can be used anywhere that a function of type `fn(p: Person)`
is required:

```
fn passAlice(f: fn(p: Person)) {
  f(alice)
}

passAlice(describePerson)
```

Teko restricts all arguments with default values to appearing in a function definition
after all arguments with no default value.

A function type may have an unnamed argument with a default value. That is, the following is a
valid function type:

```
fn(int?)
```

An example of its use:

```
fn passZeroThenNothing(f: fn(int?)) {
  f(0)
  f()
}
```

#### Argument type and subtyping

For a function type `F` to be a subtype of another function type `G`, all argument types of `F` must be
equal to or *supertypes* of the types of corresponding arguments for `G`. For example, `fn(Named)` is
a subtype of `fn(Person)`, in that a function of the former type may be used anywhere that a function
of the latter type is required.

```
fn sayName(n: Named) {
  println(n.name)
}

fn passAlice(f: fn(Person)) {
  f(alice)
}

passAlice(sayName)
```

#### Return type

For a function type `F` to be a subtype of another function type `G`, `F`'s return type must be
a subtype of `G`'s return type:

```
fn makePerson(n: int): Person {
  {
    name: "Alice",
    age: n*5,
  }
}

type IntToNamed = fn(int): Named

fn printNameFromInt(f: IntToNamed) {
  println(f(5).name)
}

printNameFromInt(makePerson)
```

#### Summary of function subtyping criteria

Because there are many validity and subtyping rules for function types, I will present them all
again in brief for clarity, without repeating the justification for each.

For a function type `F` to be valid:

* All of its arguments must have a type
* All of its unnamed arguments must precede all of its named arguments
* All of its arguments with no default value must precede all of its arguments with default values

For a function type `F` with `m` arguments and a function type `G` with `n` arguments, `F` is a
subtype of `G` if and only if:

(in the rules below, I refer to argument positions with zero-indexing)

* `m >= n`
* For all `i, 0 <= i < n`, the `i`th argument type of `F` is a supertype of the `i`th argument type of `G`
* For all `i, 0 <= i < n`, the `i`th argument of `G` either has no name, or has the same name as the `i`th argument of `F`
* For all `i, i <= n < m`, the `i`th argument of `F` has a default value
* The return type of `F` is a subtype of the return type of `G`

### Tuple types

Tuple types are defined with a sequence of types.

```
type TupType = (int, str, bool)
```

The only way to use a value of tuple type in Teko is by unpacking it:

```
tup: TupType = (0, "", false)

n, s, b := tup

println(s)
```

Tuples and tuple types must have length of at least two.

#### Tuple subtyping

A tuple type `S` is a subtype of another tuple type `T` if they have the same number of elements,
and if each element of `S` is a subtype of the corresponding element of `T`:

```
fn doSomeThings(t: (int, Named)) {
  n, named := t

  n2 : int = n + 1

  println(named.name)
}

tup : (int, Person) = (0, alice)

doSomeThings(tup)
```

### Union Types

A union type in Teko is simply defined with a set of types. This is expressed in Teko by
listing types separated by a pipe `|`. For example, `int | str`:

```
fn addZero(x: int | str): int | str {
  switch x {
    case int {
      x + 0
    }

    case str {
      x + "0"
    }
  }
}

addZero(5) // 5
addZero("5") // "50"

seven : int | str = if condition then 7 else "7"

addZero(seven)
```

Union types are represented in the Teko type checker with the Go struct

```go
type UnionType struct {
	types []TekoType
}
```

Although they are stored with an ordered slice in Go, the order of constituent types does not matter.

The constituent types of a union type cannot themselves be be union types
(cf. [Union operation](#union-operation)).

#### Union subtyping

A union type `U` is a subtype of another union type `V` if each constituent type of `U` is a
subtype of `V`; that is, if each constituent type of `U` is a subtype of at least one constituent
type of `V` (cf. [Supertype is a union type](#supertype-is-a-union-type)). This guarantees that 
regardless of the underlying value of an instance of `U`, it can be used as an instance of `V`.

```
type U = int | Person

type V = int | str | Named

fn handleUnion(v: V) {
  switch v {
    case int {
      println((v + 5)$)
    }

    case str {
      println(v)
    }

    case Named {
      println(v.name)
    }
  }
}

u1 : U = 3
u2 : U = alice

handleUnion(u1) // prints "8"
handleUnion(u2) // prints "Alice"
```

## Subtyping between kinds

I've organized this section primarily by the kind of supertype, and secondarily by the kind of subtype.

### Supertype is an object type

#### Subtype is a function or tuple type

Function types and tuple type are only subtypes of the empty object type `{}`. In fact, the
empty object type is a supertype of all types in Teko. Function types and tuple types cannot
be subtypes of any other object types, because functions and tuples do not have attributes.

#### Subtype is a union type

A union type is a subtype of an object type if and only if each type in the union is a
subtype of the object type. This is because, regardless of the underlying value of an
instance of the union type, it must be able to be used as the object type.

```
type Foo = {
  a: int,
}

fn printFooA(foo: Foo) {
  println(foo.a$)
}

type Bar = {
  a: int,
  b: int,
}

type Baz = {
  a: int,
  c: int,
}

b : (Bar | Baz) = if condition {
  {
    a: 1,
    b: 2,
  }
} else {
  {
    a: 1,
    c: 3,
  }
}

printFooA(b)
```

### Supertype is a function type

No other kind of type may be a subtype of a function type. This is because an object of
function type must be permitted to be called using function call syntax, and no value of a
non-function type may do so.

### Supertype is a tuple type

No other kind of type may be a subtype of a tuple type. This is because an object of a tuple
type must be permitted to be unpacked, and no value of non-tuple type may do so.

### Supertype is a union type

A non-union type `T` is a subtype of a union type `U` if and only if `T` is a subtype
of at least one of `U`'s constituent types.

```
type U = int | Person

fn handleUnion(u: U) {
  switch u {
    case int {
      println((u + 5)$)
    }

    case Person {
      println(v.name)
    }
  }
}

t : int = 3

handleUnion(t)
```
