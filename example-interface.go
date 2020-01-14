// Go offers built-in support for JSON encoding and
// decoding, including to and from built-in and custom
// data types.

package main

import (
	"encoding/json"
	"fmt"
	//"os"
)

// We'll use these two structs to demonstrate encoding and
// decoding of custom types below.
type response1 struct {
	Page   int
	Fruits []string
}

// Only exported fields will be encoded/decoded in JSON.
// Fields must start with capital letters to be exported.
type response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func main() {

	byt := []byte(`{
  "URL": "ampq://rabbitMQ:5032",
  "driver": "rabbitMQ",
  "args": {
    "exchange": {
      "autoDelete": false,
      "name": "api",
      "durable": true,
      "internal": false,
      "noWait": false
    },
    "queues": [
      {
        "name": "escola",
        "durable": true,
        "autoDelete": false,
        "exclusive": false,
        "noWait": false,
        "bindingKeys": [
          "escola",
          "turma",
          "aluno"
        ]
      }
    ]
  }
}
`)

	// We need to provide a variable where the JSON
	// package can put the decoded data. This
	// `map[string]interface{}` will hold a map of strings
	// to arbitrary data types.
	var dat map[string]interface{}

	// Here's the actual decoding, and a check for
	// associated errors.
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

	// In order to use the values in the decoded map,
	// we'll need to convert them to their appropriate type.
	// For example here we convert the value in `num` to
	// the expected `float64` type.
	num := dat["args"].(map[string]interface{})
	fmt.Println(num)
	q := num["queues"].([]interface{})
	fmt.Println(q)
	b := q[0].(map[string]interface{})
	fmt.Println(b["name"].())
	// Accessing nested data requires a series of ["bindingKeys"]["bindingKeys"].(interface{})
	// conversions.
	//strs := dat["strs"].([]interface{})
	//str1 := strs[0].(string)
	//fmt.Println(str1)


}
