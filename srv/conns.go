package srv

import (
	"log"
	"net"
	"bufio"
	"strings"
	"encoding/json"
)

type conns struct {
	login	chan *login
	signUp	chan *login
	logout	chan uid
	start	chan [4]uid
    peer    chan *Mail
	dao		*dao
	users	map[uid]*user
	conns	map[uid]net.Conn
	books	*books
	tables	*tables
}

func newConns(dao *dao) *conns {
	conns := new(conns)

	conns.login = make(chan *login)
	conns.signUp = make(chan *login)
	conns.logout = make(chan uid)
	conns.start = make(chan [4]uid)
	conns.peer = make(chan *Mail)

	conns.dao = dao
	conns.users = make(map[uid]*user)
	conns.conns = make(map[uid]net.Conn)
	conns.books = newBooks(conns)
	conns.tables = newTables(conns)

	return conns
}

func (conns *conns) Loop() {
	go conns.books.Loop()
	go conns.tables.Loop()

	for {
		select {
		case login := <-conns.login:
			conns.recvLogin(login)
		case sign := <-conns.signUp:
			conns.recvSignUp(sign)
		case uid := <-conns.logout:
			conns.sub(uid)
		case uids := <-conns.start:
			conns.tables.Create() <- uids
		case mail := <-conns.peer:
			conns.send(mail.To, mail.Msg)
		}
	}
}

func (conns *conns) Login() chan<- *login {
	return conns.login
}

func (conns *conns) SignUp() chan<- *login {
	return conns.signUp
}

func (conns *conns) Logout() chan<- uid {
	return conns.logout
}

func (conns *conns) Start() chan<- [4]uid {
	return conns.start
}

func (conns *conns) Peer() chan<- *Mail {
	return conns.peer
}

func (conns *conns) recvLogin(login *login) {
	if login.Version != Version {
		str := "客户端版本过旧"
		conns.reject(login.conn, newRespAuthFail(str))
	} else {
		user := conns.dao.login(login)
		if user != nil {
			conns.add(user, login.conn)
		} else {
			str := "用户名或密码错误"
			conns.reject(login.conn, newRespAuthFail(str))
		}
	}
}

func (conns *conns) recvSignUp(sign *login) {
	if sign.Version != Version {
		str := "客户端版本过旧"
		conns.reject(sign.conn, newRespAuthFail(str))
	} else {
		user := conns.dao.signUp(sign)
		if user != nil {
			conns.add(user, sign.conn)
		} else {
			str := "用户名已存在"
			conns.reject(sign.conn, newRespAuthFail(str))
		}
	}
}

func (conns *conns) add(user *user, conn net.Conn) {
	// prevent dup login
	if _, ok := conns.users[user.Id]; ok {
		str := "该用户已登录"
		conns.reject(conn, newRespAuthFail(str));
		return
	}

	conns.users[user.Id] = user
	conns.conns[user.Id] = conn
	conns.send(user.Id, newRespAuthOk(user))

	go conns.readLoop(user.Id)
}

func (conns *conns) sub(uid uid) {
	conn, found := conns.conns[uid]
	if found {
		conns.books.Unbook() <- uid
		log.Println(uid, "----")
		conn.Close()
	}

	delete(conns.conns, uid)
	delete(conns.users, uid)
}

func (conns *conns) readLoop(uid uid) {
	conn := conns.conns[uid]
	for {
		breq, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			conns.logout <- uid
			return
		}

		log.Print(uid, " ---> ", string(breq))
		var req reqTypeOnly
		if err := json.Unmarshal(breq, &req); err != nil {
			log.Fatalln("E conns.readLoop", err)
			return
		}
		conns.switchRead(uid, req.Type, breq)
	}
}

func (conns *conns) switchRead(uid uid, t string, breq []byte) {
	switch {
	case t == "look-around":
		conns.sendLookAround(uid)
	case t == "book":
		conns.books.Book() <- uid
	case t == "unbook":
		conns.books.Unbook() <- uid
	case t == "ready":
		conns.tables.Ready() <- uid
	case strings.HasPrefix(t, "t-"):
		act := reqAction{uid: uid}
		if err := json.Unmarshal(breq, &act); err != nil {
			log.Println("E conns.switchRead", err)
			return
		}
		conns.tables.Action() <- &act
	}
}

func (conns *conns) send(uid uid, msg interface{}) {
	conn, found := conns.conns[uid]
	if !found {
		log.Println("E conns.send user", uid, "not found")
		return
	}

    var jsonb []byte
    if str, ok := msg.(string); ok {
        jsonb = []byte(str)
    } else {
        var err error
        jsonb, err = json.Marshal(msg)
        if err != nil {
            log.Fatalln("conns.send", err)
        }
    }

	if _, err := conn.Write(append(jsonb, '\n')); err != nil {
		log.Println("conns.send", err)
	} else {
		log.Println(uid, "<---", string(jsonb))
	}
}

func (conns *conns) reject(conn net.Conn, msg interface{}) {
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("conns.reject", err)
	}

	if _, err := conn.Write(append(jsonb, '\n')); err != nil {
		log.Println("conns.reject", err)
	} else {
		log.Println(conn.RemoteAddr(), "<---", string(jsonb))
	}

	conn.Close()
}

func (conns *conns) sendLookAround(uid uid) {
	bookable := !conns.tables.HasUser(uid)
	connCt := len(conns.conns)
	playCt := 4 * len(conns.tables.sessions)
	bookCt := conns.books.wait
	conns.send(uid, newRespLookAround(bookable, connCt, bookCt, playCt))
}

