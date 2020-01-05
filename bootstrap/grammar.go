package bootstrap

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/arr-ai/arrai/grammar/parse"
)

var (
	grammarR = Rule("grammar")
	stmt     = Rule("stmt")
	prod     = Rule("prod")
	term     = Rule("term")
	named    = Rule("named")
	atom     = Rule("atom")
	quant    = Rule("quant")
	ident    = Rule("IDENT")
	str      = Rule("STR")
	intR     = Rule("INT")
	re       = Rule("RE")
	comment  = Rule("COMMENT")

	// WrapRE is a special rule to indicate a wrapper around all regexps and
	// strings. When supplied in the form "pre()post", then all regexes will be
	// wrapped in "pre(?:" and ")post" and all strings will be escaped using
	// regexp.QuoteMeta then likewise wrapped.
	WrapRE = Rule(".wrapRE")
)

const grammarGrammarSrc = `
// Non-terminals
grammar -> stmt+;
stmt    -> COMMENT | prod;
prod    -> IDENT "->" term+ ";";
term    -> term:"^"
         ^ term:"|"
         ^ named+;
named   -> (IDENT "=")? named
         ^ atom quant?;
atom    -> IDENT | STR | RE | "(" term ")";
quant   -> /{[?*+]}
         | "{" INT? "," INT? "}"
         | /{<:|:>?} "!"? named "!"?;

// Terminals
IDENT   -> /{[A-Za-z_\.]\w*};
STR     -> /{"((?:\\.|[^\\"])*)"};
INT     -> /{\d+};
RE      -> /{/{((?:\\.|[^\\\}])*)\}};
COMMENT -> /{//.*$|(?s:/\*(?:[^*]|\*+[^*/])\*/)};

// Special
.wrapRE -> /{\s*()\s*};
`

var grammarGrammar = Grammar{
	// Non-terminals
	grammarR: Some(stmt),
	stmt:     Oneof{comment, prod},
	prod:     Seq{ident, S("->"), Some(term), S(";")},
	term: Stack{
		Delim{Term: term, Sep: S("^")},
		Delim{Term: term, Sep: S("|")},
		Some(named),
	},
	named: Stack{
		Seq{Opt(Seq{ident, S("=")}), named},
		Seq{atom, Opt(quant)},
	},
	atom: Oneof{ident, str, re, Seq{S("("), term, S(")")}},
	quant: Oneof{
		RE(`[?*+]`),
		Seq{S("{"), Opt(intR), S(","), Opt(intR), S("}")},
		Seq{RE(`<:|:>?`), Opt(S("!")), named, Opt(S("!"))},
	},

	// Terminals
	ident:   RE(`[A-Za-z_\.]\w*`),
	str:     RE(`"((?:\\.|[^\\"])*)"`),
	intR:    RE(`\d+`),
	re:      RE(`/{((?:\\.|[^\\\}])*)\}`),
	comment: RE(`//.*$|(?s:/\*(?:[^*]|\*+[^*/])\*/)`),

	// Special
	WrapRE: RE(`\s*()\s*`),
}

func nodeRule(v interface{}) Rule {
	tag := v.(parse.Node).Tag
	backslash := strings.IndexRune(tag, '\\')
	return Rule(tag[:backslash])
}

type Grammar map[Rule]Term

// Build the grammar grammar from grammarGrammarSrc and check that it matches
// GrammarGrammar.
var core = func() Parsers {
	parsers := grammarGrammar.Compile()

	r := parse.NewScanner(grammarGrammarSrc)
	v, err := parsers.Parse(grammarR, r)
	if err != nil {
		panic(err)
	}
	if err := parsers.Grammar().ValidateParse(v); err != nil {
		panic(err)
	}
	g := v.(parse.Node)

	newGrammarGrammar := NewFromNode(g)

	if diff := DiffGrammars(grammarGrammar, newGrammarGrammar); !diff.Equal() {
		panic(fmt.Errorf(
			"mismatch between parsed and hand-crafted core grammar"+
				"\nold: %v"+
				"\nnew: %v"+
				"\ndiff: %#v",
			grammarGrammar, newGrammarGrammar, diff,
		))
	}

	return newGrammarGrammar.Compile()
}()

func Core() Parsers {
	return core
}

// ValidateParse performs numerous checks on a generated AST to ensure it
// conforms to the parser that generated it. It is useful for testing the
// parser engine, but also for any tools that synthesise parser output.
func (g Grammar) ValidateParse(v interface{}) error {
	rule := nodeRule(v)
	return g[rule].ValidateParse(g, rule, v)
}

// Unparse inverts the action of a parser, taking a generated AST and producing
// the source it came from. Currently, it doesn't quite do that, and is only
// being used for quick eyeballing to validate output.
func (g Grammar) Unparse(v interface{}, w io.Writer) (n int, err error) {
	rule := nodeRule(v)
	return g[rule].Unparse(g, v, w)
}

// Parsers holds Parsers generated by Grammar.Compile.
type Parsers struct {
	parsers map[Rule]parse.Parser
	grammar Grammar
}

func (p Parsers) Grammar() Grammar {
	return p.grammar
}

func (p Parsers) ValidateParse(v interface{}) error {
	return p.grammar.ValidateParse(v)
}

func (p Parsers) Unparse(v interface{}, w io.Writer) (n int, err error) {
	return p.grammar.Unparse(v, w)
}

// Parse parses some source per a given rule.
func (p Parsers) Parse(rule Rule, input *parse.Scanner) (interface{}, error) {
	var v interface{}
	if p.parsers[rule].Parse(input, &v) {
		if input.String() == "" {
			return v, nil
		}
		return nil, fmt.Errorf("unconsumed input: %v", input.Context())
	}
	return nil, fmt.Errorf("failed to parse %s", rule)
}

// Term represents the terms of a grammar specification.
type Term interface {
	fmt.Stringer
	Parser(name Rule, c cache) parse.Parser
	ValidateParse(g Grammar, rule Rule, v interface{}) error
	Unparse(g Grammar, v interface{}, w io.Writer) (n int, err error)
	Resolve(oldRule, newRule Rule) Term
}

type Associativity int

func NewAssociativity(s string) Associativity {
	switch s {
	case ":":
		return NonAssociative
	case ":>":
		return LeftToRight
	case "<:":
		return RightToLeft
	}
	panic(BadInput)
}

func (a Associativity) String() string {
	switch {
	case a < 0:
		return "<:"
	case a > 0:
		return ":>"
	}
	return ":"
}

const (
	RightToLeft Associativity = iota - 1
	NonAssociative
	LeftToRight
)

type (
	Rule  string
	S     string
	RE    string
	Seq   []Term
	Oneof []Term
	Stack []Term
	Delim struct {
		Term            Term
		Sep             Term
		Assoc           Associativity
		CanStartWithSep bool
		CanEndWithSep   bool
	}
	Quant struct {
		Term Term
		Min  int
		Max  int // 0 = infinity
	}
	Named struct {
		Name string
		Term Term
	}
)

func NonAssoc(term, sep Term) Delim { return Delim{Term: term, Sep: sep, Assoc: NonAssociative} }
func L2R(term, sep Term) Delim      { return Delim{Term: term, Sep: sep, Assoc: LeftToRight} }
func R2L(term, sep Term) Delim      { return Delim{Term: term, Sep: sep, Assoc: RightToLeft} }

func Opt(term Term) Quant  { return Quant{Term: term, Max: 1} }
func Any(term Term) Quant  { return Quant{Term: term} }
func Some(term Term) Quant { return Quant{Term: term, Min: 1} }

func Name(name string, term Term) Named {
	return Named{Name: name, Term: term}
}

func join(terms []Term, sep string) string {
	s := []string{}
	for _, t := range terms {
		s = append(s, t.String())
	}
	return strings.Join(s, sep)
}

func (g Grammar) String() string {
	keys := make([]string, 0, len(g))
	for key := range g {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)

	var sb strings.Builder
	count := 0
	for _, key := range keys {
		if count > 0 {
			sb.WriteString("; ")
		}
		fmt.Fprintf(&sb, "%s -> %v", key, g[Rule(key)])
		count++
	}
	return sb.String()
}

func (t Rule) String() string  { return string(t) }
func (t S) String() string     { return fmt.Sprintf("%q", string(t)) }
func (t RE) String() string    { return fmt.Sprintf("/%v/", string(t)) }
func (t Seq) String() string   { return join(t, " ") }
func (t Oneof) String() string { return join(t, " | ") }
func (t Stack) String() string { return join(t, " >> ") }
func (t Delim) String() string { return fmt.Sprintf("%v%s%v", t.Term, t.Assoc, t.Sep) }
func (t Quant) String() string { return fmt.Sprintf("%v{%d,%d}", t.Term, t.Min, t.Max) }
func (t Named) String() string { return fmt.Sprintf("<%s>%v", t.Name, t.Term) }
