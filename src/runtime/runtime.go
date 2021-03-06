package runtime

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Condition func() bool
type Machine interface {
	Init()
	InState(string) bool
	FireEvent(string) bool
	ValidateConditions()
}

var MachineRegistry = map[string]Machine{}
var Bind = ":" + strconv.Itoa(os.Getpid())

func AddMachine(name string, m Machine) {
	m.Init()
	MachineRegistry[name] = m
}

func Validate() {
	for _, v := range MachineRegistry {
		v.ValidateConditions()
	}
}

func Run() {
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			Validate()
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if m, ok := MachineRegistry[component[0]]; ok {
			m.FireEvent(component[1])
		}
	})
	http.ListenAndServe(Bind, nil)

	os.Exit(0)
}
