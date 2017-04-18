/*
 * Copyright (C) 2017 Meng Shi
 */

package yaml

import (
      "gopkg.in/yaml.v2"
    . "github.com/rookie-xy/worker/types"
)

type yamlConfigure struct {
    *Configure
}

func NewYamlConfigure(configure *Configure) *yamlConfigure {
    return &yamlConfigure{ configure }
}

func (yc *yamlConfigure) SetConfigure(configure *Configure) int {
    if configure == nil {
        return Error
    }

    yc.Configure = configure

    return Ok
}

func (yc *yamlConfigure) GetConfigure() *Configure {
    return yc.Configure
}

func (yc *yamlConfigure) Marshal(in interface{}) ([]byte, error) {
    return nil, nil
}

func (yc *yamlConfigure) Unmarshal(in []byte, out interface{}) int {
    if err := yaml.Unmarshal(in, out); err != nil {
        return Error
    }

    return Ok
}

func yamlConfigureInit(cycle *Cycle) int {
    log := cycle.Log

    configure := cycle.GetConfigure()
    if configure == nil {
        configure = NewConfigure(cycle)
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

    if configure.SetParser(yamlConfigure) == Error {
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
    yamlConfigureInit,
    nil,
    nil,
}

func init() {
    Modules = Load(Modules, &YamlConfigureModule)
}
