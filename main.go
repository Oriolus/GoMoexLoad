package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	secs "github.com/oriolus/moex/securities"
)

func main() {

	secs, err := secs.TransformAll("securities_info")

	if err != nil {
		fmt.Println(err)
	}

	buf, err := json.Marshal(secs)

	if err != nil {
		panic(err.Error())
	}

	writeErr := ioutil.WriteFile("securitites_all.json", buf, 0777)

	if writeErr != nil {
		fmt.Println(writeErr)
	}

	fmt.Println(len(secs))

}
