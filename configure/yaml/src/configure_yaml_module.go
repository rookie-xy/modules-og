/*
 * Copyright (C) 2017 Meng Shi
 */

package yaml

import (
      "gopkg.in/yaml.v2"
    . "github.com/rookie-xy/worker/types"
    "fmt"
)

type yamlConfigure struct {
    *Configure
     configure ConfigureIf
}

func NewYamlConfigure(configure *Configure) *yamlConfigure {
    return &yamlConfigure{ configure, nil }
}

func (yc *yamlConfigure) SetConfigure(configure *Configure) int {
    if configure == nil {
        return Error
    }

    yc.Configure = configure

    return Ok
}

func (so *yamlConfigure) GetConfigure() *Configure {
    return so.Configure
}

func (yc *yamlConfigure) Parser(in []byte, out interface{}) int {
    if err := yaml.Unmarshal(in, out); err != nil {
        fmt.Println(err)
				    return Error
    }

    return Ok
}

func initYamlConfigureModule(cycle *Cycle) int {
    log := cycle.Log

    configure := cycle.GetConfigure()
    if configure == nil {
        configure = NewConfigure(log)
        if configure == nil {
            return Error
        }

        if cycle.SetConfigure(configure) == Error {
            return Error
        }
    }

    yamlConfigure := NewYamlConfigure(configure)
    if yamlConfigure == nil {
        log.Error("new yaml configure error")
        return Error
    }

    if configure.Set(yamlConfigure) == Error {
        log.Error("set configure interface error")
        return Error
    }

    return Ok
}

var YamlConfigureModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    nil,
    nil,
    SYSTEM_MODULE,
    initYamlConfigureModule,
    nil,
}

func init() {
    Modules = append(Modules, &YamlConfigureModule)
}
