package html

import (
	ghtml "code.google.com/p/go.net/html"
	"os"
	"strings"
)

type HTMLDocument struct {
	Root       *ghtml.Node
	cachedPath map[string]*ghtml.Node
}

// Find a particular node in the tree.
// The input is a slash-separated list of elements,
// #names, and .classes to search for. Elements and classes
// may be combined, such as div#contentName or ul.listClass.
func (d HTMLDocument) Find(path string) (node *gthml.Node, ok bool) {
	tokens := strings.Split(path, "/")

	// First see if this is a cached path, to speed up the search.
	// Start with the longest path and go backward.
	for i := len(tokens); i > 0; i-- {
		tryCacheTokens := tokens[0:i]
		tryString := strings.Join(tryCacheTokens, "/")
		if node, ok = d.cachedPath[tryString]; ok {
			tokens = tokens[i:len(tokens)]
			break
		}
	}

	if !ok {
		node = d.Root
	}

	node, ok = findHelper(tokens, node)
	if ok {
		d.cachedPath[path] = node
	}
	return
}

func findHelper(tokens []string, parent *ghtml.Node) (*ghtml.Node, bool) {
	element := tokens[0]
	class := ""
	name := ""

	// Get the class, if given.
	if s := strings.Split(element, "."); len(s) == 2 {
		element = s[0]
		class = s[1]
	}

	// Get the name, if given, which could be either in the element or in the class.
	if s := strings.Split(element, "#"); len(s) == 2 {
		element = s[0]
		name = s[1]
	} else if s := strings.Split(class, "#"); len(s) == 2 {
		class = s[0]
		name = s[1]
	}

	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != ghtml.ElementNode {
			continue
		}
		if c.Data == element || element == "" {
			nameMatch := name == ""
			classMatch := class == ""

			if !nameMatch || !classMatch {
				for _, attr := range c.Attr {
					if nameMatch && classMatch {
						break
					}

					if !nameMatch && attr.Key == "name" {
						if attr.Val == name {
							nameMatch = true
						}
						break
					}

					if !classMatch && attr.Key == "class" {
						classes = strings.Split(attr.Val, " ")
						for _, oneClass := range classes {
							if oneClass == class {
								classMatch = true
								break
							}
						}
						break
					}
				}
			}

			if nameMatch && classMatch {
				if len(tokens) == 1 {
					// This is the final token the user is looking for.
					return c, true
				} else {
					// The user wants more nested tokens. Call the helper and return
					// the first one that found the final token.
					node, ok := findHelper(tokens[1:len(tokens)], c)
					if ok {
						return node, true
					}
				}
			}
		}
	}

	// If we get down to here, then nothing was found
	return nil, false
}

func LoadFile(filename string) (HTMLDocument, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

func Parse(reader io.Reader) (HTMLDocument, error) {
	root, err := gthml.Parse(reader)
	return HTMLDocument{Root: root, cachedPath: make(map[string]*ghtml.Node)}, err
}
