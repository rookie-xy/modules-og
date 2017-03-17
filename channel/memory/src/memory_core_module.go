/*
 * Copyright (C) 2017 Meng Shi
 */

package memory_modules

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

type MemoryCore struct {
    *Channal
    *Cycle

     name     string
     size     int
}

func NewMemoryCore() *MemoryCore {
    return &MemoryCore{}
}

var memoryCore = String{ len("memory_core"), "memory_core" }
var coreMemoryContext = &Context{
    memoryCore,
    coreContextCreate,
    coreContextInit,
}

func coreContextCreate(cycle *Cycle) unsafe.Pointer {
    memoryCore := NewMemoryCore()
    if memoryCore == nil {
        return nil
    }

    memoryCore.name = "memory test"
    memoryCore.size = 1024

    return unsafe.Pointer(memoryCore)
}

func coreContextInit(cycle *Cycle, context *unsafe.Pointer) string {
    log := cycle.GetLog()
    this := (*MemoryCore)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        log.Error("coreStdinContextInit error")
        return "0"
    }

    return "0"
}

var (
    name = String{ len("name"), "name" }
    size = String{ len("size"), "size" }
    coreMemory MemoryCore
)

var coreMemoryCommands = []Command{

    { name,
      MEMORY_CONFIG,
      SetString,
      0,
      unsafe.Offsetof(coreMemory.name),
      nil },

    { size,
      MEMORY_CONFIG,
      SetNumber,
      0,
      unsafe.Offsetof(coreMemory.size),
      nil },

    NilCommand,
}

var coreMemoryModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(coreMemoryContext),
    coreMemoryCommands,
    MEMORY_MODULE,
    coreMemoryInit,
    coreMemoryMain,
}

func coreMemoryInit(cycle *Cycle) int {
    return Ok
}

func coreMemoryMain(cycle *Cycle) int {
    return Ok
}

func init() {
    Modules = append(Modules, &coreMemoryModule)
}