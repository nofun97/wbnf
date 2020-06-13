// Code generated by "ωBNF gen" DO NOT EDIT.
// $ wbnf gen --grammar ../examples/wbnf.wbnf --start grammar --pkg wbnf
package wbnf

import (
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
)

func Grammar() parser.Parsers {
	return parser.Grammar{".wrapRE": parser.RE(`\s*()\s*`),
		"COMMENT": parser.RE(`//.*$|(?s:/\*(?:[^*]|\*+[^*/])\*/)`),
		"IDENT":   parser.RE(`@|[A-Za-z_\.]\w*`),
		"INT":     parser.RE(`\d+`),
		"RE":      parser.RE(`/{(?:\\.|{(?:(?:\d+(?:,\d*)?|,\d+)\})?|\[(?:\\.|\[:^?[a-z]+:\]|[^\]])+]|[^\\{\}])*\}|(?:(?:\[(?:\\.|\[:^?[a-z]+:\]|[^\]])+]|\\[pP](?:[a-z]|\{[a-zA-Z_]+\})|\\[a-zA-Z]|[.^$])(?:(?:[+*?]|\{\d+,?\d?\})\??)?)+`),
		"REF": parser.Seq{parser.CutPoint{Term: parser.S(`%`)},
			parser.Rule(`IDENT`),
			parser.Opt(parser.Seq{parser.S(`=`),
				parser.Eq(`default`,
					parser.Rule(`STR`))})},
		"STR": parser.RE(`"(?:\\.|[^\\"])*"|'(?:\\.|[^\\'])*'|` + "`" + `(?:` + "`" + `` + "`" + `|[^` + "`" + `])*` + "`" + ``),
		"atom": parser.Oneof{parser.Rule(`IDENT`),
			parser.Rule(`STR`),
			parser.Rule(`RE`),
			parser.Rule(`macrocall`),
			parser.Eq(`ExtRef`,
				parser.Seq{parser.CutPoint{Term: parser.S(`%%`)},
					parser.Rule(`IDENT`)}),
			parser.Rule(`REF`),
			parser.Seq{parser.S(`(`),
				parser.Rule(`term`),
				parser.S(`)`)},
			parser.Seq{parser.S(`(`),
				parser.S(`)`)}},
		"grammar": parser.Some(parser.Rule(`stmt`)),
		"macrocall": parser.Seq{parser.CutPoint{Term: parser.S(`%!`)},
			parser.Eq(`name`,
				parser.Rule(`IDENT`)),
			parser.S(`(`),
			parser.Delim{Term: parser.Opt(parser.Rule(`term`)),
				Sep: parser.S(`,`)},
			parser.S(`)`)},
		"named": parser.Seq{parser.Opt(parser.Seq{parser.Rule(`IDENT`),
			parser.Eq(`op`,
				parser.S(`=`))}),
			parser.Rule(`atom`)},
		"pragma": parser.ScopedGrammar{Term: parser.Oneof{parser.Rule(`import`),
			parser.Rule(`macrodef`)},
			Grammar: parser.Grammar{".wrapRE": parser.RE(`\s*()\s*`),
				"import": parser.Seq{parser.CutPoint{Term: parser.S(`.import`)},
					parser.Eq(`path`,
						parser.Delim{Term: parser.Oneof{parser.CutPoint{Term: parser.S(`..`)},
							parser.CutPoint{Term: parser.S(`.`)},
							parser.RE(`[a-zA-Z0-9.:]+`)},
							Sep:             parser.S(`/`),
							CanStartWithSep: true}),
					parser.Opt(parser.CutPoint{Term: parser.S(`;`)})},
				"macrodef": parser.Seq{parser.CutPoint{Term: parser.S(`.macro`)},
					parser.Eq(`name`,
						parser.Rule(`IDENT`)),
					parser.S(`(`),
					parser.Delim{Term: parser.Opt(parser.Eq(`args`,
						parser.Rule(`IDENT`))),
						Sep: parser.S(`,`)},
					parser.S(`)`),
					parser.S(`{`),
					parser.Rule(`term`),
					parser.S(`}`),
					parser.Opt(parser.CutPoint{Term: parser.S(`;`)})}}},
		"prod": parser.Seq{parser.Rule(`IDENT`),
			parser.CutPoint{Term: parser.S(`->`)},
			parser.Some(parser.Rule(`term`)),
			parser.CutPoint{Term: parser.S(`;`)}},
		"quant": parser.Oneof{parser.Eq(`op`,
			parser.RE(`[?*+]`)),
			parser.Seq{parser.S(`{`),
				parser.Opt(parser.Eq(`min`,
					parser.Rule(`INT`))),
				parser.S(`,`),
				parser.Opt(parser.Eq(`max`,
					parser.Rule(`INT`))),
				parser.S(`}`)},
			parser.Seq{parser.Eq(`op`,
				parser.RE(`<:|:>?`)),
				parser.Opt(parser.Eq(`opt_leading`,
					parser.S(`,`))),
				parser.Rule(`named`),
				parser.Opt(parser.Eq(`opt_trailing`,
					parser.S(`,`)))}},
		"stmt": parser.Oneof{parser.Rule(`COMMENT`),
			parser.Rule(`prod`),
			parser.Rule(`pragma`)},
		"term": parser.Stack{parser.Delim{Term: parser.Seq{parser.At,
			parser.Opt(parser.Seq{parser.S(`{`),
				parser.Rule(`grammar`),
				parser.S(`}`)})},
			Sep: parser.Eq(`op`,
				parser.S(`>`))},
			parser.Delim{Term: parser.At,
				Sep: parser.Eq(`op`,
					parser.S(`|`))},
			parser.Some(parser.At),
			parser.Seq{parser.Rule(`named`),
				parser.Any(parser.Rule(`quant`))}}}.Compile(nil)
}

type Stopper interface {
	ExitNode() bool
	Abort() bool
}
type nodeExiter struct{}

func (n *nodeExiter) ExitNode() bool { return true }
func (n *nodeExiter) Abort() bool    { return false }

type aborter struct{}

func (n *aborter) ExitNode() bool { return true }
func (n *aborter) Abort() bool    { return true }

var (
	NodeExiter = &nodeExiter{}
	Aborter    = &aborter{}
)

type IsWalkableType interface{ isWalkableType() }

type AtomExtRefNode struct{ ast.Node }

func (AtomExtRefNode) isWalkableType() {}

func (c AtomExtRefNode) OneIdent() *IdentNode {
	if child := ast.First(c.Node, "IDENT"); child != nil {
		return &IdentNode{child}
	}
	return nil
}

func (c AtomExtRefNode) OneToken() string {
	if child := ast.First(c.Node, ""); child != nil {
		return child.Scanner().String()
	}
	if b, ok := c.Node.(ast.Branch); ok && len(b) == 1 {
		for _, c := range b {
			if child := ast.First(c.(ast.One).Node, ""); child != nil {
				return child.Scanner().String()
			}
		}
	}
	return ""
}

type AtomNode struct{ ast.Node }

func (AtomNode) isWalkableType() {}
func (c AtomNode) Choice() int   { return ast.Choice(c.Node) }

func (c AtomNode) OneExtRef() *AtomExtRefNode {
	if child := ast.First(c.Node, "ExtRef"); child != nil {
		return &AtomExtRefNode{child}
	}
	return nil
}

func (c AtomNode) OneIdent() *IdentNode {
	if child := ast.First(c.Node, "IDENT"); child != nil {
		return &IdentNode{child}
	}
	return nil
}

func (c AtomNode) OneMacrocall() *MacrocallNode {
	if child := ast.First(c.Node, "macrocall"); child != nil {
		return &MacrocallNode{child}
	}
	return nil
}

func (c AtomNode) OneRe() *ReNode {
	if child := ast.First(c.Node, "RE"); child != nil {
		return &ReNode{child}
	}
	return nil
}

func (c AtomNode) OneRef() *RefNode {
	if child := ast.First(c.Node, "REF"); child != nil {
		return &RefNode{child}
	}
	return nil
}

func (c AtomNode) OneStr() *StrNode {
	if child := ast.First(c.Node, "STR"); child != nil {
		return &StrNode{child}
	}
	return nil
}

func (c AtomNode) OneTerm() *TermNode {
	if child := ast.First(c.Node, "term"); child != nil {
		return &TermNode{child}
	}
	return nil
}

func (c AtomNode) OneToken() string {
	if child := ast.First(c.Node, ""); child != nil {
		return child.Scanner().String()
	}
	if b, ok := c.Node.(ast.Branch); ok && len(b) == 1 {
		for _, c := range b {
			if child := ast.First(c.(ast.One).Node, ""); child != nil {
				return child.Scanner().String()
			}
		}
	}
	return ""
}

func (c AtomNode) AllToken() []string {
	var out []string
	for _, child := range ast.All(c.Node, "") {
		out = append(out, child.Scanner().String())
	}
	return out
}

type CommentNode struct{ ast.Node }

func (CommentNode) isWalkableType() {}
func (c *CommentNode) String() string {
	if c == nil || c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}

type GrammarNode struct{ ast.Node }

func (GrammarNode) isWalkableType() {}
func (c GrammarNode) AllStmt() []StmtNode {
	var out []StmtNode
	for _, child := range ast.All(c.Node, "stmt") {
		out = append(out, StmtNode{child})
	}
	return out
}

type IdentNode struct{ ast.Node }

func (IdentNode) isWalkableType() {}
func (c *IdentNode) String() string {
	if c == nil || c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}

type IntNode struct{ ast.Node }

func (IntNode) isWalkableType() {}
func (c *IntNode) String() string {
	if c == nil || c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}

type MacrocallNode struct{ ast.Node }

func (MacrocallNode) isWalkableType() {}

func (c MacrocallNode) OneName() *IdentNode {
	if child := ast.First(c.Node, "name"); child != nil {
		return &IdentNode{child}
	}
	return nil
}

func (c MacrocallNode) AllTerm() []TermNode {
	var out []TermNode
	for _, child := range ast.All(c.Node, "term") {
		out = append(out, TermNode{child})
	}
	return out
}

func (c MacrocallNode) OneToken() string {
	if child := ast.First(c.Node, ""); child != nil {
		return child.Scanner().String()
	}
	if b, ok := c.Node.(ast.Branch); ok && len(b) == 1 {
		for _, c := range b {
			if child := ast.First(c.(ast.One).Node, ""); child != nil {
				return child.Scanner().String()
			}
		}
	}
	return ""
}

func (c MacrocallNode) AllToken() []string {
	var out []string
	for _, child := range ast.All(c.Node, "") {
		out = append(out, child.Scanner().String())
	}
	return out
}

type NamedNode struct{ ast.Node }

func (NamedNode) isWalkableType() {}

func (c NamedNode) OneAtom() *AtomNode {
	if child := ast.First(c.Node, "atom"); child != nil {
		return &AtomNode{child}
	}
	return nil
}

func (c NamedNode) OneIdent() *IdentNode {
	if child := ast.First(c.Node, "IDENT"); child != nil {
		return &IdentNode{child}
	}
	return nil
}

func (c NamedNode) OneOp() string {
	if child := ast.First(c.Node, "op"); child != nil {
		return ast.First(child, "").Scanner().String()
	}
	return ""
}

type PragmaImportNode struct{ ast.Node }

func (PragmaImportNode) isWalkableType() {}

func (c PragmaImportNode) OnePath() *PragmaImportPathNode {
	if child := ast.First(c.Node, "path"); child != nil {
		return &PragmaImportPathNode{child}
	}
	return nil
}

func (c PragmaImportNode) OneToken() string {
	if child := ast.First(c.Node, ""); child != nil {
		return child.Scanner().String()
	}
	if b, ok := c.Node.(ast.Branch); ok && len(b) == 1 {
		for _, c := range b {
			if child := ast.First(c.(ast.One).Node, ""); child != nil {
				return child.Scanner().String()
			}
		}
	}
	return ""
}

type PragmaImportPathNode struct{ ast.Node }

func (PragmaImportPathNode) isWalkableType() {}
func (c PragmaImportPathNode) Choice() int   { return ast.Choice(c.Node) }

func (c PragmaImportPathNode) AllToken() []string {
	var out []string
	for _, child := range ast.All(c.Node, "") {
		out = append(out, child.Scanner().String())
	}
	return out
}

type PragmaMacrodefNode struct{ ast.Node }

func (PragmaMacrodefNode) isWalkableType() {}
func (c PragmaMacrodefNode) AllArgs() []IdentNode {
	var out []IdentNode
	for _, child := range ast.All(c.Node, "args") {
		out = append(out, IdentNode{child})
	}
	return out
}

func (c PragmaMacrodefNode) OneName() *IdentNode {
	if child := ast.First(c.Node, "name"); child != nil {
		return &IdentNode{child}
	}
	return nil
}

func (c PragmaMacrodefNode) OneTerm() *TermNode {
	if child := ast.First(c.Node, "term"); child != nil {
		return &TermNode{child}
	}
	return nil
}

func (c PragmaMacrodefNode) OneToken() string {
	if child := ast.First(c.Node, ""); child != nil {
		return child.Scanner().String()
	}
	if b, ok := c.Node.(ast.Branch); ok && len(b) == 1 {
		for _, c := range b {
			if child := ast.First(c.(ast.One).Node, ""); child != nil {
				return child.Scanner().String()
			}
		}
	}
	return ""
}

func (c PragmaMacrodefNode) AllToken() []string {
	var out []string
	for _, child := range ast.All(c.Node, "") {
		out = append(out, child.Scanner().String())
	}
	return out
}

type PragmaNode struct{ ast.Node }

func (PragmaNode) isWalkableType() {}
func (c PragmaNode) Choice() int   { return ast.Choice(c.Node) }

func (c PragmaNode) OneImport() *PragmaImportNode {
	if child := ast.First(c.Node, "import"); child != nil {
		return &PragmaImportNode{child}
	}
	return nil
}

func (c PragmaNode) OneMacrodef() *PragmaMacrodefNode {
	if child := ast.First(c.Node, "macrodef"); child != nil {
		return &PragmaMacrodefNode{child}
	}
	return nil
}

type ProdNode struct{ ast.Node }

func (ProdNode) isWalkableType() {}

func (c ProdNode) OneIdent() *IdentNode {
	if child := ast.First(c.Node, "IDENT"); child != nil {
		return &IdentNode{child}
	}
	return nil
}

func (c ProdNode) AllTerm() []TermNode {
	var out []TermNode
	for _, child := range ast.All(c.Node, "term") {
		out = append(out, TermNode{child})
	}
	return out
}

func (c ProdNode) OneToken() string {
	if child := ast.First(c.Node, ""); child != nil {
		return child.Scanner().String()
	}
	if b, ok := c.Node.(ast.Branch); ok && len(b) == 1 {
		for _, c := range b {
			if child := ast.First(c.(ast.One).Node, ""); child != nil {
				return child.Scanner().String()
			}
		}
	}
	return ""
}

func (c ProdNode) AllToken() []string {
	var out []string
	for _, child := range ast.All(c.Node, "") {
		out = append(out, child.Scanner().String())
	}
	return out
}

type QuantNode struct{ ast.Node }

func (QuantNode) isWalkableType() {}
func (c QuantNode) Choice() int   { return ast.Choice(c.Node) }

func (c QuantNode) OneMax() *IntNode {
	if child := ast.First(c.Node, "max"); child != nil {
		return &IntNode{child}
	}
	return nil
}

func (c QuantNode) OneMin() *IntNode {
	if child := ast.First(c.Node, "min"); child != nil {
		return &IntNode{child}
	}
	return nil
}

func (c QuantNode) OneNamed() *NamedNode {
	if child := ast.First(c.Node, "named"); child != nil {
		return &NamedNode{child}
	}
	return nil
}

func (c QuantNode) OneOp() string {
	if child := ast.First(c.Node, "op"); child != nil {
		return ast.First(child, "").Scanner().String()
	}
	return ""
}

func (c QuantNode) OneOptLeading() string {
	if child := ast.First(c.Node, "opt_leading"); child != nil {
		return ast.First(child, "").Scanner().String()
	}
	return ""
}

func (c QuantNode) OneOptTrailing() string {
	if child := ast.First(c.Node, "opt_trailing"); child != nil {
		return ast.First(child, "").Scanner().String()
	}
	return ""
}

func (c QuantNode) OneToken() string {
	if child := ast.First(c.Node, ""); child != nil {
		return child.Scanner().String()
	}
	if b, ok := c.Node.(ast.Branch); ok && len(b) == 1 {
		for _, c := range b {
			if child := ast.First(c.(ast.One).Node, ""); child != nil {
				return child.Scanner().String()
			}
		}
	}
	return ""
}

func (c QuantNode) AllToken() []string {
	var out []string
	for _, child := range ast.All(c.Node, "") {
		out = append(out, child.Scanner().String())
	}
	return out
}

type ReNode struct{ ast.Node }

func (ReNode) isWalkableType() {}
func (c *ReNode) String() string {
	if c == nil || c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}

type RefNode struct{ ast.Node }

func (RefNode) isWalkableType() {}

func (c RefNode) OneDefault() *StrNode {
	if child := ast.First(c.Node, "default"); child != nil {
		return &StrNode{child}
	}
	return nil
}

func (c RefNode) OneIdent() *IdentNode {
	if child := ast.First(c.Node, "IDENT"); child != nil {
		return &IdentNode{child}
	}
	return nil
}

func (c RefNode) OneToken() string {
	if child := ast.First(c.Node, ""); child != nil {
		return child.Scanner().String()
	}
	if b, ok := c.Node.(ast.Branch); ok && len(b) == 1 {
		for _, c := range b {
			if child := ast.First(c.(ast.One).Node, ""); child != nil {
				return child.Scanner().String()
			}
		}
	}
	return ""
}

type StmtNode struct{ ast.Node }

func (StmtNode) isWalkableType() {}
func (c StmtNode) Choice() int   { return ast.Choice(c.Node) }

func (c StmtNode) OneComment() *CommentNode {
	if child := ast.First(c.Node, "COMMENT"); child != nil {
		return &CommentNode{child}
	}
	return nil
}

func (c StmtNode) OnePragma() *PragmaNode {
	if child := ast.First(c.Node, "pragma"); child != nil {
		return &PragmaNode{child}
	}
	return nil
}

func (c StmtNode) OneProd() *ProdNode {
	if child := ast.First(c.Node, "prod"); child != nil {
		return &ProdNode{child}
	}
	return nil
}

type StrNode struct{ ast.Node }

func (StrNode) isWalkableType() {}
func (c *StrNode) String() string {
	if c == nil || c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}

type TermNode struct{ ast.Node }

func (TermNode) isWalkableType() {}
func (c TermNode) AllGrammar() []GrammarNode {
	var out []GrammarNode
	for _, child := range ast.All(c.Node, "grammar") {
		out = append(out, GrammarNode{child})
	}
	return out
}

func (c TermNode) OneNamed() *NamedNode {
	if child := ast.First(c.Node, "named"); child != nil {
		return &NamedNode{child}
	}
	return nil
}

func (c TermNode) OneOp() string {
	if child := ast.First(c.Node, "op"); child != nil {
		return ast.First(child, "").Scanner().String()
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

func (c TermNode) AllTerm() []TermNode {
	var out []TermNode
	for _, child := range ast.All(c.Node, "term") {
		out = append(out, TermNode{child})
	}
	return out
}

func (c TermNode) AllToken() []string {
	var out []string
	for _, child := range ast.All(c.Node, "") {
		out = append(out, child.Scanner().String())
	}
	return out
}

type WrapReNode struct{ ast.Node }

func (WrapReNode) isWalkableType() {}
func (c *WrapReNode) String() string {
	if c == nil || c.Node == nil {
		return ""
	}
	return c.Node.Scanner().String()
}

type WalkerOps struct {
	EnterAtomExtRefNode       func(AtomExtRefNode) Stopper
	ExitAtomExtRefNode        func(AtomExtRefNode) Stopper
	EnterAtomNode             func(AtomNode) Stopper
	ExitAtomNode              func(AtomNode) Stopper
	EnterCommentNode          func(CommentNode) Stopper
	ExitCommentNode           func(CommentNode) Stopper
	EnterGrammarNode          func(GrammarNode) Stopper
	ExitGrammarNode           func(GrammarNode) Stopper
	EnterIdentNode            func(IdentNode) Stopper
	ExitIdentNode             func(IdentNode) Stopper
	EnterIntNode              func(IntNode) Stopper
	ExitIntNode               func(IntNode) Stopper
	EnterMacrocallNode        func(MacrocallNode) Stopper
	ExitMacrocallNode         func(MacrocallNode) Stopper
	EnterNamedNode            func(NamedNode) Stopper
	ExitNamedNode             func(NamedNode) Stopper
	EnterPragmaImportNode     func(PragmaImportNode) Stopper
	ExitPragmaImportNode      func(PragmaImportNode) Stopper
	EnterPragmaImportPathNode func(PragmaImportPathNode) Stopper
	ExitPragmaImportPathNode  func(PragmaImportPathNode) Stopper
	EnterPragmaMacrodefNode   func(PragmaMacrodefNode) Stopper
	ExitPragmaMacrodefNode    func(PragmaMacrodefNode) Stopper
	EnterPragmaNode           func(PragmaNode) Stopper
	ExitPragmaNode            func(PragmaNode) Stopper
	EnterProdNode             func(ProdNode) Stopper
	ExitProdNode              func(ProdNode) Stopper
	EnterQuantNode            func(QuantNode) Stopper
	ExitQuantNode             func(QuantNode) Stopper
	EnterReNode               func(ReNode) Stopper
	ExitReNode                func(ReNode) Stopper
	EnterRefNode              func(RefNode) Stopper
	ExitRefNode               func(RefNode) Stopper
	EnterStmtNode             func(StmtNode) Stopper
	ExitStmtNode              func(StmtNode) Stopper
	EnterStrNode              func(StrNode) Stopper
	ExitStrNode               func(StrNode) Stopper
	EnterTermNode             func(TermNode) Stopper
	ExitTermNode              func(TermNode) Stopper
	EnterWrapReNode           func(WrapReNode) Stopper
	ExitWrapReNode            func(WrapReNode) Stopper
}

func (w WalkerOps) Walk(tree IsWalkableType) Stopper {
	switch node := tree.(type) {
	case AtomExtRefNode:
		return w.WalkAtomExtRefNode(node)

	case AtomNode:
		return w.WalkAtomNode(node)

	case CommentNode:
		if fn := w.EnterCommentNode; fn != nil {
			return fn(node)
		}

	case GrammarNode:
		return w.WalkGrammarNode(node)

	case IdentNode:
		if fn := w.EnterIdentNode; fn != nil {
			return fn(node)
		}

	case IntNode:
		if fn := w.EnterIntNode; fn != nil {
			return fn(node)
		}

	case MacrocallNode:
		return w.WalkMacrocallNode(node)

	case NamedNode:
		return w.WalkNamedNode(node)

	case PragmaImportNode:
		return w.WalkPragmaImportNode(node)

	case PragmaImportPathNode:
		return w.WalkPragmaImportPathNode(node)

	case PragmaMacrodefNode:
		return w.WalkPragmaMacrodefNode(node)

	case PragmaNode:
		return w.WalkPragmaNode(node)

	case ProdNode:
		return w.WalkProdNode(node)

	case QuantNode:
		return w.WalkQuantNode(node)

	case ReNode:
		if fn := w.EnterReNode; fn != nil {
			return fn(node)
		}

	case RefNode:
		return w.WalkRefNode(node)

	case StmtNode:
		return w.WalkStmtNode(node)

	case StrNode:
		if fn := w.EnterStrNode; fn != nil {
			return fn(node)
		}

	case TermNode:
		return w.WalkTermNode(node)

	case WrapReNode:
		if fn := w.EnterWrapReNode; fn != nil {
			return fn(node)
		}

	}
	return nil
}
func (w WalkerOps) WalkAtomExtRefNode(node AtomExtRefNode) Stopper {
	if fn := w.EnterAtomExtRefNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneIdent(); child != nil {
		child := *child
		if fn := w.EnterIdentNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}

	if fn := w.ExitAtomExtRefNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkAtomNode(node AtomNode) Stopper {
	if fn := w.EnterAtomNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneExtRef(); child != nil {
		child := *child
		if s := w.WalkAtomExtRefNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneIdent(); child != nil {
		child := *child
		if fn := w.EnterIdentNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	if child := node.OneMacrocall(); child != nil {
		child := *child
		if s := w.WalkMacrocallNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneRe(); child != nil {
		child := *child
		if fn := w.EnterReNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	if child := node.OneRef(); child != nil {
		child := *child
		if s := w.WalkRefNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneStr(); child != nil {
		child := *child
		if fn := w.EnterStrNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	if child := node.OneTerm(); child != nil {
		child := *child
		if s := w.WalkTermNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitAtomNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkGrammarNode(node GrammarNode) Stopper {
	if fn := w.EnterGrammarNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	for _, child := range node.AllStmt() {
		if s := w.WalkStmtNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitGrammarNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkMacrocallNode(node MacrocallNode) Stopper {
	if fn := w.EnterMacrocallNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneName(); child != nil {
		child := *child
		if fn := w.EnterIdentNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	for _, child := range node.AllTerm() {
		if s := w.WalkTermNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitMacrocallNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkNamedNode(node NamedNode) Stopper {
	if fn := w.EnterNamedNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneAtom(); child != nil {
		child := *child
		if s := w.WalkAtomNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneIdent(); child != nil {
		child := *child
		if fn := w.EnterIdentNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}

	if fn := w.ExitNamedNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkPragmaImportNode(node PragmaImportNode) Stopper {
	if fn := w.EnterPragmaImportNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OnePath(); child != nil {
		child := *child
		if s := w.WalkPragmaImportPathNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitPragmaImportNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkPragmaImportPathNode(node PragmaImportPathNode) Stopper {
	if fn := w.EnterPragmaImportPathNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitPragmaImportPathNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkPragmaMacrodefNode(node PragmaMacrodefNode) Stopper {
	if fn := w.EnterPragmaMacrodefNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	for _, child := range node.AllArgs() {
		if fn := w.EnterIdentNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	if child := node.OneName(); child != nil {
		child := *child
		if fn := w.EnterIdentNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	if child := node.OneTerm(); child != nil {
		child := *child
		if s := w.WalkTermNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitPragmaMacrodefNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkPragmaNode(node PragmaNode) Stopper {
	if fn := w.EnterPragmaNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneImport(); child != nil {
		child := *child
		if s := w.WalkPragmaImportNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneMacrodef(); child != nil {
		child := *child
		if s := w.WalkPragmaMacrodefNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitPragmaNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkProdNode(node ProdNode) Stopper {
	if fn := w.EnterProdNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneIdent(); child != nil {
		child := *child
		if fn := w.EnterIdentNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	for _, child := range node.AllTerm() {
		if s := w.WalkTermNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitProdNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkQuantNode(node QuantNode) Stopper {
	if fn := w.EnterQuantNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneMax(); child != nil {
		child := *child
		if fn := w.EnterIntNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	if child := node.OneMin(); child != nil {
		child := *child
		if fn := w.EnterIntNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	if child := node.OneNamed(); child != nil {
		child := *child
		if s := w.WalkNamedNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitQuantNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkRefNode(node RefNode) Stopper {
	if fn := w.EnterRefNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneDefault(); child != nil {
		child := *child
		if fn := w.EnterStrNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	if child := node.OneIdent(); child != nil {
		child := *child
		if fn := w.EnterIdentNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}

	if fn := w.ExitRefNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkStmtNode(node StmtNode) Stopper {
	if fn := w.EnterStmtNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneComment(); child != nil {
		child := *child
		if fn := w.EnterCommentNode; fn != nil {
			if s := fn(child); s != nil {
				if s.ExitNode() {
					return nil
				} else if s.Abort() {
					return s
				}
			}
		}
	}
	if child := node.OnePragma(); child != nil {
		child := *child
		if s := w.WalkPragmaNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneProd(); child != nil {
		child := *child
		if s := w.WalkProdNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitStmtNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

func (w WalkerOps) WalkTermNode(node TermNode) Stopper {
	if fn := w.EnterTermNode; fn != nil {
		if s := fn(node); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	for _, child := range node.AllGrammar() {
		if s := w.WalkGrammarNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	if child := node.OneNamed(); child != nil {
		child := *child
		if s := w.WalkNamedNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	for _, child := range node.AllQuant() {
		if s := w.WalkQuantNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}
	for _, child := range node.AllTerm() {
		if s := w.WalkTermNode(child); s != nil {
			if s.ExitNode() {
				return nil
			} else if s.Abort() {
				return s
			}
		}
	}

	if fn := w.ExitTermNode; fn != nil {
		if s := fn(node); s != nil && s.Abort() {
			return s
		}
	}
	return nil
}

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
stmt    -> COMMENT | prod | pragma;
prod    -> IDENT "->" term+ ";";
term    -> (@ ("{" grammar "}")? ):op=">"
         > @:op="|"
         > @+
         > named quant*;
named   -> (IDENT op="=")? atom;
quant   -> op=[?*+]
         | "{" min=INT? "," max=INT? "}"
         | op=/{<:|:>?} opt_leading=","? named opt_trailing=","?;
atom    -> IDENT | STR | RE | macrocall | ExtRef=("%%" IDENT) | REF | "(" term ")" | "(" ")";

macrocall   -> "%!" name=IDENT "(" term:","? ")";
REF         -> "%" IDENT ("=" default=STR)?;

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

// Special
pragma  -> import | macrodef {
                import   -> ".import" path=((".."|"."|[a-zA-Z0-9.:]+):,"/") ";"?;
                macrodef -> ".macro" name=IDENT "(" args=IDENT:","? ")" "{" term "}" ";"?;
            };

.wrapRE -> /{\s*()\s*};
`)
