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
    STDOUT_CONFIG = STDOUT_MODULE|0X00000001
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
    cycle.Configure.Block(STDOUT_MODULE, STDOUT_CONFIG)
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
   Modules = Load(Modules, &outputStdoutModule)
}