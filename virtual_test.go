package rice

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/GeertJohan/go.rice/embedded"
	"github.com/stretchr/testify/require"
)

func TestVirtualFileRead(t *testing.T) {
	// define virtual file
	data, err := ioutil.ReadFile("testdata/test.txt")
	require.NoError(t, err)
	virtFile := newVirtualFile(&embedded.EmbeddedFile{
		Filename:    "test.txt",
		FileModTime: time.Unix(1594051142, 0),
		Content:     string(data),
	})

	realFile, err := os.Open("testdata/test.txt")
	require.NoError(t, err)

	checkRead := func(buff []byte) {
		// Compare results with a real *os.File
		realN, realErr := realFile.Read(buff)
		real := fmt.Sprintf("n=%v err=%v buff=%v", realN, realErr, buff)
		virtN, virtErr := virtFile.read(buff)
		virt := fmt.Sprintf("n=%v err=%v buff=%v", virtN, virtErr, buff)
		require.Equal(t, real, virt)
	}

	buff := make([]byte, 4)
	checkRead(buff)
	checkRead(buff)
	checkRead(buff)
}
