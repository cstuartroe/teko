# Teko

Teko is a programming language of an unusual sort - a statically typed scripting language. The plan is for Teko to be both interpreted and compiled, à la Python; its interpretation-oriented organization and built-in types and data structures give it the ease of use of a scripting language like JS, Python, or Ruby, while its static typing and "functional-by-default" semantics will provide better safety and scalability.

**Major design features include:**

 * Static typing
 * The freedom of organization of a scripting language - put variable declarations, control blocks, or whatever you want right in the body of a file, no need to wrap it in a class or `main()`!
 * Zero type coercion - made up for with a variety of lightweight conversion syntax (`5$` means `5.to_str()`)
 * Everything is a method: e.g., `+` is *always* an alias for `.add()`; `3 + 4` *is* `(3).add(4)`
 * Everything is an expression - assignments and control blocks evaluate to a value
 * Out of the box data structures, including dynamic arrays, sets, maps, and algebraic data types
 * Elements of functional or functional-by-default style programming:
   * First-order functions and closures
   * Fastidiousness about mutability - mutable variables must be declared with the `var` keyword
   * Type inference - type may be omitted from declarations
   * Lightweight currying with ellipsis `..` syntax
 * Purely structural types
   * Objects can be created without a constructor, by simply listing all members
   * Class constructors are just syntactic sugar and you often don't need them - they're actually just functions!

## Installation

1. Make sure you have the [Go programming language](https://golang.org/dl/) installed.

3. Clone this repo and `cd` into it.

From here your road diverges in a yellow wood:

### Using `go install`

3. Make sure `go install` works as expected. If you're not sure how to do that, try adding `export PATH=$PATH:$(go env GOPATH)/bin` to your `.bashrc` or something similar.

4. Run `go install`

5. Try running `teko tests/simple.to`. Hopefully you don't get something that looks like

```
Command 'teko' not found, did you mean:
```

If you do, revisit step 3.

### Manually moving executable to `PATH`

3. Run `go build`

4. Move the newly created `teko` executable into any directory in your `PATH`

5. Try running `teko tests/simple.to`. Hopefully you don't get something that looks like

```
Command 'teko' not found, did you mean:
```

If you do, revisit step 4.

### Running executable from working directory

3. Run `go build`

4. Try running `./teko tests/simple.to`. If something doesn't look right, that's on us.

## Acknowledgements

I want to thank [@iafisher](https://github.com/iafisher), [@Pierre-vh](https://github.com/Pierre-vh), [@goldfirere](https://github.com/goldfirere/), and [@lwpulsifer](https://github.com/lwpulsifer) for their assistance so far in thinking through the design and implementation of Teko.
