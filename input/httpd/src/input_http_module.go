/*
 * Copyright (C) 2017 Meng Shi
 */

package modules

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

const (
    HTTP_MODULE = 0x0001
    HTTP_CONFIG = 0x00010000
)

var httpModule = String{ len("http_module"), "http_module" }
var inputHttpContext = &AbstractContext{
    httpModule,
    nil,
    nil,
}

var httpd = String{ len("httpd"), "httpd" }
var inputHttpCommands = []Command{

    { httpd,
      USER_CONFIG|CONFIG_BLOCK,
      httpBlock,
      0,
      0,
      nil },

    NilCommand,
}

func httpBlock(configure *AbstractConfigure, command *Command, cycle *AbstractCycle, config *unsafe.Pointer) string {
    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != HTTP_MODULE {
            continue
        }

        module.CtxIndex++
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != HTTP_MODULE {
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

    if configure.SetModuleType(HTTP_MODULE) == Error {
        return "0"
    }

    if configure.SetCommandType(HTTP_CONFIG) == Error {
        return "0"
    }

    if configure.Parse(cycle) == Error {
        return "0"
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != HTTP_MODULE {
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

var inputHttpModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(inputHttpContext),
    inputHttpCommands,
    INPUT_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &inputHttpModule)
}
