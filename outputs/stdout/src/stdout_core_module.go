/*
 * Copyright (C) 2017 Meng Shi
 */

package stdout

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
"fmt"
)

type StdoutCore struct {
    *Cycle
    *File

     status   bool
     channel  string
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
    stdoutCore.channel = "zhangyue"

    return unsafe.Pointer(stdoutCore)
}

func coreStdoutContextInit(cycle *Cycle, context *unsafe.Pointer) string {
    log := cycle.GetLog()
    this := (*StdoutCore)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        log.Error("coreStdoutContextInit error")
        return "0"
    }

    fmt.Println(this.channel)

    return "0"
}

var (
    coreStatus = String{ len("status"), "status" }
    coreChannel = String{ len("channel"), "channel" }
    coreStdout StdoutCore
)

var coreStdoutCommands = []Command{

    { coreStatus,
      STDOUT_CONFIG|CONFIG_VALUE,
      SetFlag,
      0,
      unsafe.Offsetof(coreStdout.status),
      nil },

    { coreChannel,
      STDOUT_CONFIG|CONFIG_VALUE,
      SetString,
      0,
      unsafe.Offsetof(coreStdout.channel),
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
    Modules = Load(Modules, &coreStdoutModule)
}
