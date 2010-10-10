package stick

import (
//	"github.com/thoj/Go-IRC-Client-Library"
         irc "../Go-IRC-Client-Library/_obj/irc"
	"fmt"
	"os"
	"regexp"
	"json"
)



func msgDispatcher (e *irc.IRCEvent, conn *irc.IRCConnection, info *NetConf) {
    var chancfg ChanConf
    for _, chancfg = range info.Channels {
        if chancfg.Name == e.Arguments[0] {
            break
        }
    }
    for _, action := range chancfg.Actions {
        re, err := regexp.Compile(action.Match)
        if err == nil && re.MatchString(e.Message) {
            actionDispatcher(&action, e, conn, &chancfg)
        }
    }
    fmt.Printf("%#v\n", e)
}

func ircAction(c string, parms string, conn *irc.IRCConnection) {
    conn.Privmsg(c, "\001ACTION " + parms + "\001")
}

func learn(c string, p string, conn *irc.IRCConnection, acts map[string]ChanActConf){
    act := new(ChanActConf)
    err := json.Unmarshal([]byte(p), act)
    if err != nil {
        conn.Privmsg(c, "Sorry didn't understand: " + p)
        conn.Privmsg(c, err.String())
    }
    conn.Privmsg(c, "Learning: " + p)
    acts[string(len(acts))] = *act
}

func actionDispatcher(act *ChanActConf, e *irc.IRCEvent, conn *irc.IRCConnection, cfg *ChanConf) {
    re, err := regexp.Compile(`{\$.+}`)
    if err != nil {
        return
    }
    tmpre, err := regexp.Compile(act.Match)
    if err != nil {
        return
    }
    parms := re.ReplaceAllStringFunc(act.Parms, func(s string) string{
          switch s {
            case "{$victim}":
                s = e.Nick
            case "{$message}":
                s = e.Message
            case "{$message-match}":
                s = tmpre.ReplaceAllString(e.Message, "")
          }
          return s
        })
    switch act.Action {
        case "say":
            conn.Privmsg(e.Arguments[0], parms)
        case "action":
            ircAction(e.Arguments[0], parms, conn)
        case "learn":
            learn(e.Arguments[0], parms, conn, cfg.Actions)
        default:
            break
    }
}

func Init(confpath *string) map[string]*irc.IRCConnection {
    conf := parseconf(*confpath)
    conns := make(map[string]*irc.IRCConnection)
     for net, info := range conf.Networks {
        conns[net] = irc.IRC(info.Nick, info.Realname)
        if err :=conns[net].Connect(net); err != nil {
            fmt.Printf("%s\n", err)
            fmt.Printf("%#v\n", conns[net])
            os.Exit(1)
        }
        for _, cn := range info.Channels {
            conns[net].Join(cn.String())
        }
        conns[net].AddCallback("PRIVMSG", func(e *irc.IRCEvent){ msgDispatcher(e, conns[net], &info)})
        conns[net].Loop();
    }
    return conns
}
