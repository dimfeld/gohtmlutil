// gohtmlutil provides helper functions for use with the package `code.google.com/p/go.net/html`.
package gohtmlutil

import (
	"code.google.com/p/go.net/html"
	"strconv"
	"strings"
)

// Find a particular node in the tree.
// The input is a slash-separated list of elements,
// #names, and .classes to search for. Elements and classes
// may be combined, such as div#contentName or ul.listClass.
// A token may also be prefixed with a count N and asterisk (e.g. 2*span),
// which will find the Nth match.
func Find(root *html.Node, path string) (node *html.Node, ok bool) {
	tokens := strings.Split(path, "/")
	node, ok = findHelper(tokens, root)
	return
}

func findHelper(tokens []string, parent *html.Node) (*html.Node, bool) {
	index := -1
	element := tokens[0]
	class := ""
	name := ""

	if s := strings.Split(element, "*"); len(s) == 2 {
		var err error
		index, err = strconv.Atoi(s[0])
		if err != nil {
			// Better error handling here?
			index = -1
		}
		element = s[1]
	}

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

	matchElementCount := 0
	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
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
						} else {
							break
						}
					}

					if !classMatch && attr.Key == "class" {
						classes := strings.Split(attr.Val, " ")
						for _, oneClass := range classes {
							if oneClass == class {
								classMatch = true
							}
						}
						if !classMatch {
							break
						}
					}
				}
			}

			if nameMatch && classMatch {
				matchElementCount++
				if index != -1 && index != matchElementCount {
					// We're looking for the nth element, and haven't reached n yet.
					continue
				}

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
