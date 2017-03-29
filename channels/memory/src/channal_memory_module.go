/*
 * Copyright (C) 2017 Meng Shi
 */

package memory

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
    . "github.com/rookie-xy/worker/modules"
)

const (
    MEMORY_MODULE = CHANNEL_MODULE|0x01000000
    MEMORY_CONFIG = 0x00000001
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
    cycle.Configure.Block(MEMORY_MODULE, MEMORY_CONFIG)
    return Ok
}

var channalMemoryModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(channalMemoryContext),
    channalMemoryCommands,
    CHANNEL_MODULE,
    nil,
    nil,
}

func init() {
    Modules = Load(Modules, &channalMemoryModule)
}
