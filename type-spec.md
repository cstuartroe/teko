# a specification of teko's type system

## overview

teko has a unified type system, in which every type is a subtype of `obj`, 
and fully first-class types, which are objects and can used in expressions.

[click here](http://www.conorstuartroe.com/tekotypes) for an interactive diagram of teko's types.

## types and their fields

**obj**

parent: `obj`

fields:

`str to_str()` - returns a string representation of the object

`bool is(obj other)`
 
---

**type**

parent: `obj`

fields:

`{str: type} fields` - a map of all the fields of the type that all instances have

`type parent` - the direct parent type 

`bool eq(type other)` - evaluates whether a type is equivalent to another type
 
---

**bool**

parent: `obj`

fields:

`bool and(bool other)`

`bool or(bool other)`

`bool eq(bool other)`

`bool not()`

---

**real**

parent: `obj`

fields:

`real add(real other)`

`real sub(real other)`

`real mult(real other)`

`real div(real other)`

`real exp(real other)`

`int compare(real other)`

---

**int**

parent: `real`

fields:

`int add(int other)`

`int sub(int other)`

`int mult(int other)`

`int div(int other)`

`int exp(int other)`

`int mod(int other)`

*because `int` is a subtype of real, it still supports the operations with `real`s specified in the fields of `real`.*

---

**char**

parent: `obj`

fields: 

`int compare(char other)`

---

**label**

parent: `obj`

fields:

`bool eq(label other)`

---

**unit**

parent: `obj`

fields:

`unit mult(real x)`

`unit div(real x)`

`unit mult(unit other)`

`unit div(unit other)`

`unit exp(int e)`

`bool eq(unit other)`

`enum(unit) units`

---

**struct**

parent: `type`

fields:

`str[] labels`

`type[] types`

`obj[] defaults`

`bool[] by_ref`

`bool[] mutable`

## type families

*parameterized with generic types `X` and `Y`*

---

**~X**

parent: `X`

fields:

---

**iterable(X)**

parent: `obj`

fields:

`bool contains(X e)`

---

**X[]**

parent: `iterable(X)`

fields:

`X get(int i)`

`bool eq(X[] other)`

`void append(X e)`

*`str` is char[]*

---

**X{}**

parent: `iterable(X)`

fields:

`bool eq(X{} other)`

`void add(X e)`

---

**enum(X)**

parents: `iterable(X)`, `type`

fields:

---

**instances of enum(X)**

parent: `X`

fields:

---

**instances of struct**

parent: `obj`

fields:

*the fields are declared as part of the value of the instance.*

*e.g., `(int n, str s)` has the fields `int n` and `str s`.*

## de-sugaring

**arithmetic**

`a + b` becomes `a.add(b)`, or `b.add(a)` if the former fails a type check.

 * likewise for `-` to `sub`, `*` to `mult`
 * `/` similarly becomes `div`, `^` to `exp`, and `%` to `mod`, except that commutativity is not assumed
 
**comparison**

a type `X` may implement either `int compare(X other)` or `bool eq(X other)`

for `X a, b`,

if `X` implements `eq`:
 * `a == b` becomes `a.eq(b)`, `a != b` becomes `!(a.eq(b))`
 * `<`, `>`, `<=`, `>=` are invalid operations
 
if `X` implements `compare`:
 * `a == b` becomes `a.compare(b) == 0`
 * similar for other comparators
 
**exploding**

for any expression `(e)`, `~e` becomes `explode(e)`

teko core has two methods called `explode`:
 * `type explode(type X)` returns an exploded type
   * e.g., `~int` is the type of `1, 2, 3`
 * `~X explode(X[] xarr)` returns an exploded object
   * e.g., `~[1, 2, 3]` returns `1, 2, 3`
