package html

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var document string = `<html><body>
    <h1 name="header">Header</h1>
	<p name="content" class="pclass">PClass
		<div>AnonymousDiv</div>
		<div name="abc">AbcName</div>
		<div class="abc">AbcClass</div>
		<div name="abc" class="abc">AbcNameClass</div>
		<div>AnotherAnonymousDiv
			<ul>
			<li>123</li>
			<li class="gray">GrayClass</li>
			<li>456</li>
			</ul>
			<ul>
			<li name="SecondList">SecondList</li>
			</ul>
		</div>
	</p>
	</body>
	</html>`

func testFind(t *testing.T, doc HTMLDocument) {

}

func TestLoadFile(t *testing.T) {
	tempFile, err = ioutil.TempFile("", "htmldoc")
	if err != nil {
		t.Error("Could not create temporary file")
		t.FailNow()
	}
	filename := tempFile.Name()
	defer os.Remove(filename)

	buf := bytes.NewBufferString(document)
	_, err = buf.WriteTo(tempFile)
	tempFile.Close()
	if err != nil {
		t.Error("Could not write data to temporary file")
		t.FailNow()
	}

	document, err := LoadFile(tempFile)
	testFind(t, document)
}

func TestParse(t *testing.T) {
	testFind(t, document)
}

// BenchmarkFind ensures that repeated finds on the same complex path
// is fast.
func BenchmarkFind(b *testing.B) {

}
