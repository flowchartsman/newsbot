package main

import (
    "log"
    "github.com/anaxagoras/toml"
    "os"
    "path/filepath"
)

var BinPath string;
func init() {
    if _, err := toml.DecodeFile("newsbot.conf", &config); err != nil {
        log.Fatal(err)
    }
    
    //TODO: Find a way to actually assign to the global
    path, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
    }
    BinPath = path
}

