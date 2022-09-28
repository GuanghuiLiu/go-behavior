package gateway

import (
	"net"
	"strconv"

	"github.com/GuanghuiLiu/behavior"
	"github.com/GuanghuiLiu/behavior/model"
	proto_game "github.com/GuanghuiLiu/behavior/tcp/long/game"
	pb "github.com/GuanghuiLiu/behavior/tcp/long/game/pb"
	"github.com/GuanghuiLiu/behavior/utils"
	"google.golang.org/protobuf/proto"
)

// MaxConnID for gateway process name
var MaxConnID uint64

type Reader struct {
	model.QuickHandel
	conn      *net.TCPConn
	toProcess string
	writer    string
	doPack    *proto_game.Packer
	router    behavior.IRouter
}

func newReader(conn *net.TCPConn, router behavior.CreateRouter) (*Reader, error) {
	r := &Reader{
		conn: conn,
	}
	connID, err := model.GetId("ConnID")
	if err != nil {
		return r, err
	}
	r.Name = NameReaderPrefix + strconv.FormatUint(connID, 10)
	r.toProcess = NullString
	r.doPack = proto_game.NewPack(r.Name, conn)
	r.router = router()
	return r, nil
}

func (r *Reader) receiveFunc() bool {
	clientMsg, err := r.doPack.CircleUnpack()
	// clientMsg, err := r.doPack.Unpack(r.conn)
	if err != nil {
		return true
	}
	// login 在当前process处理；其他业务，在toProcess处理
	switch clientMsg.ProtoID {
	case proto_game.Login:
		// 已登录，不做任何处理
		if r.toProcess != NullString {
			break
		}

		login, err2 := r.unPackLogin(clientMsg.Data)
		if err2 != nil {
			return true
		}

		uid := r.getUid(login.Name, login.Password)

		// start writer
		w := newWriter(r.Name, r.conn)

		w.Run(w, model.SetRetry(0, 0))
		r.writer = w.Name
		r.toProcess = uid

		if ok, _ := model.OnThisNode(uid); ok {
			if err = r.SendInfo(uid, proto_game.ResetConn, utils.Encode(w.Name)); err == nil {
				break
			}
		}

		role := behavior.NewActor(uid, w.Name, r.router)

		if err3 := role.Run(role, model.SetIsGlobal(true), model.SetLiveTime(RoleLiveTime)); err3 != nil {
			r.reLogin()
			return true
		}
		if e := r.loginSuccess(); e != nil {
			return true
		}

	default:
		if r.toProcess == NullString {
			r.reLogin()
			return true
		}
		if e := r.SendInfo(r.toProcess, clientMsg.ProtoID, clientMsg.Data); e != nil {
			r.reLogin()
			return true
		}
	}
	return false
}

func (r *Reader) HandlerStop() model.Handler {
	if r.writer != NullString {
		r.SendStop(r.writer)
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return r
}

func (r *Reader) HandlerInfo(msg *model.Message) model.Handler {
	return r
}

func (r *Reader) unPackLogin(msg []byte) (*pb.Login, error) {

	login := &pb.Login{}
	if err := proto.Unmarshal(msg, login); err != nil {
		return nil, err
	}

	return login, nil
}

func (r *Reader) send2Client(protoID uint64, data []byte) error {
	return r.doPack.Send2Client(r.conn, protoID, data)
}

func (r *Reader) reLogin() {
	m := &pb.LoginInfo{
		Text: "not login,place re login",
	}
	binaryData, _ := proto.Marshal(m)
	r.send2Client(proto_game.LoginInfo, binaryData)
}

func (r *Reader) loginSuccess() error {
	m := &pb.LoginInfo{
		Text: "login success",
	}
	binaryData, _ := proto.Marshal(m)
	return r.send2Client(proto_game.LoginInfo, binaryData)
}

func (r *Reader) getUid(name, passWprd string) string {
	// todo DB 校验
	return name
}
