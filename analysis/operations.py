#!/usr/bin/env python3

fset = lambda *args: frozenset(args)

alphabet = set("abc")
transitions = {
    l + r
    for l in alphabet|{"+"}
    for r in alphabet|{"-"}
}

for i in range(2**len(transitions)):
    
print(len(transitions))
