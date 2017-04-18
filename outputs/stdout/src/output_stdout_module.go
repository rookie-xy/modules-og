/*
 * Copyright (C) 2017 Meng Shi
 */

package stdout

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
    . "github.com/rookie-xy/worker/modules"
)

const (
    STDOUT_MODULE = OUTPUT_MODULE|0x01000000
    STDOUT_CONFIG = USER_CONFIG|CONFIG_ARRAY
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
      STDOUT_CONFIG,
      stdoutBlock,
      0,
      0,
      nil },

    NilCommand,
}

func stdoutBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if nil == cycle {
        return Error
    }

    flag := STDOUT_CONFIG|CONFIG_VALUE
    cycle.Block(cycle, STDOUT_MODULE, flag)

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
    nil,
}

func init() {
    Modules = Load(Modules, &outputStdoutModule)
}