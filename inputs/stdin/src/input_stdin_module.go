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
    STDIN_MODULE = 0x000000000001
    STDIN_CONFIG = 0x00020000
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
    cycle.Configure.Block(INPUT_MODULE|STDIN_MODULE, STDIN_CONFIG)
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