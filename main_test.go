package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestCmp(t *testing.T) {
	left, err := os.Open("testdata/1-left")
	require.NoError(t, err)
	right, err := os.Open("testdata/1-right")
	require.NoError(t, err)
	wantOut, err := ioutil.ReadFile("testdata/1-cmp")
	require.NoError(t, err)

	out := strings.Builder{}
	cmp(&out, left, right)

	assert.Equal(t, string(wantOut), out.String())
}
