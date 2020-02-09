package wbnf

import (
	"io"

	"github.com/arr-ai/wbnf/errors"
	"github.com/arr-ai/wbnf/parser"
)

// The following methods assume a valid parse. Call (Term).ValidateParse first if
// unsure.

func (t S) Unparse(g Grammar, e parser.TreeElement, w io.Writer) (n int, err error) {
	return w.Write([]byte(e.(parser.Scanner).String()))
}

func (t RE) Unparse(g Grammar, e parser.TreeElement, w io.Writer) (n int, err error) {
	return w.Write([]byte(e.(parser.Scanner).String()))
}
func (t REF) Unparse(g Grammar, e parser.TreeElement, w io.Writer) (n int, err error) {
	return w.Write([]byte("\\" + e.(parser.Scanner).String()))
}

func unparse(g Grammar, term Term, e parser.TreeElement, w io.Writer, N *int) error {
	n, err := term.Unparse(g, e, w)
	if err == nil {
		*N += n
	}
	return err
}

func (t Seq) Unparse(g Grammar, e parser.TreeElement, w io.Writer) (n int, err error) {
	node := e.(parser.Node)
	for i, term := range t {
		if err = unparse(g, term, node.Children[i], w, &n); err != nil {
			return
		}
	}
	return n, nil
}

func (t Oneof) Unparse(g Grammar, e parser.TreeElement, w io.Writer) (n int, err error) {
	node := e.(parser.Node)
	return t[node.Extra.(Choice)].Unparse(g, node.Children[0], w)
}

func (t Delim) Unparse(g Grammar, e parser.TreeElement, w io.Writer) (n int, err error) {
	node := e.(parser.Node)
	tgen := t.LRTerms(node)
	for _, child := range node.Children {
		if err = unparse(g, tgen.Next(), child, w, &n); err != nil {
			return
		}
	}
	return
}

func (t Quant) Unparse(g Grammar, e parser.TreeElement, w io.Writer) (n int, err error) {
	for _, child := range e.(parser.Node).Children {
		if err = unparse(g, t.Term, child, w, &n); err != nil {
			return
		}
	}
	return
}

//-----------------------------------------------------------------------------

func (t Rule) Unparse(g Grammar, e parser.TreeElement, w io.Writer) (n int, err error) {
	return g[t].Unparse(g, e, w)
}

//-----------------------------------------------------------------------------

func (t Stack) Unparse(_ Grammar, _ parser.TreeElement, _ io.Writer) (int, error) {
	panic(errors.Inconceivable)
}

//-----------------------------------------------------------------------------

func (t Named) Unparse(g Grammar, e parser.TreeElement, w io.Writer) (n int, err error) {
	err = unparse(g, t.Term, e, w, &n)
	return
}
