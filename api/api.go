package api

import (
	"bytes"
	"fmt"
	"log"
)

func die(err error) {
	if err != nil {
		fmt.Println("Error in module installtion")
		log.Fatal(err)
	}
}

func Concat(path ...string) string {
	var b bytes.Buffer

	for _, p := range path {
		b.WriteString(p)
	}

	return b.String()
}
