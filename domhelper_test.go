package gohtmlutil

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"fmt"
	"strings"
	"testing"
)

var document string = `<html><body>
    <h1 name="header">Header</h1>
	<p name="content" class="pclass">PClass</p>
	<div>
		<div>AnonymousDiv</div>
		<div name="abc">AbcName</div>
		<div class="abc">AbcClass</div>
		<div name="abc" class="abc">AbcNameClass</div>
		<div>AnotherAnonymousDiv
			<ul>
			<li>123</li>
			<li class="gray red">GrayRedClass</li>
			<li>456</li>
			</ul>
			<ul>
			<li name="SecondList">SecondList</li>
			</ul>
		</div>
	</div>
	</body>
	</html>`

func nodeDesc(node *html.Node) string {
	desc := node.Data
	for _, attr := range node.Attr {
		if attr.Key == "name" {
			desc = desc + "#" + attr.Val
		}

		if attr.Key == "class" {
			desc = desc + "." + attr.Val
		}
	}
	return desc
}

func pathToRoot(node *html.Node) string {
	path := node.Data
	for node.Type != html.DocumentNode {
		node = node.Parent
		path = node.Data + "/" + path
	}
	return path
}

func TestFind(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(document))
	if err != nil {
		t.Error("Failed to parse data with error", err)
		t.FailNow()
	}

	// Look for a node, and use the first child's data to verify that it
	// is the right one.
	testNode := func(path, childData string) {
		node, ok := Find(doc, path)
		if !ok {
			t.Error("No match for path", path)
			return
		}
		child := node.FirstChild
		if child == nil || child.Data != childData {
			t.Errorf("Path %s found incorrect node %s at %s",
				path, nodeDesc(node), pathToRoot(node))
			if child == nil {
				t.Error("Child was nil")
			} else {
				t.Error("Child data was", child.Data)
			}
		}
	}

	testNodeMiss := func(path string) {
		node, ok := Find(doc, path)
		if ok {
			t.Errorf("Expected miss for path %s but found %s at %s",
				path, nodeDesc(node), pathToRoot(node))
		}
	}

	testNode("html/body/h1", "Header")
	testNode("html/body/h1", "Header")
	testNode("html/body/#header", "Header")
	testNode("html/body/h1#header", "Header")
	testNode("html/body/p", "PClass")
	testNode("html/body/p#content", "PClass")
	testNode("html/body/p.pclass", "PClass")
	testNode("html/body/div/div", "AnonymousDiv")
	testNode("html/body/div/div#abc", "AbcName")
	testNode("html/body/div/div.abc", "AbcClass")
	testNode("html/body/div/2*div.abc", "AbcNameClass")
	testNode("html/body/div/div.abc#abc", "AbcNameClass")
	testNode("html/body/div/div#abc.abc", "AbcNameClass")
	testNode("html/body/div/div/ul/li", "123")
	testNode("html/body/div/div/ul/li.gray", "GrayRedClass")
	testNode("html/body/div/div/ul/li.red", "GrayRedClass")
	testNode("html/body/div/div/ul/.red", "GrayRedClass")
	testNode("html/body/div/div/ul/li#SecondList", "SecondList")
	testNode("html/body/div/div/ul/#SecondList", "SecondList")
	testNode("html/body/div/2*div", "AbcName")
	testNode("html/body/1*div/2*div", "AbcName")
	testNode("html/body/div/5*div/ul/li", "123")
	testNode("html/body/div/1*div", "AnonymousDiv")
	testNode("html/body/div/div/1*ul/.gray", "GrayRedClass")
	testNodeMiss("html/body/div/4*div/ul/li")
	testNodeMiss("html/body/2*div")
	testNodeMiss("html/body/0*div")
	testNodeMiss("html/body/div/div/2*ul/.gray")
	testNodeMiss("html/body/h1#abc")
	testNodeMiss("h1")
	testNodeMiss("html/body/div/div/ul/li#SecondList.SecondList")
	testNodeMiss("html/body/div/div/ul/li.SecondList")

	// Make sure it works when not at the document root
	doc, _ = Find(doc, "html/body/div")
	testNode("div/ul/li", "123")

	buf := &bytes.Buffer{}
	html.Render(buf, doc)
	t.Log(string(buf.Bytes()))
}

func BenchmarkFind(b *testing.B) {
	buf := bytes.NewBufferString(document)
	doc, err := html.Parse(buf)
	if err != nil {
		b.Error("Failed to parse data with error", err)
		b.FailNow()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Find(doc, "html/body/div/div/ul/li#SecondList")
	}
}

func ExampleFind() {
	document := `<html><body><div>
		<span>Some text</span>
		<span name="abc">ABC</span>
		<span class="fancytext">Fancy Text</span>
		</div>
		</body>
		</html>`
	root, err := html.Parse(strings.NewReader(document))
	if err != nil {
		fmt.Println("Error parsing document:", err)
		return
	}

	node, _ := Find(root, "html/body/div/#abc")
	fmt.Println("Text for #abc is", node.FirstChild.Data)

	node, _ = Find(root, "html/body/div/span.fancytext")
	fmt.Println("Text for span.fancytext is", node.FirstChild.Data)

	node, _ = Find(root, "html/body/div/2*span")
	fmt.Println("Text for 2nd span element is", node.FirstChild.Data)

	divNode, _ := Find(root, "html/body/div")
	node, _ = Find(divNode, "3*span")
	fmt.Println("Text for 3rd span element is", node.FirstChild.Data)

	// Output:
	// Text for #abc is ABC
	// Text for span.fancytext is Fancy Text
	// Text for 2nd span element is ABC
	// Text for 3rd span element is Fancy Text
}
