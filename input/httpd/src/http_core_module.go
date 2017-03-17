/*
 * Copyright (C) 2017 Meng Shi
 */

package modules

import (
      "net/http"
      "unsafe"

      "github.com/gorilla/mux"
    . "github.com/rookie-xy/worker/types"
    "fmt"
    "strings"
)

const (
    LOCATION_MODULE = 0x0005
    LOCATION_CONFIG = 0x00050000
)

type AbstractHttpCore struct {
    *AbstractCycle
    *AbstractFile

     listen    string
     timeout   int
     location  *AbstractLocationHttp
}

func NewHttpCore() *AbstractHttpCore {
    return &AbstractHttpCore{}
}

var httpCore = String{ len("http_core"), "http_core" }
var coreHttpContext = &AbstractContext{
    httpCore,
    coreHttpContextCreate,
    coreHttpContextInit,
}

func coreHttpContextCreate(cycle *AbstractCycle) unsafe.Pointer {
    coreHttp := NewHttpCore()
    if coreHttp == nil {
        return nil
    }

    coreHttp.listen = "127.0.0.1:9800"
    coreHttp.timeout = 3
    coreHttp.location = nil

    return unsafe.Pointer(coreHttp)
}

func coreHttpContextInit(cycle *AbstractCycle, context *unsafe.Pointer) string {
    log := cycle.GetLog()

    this := (*AbstractHttpCore)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        log.Error("error")
        return "0"
    }

    coreHttp = *this

    return "0"
}

var (
    listen   = String{ len("listen"), "listen" }
    timeout  = String{ len("timeout"), "timeout" }
    location = String{ len("location"), "location" }

    coreHttp  AbstractHttpCore
)

var coreHttpCommands = []Command{

    { listen,
      HTTP_CONFIG,
      SetString,
      0,
      unsafe.Offsetof(coreHttp.listen),
      nil },

    { timeout,
      HTTP_CONFIG,
      SetNumber,
      0,
      unsafe.Offsetof(coreHttp.timeout),
      nil },

    { location,
      HTTP_CONFIG,
      locationBlock,
      0,
      unsafe.Offsetof(coreHttp.location),
      nil },

    NilCommand,
}

func locationBlock(configure *AbstractConfigure, command *Command, cycle *AbstractCycle, config *unsafe.Pointer) string {
    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != LOCATION_MODULE {
            continue
        }

        module.CtxIndex++
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != LOCATION_MODULE {
            continue
        }

        context := (*AbstractContext)(unsafe.Pointer(module.Context))
        if context == nil {
            continue
        }

        if handle := context.Create; handle != nil {
            this := handle(cycle)
            if cycle.SetContext(module.Index, &this) == Error {
                return "0"
            }
        }
    }

    if configure.SetModuleType(LOCATION_MODULE) == Error {
        return "0"
    }

    if configure.SetCommandType(LOCATION_CONFIG) == Error {
        return "0"
    }

    if configure.Parse(cycle) == Error {
        return "0"
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != LOCATION_MODULE {
            continue
        }

        this := (*AbstractContext)(unsafe.Pointer(module.Context))
        if this == nil {
            continue
        }

        context := cycle.GetContext(module.Index)
        if context == nil {
            continue
        }

        if init := this.Init; init != nil {
            if init(cycle, context) == "-1" {
                return "0"
            }
        }
    }

    return "0"
}

var coreHttpModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(coreHttpContext),
    coreHttpCommands,
    HTTP_MODULE,
    coreHttpInit,
    coreHttpMain,
}

func coreHttpInit(cycle *AbstractCycle) int {
    fmt.Println(coreHttp.listen)
//    fmt.Println(coreHttp.timeout)

    if coreHttp.location == nil {
        coreHttp.location = &httpLocation
    }

//    fmt.Println(coreHttp.location.document)
//    fmt.Println(coreHttp.location.bufsize)

    return Ok
}

type SwitchHandler struct {
    mux http.Handler
}

func (s *SwitchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.mux.ServeHTTP(w, r)
}

func coreHttpMain(cycle *AbstractCycle) int {
    fmt.Println("http main")

    doccument := coreHttp.location.document
    path := doccument[strings.LastIndex(doccument, "/") : ] + "/"
    if path == "" {
        return Error
    }

    r := mux.NewRouter()

    fmt.Println(path[1 : len(path) - 1])

    s := http.StripPrefix(path, http.FileServer(http.Dir(path[1 : len(path) - 1])))

    r.PathPrefix(path).Handler(s)

    handler := &SwitchHandler{mux: r}
    http.Handle("/", handler)

    err := http.ListenAndServe(coreHttp.listen, nil)
    if err != nil {
        fmt.Println("ok")
    } else {
        fmt.Println("error")
    }

    return Ok
}

func init() {
    Modules = append(Modules, &coreHttpModule)
}
