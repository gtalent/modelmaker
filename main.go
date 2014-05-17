/*
   Copyright 2013 gtalent2@gmail.com

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/gtalent/cyborgbear/parser"
	"io/ioutil"
	"os"
)

func main() {
	out := flag.String("o", "stdout", "File or file set(languages with header files) to write the output to")
	in := flag.String("i", "", "The model file to generate JSON serialization code for")
	namespace := flag.String("n", "models", "Namespace for the models")
	outputType := flag.String("t", "cpp-jansson", "Output type(cpp-jansson, cpp-qt, go)")
	boost := flag.Bool("cpp-boost", false, "Boost serialization enabled")
	lowerCase := flag.Bool("lc", false, "Make variable names lowercase in output models")
	version := flag.Bool("v", false, "version")
	flag.Parse()

	if *version {
		fmt.Println("cyborgbear version " + cyborgbear_version)
		return
	}
	parseFile(*in, *out, *namespace, *outputType, *boost, *lowerCase)
}

func parseFile(path, outFile, namespace, outputType string, boost, lowerCase bool) {
	ss, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not find or open specified model file")
		os.Exit(1)
	}
	input := string(ss)

	ioutputType := USING_JANSSON
	switch outputType {
	case "cpp-jansson":
		ioutputType = USING_JANSSON
	case "cpp-qt":
		ioutputType = USING_QT
	case "go", "Go":
		ioutputType = USING_GO
	}

	var out Out
	switch ioutputType {
	case USING_JANSSON, USING_QT:
		out = NewCOut(namespace, ioutputType, boost, lowerCase)
	case USING_GO:
		out = NewGo(namespace)
	}

	models, err := parser.Parse(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
		return
	} else {
		for _, v := range models {
			out.addClass(v.Name)
			for _, vv := range v.Vars {
				out.addVar(vv.Name, vv.Type)
			}
			out.closeClass(v.Name)
		}

		if outFile == "stdout" {
			fmt.Print(out.write(""))
		} else {
			out.writeFile(outFile)
		}
	}
}
