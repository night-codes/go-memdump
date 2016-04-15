package memdump

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func join(bufs ...[]byte) []byte {
	return bytes.Join(bufs, nil)
}

func TestDelimitedReader_Simple(t *testing.T) {
	data := join([]byte("abc"), delim, []byte("defggg"), delim)
	r := newDelimitedReader(bytes.NewReader(data))

	buf := make([]byte, 100)

	n, err := r.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, "abc", string(buf[:n]))

	r.Next()

	n, err = r.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 6, n)
	assert.Equal(t, "defggg", string(buf[:n]))

	r.Next()
}

func TestDelimitedReader_Long(t *testing.T) {
	data := make([]byte, 132)
	for i := range data {
		data[i] = 0
	}
	copy(data[100:], delim)
	r := newDelimitedReader(bytes.NewReader(data))

	buf := make([]byte, 80)

	n, err := r.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, 80, n)

	n, err = r.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 20, n)

	r.Next()
}

func TestDelimitedReader_SimpleThenEmpty(t *testing.T) {
	data := join([]byte("abc"), delim, delim)
	r := newDelimitedReader(bytes.NewReader(data))

	buf := make([]byte, 100)

	n, err := r.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, "abc", string(buf[:n]))

	r.Next()

	n, err = r.Read(buf)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)

	r.Next()
}
func TestDelimitedReader_EmptyUnterminated(t *testing.T) {
	r := newDelimitedReader(bytes.NewReader(nil))

	buf := make([]byte, 100)

	n, err := r.Read(buf)
	assert.Equal(t, ErrUnexpectedEOF, err)
	assert.Equal(t, 0, n)
}

func TestDelimitedReader_Unterminated(t *testing.T) {
	r := newDelimitedReader(bytes.NewReader([]byte("abc")))
	buf, err := ioutil.ReadAll(r)
	assert.Equal(t, ErrUnexpectedEOF, err)
	assert.Equal(t, "abc", string(buf))
}
