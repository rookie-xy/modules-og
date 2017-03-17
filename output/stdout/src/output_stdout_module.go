/*
 * Copyright (C) 2017 Meng Shi
 */

package stdout_modules

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

const (
    STDOUT_MODULE = 0x0004
    STDOUT_CONFIG = 0X00040000
)

var stdoutModule = String{ len("stdout_module"), "stdout_module" }
var outputStdoutContext = &Context{
    stdoutModule,
    nil,
    nil,
}

var stdout = String{ len("stdout"), "stdout" }
var outputStdoutCommands = []Command{

    { stdout,
      USER_CONFIG|CONFIG_BLOCK,
      stdoutBlock,
      0,
      0,
      nil },

    NilCommand,
}

func stdoutBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if cycle == nil {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != STDOUT_MODULE {
            continue
        }

        module.CtxIndex++
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != STDOUT_MODULE {
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

    if configure.SetModuleType(STDOUT_MODULE) == Error {
				    return Error
    }

    if configure.SetCommandType(STDOUT_CONFIG) == Error {
				    return Error
    }

    if configure.Parse(cycle) == Error {
				    return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != STDOUT_MODULE {
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

var outputStdoutModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(outputStdoutContext),
    outputStdoutCommands,
    OUTPUT_MODULE,
    nil,
    nil,
}

func init() {
   Modules = append(Modules, &outputStdoutModule)
}