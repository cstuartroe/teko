package lexparse

type tokenType string

const (
	SymbolT tokenType = "symbol"

	StringT = "string"
	CharT   = "char"
	IntT    = "int"
	FloatT  = "float"
	BoolT   = "bool"

	VarT = "var"
	FnT  = "fn"

	InT    = "in"
	ForT   = "for"
	WhileT = "while"
	IfT    = "if"
	ThenT  = "then"
	ElseT  = "else"
	ScopeT = "scope"

	SetT  = "set"
	MapT  = "map"
	DataT = "data"
	TypeT = "type"
	DoT   = "do"

	BinopT      = "binary operation" // + - * / ^ % & |
	ComparisonT = "comparison"       // == != < > <= >=
	UpdaterT    = "updater"          // <- += -= *= /= ^= %= &= |=
	PrefixT     = "prefix"           // ! ~ ?
	SuffixT     = "suffix"           // $ # .

	LParT      = "("
	RParT      = ")"
	LSquareBrT = "["
	RSquareBrT = "]"
	LCurlyBrT  = "{"
	RCurlyBrT  = "}"

	DotT       = "."
	QMarkT     = "?"
	EllipsisT  = ".."
	CommaT     = ","
	SemicolonT = ";"
	ColonT     = ":"
	SubtypeT   = "<:"
	ArrowT     = "->"
	EqualT     = "="
	PipeT      = "|>"

	LineCommentT = "line comment"
	BlockCommentT = "block comment"
)

// Go doesn't have sets, which is dumb.
var punct_combos map[string]bool = map[string]bool{
	"==": true,
	"!=": true,
	"<=": true,
	">=": true,
	"+=": true,
	"-=": true,
	"*=": true,
	"/=": true,
	"^=": true,
	"%=": true,
	"&=": true,
	"|=": true,
	"->": true,
	"<-": true,
	"<:": true,
	"..": true,
	"|>": true,
}

var binops map[string]string = map[string]string{
	"&": "and",
	"|": "or",

	"+": "add",
	"-": "sub",

	"%": "mod",

	"*": "mult",
	"/": "div",

	"^": "exp",
}

var comparisons map[string]string = map[string]string{
	"==": "eq",
	"!":  "neq",
	"<":  "lt",
	">":  "gt",
	"<=": "leq",
	">=": "geq",
}

var updaters map[string]string = map[string]string{
	"+=": "add",
	"-=": "sub",
	"*=": "mult",
	"/=": "div",
	"^=": "exp",
	"%=": "mod",
	"&=": "and",
	"|=": "or",
}

var prefixes map[string]string = map[string]string{
	"!": "not",
	"~": "explode",
}

var suffixes map[string]string = map[string]string{
	"$": "to_str",
	"#": "to_int",
	".": "to_float",
	"?": "to_bool",
}

const (
	min_prec int = iota
	setter_prec
	add_sub_prec
	mult_div_prec
	exp_prec
	max_prec
)

var binop_precs map[string]int = map[string]int{
	"and": add_sub_prec,
	"or":  add_sub_prec,

	"add": add_sub_prec,
	"sub": add_sub_prec,

	"mod": mult_div_prec, // is this right??

	"mult": mult_div_prec,
	"div":  mult_div_prec,

	"exp": exp_prec,
}
