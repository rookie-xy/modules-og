/*
 * Copyright (C) 2017 Meng Shi
 */

package modules

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

type AbstractLocationHttp struct {
    document string
    bufsize  int
}

func NewLocationHttp() *AbstractLocationHttp {
    return &AbstractLocationHttp{}
}

var httpLocationContext = &AbstractContext{
    location,
    httpLocationContextCreate,
    httpLocationContextInit,
}

func httpLocationContextCreate(cycle *AbstractCycle) unsafe.Pointer {
    httpLocation := NewLocationHttp()
    if httpLocation == nil {
        return nil
    }

    httpLocation.document = "/data/service/http/mengshi"
    httpLocation.bufsize = 256

    return unsafe.Pointer(httpLocation)
}

func httpLocationContextInit(cycle *AbstractCycle, context *unsafe.Pointer) string {
    log := cycle.GetLog()
    this := (*AbstractLocationHttp)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        log.Error("coreStdinContextInit error")
        return "0"
    }

    httpLocation = *this

    return "0"
}

var (
    document = String{ len("document"), "document" }
    bufsize  = String{ len("bufsize"), "bufsize" }

    httpLocation AbstractLocationHttp
)

var httpLocationCommands = []Command{

    { document,
      LOCATION_CONFIG,
      SetString,
      0,
      unsafe.Offsetof(httpLocation.document),
      nil },

    { bufsize,
      LOCATION_CONFIG,
      SetNumber,
      0,
      unsafe.Offsetof(httpLocation.bufsize),
      nil },

    NilCommand,
}

var httpLocationModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(httpLocationContext),
    httpLocationCommands,
    LOCATION_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &httpLocationModule)
}