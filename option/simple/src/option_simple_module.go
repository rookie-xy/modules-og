/*
 * Copyright (C) 2017 Meng Shi
 */

package simple

import (
    . "github.com/rookie-xy/worker/types"
)

type simpleOption struct {
    *Option

    option OptionIf
}

func NewSimpleOption() *simpleOption {
    return &simpleOption{ nil, nil }
}

func (so *simpleOption) SetOption(option *Option) int {
    if option == nil {
        return Error
    }

    so.Option = option

    return Ok
}

func (so *simpleOption) GetOption() *Option {
    return so.Option
}

func (o *simpleOption) Parse() int {
    log := o.Log

    argv := o.GetArgv()

    for i := 1; i < o.GetArgc(); i++ {

        if argv[i][0] != '-' {
            return Error
        }

        switch argv[i][1] {

        case 'c':
	    if argv[i + 1] == "" {
	        return Error
	    }

            // file://path=/home/
            o.SetItem("configure", "file://resource=" + argv[i + 1])
            i++
            break

        case 'z':
	    if argv[i + 1] == "" {
	        return Error
	    }

            // file://path=/home/
            o.SetItem("configure", "zookeeper://resource=" + argv[i + 1])
            i++
            break

        case 't':
            o.SetItem("test", true)
	    break

        default:
            o.SetItem("invaild", "")
            log.Info("not found any option")
            //o.result["error"] = "not found any option"
            break
        }
    }

    return Ok
}

func initSimpleOptionModule(cycle *Cycle) int {


    option := cycle.GetOption()
    if option == nil {
        return Error
    }

    simpleOption := NewSimpleOption()

    //simpleOption.SetOptionIf(simpleOption)
    log := option.Log

    if simpleOption.SetOption(option) == Error {
        log.Error("set option error")
        return Error
    }

    simpleOption.option = simpleOption
/*
    fmt.Println("menggggggggggggggggggggggggggggggg")
    if simpleOption.option.Parse() == Error {
        return Error
    }
    */

    option.SetOptionIf(simpleOption.option)

    return Ok
}

var SimpleOptionModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    nil,
    nil,
    SYSTEM_MODULE,
    initSimpleOptionModule,
    nil,
}

func init() {
    Modules = append(Modules, &SimpleOptionModule)
}
