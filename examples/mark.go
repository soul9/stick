package main

import (
    "flag"
    "os"
    "fmt"
    "../_obj/stick"
)

func main() {
    confpath := flag.String("c", "conf.json", "Configuration file path")
    flag.Parse()
    conn, err := stick.Init(confpath)
    if err != nil {
        fmt.Println(conn, err.String())
        os.Exit(1)
    }
    os.Exit(0)
}
