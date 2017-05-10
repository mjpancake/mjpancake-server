package tssn

import (
	"encoding/json"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mjpancake/hisa/db"
	"github.com/mjpancake/hisa/model"
	"github.com/mjpancake/hisa/node"
	"github.com/mjpancake/hisa/saki"
)

func (tssn *tssn) Happy(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		tssn.bye(ctx)
	case *actor.ReceiveTimeout:
		tssn.sweepAll()
	case *pcAction:
		tssn.handleAction(msg.Uid, msg.Act)
	default:
		log.Printf("tssn.Ready unexpected %T\n", msg)
	}

	switch ctx.Message().(type) {
	case *actor.Stopping, *actor.Stopped:
	default:
		tssn.checkGameOver()
	}
}

func (tssn *tssn) handleAction(uid model.Uid, act *model.CsAction) {
	i, _ := tssn.findUser(uid)
	if act.ActStr == "RESUME" {
		tssn.onlines[i] = true
	} else if act.Nonce != tssn.nonces[i] {
		tssn.kick(i, "wrong action nonce")
		return
	}
	tssn.waits[i] = false
	mails := tssn.table.Action(i, act.ActStr, act.ActArg)
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *tssn) start() {
	log.Println("TSSN ****", tssn.uids[0], tssn.gids)
	tssn.table = saki.NewTableSession(
		int(tssn.gids[0]), int(tssn.gids[1]),
		int(tssn.gids[2]), int(tssn.gids[3]))

	mails := tssn.table.Start()
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *tssn) handleMails(mails saki.MailVector) {
	size := int(mails.Size())
	if size > 0 {
		var nonceInced [4]bool
		for i := 0; i < size; i++ {
			w := mails.Get(i).GetTo()
			if w != -1 && !nonceInced[w] {
				tssn.nonces[w]++
				nonceInced[w] = true
			}
		}
	}

	for i := 0; i < size; i++ {
		toWhom := mails.Get(i).GetTo()
		str := mails.Get(i).GetMsg()
		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(str), &msg); err != nil {
			log.Fatalln("unmarshal c++ str", err)
		}
		if toWhom == -1 {
			tssn.handleSystemMail(msg, str)
		} else {
			tssn.sendUserMail(toWhom, msg)
		}
	}
}

func (tssn *tssn) sendUserMail(who int, msg map[string]interface{}) {
	msg["Nonce"] = tssn.nonces[who]
	if msg["Event"] == "resume" {
		right := (who + 1) % 4
		cross := (who + 2) % 4
		left := (who + 3) % 4
		rotated := [4]*model.User{
			tssn.users[who],
			tssn.users[right],
			tssn.users[cross],
			tssn.users[left],
		}
		msg["Args"].(map[string]interface{})["users"] = rotated
	}

	err := tssn.sendPeer(who, msg)
	if err != nil && msg["Event"] == "activated" {
		if tssn.anyOnline() && !tssn.table.GameOver() {
			tssn.sweepOne(who)
		}
	}
}

func (tssn *tssn) handleSystemMail(msg map[string]interface{},
	msgStr string) {
	switch msg["Type"] {
	case "round-start-log":
		fmt := "TSSN .... %v %v.%v%s d=%v depo=%v seed=%v"
		al := ""
		if msg["allLast"].(bool) {
			al = "a"
		}
		log.Printf(fmt, tssn.uids[0], msg["round"], msg["extra"], al,
			msg["dealer"], msg["deposit"], uint(msg["seed"].(float64)))
	case "table-end-stat":
		var stat model.EndTableStat
		err := json.Unmarshal([]byte(msgStr), &stat)
		if err != nil {
			log.Fatalln("table-end-stat unmarshal", err)
		}
		tssn.injectUsers(stat.Replay)
		db.UpdateUserGirl(tssn.bookType, tssn.uids, tssn.gids, &stat)
		for w := 0; w < 4; w++ {
			node.Umgr.Tell(&node.MuUpdateInfo{Uid: tssn.uids[w]})
		}
	case "riichi-auto":
		time.Sleep(1000 * time.Millisecond)
		who := int(msg["Who"].(float64))
		act := &model.CsAction{
			Nonce:  tssn.nonces[who],
			ActStr: "SPIN_OUT",
			ActArg: "-1",
		}
		tssn.handleAction(tssn.uids[who], act)
	default:
		log.Fatalln("unknown system mail", msg)
	}
}

func (tssn *tssn) injectUsers(replay map[string]interface{}) {
	var users [4]map[string]interface{}
	for w := 0; w < 4; w++ {
		user := make(map[string]interface{})
		user["Id"] = tssn.uids[w]
		user["Username"] = tssn.users[w].Username
		user["Level"] = tssn.users[w].Level
		user["Rating"] = tssn.users[w].Rating
		users[w] = user
	}
	replay["users"] = users
}

func (tssn *tssn) sweepOne(i int) {
	mails := tssn.table.SweepOne(i)
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}

func (tssn *tssn) sweepAll() {
	var targets int
	mails := tssn.table.SweepAll(&targets)
	for w := uint(0); w < 4; w++ {
		if (targets & (1 << w)) != 0 {
			tssn.waits[w] = false
			tssn.kick(int(w), "happy timeout")
		}
	}
	defer saki.DeleteMailVector(mails)
	tssn.handleMails(mails)
}