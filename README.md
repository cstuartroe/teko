# Teko

Teko is a statically typed, mixed-paradigm, general-purpose scripting language. Its major design goal is to be a fast, statically analyzable language
that feels as breezy to write as JavaScript or Python, with a type system that gets out of your way whenever possible.

To this end, it shares many features with TypeScript, including simple structural typing; a love for generics, union types and other type operations; 
type inference; and some distinctive syntactic constructions.

Other major design features include:

 * Zero type coercion - made up for with a variety of lightweight conversion syntax (`5$` means `5.to_str()`)
 * Everything is a method: e.g., `+` is *always* an alias for `.add()`; `3 + 4` *is* `(3).add(4)`
 * Everything is an expression - assignments and control blocks evaluate to a value
 * Out of the box data structures, including dynamic arrays, sets, maps, and no-frills JS-style objects
 * Elements of functional or functional-by-default style programming:
   * First-order functions and closures
   * Immutability is the default
   * Partial function application

Check out [the test suite](/tests) for some examples of the current and planned state of the language.

## Installation

Make sure you have the [Go programming language](https://golang.org/dl/) installed.

Then, simply run 

```sh
go install github.com/cstuartroe/teko@latest
```

## Acknowledgements

I want to thank [@iafisher](https://github.com/iafisher), [@Pierre-vh](https://github.com/Pierre-vh), [@goldfirere](https://github.com/goldfirere/), and [@lwpulsifer](https://github.com/lwpulsifer) for their assistance so far in thinking through the design and implementation of Teko.
