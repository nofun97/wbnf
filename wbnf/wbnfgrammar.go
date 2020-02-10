// Code generated by "ωBNF gen" DO NOT EDIT.
// $ wbnf gen --grammar ../examples/wbnf.wbnf --start grammar --pkg wbnf
package wbnf

import (
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
)

func Grammar() parser.Parsers {
	return parser.Grammar{"grammar": parser.Some(parser.Rule(`stmt`)),
		"stmt": parser.Oneof{parser.Rule(`COMMENT`),
			parser.Rule(`prod`)},
		"prod": parser.Seq{parser.Rule(`IDENT`),
			parser.S("->"),
			parser.Some(parser.Rule(`term`)),
			parser.S(";")},
		"term": parser.Stack{parser.Delim{Term: parser.Rule(`@`),
			Sep: parser.Eq("op",
				parser.S(">")),
			Assoc: parser.NonAssociative},
			parser.Delim{Term: parser.Rule(`@`),
				Sep: parser.Eq("op",
					parser.S("|")),
				Assoc: parser.NonAssociative},
			parser.Some(parser.Rule(`@`)),
			parser.Seq{parser.Rule(`named`),
				parser.Any(parser.Rule(`quant`))}},
		"named": parser.Seq{parser.Opt(parser.Seq{parser.Rule(`IDENT`),
			parser.Eq("op",
				parser.S("="))}),
			parser.Rule(`atom`)},
		"quant": parser.Oneof{parser.Eq("op",
			parser.RE(`[?*+]`)),
			parser.Seq{parser.S("{"),
				parser.Opt(parser.Eq("min",
					parser.Rule(`INT`))),
				parser.S(","),
				parser.Opt(parser.Eq("max",
					parser.Rule(`INT`))),
				parser.S("}")},
			parser.Seq{parser.Eq("op",
				parser.RE(`<:|:>?`)),
				parser.Opt(parser.Eq("opt_leading",
					parser.S(","))),
				parser.Rule(`named`),
				parser.Opt(parser.Eq("opt_trailing",
					parser.S(",")))}},
		"atom": parser.Oneof{parser.Rule(`IDENT`),
			parser.Rule(`STR`),
			parser.Rule(`RE`),
			parser.Rule(`REF`),
			parser.Seq{parser.S("("),
				parser.Rule(`term`),
				parser.S(")")},
			parser.Seq{parser.S("("),
				parser.S(")")}},
		"COMMENT": parser.RE(`//.*$|(?s:/\*(?:[^*]|\*+[^*/])\*/)`),
		"IDENT":   parser.RE(`@|[A-Za-z_\.]\w*`),
		"INT":     parser.RE(`\d+`),
		"STR":     parser.RE(`"(?:\\.|[^\\"])*"|'(?:\\.|[^\\'])*'|` + "`" + `(?:` + "`" + `` + "`" + `|[^` + "`" + `])*` + "`" + ``),
		"RE":      parser.RE(`/{(?:\\.|{(?:(?:\d+(?:,\d*)?|,\d+)\})?|\[(?:\\.|\[:^?[a-z]+:\]|[^\]])+]|[^\\{\}])*\}|(?:(?:\[(?:\\.|\[:^?[a-z]+:\]|[^\]])+]|\\[pP](?:[a-z]|\{[a-zA-Z_]+\})|\\[a-zA-Z]|[.^$])(?:(?:[+*?]|\{\d+,?\d?\})\??)?)+`),
		"REF": parser.Seq{parser.S("%"),
			parser.Rule(`IDENT`),
			parser.Opt(parser.Seq{parser.S("="),
				parser.Eq("default",
					parser.Rule(`STR`))})},
		".wrapRE": parser.RE(`\s*()\s*`)}.Compile(nil)
}

type WrapreNode struct{ ast.Node }

func (c WrapreNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkWrapreNode(node WrapreNode, ops WalkerOps) {
	if fn := ops.EnterWrapreNode; fn != nil {
		fn(node)
	}

	if fn := ops.ExitWrapreNode; fn != nil {
		fn(node)
	}
}

type CommentNode struct{ ast.Node }

func (c CommentNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkCommentNode(node CommentNode, ops WalkerOps) {
	if fn := ops.EnterCommentNode; fn != nil {
		fn(node)
	}

	if fn := ops.ExitCommentNode; fn != nil {
		fn(node)
	}
}

type IdentNode struct{ ast.Node }

func (c IdentNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkIdentNode(node IdentNode, ops WalkerOps) {
	if fn := ops.EnterIdentNode; fn != nil {
		fn(node)
	}

	if fn := ops.ExitIdentNode; fn != nil {
		fn(node)
	}
}

type IntNode struct{ ast.Node }

func (c IntNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkIntNode(node IntNode, ops WalkerOps) {
	if fn := ops.EnterIntNode; fn != nil {
		fn(node)
	}

	if fn := ops.ExitIntNode; fn != nil {
		fn(node)
	}
}

type ReNode struct{ ast.Node }

func (c ReNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkReNode(node ReNode, ops WalkerOps) {
	if fn := ops.EnterReNode; fn != nil {
		fn(node)
	}

	if fn := ops.ExitReNode; fn != nil {
		fn(node)
	}
}

type RefNode struct{ ast.Node }

func (c RefNode) AllIdent() []IdentNode {
	var out []IdentNode
	for _, child := range ast.All(c.Node, "IDENT") {
		out = append(out, IdentNode{child})
	}
	return out
}

func (c RefNode) OneIdent() IdentNode {
	return IdentNode{ast.First(c.Node, "IDENT")}
}

func (c RefNode) AllDefault() []StrNode {
	var out []StrNode
	for _, child := range ast.All(c.Node, "default") {
		out = append(out, StrNode{child})
	}
	return out
}

func (c RefNode) OneDefault() StrNode {
	return StrNode{ast.First(c.Node, "default")}
}
func WalkRefNode(node RefNode, ops WalkerOps) {
	if fn := ops.EnterRefNode; fn != nil {
		fn(node)
	}
	for _, child := range node.AllIdent() {
		WalkIdentNode(child, ops)
	}

	if fn := ops.ExitRefNode; fn != nil {
		fn(node)
	}
}

type StrNode struct{ ast.Node }

func (c StrNode) String() string {
	if c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}
func WalkStrNode(node StrNode, ops WalkerOps) {
	if fn := ops.EnterStrNode; fn != nil {
		fn(node)
	}

	if fn := ops.ExitStrNode; fn != nil {
		fn(node)
	}
}

type AtomNode struct{ ast.Node }

func (c AtomNode) Choice() int {
	return ast.Choice(c.Node)
}

func (c AtomNode) AllIdent() []IdentNode {
	var out []IdentNode
	for _, child := range ast.All(c.Node, "IDENT") {
		out = append(out, IdentNode{child})
	}
	return out
}

func (c AtomNode) OneIdent() IdentNode {
	return IdentNode{ast.First(c.Node, "IDENT")}
}

func (c AtomNode) AllRe() []ReNode {
	var out []ReNode
	for _, child := range ast.All(c.Node, "RE") {
		out = append(out, ReNode{child})
	}
	return out
}

func (c AtomNode) OneRe() ReNode {
	return ReNode{ast.First(c.Node, "RE")}
}

func (c AtomNode) AllRef() []RefNode {
	var out []RefNode
	for _, child := range ast.All(c.Node, "REF") {
		out = append(out, RefNode{child})
	}
	return out
}

func (c AtomNode) OneRef() RefNode {
	return RefNode{ast.First(c.Node, "REF")}
}

func (c AtomNode) AllStr() []StrNode {
	var out []StrNode
	for _, child := range ast.All(c.Node, "STR") {
		out = append(out, StrNode{child})
	}
	return out
}

func (c AtomNode) OneStr() StrNode {
	return StrNode{ast.First(c.Node, "STR")}
}

func (c AtomNode) AllTerm() []TermNode {
	var out []TermNode
	for _, child := range ast.All(c.Node, "term") {
		out = append(out, TermNode{child})
	}
	return out
}

func (c AtomNode) OneTerm() TermNode {
	return TermNode{ast.First(c.Node, "term")}
}
func WalkAtomNode(node AtomNode, ops WalkerOps) {
	if fn := ops.EnterAtomNode; fn != nil {
		fn(node)
	}
	for _, child := range node.AllIdent() {
		WalkIdentNode(child, ops)
	}
	for _, child := range node.AllRe() {
		WalkReNode(child, ops)
	}
	for _, child := range node.AllRef() {
		WalkRefNode(child, ops)
	}
	for _, child := range node.AllStr() {
		WalkStrNode(child, ops)
	}
	for _, child := range node.AllTerm() {
		WalkTermNode(child, ops)
	}

	if fn := ops.ExitAtomNode; fn != nil {
		fn(node)
	}
}

type GrammarNode struct{ ast.Node }

func (c GrammarNode) AllStmt() []StmtNode {
	var out []StmtNode
	for _, child := range ast.All(c.Node, "stmt") {
		out = append(out, StmtNode{child})
	}
	return out
}

func (c GrammarNode) OneStmt() StmtNode {
	return StmtNode{ast.First(c.Node, "stmt")}
}
func WalkGrammarNode(node GrammarNode, ops WalkerOps) {
	if fn := ops.EnterGrammarNode; fn != nil {
		fn(node)
	}
	for _, child := range node.AllStmt() {
		WalkStmtNode(child, ops)
	}

	if fn := ops.ExitGrammarNode; fn != nil {
		fn(node)
	}
}

type NamedNode struct{ ast.Node }

func (c NamedNode) AllIdent() []IdentNode {
	var out []IdentNode
	for _, child := range ast.All(c.Node, "IDENT") {
		out = append(out, IdentNode{child})
	}
	return out
}

func (c NamedNode) OneIdent() IdentNode {
	return IdentNode{ast.First(c.Node, "IDENT")}
}

func (c NamedNode) AllAtom() []AtomNode {
	var out []AtomNode
	for _, child := range ast.All(c.Node, "atom") {
		out = append(out, AtomNode{child})
	}
	return out
}

func (c NamedNode) OneAtom() AtomNode {
	return AtomNode{ast.First(c.Node, "atom")}
}

func (c NamedNode) AllOp() []string {
	var out []string
	for _, child := range ast.All(c.Node, "op") {
		out = append(out, child.Scanner().String())
	}
	return out
}

func (c NamedNode) OneOp() string {
	if child := ast.First(c.Node, "op"); child != nil {
		return child.Scanner().String()
	}
	return ""
}
func WalkNamedNode(node NamedNode, ops WalkerOps) {
	if fn := ops.EnterNamedNode; fn != nil {
		fn(node)
	}
	for _, child := range node.AllIdent() {
		WalkIdentNode(child, ops)
	}
	for _, child := range node.AllAtom() {
		WalkAtomNode(child, ops)
	}

	if fn := ops.ExitNamedNode; fn != nil {
		fn(node)
	}
}

type ProdNode struct{ ast.Node }

func (c ProdNode) AllIdent() []IdentNode {
	var out []IdentNode
	for _, child := range ast.All(c.Node, "IDENT") {
		out = append(out, IdentNode{child})
	}
	return out
}

func (c ProdNode) OneIdent() IdentNode {
	return IdentNode{ast.First(c.Node, "IDENT")}
}

func (c ProdNode) AllTerm() []TermNode {
	var out []TermNode
	for _, child := range ast.All(c.Node, "term") {
		out = append(out, TermNode{child})
	}
	return out
}

func (c ProdNode) OneTerm() TermNode {
	return TermNode{ast.First(c.Node, "term")}
}
func WalkProdNode(node ProdNode, ops WalkerOps) {
	if fn := ops.EnterProdNode; fn != nil {
		fn(node)
	}
	for _, child := range node.AllIdent() {
		WalkIdentNode(child, ops)
	}
	for _, child := range node.AllTerm() {
		WalkTermNode(child, ops)
	}

	if fn := ops.ExitProdNode; fn != nil {
		fn(node)
	}
}

type QuantNode struct{ ast.Node }

func (c QuantNode) Choice() int {
	return ast.Choice(c.Node)
}

func (c QuantNode) AllMax() []IntNode {
	var out []IntNode
	for _, child := range ast.All(c.Node, "max") {
		out = append(out, IntNode{child})
	}
	return out
}

func (c QuantNode) OneMax() IntNode {
	return IntNode{ast.First(c.Node, "max")}
}

func (c QuantNode) AllMin() []IntNode {
	var out []IntNode
	for _, child := range ast.All(c.Node, "min") {
		out = append(out, IntNode{child})
	}
	return out
}

func (c QuantNode) OneMin() IntNode {
	return IntNode{ast.First(c.Node, "min")}
}

func (c QuantNode) AllNamed() []NamedNode {
	var out []NamedNode
	for _, child := range ast.All(c.Node, "named") {
		out = append(out, NamedNode{child})
	}
	return out
}

func (c QuantNode) OneNamed() NamedNode {
	return NamedNode{ast.First(c.Node, "named")}
}

func (c QuantNode) AllOp() []string {
	var out []string
	for _, child := range ast.All(c.Node, "op") {
		out = append(out, child.Scanner().String())
	}
	return out
}

func (c QuantNode) OneOp() string {
	if child := ast.First(c.Node, "op"); child != nil {
		return child.Scanner().String()
	}
	return ""
}

func (c QuantNode) AllOptLeading() []string {
	var out []string
	for _, child := range ast.All(c.Node, "opt_leading") {
		out = append(out, child.Scanner().String())
	}
	return out
}

func (c QuantNode) OneOptLeading() string {
	if child := ast.First(c.Node, "opt_leading"); child != nil {
		return child.Scanner().String()
	}
	return ""
}

func (c QuantNode) AllOptTrailing() []string {
	var out []string
	for _, child := range ast.All(c.Node, "opt_trailing") {
		out = append(out, child.Scanner().String())
	}
	return out
}

func (c QuantNode) OneOptTrailing() string {
	if child := ast.First(c.Node, "opt_trailing"); child != nil {
		return child.Scanner().String()
	}
	return ""
}
func WalkQuantNode(node QuantNode, ops WalkerOps) {
	if fn := ops.EnterQuantNode; fn != nil {
		fn(node)
	}
	for _, child := range node.AllNamed() {
		WalkNamedNode(child, ops)
	}

	if fn := ops.ExitQuantNode; fn != nil {
		fn(node)
	}
}

type StmtNode struct{ ast.Node }

func (c StmtNode) Choice() int {
	return ast.Choice(c.Node)
}

func (c StmtNode) AllComment() []CommentNode {
	var out []CommentNode
	for _, child := range ast.All(c.Node, "COMMENT") {
		out = append(out, CommentNode{child})
	}
	return out
}

func (c StmtNode) OneComment() CommentNode {
	return CommentNode{ast.First(c.Node, "COMMENT")}
}

func (c StmtNode) AllProd() []ProdNode {
	var out []ProdNode
	for _, child := range ast.All(c.Node, "prod") {
		out = append(out, ProdNode{child})
	}
	return out
}

func (c StmtNode) OneProd() ProdNode {
	return ProdNode{ast.First(c.Node, "prod")}
}
func WalkStmtNode(node StmtNode, ops WalkerOps) {
	if fn := ops.EnterStmtNode; fn != nil {
		fn(node)
	}
	for _, child := range node.AllComment() {
		WalkCommentNode(child, ops)
	}
	for _, child := range node.AllProd() {
		WalkProdNode(child, ops)
	}

	if fn := ops.ExitStmtNode; fn != nil {
		fn(node)
	}
}

type TermNode struct{ ast.Node }

func (c TermNode) AllTerm() []TermNode {
	var out []TermNode
	for _, child := range ast.All(c.Node, "term") {
		out = append(out, TermNode{child})
	}
	return out
}

func (c TermNode) OneTerm() TermNode {
	return TermNode{ast.First(c.Node, "term")}
}

func (c TermNode) AllNamed() []NamedNode {
	var out []NamedNode
	for _, child := range ast.All(c.Node, "named") {
		out = append(out, NamedNode{child})
	}
	return out
}

func (c TermNode) OneNamed() NamedNode {
	return NamedNode{ast.First(c.Node, "named")}
}

func (c TermNode) AllOp() []string {
	var out []string
	for _, child := range ast.All(c.Node, "op") {
		out = append(out, child.Scanner().String())
	}
	return out
}

func (c TermNode) OneOp() string {
	if child := ast.First(c.Node, "op"); child != nil {
		return child.Scanner().String()
	}
	return ""
}

func (c TermNode) AllQuant() []QuantNode {
	var out []QuantNode
	for _, child := range ast.All(c.Node, "quant") {
		out = append(out, QuantNode{child})
	}
	return out
}

func (c TermNode) OneQuant() QuantNode {
	return QuantNode{ast.First(c.Node, "quant")}
}
func WalkTermNode(node TermNode, ops WalkerOps) {
	if fn := ops.EnterTermNode; fn != nil {
		fn(node)
	}
	for _, child := range node.AllTerm() {
		WalkTermNode(child, ops)
	}
	for _, child := range node.AllNamed() {
		WalkNamedNode(child, ops)
	}
	for _, child := range node.AllQuant() {
		WalkQuantNode(child, ops)
	}

	if fn := ops.ExitTermNode; fn != nil {
		fn(node)
	}
}

type WalkerOps struct {
	EnterWrapreNode  func(WrapreNode)
	ExitWrapreNode   func(WrapreNode)
	EnterCommentNode func(CommentNode)
	ExitCommentNode  func(CommentNode)
	EnterIdentNode   func(IdentNode)
	ExitIdentNode    func(IdentNode)
	EnterIntNode     func(IntNode)
	ExitIntNode      func(IntNode)
	EnterReNode      func(ReNode)
	ExitReNode       func(ReNode)
	EnterRefNode     func(RefNode)
	ExitRefNode      func(RefNode)
	EnterStrNode     func(StrNode)
	ExitStrNode      func(StrNode)
	EnterAtomNode    func(AtomNode)
	ExitAtomNode     func(AtomNode)
	EnterGrammarNode func(GrammarNode)
	ExitGrammarNode  func(GrammarNode)
	EnterNamedNode   func(NamedNode)
	ExitNamedNode    func(NamedNode)
	EnterProdNode    func(ProdNode)
	ExitProdNode     func(ProdNode)
	EnterQuantNode   func(QuantNode)
	ExitQuantNode    func(QuantNode)
	EnterStmtNode    func(StmtNode)
	ExitStmtNode     func(StmtNode)
	EnterTermNode    func(TermNode)
	ExitTermNode     func(TermNode)
}

func (w WalkerOps) Walk(tree GrammarNode)  { WalkGrammarNode(tree, w) }
func (c GrammarNode) GetAstNode() ast.Node { return c.Node }

func NewGrammarNode(from ast.Node) GrammarNode { return GrammarNode{from} }

func Parse(input *parser.Scanner) (GrammarNode, error) {
	p := Grammar()
	tree, err := p.Parse("grammar", input)
	if err != nil {
		return GrammarNode{nil}, err
	}
	return GrammarNode{ast.FromParserNode(p.Grammar(), tree)}, nil
}

func ParseString(input string) (GrammarNode, error) {
	return Parse(parser.NewScanner(input))
}

var grammarGrammarSrc = unfakeBackquote(`
// Non-terminals
grammar -> stmt+;
stmt    -> COMMENT | prod;
prod    -> IDENT "->" term+ ";";
term    -> @:op=">"
         > @:op="|"
         > @+
         > named quant*;
named   -> (IDENT op="=")? atom;
quant   -> op=[?*+]
         | "{" min=INT? "," max=INT? "}"
         | op=/{<:|:>?} opt_leading=","? named opt_trailing=","?;
atom    -> IDENT | STR | RE | REF | "(" term ")" | "(" ")";

// Terminals
COMMENT -> /{ //.*$
            | (?s: /\* (?: [^*] | \*+[^*/] ) \*/ )
            };
IDENT   -> /{@|[A-Za-z_\.]\w*};
INT     -> \d+;
STR     -> /{ " (?: \\. | [^\\"] )* "
            | ' (?: \\. | [^\\'] )* '
            | ‵ (?: ‵‵  | [^‵]   )* ‵
            };
RE      -> /{
             /{
               (?:
                 \\.
                 | { (?: (?: \d+(?:,\d*)? | ,\d+ ) \} )?
                 | \[ (?: \\. | \[:^?[a-z]+:\] | [^\]] )+ ]
                 | [^\\{\}]
               )*
             \}
           | (?:
               (?:
                 \[ (?: \\. | \[:^?[a-z]+:\] | [^\]] )+ ]
               | \\[pP](?:[a-z]|\{[a-zA-Z_]+\})
               | \\[a-zA-Z]
               | [.^$]
               )(?: (?:[+*?]|\{\d+,?\d?\}) \?? )?
             )+
           };
REF     -> "%" IDENT ("=" default=(STR|RE))?;

// Special
.wrapRE -> /{\s*()\s*};
`)
