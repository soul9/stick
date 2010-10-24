package stick

import (
	"github.com/thoj/Go-IRC-Client-Library"
//         irc "../Go-IRC-Client-Library/_obj/irc"
)

func IrcAction(c string, parms string, conn *irc.IRCConnection) {
    conn.Privmsg(c, "\001ACTION " + parms + "\001")
}