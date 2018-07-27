package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	goyaml "gopkg.in/yaml.v2"
)

var version = "latest"
var versionFlag = flag.Bool("version", false, fmt.Sprintf("prints current %s version", os.Args[1:]))
var yaml2jsonFlag = flag.Bool("yaml2json", false, "yaml2json")

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
	if *yaml2jsonFlag {
		yaml2json()
	} else {
		json2yaml()
	}
}

func yaml2json() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	decoder := goyaml.NewDecoder(os.Stdin)

	var data interface{}
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	err = parse(&data)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	_, err = out.WriteTo(os.Stdout)
	if err == io.EOF {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
}

func parse(i *interface{}) (err error) {
	switch in := (*i).(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(in))
		for k, v := range in {
			if err = parse(&v); err != nil {
				return err
			}
			var s string
			switch k.(type) {
			case string:
				s = k.(string)
			case int:
				s = strconv.Itoa(k.(int))
			default:
				return fmt.Errorf("type mismatch: expect map key string or int; got: %T", k)
			}
			m[s] = v
		}
		*i = m
	case []interface{}:
		for i := len(in) - 1; i >= 0; i-- {
			if err = parse(&in[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func json2yaml() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	encoder := goyaml.NewEncoder(os.Stdout)

	var data interface{}

	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Fatal(err)
	}

	err = encoder.Encode(&data)
	if err != nil {
		log.Fatal(err)
	}
}
