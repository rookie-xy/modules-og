/*
 * Copyright (C) 2017 Meng Shi
 */

package memory

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
"fmt"
)

type MemoryCore struct {
    *Channel
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

    fmt.Println(this.name)

    return "0"
}

var (
    name = String{ len("name"), "name" }
    size = String{ len("size"), "size" }
    coreMemory MemoryCore
)

var coreMemoryCommands = []Command{

    { name,
      MEMORY_CONFIG|CONFIG_VALUE,
      SetString,
      0,
      unsafe.Offsetof(coreMemory.name),
      nil },

    { size,
      MEMORY_CONFIG|CONFIG_VALUE,
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
    coreMemoryExit,
}

func coreMemoryInit(cycle *Cycle) int {
    return Ok
}

func coreMemoryMain(cycle *Cycle) int {
    return Ok
}

func coreMemoryExit(cycle *Cycle) int {
    return Ok
}

func init() {
    Modules = Load(Modules, &coreMemoryModule)
}