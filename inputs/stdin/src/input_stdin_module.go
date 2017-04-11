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
    STDIN_CONFIG = USER_CONFIG|CONFIG_ARRAY
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
      STDIN_CONFIG,
      stdinBlock,
      0,
      0,
      nil },

    NilCommand,
}

func stdinBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    if nil == cycle {
        return Error
    }

    flag := STDIN_CONFIG|CONFIG_VALUE
    cycle.Block(cycle, STDIN_MODULE, flag)

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