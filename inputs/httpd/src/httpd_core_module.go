/*
 * Copyright (C) 2017 Meng Shi
 */

package httpd

import (
      "net/http"
      "unsafe"

      "github.com/gorilla/mux"
    . "github.com/rookie-xy/worker/types"
    "strings"
"net"
    "fmt"
    "time"
)

const (
    LOCATION_MODULE = HTTPD_MODULE|0x00100000
    LOCATION_CONFIG = HTTPD_CONFIG|CONFIG_MAP
)

type HttpdCore struct {
    //*Log
    *Cycle
    //*Routine

     listen    string
     timeout   int
     location  *LocationHttpd
     gid       []int64

     listener   net.Listener
}

func NewHttpdCore() *HttpdCore {
    return &HttpdCore{}
}

func (hc *HttpdCore) Init() int {
    document := hc.location.document
    fmt.Println(document)
    path := document[strings.LastIndex(document, "/") : ] + "/"
    fmt.Println(path)
    if path == "" {
        hc.Error("paht is null")
        return Error
    }

    router := mux.NewRouter()

    file := http.StripPrefix(path, http.FileServer(http.Dir(path[1 : len(path) - 1])))

    router.PathPrefix(path).Handler(file)

    handler := &SwitchHandler{mux: router}
    http.Handle("/", handler)

    listener, error := net.Listen("tcp", hc.listen)
    if error != nil {
        return Error
    }

    hc.listener = listener

    return Ok
}
/*
func httpServer(args ...interface{}) int {
    var i interface{}

    l := len(args)
    switch l {
    case 1:
        i = args[l-1]
    default:

    }

    e := i.([]interface{})
    listener := e[0].(net.Listener)
    //name := e[1].(string)

    http.Serve(listener, nil)

    return Ok
}
*/

func httpServer(argv interface{}) int {
    if argv == nil {
        return Error
    }

    listener := argv.(net.Listener)

    http.Serve(listener, nil)

    return Ok
}

func (hc *HttpdCore) Run() int {
    if hc.Routine == nil {
        hc.Error("routine is null")
        return Error
    }

    if gid := hc.Routine.Go("httpServer", httpServer,
                            hc.listener); gid != -1 {

        hc.gid = append(hc.gid, gid)

    } else {
        return Error
    }

    select {

    }

    return Ok
}

func (hc *HttpdCore) Clear() int {
    if hc.listener != nil {
        hc.listener.Close()
    }

    return Ok
}

var httpdCore = String{ len("httpd_core"), "httpd_core" }
var coreHttpdContext = &Context{
    httpdCore,
    coreHttpdContextCreate,
    coreHttpdContextInit,
}

func coreHttpdContextCreate(cycle *Cycle) unsafe.Pointer {
    coreHttpd := NewHttpdCore()
    if coreHttpd == nil {
        return nil
    }

    coreHttpd.listen = "127.0.0.1:9800"
    coreHttpd.timeout = 3
    coreHttpd.location = nil

    return unsafe.Pointer(coreHttpd)
}

func coreHttpdContextInit(cycle *Cycle, context *unsafe.Pointer) string {
    log := cycle.GetLog()

    this := (*HttpdCore)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        log.Error("error")
        return "0"
    }

    coreHttpd = *this

    return "0"
}

var (
    listen   = String{ len("listen"), "listen" }
    timeout  = String{ len("timeout"), "timeout" }
    location = String{ len("location"), "location" }

    coreHttpd HttpdCore
)

var coreHttpdCommands = []Command{

    { listen,
      HTTPD_CONFIG|CONFIG_VALUE,
      SetString,
      0,
      unsafe.Offsetof(coreHttpd.listen),
      nil },

    { timeout,
      HTTPD_CONFIG|CONFIG_VALUE,
      SetNumber,
      0,
      unsafe.Offsetof(coreHttpd.timeout),
      nil },

    { location,
      HTTPD_CONFIG|CONFIG_VALUE,
      locationBlock,
      0,
      unsafe.Offsetof(coreHttpd.location),
      nil },

    NilCommand,
}

func locationBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if nil == cycle {
        return Error
    }

    flag := LOCATION_CONFIG|CONFIG_VALUE
    cycle.Block(cycle, LOCATION_MODULE, flag)

    return Ok
}

var coreHttpdModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(coreHttpdContext),
    coreHttpdCommands,
    HTTPD_MODULE,
    coreHttpdInit,
    coreHttpdMain,
    coreHttpdExit,
}

func coreHttpdInit(cycle *Cycle) int {
    coreHttpd.Cycle = cycle
    //coreHttpd.Log = cycle.Log
    //coreHttpd.Routine = cycle.Routine
    fmt.Println(coreHttpd.listen)

    if coreHttpd.location == nil {
        coreHttpd.location = &httpdLocation
    }

    return Ok
}

type SwitchHandler struct {
    mux http.Handler
}

func (s *SwitchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.mux.ServeHTTP(w, r)
}

func coreHttpdMain(cycle *Cycle) int {
    if coreHttpd.Init() == Error {
        cycle.Error("init error")
        return Error
    }

    if coreHttpd.Run() == Error {
        cycle.Error("run error")
        return Error
    }

    time.Sleep(1000)

    //coreHttpd.Routine.Check(coreHttpd.gid, BLOCKING)
    /*
        for _, id := range coreHttpd.gid {
            event := coreHttpd.Routine.GetEvent(id)
            if event != nil {
                select {

                case e := <-event.This:
                    op := e.GetOpcode()
                    switch op {
                    case RELOAD:
                        goto QUIT
                    }
                }
            }
        }

QUIT:
*/
    coreHttpd.Clear()

    return Ok
}

func coreHttpdExit(_ *Cycle) int {
    for _, id := range coreHttpd.gid {
        event := coreHttpd.Routine.GetEvent(id)
        if event != nil {
            event.SetOpcode(RELOAD)
            event.This <- event
        }
    }

    return Ok
}

func init() {
    Modules = Load(Modules, &coreHttpdModule)
}
