/*
 * Copyright (C) 2017 Meng Shi
 */

package mlog

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
    "fmt"
)

const (
    MLOG_MODULE = 0x0002
    MLOG_CONFIG = 0x00020000
)

var mlogModule = String{ len("mlog_module"), "mlog_module" }
var logMlogContext = &Context{
    mlogModule,
    nil,
    nil,
}

var mlog = String{ len("mlog"), "mlog" }
var logMlogCommands = []Command{

    { mlog,
      USER_CONFIG|CONFIG_BLOCK,
      mlogBlock,
      0,
      0,
      nil },

    NilCommand,
}

func mlogBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if cycle == nil {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != MLOG_MODULE {
            continue
        }

        module.CtxIndex++
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != MLOG_MODULE {
            continue
        }

        context := (*Context)(unsafe.Pointer(module.Context))
        if context == nil {
            continue
        }

        if handle := context.Create; handle != nil {
            this := handle(cycle)
            fmt.Println(module.Index)
            if cycle.SetContext(module.Index, &this) == Error {
												    return Error
            }
        }
    }

    configure := cycle.GetConfigure()
    if configure == nil {
        return Error
    }

    if configure.SetModuleType(MLOG_MODULE) == Error {
				    return Error
    }

    if configure.SetCommandType(MLOG_CONFIG) == Error {
				    return Error
    }

    if configure.Materialized(cycle) == Error {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != MLOG_MODULE {
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

var logMlogModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(logMlogContext),
    logMlogCommands,
    LOG_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &logMlogModule)
}