package main

import (
	irc "github.com/thoj/Go-IRC-Client-Library"
	"fmt"
	"os"
	"json"
	"flag"
	"reflect"
)

func parseconf (conffile string) map[string]interface{} {
    f, err := os.Open(conffile, os.O_RDONLY, 0)
    if err != nil {
        fmt.Println("An error occurred while opening the configuration file for reading: ", err.String())
    }
    defer f.Close()
    d := json.NewDecoder(f)
    c := make(map[string]interface{})
    err = d.Decode(&c);
    if  err != nil {
        fmt.Println("An error occurred while parsing the config file: ", err.String())
    }
    fmt.Printf("%#V\n", c["networks"])
    for key, val := range c["networks"].(map[string]interface{}) {
        fmt.Printf("%s :\n", key)
        for key, val := range val.(map[string]interface{}) {
            switch t := val.(type) {
                case nil:
                    fmt.Printf("%s : nil\n", key)
                case string:
                    fmt.Printf("%s : %s\n", key, val)
                default:
                    if v, ok := val.(*reflect.ArrayValue); ok == true {
                        fmt.Printf("%s: %V", key, v)
/*                        for _, elm := range v.Elem {
                            fmt.Printf("%s, ", elm)
                        }*/
                        fmt.Printf("\n")
                    } else {
                        fmt.Printf("Don't know %V\n", val)
                    }
            }
        }
        fmt.Printf("\n")
    }
    return c
}

func dispatcher (e *irc.IRCEvent) {

}

func main() {
    confpath := flag.String("c", "-c /path/to/conf/file.json", "Configuration file path")
    flag.Parse()
    _ = parseconf(*confpath)
    os.Exit(0)
    irccon := irc.IRC("nsfw", "nsfw")
    err := irccon.Connect("chat.freenode.net:6667")
    if err != nil {
    	fmt.Printf("%s\n", err)
    	fmt.Printf("%#v\n", irccon)
    	os.Exit(1)
    }
    irccon.AddCallback("001", func(e *irc.IRCEvent) { irccon.Join("#sabayon-hu") })
    irccon.AddCallback("JOIN", func(e *irc.IRCEvent) {
                                       irccon.Privmsg(e.Message, "mi a fakk van?!")
                                   })
    irccon.AddCallback("PRIVMSG", func(e *irc.IRCEvent) {
                                   if e.Message == "prout" {
                                               for _, s := range e.Arguments {
                                                   if s[0] == '#' {
                                                       irccon.Privmsg(s, "francharb: PROUTZOR")
                                                   }
                                               }
                                   }})
    irccon.Loop();
}
