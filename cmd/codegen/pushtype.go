package codegen

func (tm *TypeMap) findType(name string) GrammarType {
	return (*tm)[name]
}

func (tm *TypeMap) pushType(name, parent string, child GrammarType) {
	if name == "" {
		tm.createOrAddParent(parent, child)
		return
	}

	switch child.(type) {
	case namedToken, namedRule, backRef, stackBackRef:
		// Dont' create a type it if doesnt exist
		tm.createOrAddParent(parent, child)
	default:
		panic("It doesnt make sense to add this type with a name param!")
	}
}

func (tm *TypeMap) createOrAddParent(parent string, child GrammarType) {
	parentTypeName := GoTypeName(parent)
	var val GrammarType
	if p := tm.findType(parentTypeName); p != nil {
		var children []GrammarType
		// Check if the parent needs to be upgraded to a rule() instead of a basicrule()
		if basic, ok := p.(basicRule); ok {
			children = []GrammarType{basic.Upgrade()}
		} else if _, ok := p.(rule); ok {
			children = p.Children()
		} else {
			panic("Only rule{} or basicrule() can be a root node")
		}
		val = rule{name: parentTypeName, childs: checkForDupes(children, child)}
	} else {
		// Need a new parent
		if v, ok := child.(unnamedToken); ok && v.count.wantOne() {
			val = basicRule(parentTypeName)
		} else {
			val = rule{name: parentTypeName, childs: []GrammarType{child}}
		}
	}
	(*tm)[parentTypeName] = val
}

func getNewCount(old countManager, new GrammarType) countManager {
	switch t := new.(type) {
	case unnamedToken:
		return old.merge(t.count)
	case namedToken:
		return old.merge(t.count)
	case namedRule:
		return old.merge(t.count)
	case stackBackRef, backRef:
		return old.forceMany()
	}
	return old
}

func checkForDupes(children []GrammarType, next GrammarType) []GrammarType {
	if next == nil {
		return children
	}
	result := make([]GrammarType, 0, len(children)+1)
	appendNext := true
	for _, c := range children {
		if next.Ident() != c.Ident() {
			result = append(result, c)
			continue
		}
		switch child := c.(type) {
		case unnamedToken:
			child.count = getNewCount(child.count, next)
			c = child
			appendNext = false
		case namedToken:
			child.count = getNewCount(child.count, next)
			c = child
			appendNext = false
		case namedRule:
			child.count = getNewCount(child.count, next)
			c = child
			appendNext = false
		case stackBackRef:
			if _, ok := next.(stackBackRef); ok {
				return children
			}
		case choice:
			if _, ok := next.(choice); ok {
				return children
			}
		}
		result = append(result, c)
	}
	if appendNext {
		return append(result, next)
	}
	return result
}
