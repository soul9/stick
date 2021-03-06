package stick

import (
	"github.com/thoj/Go-IRC-Client-Library"
//         irc "../Go-IRC-Client-Library/_obj/irc"
	"fmt"
	"os"
	"regexp"
)

type Stick struct {
    Conf *StickConf
    Conns map[string]*irc.IRCConnection
    
}

func (st *Stick) msgDispatcher (e *irc.IRCEvent, net string) {
    //make a copy of the event so it doesn't get modified (when replying to (i.e.) learn messages, the action is auto-triggered)
    newevent := *e
    e = &newevent
    if e.Nick == st.Conf.Conf.Networks[net].Nick {
        return
    }
    var chancfg ChanConf
    for _, chancfg = range st.Conf.Conf.Networks[net].Channels {
        if chancfg.Name == e.Arguments[0] {
            break
        }
    }
    for _, action := range chancfg.Actions {
        action = st.replaceVars(action, e, net)
        re, err := regexp.Compile(action.Match)
        if err == nil && re.MatchString(e.Message) {
            st.actionDispatcher(&action, e, net)
        }
    }
    fmt.Printf("%#v\n", e)
}


func (st *Stick) replaceVars(c ChanActConf, e *irc.IRCEvent, net string) ChanActConf {
    info := st.Conf.Conf.Networks[net]
    re, _ := regexp.Compile(`{\$.+}`)
    tmpre, err := regexp.Compile(c.Match)
    if err != nil {
        st.Conns[net].Privmsg(e.Arguments[0], "oops in replaceVars: " + err.String())
        return c
    }
    f := func(s string) string{
         switch s {
            case "{$victim}":
                s = e.Nick
            case "{$message}":
                s = e.Message
            case "{$message-match}":
                s = tmpre.ReplaceAllString(e.Message, "")
            case "{$mynick}":
                s = info.Nick
          }
          return s
        }
    match2 := re.ReplaceAllStringFunc(c.Match, f)
    tmpre, err = regexp.Compile(match2)
    if err == nil {
        c.Match = match2
    }
    c.Parms = re.ReplaceAllStringFunc(c.Parms, f)
    return c
}

func (st *Stick) actionDispatcher(act *ChanActConf, e *irc.IRCEvent, net string) {
    switch act.Action {
        case "say":
            st.Conns[net].Privmsg(e.Arguments[0], act.Parms)
        case "action":
            IrcAction(e.Arguments[0], act.Parms, st.Conns[net])
        case "learn":
            st.Learn(act.Parms, e, net)
        case "part":
            st.Conns[net].Part(e.Arguments[0])
        case "quit":
            st.Conns[net].Quit()
        default:
            break
    }
}

func Init(confpath *string) (*Stick, os.Error) {
    st := new(Stick)
    var err os.Error
    st.Conf, err = InitConf(*confpath)
    if err != nil {
        return st, os.NewError("Couldn't read configuration: " + err.String())
    }
    st.Conns = make(map[string]*irc.IRCConnection)
    for net, info := range st.Conf.Conf.Networks {
        st.Conns[net] = irc.IRC(info.Nick, info.Realname)
        if err :=st.Conns[net].Connect(net); err != nil {
            fmt.Printf("%s\n", err)
            fmt.Printf("%#v\n", st.Conns[net])
            os.Exit(1)
        }
        for _, cn := range info.Channels {
            st.Conns[net].Join(cn.String())
        }
        st.Conns[net].AddCallback("PRIVMSG", func(e *irc.IRCEvent){ st.msgDispatcher(e, net)})
        st.Conns[net].Loop();
    }
    return st, nil
}
