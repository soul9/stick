package stick
import (
	"github.com/thoj/Go-IRC-Client-Library"
	"os"
	"json"
	"regexp"
)

func (st *Stick) Learn(p string, e *irc.IRCEvent, net string) os.Error {
    var act ChanActConf
    err := json.Unmarshal([]byte(p), &act)
    if err != nil {
        return os.NewError("Couldn't unmarshal json: " + p + ": " + err.String())
    }
    conn := st.Conns[net]
    //check if the match regexp is a valid regexp
    _, err = regexp.Compile(st.replaceVars(act, e, net).Match)
    if err != nil {
        return os.NewError("Action match isn't a valid regexp: " + act.Match + ": " + err.String())
    }
    err = st.Conf.AddActionChan(e.Arguments[0], act, net)
    if err != nil {
        return os.NewError("Couldn't add action: "+err.String())
    }
    conn.Privmsg(e.Arguments[0], "Learned: " + p)
    return nil
}

func (st *Stick) Forget(rlno int) os.Error {
}

func (st *Stick) ListActions() os.Error {
}