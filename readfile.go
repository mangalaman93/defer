package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	jsonFile, err := os.Open("priv/states.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	states := make(map[string]string)
	if err := json.Unmarshal(data, &states); err != nil {
		panic(err)
	}

	for shortName, fullName := range states {
		fmt.Println(fullName, "->", shortName)
	}
}
