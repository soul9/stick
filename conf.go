package stick

import (
  "json"
  "os"
  "fmt"
)

type ChanActConf struct {
    Match string "match"
    Action string "action"
    Parms string "parms"
}

type ChanConf struct {
    Name string "name"
    Actions map[string]ChanActConf "actions"
}
func (c ChanConf) String() string {
    return string(c.Name)
}

type NetConf struct {
    Nick string "nick"
    Realname string "realname"
    Channels []ChanConf "channels"
}

type NetsConf map[string] NetConf

type PlugConf []string
type ActConf []string

type Conf struct {
    Networks NetsConf "networks"
    Plugins PlugConf "plugins"
    Actions ActConf "actions"
}


func parseconf (conffile string) Conf {
    f, err := os.Open(conffile, os.O_RDONLY, 0)
    if err != nil {
        fmt.Println("An error occurred while opening the configuration file for reading: ", err.String())
    }
    defer f.Close()
    d := json.NewDecoder(f)
    var c Conf
    err = d.Decode(&c);
    if  err != nil {
        fmt.Println("An error occurred while parsing the config file: ", err.String())
    }
    return c
}
