/*
 * Copyright (C) 2017 Meng Shi
 */

package file

import (
      "strings"
      "github.com/fsnotify/fsnotify"

    . "github.com/rookie-xy/worker/types"

    "unsafe"
    "log"
    "fmt"
)

const (
    RESOURCE = "/data/service"
    FILENAME = "configure"
)

type FileConfigure struct {
    *Configure

     watcher     *fsnotify.Watcher

     resource     string
     fileName     string

     Notice       chan *Event
}

func NewFileConfigure(configure *Configure) *FileConfigure {
    return &FileConfigure{
        Configure: configure,
        resource:  RESOURCE,
        fileName:  FILENAME,
        Notice:    make(chan *Event),
    }
}

var fileConfig = String{ len("file_configure"), "file_configure" }
var fileConfigureCtx = &Context{
    fileConfig,
    nil,
    nil,
}

func (fc *FileConfigure) SetConfigure(configure *Configure) int {
    if configure == nil {
        return Error
    }

    fc.Configure = configure

    return Ok
}

func (fc *FileConfigure) GetConfigure() *Configure {
    return fc.Configure
}

func (fc *FileConfigure) SetResource(resource string) int {
    if resource == "" {
        return Error
    }

    fc.resource = resource

    return Ok
}

func (fc *FileConfigure) GetResource() string {
    if fc.resource == "" {
        return ""
    }

    return fc.resource
}

func (fc *FileConfigure) SetFileName(fileName string) int {
    if fileName == "" {
        return Error
    }

    fc.fileName = fileName

    return Ok
}

func (fc *FileConfigure) GetFileName() string {
    return fc.fileName
}

func (fc *FileConfigure) Set() int {
    if fc.AbstractFile == nil {
        fc.Error("file configure set error")
        return Error
    }

    flag := false
    if fc.Reader() == Error {
        fc.Error("configure read file error")
        flag = true
    }

    if fc.Closer() == Error {
        fc.Warn("file close error: %d\n", 10)
        return Error
    }

    if flag {
        return Error
    }

    return Ok
}

func (fc *FileConfigure) Get() int {
    if fc.AbstractFile == nil {
        fc.AbstractFile = NewAbstractFile(fc.Log)
    }

    resource := fc.GetResource()
    if resource == "" {
        return Error
    }

    if fc.Open(resource) == Error {
        fc.Error("configure open file error")
        return Error
    }

    return Ok
}

func (fc *FileConfigure) GetType() unsafe.Pointer {
    return unsafe.Pointer(fc)
}

func (fc *FileConfigure) Clear() {
    return
}

var fileConfigure *FileConfigure

func fileConfigureInit(cycle *Cycle) int {
    option := cycle.GetOption()
    if option == nil {
        return Error
    }

    item := option.GetItem("configure")
    if item == nil {
        return Error
    }

    file := item.(string)

    fileType := file[0 : strings.Index(file, ":")]
    if fileType == "" {
        return Error
    }

    if fileType != "file" {
        return Ignore
    }

    configure := cycle.GetConfigure()
    if configure == nil {
        if configure = NewConfigure(cycle); configure == nil {
            return Error
        }
    }

    fileConfigure = NewFileConfigure(configure)
    if fileConfigure == nil {
        return Error
    }

    if watcher, error := fsnotify.NewWatcher(); error != nil {
        return Error
    } else {
        fileConfigure.watcher = watcher
    }

    if fileConfigure.SetName(fileType) == Error {
        return Error
    }

    fileName := file[strings.LastIndex(file, "/") + 1 : ]
    if fileName == "" {
        return Error
    }

    if fileConfigure.SetFileName(fileName) == Error {
        return Error
    }

    resource := file[strings.Index(file, "=") + 1 : ]
    if resource == "" {
        return Error
    }

    if fileConfigure.SetResource(resource) == Error {
        return Error
    }

    if error := fileConfigure.watcher.Add(resource); error != nil {
        return Error
    }

    if cycle.SetConfigure(configure) == Error {
        return Error
    }

    if configure.SetHandle(fileConfigure) == Error {
        return Error
    }

    return Ok
}

//func fileConfigureMain(cycle *Cycle) int {
func fileConfigureMain(args ...interface{}) int {
    cycle := args.(Cycle)
    flag := Error

    configure := cycle.GetConfigure()
    if configure == nil {
        return flag
    }

    content := configure.GetHandle()
    if content == nil {
        return flag
    }

    if content.Get() == Error {
        return flag
    }

    if content.Set() == Error {
        return flag
    }

    notice := NewEvent()

    if flag == Error {
        notice.SetOpcode(LOAD)
        notice.SetName("load")
        configure.Event <- notice
    }

    fcp := content.GetType()
    if fcp == nil {
        return Error
    }

    fc := (*FileConfigure)(unsafe.Pointer(uintptr(fcp)))
    if fc == nil {
        return Error
    }

    defer fc.watcher.Close()

    quit := false

    for {
        select {

        case event := <-fc.watcher.Events:
            if event.Op & fsnotify.Write == fsnotify.Write {
                fmt.Println("mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm")
                notice.SetOpcode(RELOAD)
                notice.SetName("reload")
                configure.Event <- notice
            }

        case err := <-fc.watcher.Errors:
            log.Println("error:", err)
        /*
        case e := <-fc.Notice.GetNotice():
            if op := e.GetOpcode(); op == SYSTEM_MODULE {
                quit = true
            }
            */
        }

        if quit {
            break
        }
    }

    fc.Clear()

    return Ok
}

func fileConfigureExit(cycle *Cycle) int {
    //fileConfigure.Event <- 1

    fileConfigure.Clear()
//    fileConfigure.Quit()
    return Ok
}

var FileConfigureModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(fileConfigureCtx),
    nil,
    SYSTEM_MODULE,
    fileConfigureInit,
    fileConfigureMain,
    fileConfigureExit,
}

func init() {
    Modules = Load(Modules, &FileConfigureModule)
}