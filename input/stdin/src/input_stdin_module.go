/*
 * Copyright (C) 2017 Meng Shi
 */

package stdin_modules

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
    "fmt"
)

const (
    STDIN_MODULE = 0x0002
    STDIN_CONFIG = 0x00020000
)

var stdinModule = String{ len("stdin_module"), "stdin_module" }
var inputStdinContext = &Context{
    stdinModule,
    nil,
    nil,
}

var stdin = String{ len("stdin"), "stdin" }
var inputStdinCommands = []Command{

    { stdin,
      USER_CONFIG|CONFIG_BLOCK,
      stdinBlock,
      0,
      0,
      nil },

    NilCommand,
}

func stdinBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if cycle == nil {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != STDIN_MODULE {
            continue
        }

        module.CtxIndex++
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != STDIN_MODULE {
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

    if configure.SetModuleType(STDIN_MODULE) == Error {
				    return Error
    }

    if configure.SetCommandType(STDIN_CONFIG) == Error {
				    return Error
    }

    if configure.Parse(cycle) == Error {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != STDIN_MODULE {
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

var inputStdinModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(inputStdinContext),
    inputStdinCommands,
    INPUT_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &inputStdinModule)
}