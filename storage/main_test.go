package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunc(t *testing.T) {
	cases := []struct {
		desc string
		a    int
		b    int
		exp  int
	}{
		{desc: "test 1", a: 1, b: 2, exp: 3},
		{desc: "test 2", a: 2, b: 2, exp: 4},
		{desc: "test 3", a: 1, b: 2, exp: 33},
	}

	for _, v := range cases {
		res := v.a + v.b
		t.Run(v.desc, func(t *testing.T) {
			assert.Equal(t, v.exp, res)
		})
	}
}
