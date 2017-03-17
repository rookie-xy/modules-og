/*
 * Copyright (C) 2017 Meng Shi
 */

package stdin_modules

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
"fmt"
)

type StdinCore struct {
    *Cycle
    *File

     status   bool
     channal  string
}

func NewStdinCore() *StdinCore {
    return &StdinCore{}
}

var stdinCore = String{ len("stdin_core"), "stdin_core" }
var coreStdinContext = &Context{
    stdinCore,
    coreStdinContextCreate,
    coreStdinContextInit,
}

func coreStdinContextCreate(cycle *Cycle) unsafe.Pointer {
    stdinCore := NewStdinCore()
    if stdinCore == nil {
        return nil
    }

    stdinCore.status = false
    stdinCore.channal = "mengshi"

    return unsafe.Pointer(stdinCore)
}

func coreStdinContextInit(cycle *Cycle, context *unsafe.Pointer) string {
    log := cycle.GetLog()
    this := (*StdinCore)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        log.Error("coreStdinContextInit error")
        return "0"
    }

    return "0"
}

var (
    status = String{ len("status"), "status" }
    channal = String{ len("channal"), "channal" }
    coreStdin StdinCore
)

var coreStdinCommands = []Command{

    { status,
      STDIN_CONFIG,
      SetFlag,
      0,
      unsafe.Offsetof(coreStdin.status),
      nil },

    { channal,
      STDIN_CONFIG,
      SetString,
      0,
      unsafe.Offsetof(coreStdin.channal),
      nil },

    NilCommand,
}

var coreStdinModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(coreStdinContext),
    coreStdinCommands,
    STDIN_MODULE,
    coreStdinInit,
    coreStdinMain,
}

<<<<<<< HEAD
func coreStdinInit(cycle *AbstractCycle) int {
    return Ok
}

func coreStdinMain(cycle *AbstractCycle) int {
    fmt.Println("stdin main")
=======
func coreStdinInit(cycle *Cycle) int {
    /*
    context := cycle.GetContext(10)
    if context == nil {
        return Error
    }

    this := (*AbstractStdinCore)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        return Error
    }

    fmt.Println(this.channal)
    */

    return Ok
}

func coreStdinMain(cycle *Cycle) int {
>>>>>>> 6a7781c2aa20c39bae6abd9cfb65bdcb3653261a
    return Ok
}

func init() {
    Modules = append(Modules, &coreStdinModule)
}