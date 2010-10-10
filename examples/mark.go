package main

import (
    "flag"
    "os"
    "../_obj/stick"
)

func main() {
    confpath := flag.String("c", "../conf.json", "Configuration file path")
    flag.Parse()
    stick.Init(confpath)
    os.Exit(0)
}
