/*
 * Copyright (C) 2017 Meng Shi
 */

package httpd

import (
      "net/http"
      "unsafe"

      "github.com/gorilla/mux"
    . "github.com/rookie-xy/worker/types"
    "fmt"
    "strings"
)

const (
    LOCATION_MODULE = HTTPD_MODULE|0x00100000
    LOCATION_CONFIG = LOCATION_MODULE|0x00000001
)

type HttpdCore struct {
    *Cycle
    *File

     listen    string
     timeout   int
     location  *LocationHttpd
}

func NewHttpdCore() *HttpdCore {
    return &HttpdCore{}
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

    coreHttpd  HttpdCore
)

var coreHttpdCommands = []Command{

    { listen,
      HTTPD_CONFIG,
      SetString,
      0,
      unsafe.Offsetof(coreHttpd.listen),
      nil },

    { timeout,
      HTTPD_CONFIG,
      SetNumber,
      0,
      unsafe.Offsetof(coreHttpd.timeout),
      nil },

    { location,
      HTTPD_CONFIG|CONFIG_BLOCK,
      locationBlock,
      0,
      unsafe.Offsetof(coreHttpd.location),
      nil },

    NilCommand,
}

func locationBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    cycle.Configure.Block(LOCATION_MODULE, LOCATION_CONFIG)
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
}

func coreHttpdInit(cycle *Cycle) int {
//    fmt.Println(coreHttpd.listen)
//    fmt.Println(coreHttpd.timeout)

    if coreHttpd.location == nil {
        coreHttpd.location = &httpdLocation
    }

//    fmt.Println(coreHttpd.location.document)
//    fmt.Println(coreHttpd.location.bufsize)

    return Ok
}

type SwitchHandler struct {
    mux http.Handler
}

func (s *SwitchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.mux.ServeHTTP(w, r)
}

func coreHttpdMain(cycle *Cycle) int {
    fmt.Println("httpd main")

    log := cycle.Log
    fmt.Println(log.GetPath())

    document := coreHttpd.location.document
    path := document[strings.LastIndex(document, "/") : ] + "/"
    if path == "" {
        return Error
    }

    r := mux.NewRouter()

    fmt.Println(path[1 : len(path) - 1])

    s := http.StripPrefix(path, http.FileServer(http.Dir(path[1 : len(path) - 1])))

    r.PathPrefix(path).Handler(s)

    handler := &SwitchHandler{mux: r}
    http.Handle("/", handler)

    err := http.ListenAndServe(coreHttpd.listen, nil)
    if err != nil {
        fmt.Println("ok")
    } else {
        fmt.Println("error")
    }

    return Ok
}

func init() {
    Modules = Load(Modules, &coreHttpdModule)
}
