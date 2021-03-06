package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"flag"
	"os"
	"fmt"
)

var (
	coregroupFile = flag.String("file", "", "")
)

var usage = `Usage: coregroups [options...]

Options:

  -file  		JSON-file containing your endpoints
`

type coregroup struct {
	Application   string `json:"application"`
	CoregroupName string `json:"coregroupName"`
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func viewHandler(coregroups *[]coregroup) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var applicationName string
		if r.URL.Query().Get("application") == "" {
			for _, cg := range *coregroups {
				var line string
				line = cg.Application + ":" + cg.CoregroupName + "\n"
				w.Write([]byte(line))
			}
		} else {
			applicationName = r.URL.Query().Get("application")
			for _, correctCoregroup := range *coregroups {
				if applicationName == correctCoregroup.Application {
					w.Write([]byte(correctCoregroup.CoregroupName))
					return
				}
			}
			w.Write([]byte("DefaultCoreGroup"))
		}
	})
}

func main() {

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	if flag.NFlag() < 1 {
		usageAndExit("You did not supply enough arguments")
	}

	data, err := ioutil.ReadFile(*coregroupFile)

	if err != nil {
		log.Fatal("unable to read file ", *coregroupFile)
		panic(err)
	}

	coregroups := []coregroup{}
	err = json.Unmarshal(data, &coregroups)

	if err != nil {
		log.Fatal("Couldn't parse JSON: ", string(data))
		panic(err)
	}

	mux := http.NewServeMux()
	vh := viewHandler(&coregroups)
	mux.HandleFunc("/isAlive", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "")
	})
	mux.HandleFunc("/isReady", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "")
	})
	mux.Handle("/", vh)
	err = http.ListenAndServe(":80", mux)

	if err != nil {
		log.Fatal("Couldn't start application. ", err)
		panic(err)
	}
}
