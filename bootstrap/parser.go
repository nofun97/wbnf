package bootstrap

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arr-ai/wbnf/parser"
	parse "github.com/arr-ai/wbnf/parser"
)

const (
	stackDelim = "#"

	seqTag   = "_"
	oneofTag = "|"
	delimTag = ":"
	quantTag = "?"
)

type cache struct {
	parsers    map[Rule]parse.Parser
	grammar    Grammar
	rulePtrses map[Rule][]*parse.Parser
}

func (c cache) registerRule(parser *parse.Parser) {
	if rule, ok := (*parser).(ruleParser); ok {
		c.rulePtrses[rule.t] = append(c.rulePtrses[rule.t], parser)
	}
}

func (c cache) registerRules(parsers []parse.Parser) {
	for i := range parsers {
		c.registerRule(&parsers[i])
	}
}

func (c cache) makeParsers(terms []Term) []parse.Parser {
	parsers := make([]parse.Parser, 0, len(terms))
	for _, t := range terms {
		parsers = append(parsers, t.Parser("", c))
	}
	c.registerRules(parsers)
	return parsers
}

func ruleOrAlt(rule Rule, alt Rule) Rule {
	if rule == "" {
		return alt
	}
	return rule
}

type putter func(output interface{}, extra interface{}, children ...interface{}) bool

func tag(rule Rule, alt Rule) putter {
	rule = ruleOrAlt(rule, alt)
	return func(output interface{}, extra interface{}, children ...interface{}) bool {
		parse.PtrAssign(output, parse.Node{
			Tag:      string(rule),
			Extra:    extra,
			Children: children,
		})
		return true
	}
}

func (g Grammar) clone() Grammar {
	clone := make(Grammar, len(g))
	for rule, term := range g {
		clone[rule] = term
	}
	return clone
}

func (g Grammar) resolveStacks() {
	for rule, term := range g {
		if stack, ok := term.(Stack); ok {
			oldRule := rule
			for i, layer := range stack {
				newRule := rule
				if j := (i + 1) % len(stack); j > 0 {
					newRule = Rule(fmt.Sprintf("%s%s%d", rule, stackDelim, j))
				}
				g[oldRule] = layer.Resolve(rule, newRule)
				oldRule = newRule
			}
		}
	}
}

// Compile prepares a grammar for parsing. The parser holds a copy of the
// grammar modified to support parser execution.
func (g Grammar) Compile() Parsers {
	for _, term := range g {
		if _, ok := term.(Stack); ok {
			g = g.clone()
			g.resolveStacks()
			break
		}
	}

	c := cache{
		parsers:    map[Rule]parse.Parser{},
		grammar:    g,
		rulePtrses: map[Rule][]*parse.Parser{},
	}
	for rule, term := range g {
		c.parsers[rule] = term.Parser(rule, c)
	}

	for rule, rulePtrs := range c.rulePtrses {
		term := c.parsers[rule]
		for _, rulePtr := range rulePtrs {
			*rulePtr = term
		}
	}

	return Parsers{
		parsers:    c.parsers,
		grammar:    g,
		singletons: g.singletons(),
	}
}

func Compile(grammar string) (Parsers, error) {
	v, err := Core().Parse(grammarR, parser.NewScanner(grammar))
	if err != nil {
		return Parsers{}, err
	}
	return NewFromNode(v.(parser.Node)).Compile(), nil
}

//-----------------------------------------------------------------------------

type ruleParser struct {
	rule Rule
	t    Rule
}

func (p ruleParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	panic(Inconceivable)
}

func (t Rule) Parser(rule Rule, c cache) parse.Parser {
	return ruleParser{
		rule: rule,
		t:    t,
	}
}

//-----------------------------------------------------------------------------

func eatRegexp(input *parse.Scanner, re *regexp.Regexp, output interface{}) bool {
	var eaten [2]parse.Scanner
	if n, ok := input.EatRegexp(re, nil, eaten[:]); ok {
		parse.PtrAssign(output, eaten[n-1])
		return true
	}
	return false
}

type sParser struct {
	rule Rule
	t    S
	re   *regexp.Regexp
}

func (p *sParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	return eatRegexp(input, p.re, output)
}

func (t S) Parser(rule Rule, c cache) parse.Parser {
	re := "(" + regexp.QuoteMeta(string(t)) + ")"
	if wrap, has := c.grammar[WrapRE]; has {
		re = strings.Replace(string(wrap.(RE)), "()", "(?:"+re+")", 1)
	}
	return &sParser{
		rule: rule,
		t:    t,
		re:   regexp.MustCompile(`(?m)\A` + re),
	}
}

type reParser struct {
	rule Rule
	t    RE
	re   *regexp.Regexp
}

func (p *reParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	return eatRegexp(input, p.re, output)
}

func (t RE) Parser(rule Rule, c cache) parse.Parser {
	re := "(" + string(t) + ")"
	if wrap, has := c.grammar[WrapRE]; has {
		re = strings.Replace(string(wrap.(RE)), "()", "(?:"+re+")", 1)
	}
	return &reParser{
		rule: rule,
		t:    t,
		re:   regexp.MustCompile(`(?m)\A` + re),
	}
}

//-----------------------------------------------------------------------------

type seqParser struct {
	rule    Rule
	t       Seq
	parsers []parse.Parser
	put     putter
}

func (p *seqParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	defer enterf("%s: %T %[2]v", p.rule, p.t).exitf("%v %v", &out, output)
	result := make([]interface{}, 0, len(p.parsers))
	furthest := *input
	for _, parser := range p.parsers {
		var v interface{}
		if !parser.Parse(input, &v) {
			*input = furthest
			return false
		}
		furthest = *input
		result = append(result, v)
	}
	return p.put(output, nil, result...)
}

func (t Seq) Parser(rule Rule, c cache) parse.Parser {
	return &seqParser{
		rule:    rule,
		t:       t,
		parsers: c.makeParsers(t),
		put:     tag(rule, seqTag),
	}
}

//-----------------------------------------------------------------------------

type delimParser struct {
	rule Rule
	t    Delim
	term parse.Parser
	sep  parse.Parser
	put  putter
}

func parseAppend(p parse.Parser, input *parse.Scanner, slice *[]interface{}) bool {
	var v interface{}
	if p.Parse(input, &v) {
		*slice = append(*slice, v)
		return true
	}
	return false
}

type Empty struct{}

func (p *delimParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	defer enterf("%s: %T %[2]v", p.rule, p.t).exitf("%v %v", &out, output)
	var result []interface{}

	switch {
	case parseAppend(p.term, input, &result):
	case p.t.CanStartWithSep:
		result = append(result, Empty{})
		if !parseAppend(p.sep, input, &result) {
			return false
		}
		if !parseAppend(p.term, input, &result) {
			result = append(result, Empty{})
			return p.put(output, Associativity(0), result...)
		}
	default:
		return false
	}

	start := *input
	for parseAppend(p.sep, input, &result) {
		start = *input
		if !parseAppend(p.term, input, &result) {
			break
		}
		start = *input
	}
	*input = start

	if p.t.CanEndWithSep {
		if parseAppend(p.sep, input, &result) {
			result = append(result, Empty{})
		}
	}

	if n := len(result); n > 1 {
		switch p.t.Assoc {
		case LeftToRight:
			v := result[0]
			for i := 1; i < n; i += 2 {
				p.put(&v, Associativity(i/2), v, result[i], result[i+1])
			}
			*output.(*interface{}) = v
		case RightToLeft:
			v := result[n-1]
			for i := 1; i < n; i += 2 {
				j := n - 1 - i
				p.put(&v, Associativity(-j/2), result[j-1], result[j], v)
			}
			*output.(*interface{}) = v
		}
	}

	return p.put(output, Associativity(0), result...)
}

func (t Delim) Parser(rule Rule, c cache) parse.Parser {
	p := &delimParser{
		rule: rule,
		t:    t,
		term: t.Term.Parser("", c),
		sep:  t.Sep.Parser("", c),
		put:  tag(rule, delimTag),
	}
	c.registerRule(&p.term)
	c.registerRule(&p.sep)
	return p
}

func (t Delim) LRTerms(node parse.Node) (left, right Term) {
	associativity := node.Extra.(Associativity)
	switch {
	case associativity < 0:
		return t.Term, t
	case associativity > 0:
		return t, t.Term
	}
	return t.Term, t.Term
}

//-----------------------------------------------------------------------------

type quantParser struct {
	rule Rule
	t    Quant
	term parse.Parser
	put  putter
}

func (p *quantParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	defer enterf("%s: %T %[2]v", p.rule, p.t).exitf("%v %v", &out, output)
	result := make([]interface{}, 0, p.t.Min)
	var v interface{}
	start := *input
	for i := 0; (p.t.Max == 0 || i < p.t.Max) && p.term.Parse(&start, &v); i++ {
		result = append(result, v)
		*input = start
	}
	if len(result) >= p.t.Min {
		return p.put(output, nil, result...)
	}
	return false
}

func (t Quant) Parser(rule Rule, c cache) parse.Parser {
	p := &quantParser{
		rule: rule,
		t:    t,
		term: t.Term.Parser("", c),
		put:  tag(rule, quantTag),
	}
	c.registerRule(&p.term)
	return p
}

//-----------------------------------------------------------------------------

type oneofParser struct {
	rule    Rule
	t       Oneof
	parsers []parse.Parser
	put     putter
}

func (p *oneofParser) Parse(input *parse.Scanner, output interface{}) (out bool) {
	defer enterf("%s: %T %[2]v", p.rule, p.t).exitf("%v %v", &out, output)
	furthest := *input
	for i, parser := range p.parsers {
		var v interface{}
		start := *input
		if parser.Parse(&start, &v) {
			*input = start
			return p.put(output, i, v)
		}
		if furthest.Offset() < start.Offset() {
			furthest = start
		}
	}
	*input = furthest
	return false
}

func (t Oneof) Parser(rule Rule, c cache) parse.Parser {
	return &oneofParser{
		rule:    rule,
		t:       t,
		parsers: c.makeParsers(t),
		put:     tag(rule, oneofTag),
	}
}

//-----------------------------------------------------------------------------

func (t Stack) Parser(_ Rule, _ cache) parse.Parser {
	panic(Inconceivable)
}

//-----------------------------------------------------------------------------

func (t Named) Parser(rule Rule, c cache) parse.Parser {
	return t.Term.Parser(Rule(t.Name), c)
}
