// A simple binary to link together the mapper servers.
package plugin_test

import (
	"github.com/jackskj/protoc-gen-map/examples"
	"github.com/jackskj/protoc-gen-map/testdata"
	"github.com/jackskj/protoc-gen-map/testdata/gentest"
	"github.com/jackskj/protoc-gen-map/testdata/initdb"
	"testing"
)

func TestLink(t *testing.T) {
	_ = examples.BlogRequest{}
	_ = testdata.TestReflectServiceMapServer{}
	_ = gentest.EmptyRequest{}
	_ = initdb.InsertTagRequest{}
}
