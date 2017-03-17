/*
 * Copyright (C) 2017 Meng Shi
 */

package grok

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

const (
    GROK_MODULE = 0x0007
    GROK_CONFIG = 0x00070000
)

var grokModule = String{ len("grok_module"), "grok_module" }
var filterGrokContext = &Context{
    grokModule,
    nil,
    nil,
}

var	grok = String{ len("grok"), "grok" }
var filterGrokCommands = []Command{

    { grok,
      USER_CONFIG|CONFIG_BLOCK,
      grokBlock,
      0,
      0,
      nil },

    NilCommand,
}

func grokBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if cycle == nil {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != GROK_MODULE {
            continue
        }

        module.CtxIndex++
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != GROK_MODULE {
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

    if configure.SetModuleType(GROK_MODULE) == Error {
        return Error
    }

    if configure.SetCommandType(GROK_CONFIG) == Error {
        return Error
    }

    if configure.Parse(cycle) == Error {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != GROK_MODULE {
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

var filterGrokModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(filterGrokContext),
    filterGrokCommands,
    FILTER_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &filterGrokModule)
}