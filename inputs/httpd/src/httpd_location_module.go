/*
 * Copyright (C) 2017 Meng Shi
 */

package httpd

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

type LocationHttpd struct {
    document string
    bufsize  int
}

func NewLocationHttpd() *LocationHttpd {
    return &LocationHttpd{}
}

var httpdLocationContext = &Context{
    location,
    httpdLocationContextCreate,
    httpdLocationContextInit,
}

func httpdLocationContextCreate(cycle *Cycle) unsafe.Pointer {
    httpdLocation := NewLocationHttpd()
    if httpdLocation == nil {
        return nil
    }

    httpdLocation.document = "/data/service/httpd/mengshi"
    httpdLocation.bufsize = 256

    return unsafe.Pointer(httpdLocation)
}

func httpdLocationContextInit(cycle *Cycle, context *unsafe.Pointer) string {
    log := cycle.GetLog()

    this := (*LocationHttpd)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        log.Error("coreStdinContextInit error")
        return "0"
    }

    httpdLocation = *this

    return "0"
}

var (
    document = String{ len("document"), "document" }
    bufsize  = String{ len("bufsize"), "bufsize" }

    httpdLocation LocationHttpd
)

var httpdLocationCommands = []Command{

    { document,
      LOCATION_CONFIG,
      SetString,
      0,
      unsafe.Offsetof(httpdLocation.document),
      nil },

    { bufsize,
      LOCATION_CONFIG,
      SetNumber,
      0,
      unsafe.Offsetof(httpdLocation.bufsize),
      nil },

    NilCommand,
}

var httpdLocationModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(httpdLocationContext),
    httpdLocationCommands,
    HTTPD_MODULE|LOCATION_MODULE,
    nil,
    nil,
}

func init() {
    Modules = Load(Modules, &httpdLocationModule)
}