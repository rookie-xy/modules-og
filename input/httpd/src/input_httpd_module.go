/*
 * Copyright (C) 2017 Meng Shi
 */

package httpd

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

const (
    HTTPD_MODULE = 0x0001
    HTTPD_CONFIG = 0x00010000
)

var httpdModule = String{ len("httpd_module"), "httpd_module" }
var inputHttpdContext = &Context{
    httpdModule,
    nil,
    nil,
}

var httpdd = String{ len("httpd"), "httpd" }
var inputHttpdCommands = []Command{

    { httpdd,
      USER_CONFIG|CONFIG_BLOCK,
      httpdBlock,
      0,
      0,
      nil },

    NilCommand,
}

func httpdBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if cycle == nil {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != HTTPD_MODULE {
            continue
        }

        module.CtxIndex++
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != HTTPD_MODULE {
            continue
        }

        context := (*Context)(unsafe.Pointer(module.Context))
        if context == nil {
            continue
        }

        if handle := context.Create; handle != nil {
            this := handle(cycle)
            if cycle.SetContext(module.Index, &this) == Error {
                return Error
            }
        }
    }

    configure := cycle.GetConfigure()
    if configure == nil {
        return Error
    }

    if configure.SetModuleType(HTTPD_MODULE) == Error {
        return Error
    }

    if configure.SetCommandType(HTTPD_CONFIG) == Error {
        return Error
    }

    if configure.Parse(cycle) == Error {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != HTTPD_MODULE {
            continue
        }

        this := (*Context)(unsafe.Pointer(module.Context))
        if this == nil {
            continue
        }

        context := cycle.GetContext(module.Index)
        if context == nil {
            continue
        }

        if init := this.Init; init != nil {
            if init(cycle, context) == "-1" {
                return Error
            }
        }
    }

    return Ok
}

var inputHttpdModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(inputHttpdContext),
    inputHttpdCommands,
    INPUT_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &inputHttpdModule)
}
