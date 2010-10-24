package stick

import (
  "json"
  "os"
  "strconv"
)

type ChanActConf struct {
    Match string "match"
    Action string "action"
    Parms string "parms"
}

type ChanConf struct {
    Name string "name"
    Actions map[string]ChanActConf "actions"
    Users map[string]UserConf "users"
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

type UserConf struct {
    Role string "role"
    Sha1pass string "sha1pass"
    Email string "e-mail"
}

type StickConfJson struct {
    Networks NetsConf "networks"
    Plugins PlugConf "plugins"
    Actions ActConf "actions"
    Users map[string]UserConf "users"
}

type StickConf struct {
    Conf StickConfJson
    conffile string
}

func InitConf (conffile string) (*StickConf, os.Error) {
    var c StickConf
    c.conffile = conffile
    err := c.RdConf()
    return &c, err
}

func (c *StickConf) RdConf() os.Error {
    f, err := os.Open(c.conffile, os.O_RDONLY, 0)
    if err != nil {
        return os.NewError("An error occurred while parsing the config file: " + err.String())
    }
    defer f.Close()
    d := json.NewDecoder(f)
    err = d.Decode(&c.Conf)
    if  err != nil {
        return os.NewError("An error occurred while parsing the config file: " + err.String())
    }
    return nil
}

func (c *StickConf) WrConf() os.Error {
    f, err := os.Open(c.conffile, os.O_WRONLY, 600)
    if err != nil {
        return os.NewError("An error occurred while parsing the config file: " + err.String())
    }
    defer f.Close()
    jsonbytes, err := json.MarshalIndent(c.Conf, "", "    ")
    if  err != nil {
        return os.NewError("An error occurred while marshaling struct: " + err.String())
    }
    _, err = f.Write(jsonbytes)
    if  err != nil {
        return os.NewError("An error occurred while parsing the config file: " + err.String())
    }
    return nil
}

func (c *StickConf) AddActionChan(ch string, action ChanActConf, net string) os.Error {
    var i int
    for i, _ =  range c.Conf.Networks[net].Channels {
        if c.Conf.Networks[net].Channels[i].Name == ch {
            break
        }
    }
    c.Conf.Networks[net].Channels[i].Actions[strconv.Itoa(len(c.Conf.Networks[net].Channels[i].Actions)+1)] = action
    err := c.WrConf()
    if err != nil {
        return os.NewError("Couldn't write config to file: " + err.String())
    }
    return nil
}