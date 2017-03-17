/*
 * Copyright (C) 2017 Meng Shi
 */

package multiline

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

const (
    MULTILINE_MODULE = 0x0007
    MULTILINE_CONFIG = 0x00070000
)

var multilineModule = String{ len("multiline_module"), "multiline_module" }
var codecMultilineContext = &Context{
    multilineModule,
    nil,
    nil,
}

var	multiline = String{ len("multiline"), "multiline" }
var codecMultilineCommands = []Command{

    { multiline,
      USER_CONFIG|CONFIG_BLOCK,
      multilineBlock,
      0,
      0,
      nil },

    NilCommand,
}

func multilineBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if cycle == nil {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != MULTILINE_MODULE {
            continue
        }

        module.CtxIndex++
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != MULTILINE_MODULE {
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

    if configure.SetModuleType(MULTILINE_MODULE) == Error {
        return Error
    }

    if configure.SetCommandType(MULTILINE_CONFIG) == Error {
        return Error
    }

    if configure.Parse(cycle) == Error {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != MULTILINE_MODULE {
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

var codecMultilineModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(codecMultilineContext),
    codecMultilineCommands,
    CODEC_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &codecMultilineModule)
}