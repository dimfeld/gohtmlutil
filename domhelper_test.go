package gohtmlutil

import (
	"bytes"
	"code.google.com/p/go.net/html"
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
	buf := bytes.NewBufferString(document)
	doc, err := html.Parse(buf)
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
	// Test cache by reading again.
	testNode("html/body/h1", "Header")
	testNode("html/body/#header", "Header")
	testNode("html/body/h1#header", "Header")
	testNode("html/body/p", "PClass")
	testNode("html/body/p#content", "PClass")
	testNode("html/body/p.pclass", "PClass")
	testNode("html/body/div/div", "AnonymousDiv")
	testNode("html/body/div/div#abc", "AbcName")
	testNode("html/body/div/div.abc", "AbcClass")
	testNode("html/body/div/div.abc#abc", "AbcNameClass")
	testNode("html/body/div/div#abc.abc", "AbcNameClass")
	testNode("html/body/div/div/ul/li", "123")
	testNode("html/body/div/div/ul/li.gray", "GrayRedClass")
	testNode("html/body/div/div/ul/li.red", "GrayRedClass")
	testNode("html/body/div/div/ul/.red", "GrayRedClass")
	testNode("html/body/div/div/ul/li#SecondList", "SecondList")
	testNode("html/body/div/div/ul/#SecondList", "SecondList")
	testNodeMiss("html/body/h1#abc")
	testNodeMiss("h1")
	testNodeMiss("html/body/div/div/ul/li#SecondList.SecondList")
	testNodeMiss("html/body/div/div/ul/li.SecondList")

	// Make sure it works when not at the document root
	doc, _ = Find(doc, "html/body/div")
	testNode("div/ul/li", "123")

	buf = &bytes.Buffer{}
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
