/*
 * Copyright (C) 2017 Meng Shi
 */

package stdout_modules

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

type StdoutCore struct {
    *Cycle
    *File

     status   bool
     channal  string
}

func NewStdoutCore() *StdoutCore {
    return &StdoutCore{}
}

var stdoutCore = String{ len("stdout_core"), "stdout_core" }
var coreStdoutContext = &Context{
    stdoutCore,
    coreStdoutContextCreate,
    coreStdoutContextInit,
}

func coreStdoutContextCreate(cycle *Cycle) unsafe.Pointer {
    stdoutCore := NewStdoutCore()
    if stdoutCore == nil {
        return nil
    }

    stdoutCore.status = false
    stdoutCore.channal = "zhangyue"

    return unsafe.Pointer(stdoutCore)
}

func coreStdoutContextInit(cycle *Cycle, context *unsafe.Pointer) string {
    log := cycle.GetLog()
    this := (*StdoutCore)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        log.Error("coreStdoutContextInit error")
        return "0"
    }

    return "0"
}

var (
    coreStatus = String{ len("status"), "status" }
    coreChannal = String{ len("channal"), "channal" }
    coreStdout StdoutCore
)

var coreStdoutCommands = []Command{

    { coreStatus,
      STDOUT_CONFIG,
      SetFlag,
      0,
      unsafe.Offsetof(coreStdout.status),
      nil },

    { coreChannal,
      STDOUT_CONFIG,
      SetString,
      0,
      unsafe.Offsetof(coreStdout.channal),
      nil },

    NilCommand,
}

var coreStdoutModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(coreStdoutContext),
    coreStdoutCommands,
    STDOUT_MODULE,
    coreStdoutInit,
    coreStdoutMain,
}

func coreStdoutInit(cycle *Cycle) int {
    return Ok
}

func coreStdoutMain(cycle *Cycle) int {
    return Ok
}

func init() {
    Modules = append(Modules, &coreStdoutModule)
}
