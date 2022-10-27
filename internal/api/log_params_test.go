
package api

import (
    "testing"
    "github.com/stretchr/testify/assert"
)



func TestLogParams(t *testing.T) {

    type test_struct struct {
        First string
        Second int
    }

    p := LogParams{}

    p.Add("one", 1)
    p.Add("string", "str")
    p.Add("struct", test_struct{First:"first_string",Second:1234})

    expected := []interface{}([]interface {}{
        "one",1,
        "string","str",
        "struct",test_struct{
            First:"first_string",
            Second:1234,
        },
    })

    assert.ElementsMatch(t, expected, p)
}

func TestLogParamsCombine(t *testing.T) {

    a := LogParams{"one",1,"string1","str1"}

    b := LogParams{"two",2,"string2","str2"}

    expected := LogParams{
        "one",1,
        "string1","str1",
        "two",2,
        "string2","str2",
    }

    assert.Equal(t, expected, a.Combine(b))
}