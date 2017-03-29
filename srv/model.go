package srv

import (
)

// user id
type uid uint

// girl id, signed-int for compatibility to libsaki
type gid int

// level, pt, and rating
type lpr struct {
	Level		int
	Pt			int
	Rating		float64
}

type user struct {
	Id			uid
	Username	string
	lpr
}

type girl struct {
	Id			gid
	lpr
}

type bookType int

func (b bookType) index() int {
	return int(b)
}

func (b bookType) valid() bool {
	i := int(b)
	return 0 <= i && i < 4
}

type reqTypeOnly struct {
	Type		string
}

type reqAuth struct {
	Type		string
	Version		string
	Username	string
	Password	string
}

type reqAction struct {
	Nonce		int
	ActStr		string
	ActArg		string
}

type reqBook struct {
	BookType	bookType
}

type reqChoose struct {
	GirlIndex	int
}

type respTypeOnly struct {
	Type		string
}

type respAuthFail struct {
	Type		string
	Ok			bool
	Reason		string
}

func newRespAuthFail(str string) interface{} {
	return respAuthFail{"auth", false, str}
}

type respAuthOk struct {
	Type	string
	Ok		bool
	User	*user
}

func newRespAuthOk(u *user) interface{} {
	return respAuthOk{"auth", true, u}
}

type bookEntry struct {
	Bookable	bool
	Book		int
	Play		int
}

type respLookAround struct {
	Type		string
	Conn		int
	Books		[4]bookEntry
}

func newRespLookAround(conn int) *respLookAround {
	resp := new(respLookAround)
	resp.Type = "look-around"
	resp.Conn = conn
	return resp
}

type respUpdateUser struct {
	Type		string
	User		*user
}

func newRespUpdateUser(user *user) *respUpdateUser {
	resp := new(respUpdateUser)
	resp.Type = "update-user"
	resp.User = user
	return resp
}



