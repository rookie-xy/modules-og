/*
 * Copyright (C) 2017 Meng Shi
 */

package memory_modules

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

const (
    MEMORY_MODULE = 0x0003
    MEMORY_CONFIG = 0x00030000
)

var memoryModule = String{ len("memory_module"), "memory_module" }
var channalMemoryContext = &Context{
    memoryModule,
    nil,
    nil,
}

var	memory = String{ len("memory"), "memory" }
var channalMemoryCommands = []Command{

    { memory,
      USER_CONFIG|CONFIG_BLOCK,
      memoryBlock,
      0,
      0,
      nil },

    NilCommand,
}

func memoryBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if cycle == nil {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != MEMORY_MODULE {
            continue
        }

        module.CtxIndex++
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != MEMORY_MODULE {
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

    if configure.SetModuleType(MEMORY_MODULE) == Error {
        return Error
    }

    if configure.SetCommandType(MEMORY_CONFIG) == Error {
        return Error
    }

    if configure.Parse(cycle) == Error {
        return Error
    }

    for m := 0; Modules[m] != nil; m++ {
        module := Modules[m]
        if module.Type != MEMORY_MODULE {
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

var channalMemoryModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(channalMemoryContext),
    channalMemoryCommands,
    CHANNAL_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &channalMemoryModule)
}
