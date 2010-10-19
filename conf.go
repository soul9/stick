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

type NetsConf map[string] *NetConf

type PlugConf []string
type ActConf []string

type StickConf struct {
    Networks NetsConf "networks"
    Plugins PlugConf "plugins"
    Actions ActConf "actions"
}


func InitConf (conffile string) (*StickConf, os.Error) {
    var c StickConf
    return c.ReadConf(conffile)
}

func (c *StickConf) ReadConf(conffile string) (*StickConf, os.Error) {
    f, err := os.Open(conffile, os.O_RDONLY, 0)
    if err != nil {
        fmt.Println("An error occurred while opening the configuration file for reading: ", err.String())
    }
    defer f.Close()
    d := json.NewDecoder(f)
    err = d.Decode(&c);
    if  err != nil {
        return c, os.NewError("An error occurred while parsing the config file: " + err.String())
    }
    return c, nil
}