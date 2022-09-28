package name_server

import (
	"errors"
	"fmt"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
	store "github.com/zeromicro/go-zero/core/stores/redis"
	"strings"
	"sync"
)

// go-zero的redis可以处理集群，所以用此redis
// 默认单节点，store.Type="node"，集群模式：store.Type="cluster"
// store := store.New("",func(r *store.Redis){r.Type = "node"})

var clusterData *storeInfo

type storeInfo struct {
	// model.QuickHandel
	store   *store.Redis
	appName string
	// localKeyNode  sync.Map //map[key]node
	localNodeAddr sync.Map //map[node]addr
	localNodeConn sync.Map //map[node]conn
	// deadNode      *[]string
}

func initCluster(appName string) {
	clusterData = newStore(appName)
}

func newStore(appName string) *storeInfo {

	return &storeInfo{
		store:   store.New(RedisAddr),
		appName: appName,
	}
}

// func (s *storeInfo) HandlerInfo(msg *model.Message) model.Handler {
// 	return s
// }

func (s *storeInfo) getKeyAddr(uid string, isRepair bool) (string, error) {

	node, err := s.getNodeByUser(uid)
	if err != nil {
		return NullString, err
	}
	if node == NullString {
		return NullString, errors.New("no target")
	}
	isDead, err := s.isDeadNode(node)
	if isDead || err != nil {
		return NullString, errors.New(" target dead")
	}

	if isRepair {
		s.delNodeAddr(node)
		s.putDeadNode(node)
		return NullString, errors.New("repair")
	}

	return s.getCacheNodeAddr(node)
}

//
// func (s *storeInfo) getCacheKeyAddr(modelName string) (string, error) {
//
// 	if err != nil {
// 		addr, ok := s.localNodeAddr.Load(node)
// 		if ok {
// 			return addr.(string), nil
// 		}
// 	}
// 	return s.getRepairKeyAddr(modelName)
// }
//
// func (s *storeInfo) getRepairKeyAddr(modelName string) (string, error) {
// 	s.localKeyNode.Delete(modelName)
// 	node, err := s.getNodeByUser(modelName)
// 	if err != nil || node == NullString {
// 		return NullString, err
// 	}
// 	addr, err := s.getRepairNodeAddr(node)
// 	if err != nil || node == NullString {
// 		return NullString, err
// 	}
// 	s.localKeyNode.Store(modelName, addr)
// 	return addr, nil
// }

func (s *storeInfo) putDeadNode(node string) error {
	if _, err := s.store.Sadd(s.getDeadNodeKey(), node); err != nil {
		return err
	}
	return nil
}

func (s *storeInfo) isDeadNode(node string) (bool, error) {
	return s.store.Sismember(s.getDeadNodeKey(), node)
}

func (s *storeInfo) getNodeByUser(modelName string) (string, error) {
	return s.store.Get(s.getModelKey(modelName))
}

func (s *storeInfo) getNodeAddr(nodeName string, isRepair bool) (string, error) {

	if isRepair {
		return s.getRepairNodeAddr(nodeName)
	}
	return s.getCacheNodeAddr(nodeName)
}

func (s *storeInfo) getRepairNodeAddr(node string) (string, error) {
	s.localNodeAddr.Delete(node)
	addr, err := s.store.Get(s.getNodeAddrKey(node))
	if err != nil {
		return NullString, err
	}
	if addr == NullString {
		return NullString, errors.New("no this node")
	}
	s.localNodeAddr.Store(node, addr)
	return addr, nil
}
func (s *storeInfo) getCacheNodeAddr(node string) (string, error) {
	addr, ok := s.localNodeAddr.Load(node)
	if ok {
		return addr.(string), nil
	}
	return s.getRepairNodeAddr(node)
}
// todo 做更优策略，缓存权重，不用每次都读redis
func (s *storeInfo) getHealthiestNodeAddr(isRepair bool) (string, error) {
	// node 数量不可能特别多，可以使用keys
	// 所有节点
	keys, err := s.store.Keys(s.getNodeAddrKey("*"))
	if err != nil {
		return NullString, err
	}
	min := 100000
	mapCountKey := make(map[int64]string)
	// 找到节点中，线程数量最小的节点
	for _, key := range keys {
		temp := strings.Split(key, ":")
		node := temp[len(temp)-1]
		count, err2 := s.store.Scard(s.getNodeListKey(node))
		if err2 != nil {
			return NullString, err2
		}
		mapCountKey[count] = node
		min = utils.Min(min, int(count))
	}
	node := mapCountKey[int64(min)]
	return s.getNodeAddr(node, isRepair)
}

func (s *storeInfo) deleteLocalNodeAddr(node string) {
	s.localNodeAddr.Delete(node)
	s.localNodeConn.Delete(node)
}

func (s *storeInfo) delNodeAddr(key string) error {
	s.deleteLocalNodeAddr(key)
	// 删除云端节点地址
	if _, err := s.store.Del(s.getNodeAddrKey(key)); err != nil {
		return err
	}
	// 删除云端节点进程展示
	// if _, err := s.delNodeList(key); err != nil {
	// 	return err
	// }
	return nil
}

func (s *storeInfo) delUserNode(key string) (int, error) {
	return s.store.Del(s.getModelKey(key))
}

func (s *storeInfo) delNodeList(key string) (int, error) {
	return s.store.Del(s.getNodeListKey(key))
}

func (s *storeInfo) getModelKey(key string) string {
	return fmt.Sprintf(model.ModelNodeKey, s.appName, key)
}
func (s *storeInfo) getNodeListKey(key string) string {
	return fmt.Sprintf(model.NodeListKey, s.appName, key)
}
func (s *storeInfo) getNodeAddrKey(key string) string {
	return fmt.Sprintf(model.NodeAddrKey, s.appName, model.C2S, key)
}
func (c *storeInfo) getDeadNodeKey() string {
	return fmt.Sprintf(model.DeadNodeKey, c.appName)
}
