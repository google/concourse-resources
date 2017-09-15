package temputil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	tempMan := TempFileManager{}
	defer tempMan.Cleanup()

	contents := "file contents"
	path, err := tempMan.Create("temputil-test", contents)
	assert.NoError(t, err)

	actualContents, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	assert.EqualValues(t, contents, actualContents)
}

func TestCleanup(t *testing.T) {
	tempMan := TempFileManager{}

	path, err := tempMan.Create("temputil-test", "foo")
	_, err = os.Stat(path)
	assert.NoError(t, err)

	tempMan.Cleanup()
	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))

}
