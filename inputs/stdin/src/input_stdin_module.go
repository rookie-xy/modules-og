/*
 * Copyright (C) 2017 Meng Shi
 */

package stdin

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
    . "github.com/rookie-xy/worker/modules"
)

const (
    STDIN_MODULE = INPUT_MODULE|0x01000000
    STDIN_CONFIG = 0x00030000
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
    cycle.Configure.Block(STDIN_MODULE, STDIN_CONFIG)
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
    Modules = Load(Modules, &inputStdinModule)
}