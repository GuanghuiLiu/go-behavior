package model

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	store "github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	// behavior:(app):model:(modelID)
	ModelNodeKey = "behavior:%s:process-on-node:%s"
	// behavior:(app):node:list:(nodeID)
	NodeListKey = "behavior:%s:node-process-list:%s"
	// behavior:(app):node:addr:(nodeID)
	NodeAddrKey = "behavior:%s:node-addr:%s:%s"
	// behavior:(app):dead_node
	DeadNodeKey = "behavior:%s:dead_node"
)

// for NodeAddrKey
const (
	C2S = "c2s"
	S2S = "s2s"
)

var localName string = "bj_0001_00091"

// 过期时间，必须大于，同步时间。保证不被误删
const syncTime time.Duration = 5 * time.Second // seconds

const modelExpire = 60 * 60 * 2 // seconds
const nodeExpire = 60 * 60 * 2  // seconds

var failedDelete []string
var failedRegister []string

// go-zero的redis可以处理集群，所以用它的redis
var ClusterCenter *clusterCenter

type clusterCenter struct {
	appName           string
	deleteModelName   chan string
	registerModelName chan string
	syncTk            *time.Ticker
	store             *store.Redis // 未用接口，直接用redis;暂时没有修改意愿
	localCenter       sync.Map     //map[process]addr
	localConn         sync.Map     //map[addr]conn
}

func InitCluster(app, nodeName string, s2s bool, s2sPort uint) {
	localName = nodeName
	c := clusterCenter{
		appName:           app,
		deleteModelName:   make(chan string),
		registerModelName: make(chan string),
		store:             newRedis(),
		syncTk:            time.NewTicker(syncTime),
	}
	ClusterCenter = &c
	if s2s {
		go startS2S(s2sPort)
		initS2SHeartbeat()
	}
	go c.run()
}

func newRedis() *store.Redis {

	// 默认单节点，redis.Type="node"，集群模式：redis.Type="cluster"
	// redis := store.New("",func(r *store.Redis){r.Type = "node"})

	return store.New(RedisAddr)
}
func (c *clusterCenter) run() {
	defer func() {
		err := recover()
		fmt.Println("cluster_center panic:", err, string(debug.Stack()))
		c.run()
	}()
	for {
		select {
		case rn := <-c.registerModelName:
			c.registerModel(rn)
			c.registerNode(rn)
		case dn := <-c.deleteModelName:
			c.deleteModel(c.getModelKey(dn))
		case <-c.syncTk.C:
			c.syneData()
			c.syncTk.Reset(syncTime)
		}
	}
}

func (c *clusterCenter) syneData() {
	// 同步节点信息
	c.syncNodeListProcess()

	// 重新设置model过期时间
	c.reExpireModelNode()

	// 处理失败
	c.reDelModelNode()
	c.reRegModelNode()
}

func (c *clusterCenter) reExpireModelNode() {
	for modelName, _ := range allModel {
		c.store.Expire(c.getModelKey(modelName), modelExpire)
	}
}
func (c *clusterCenter) reDelModelNode() {
	fd := []string{}
	for _, modelName := range failedDelete {
		if err := c.deleteModel(c.getModelKey(modelName)); err != nil {
			fd = append(fd, modelName)
		}
	}
	failedDelete = fd
}
func (c *clusterCenter) reRegModelNode() {
	fr := []string{}
	for _, modelName := range failedRegister {
		if err := c.registerModel(c.getModelKey(modelName)); err != nil {
			fr = append(fr, modelName)
		}
	}
	failedDelete = fr
}

func (c *clusterCenter) syncNodeListProcess() {
	c.store.Del(c.getNodeListKey(localName))
	for key, _ := range allModel {
		c.store.Sadd(c.getNodeListKey(localName), key)
	}
	c.store.Expire(c.getNodeListKey(localName), nodeExpire)
}
func (c *clusterCenter) registerNode(processName string) error {

	c.store.Sadd(c.getNodeListKey(localName), processName)
	c.store.Expire(c.getNodeListKey(localName), nodeExpire)
	return nil
}

func (c *clusterCenter) cleanNode(direction, process string) error {

	node, err := c.getNodeByModel(process)
	if err != nil {
		return err
	}
	c.deleteModelNode(process)
	if err = c.delNodeAddr(direction, process); err != nil {
		return err
	}
	isDear, err := c.isDeadNode(node)
	if node == NullString || isDear || err != nil {
		return nil
	}
	return c.putDeadNode(node)

}

func (c *clusterCenter) getAddr(direction, process string) (string, error) {
	return c.getAddrCache(direction, process)
}

func (c *clusterCenter) getAddrRepair(end, key string) (string, error) {
	c.localCenter.Delete(key)
	node, err := c.getNodeByModel(key)
	if err != nil {
		return "", err
	}
	a, e := c.getNodeAddr(end, node)
	if e != nil {
		return "", e
	}
	c.localCenter.Store(key, a)
	return a, nil
}

func (c *clusterCenter) getAddrCache(direction, process string) (string, error) {
	addr, ok := c.localCenter.Load(process)
	if ok {
		return addr.(string), nil
	}
	return c.getAddrRepair(direction, process)
}

func (c *clusterCenter) registerModel(name string) error {
	if _, err := c.store.SetnxEx(c.getModelKey(name), localName, modelExpire); err != nil {
		failedRegister = append(failedRegister, name)
		return err
	}
	return nil
}
func (c *clusterCenter) deleteModel(name string) error {
	c.store.Srem(c.getNodeListKey(localName), name)
	if _, err := c.store.Del(c.getModelKey(name)); err != nil {
		failedDelete = append(failedDelete, name)
		return err
	}
	return nil
}
func (c *clusterCenter) sentRegister(name string) {
	if c != nil {
		c.registerModelName <- name
	}

}
func (c *clusterCenter) sentDelete(name string) {
	if c != nil {
		c.deleteModelName <- name
	}

}

func (c *clusterCenter) SetNodeAddr(end, addr string) error {
	if _, err := c.deleteDeadNode(localName); err != nil {
		return err
	}
	if err := c.store.Set(c.getNodeAddrKey(end, localName), addr); err != nil {
		return err
	}
	return nil
	// return c.store.Expire(c.getNodeAddrKey(end, localName), nodeExpire)
}

func (c *clusterCenter) getNodeByModel(process string) (string, error) {
	return c.store.Get(c.getModelKey(process))
}

func (c *clusterCenter) getNodeAddr(end, node string) (string, error) {
	// c.store.Expire(c.getNodeAddrKey(end, localName), nodeExpire)
	return c.store.Get(c.getNodeAddrKey(end, node))
}

func (c *clusterCenter) deleteNodeProcessList(node string) {
	c.store.Del(c.getNodeListKey(node))
}

func (c *clusterCenter) deleteModelNode(process string) {
	c.store.Del(c.getModelKey(process))
}

func (c *clusterCenter) deleteLocalNodeAddr(node string) {
	c.localCenter.Delete(node)
}

func (c *clusterCenter) delNodeAddr(direction, key string) error {
	c.deleteLocalNodeAddr(key)
	if _, err := c.store.Del(c.getNodeAddrKey(direction, key)); err != nil {
		return err
	}
	return nil
}

func (c *clusterCenter) deleteDeadNode(node string) (int, error) {
	return c.store.Srem(c.getDeadNodeKey(), node)
}

func (c *clusterCenter) putDeadNode(node string) error {
	if _, err := c.store.Sadd(c.getDeadNodeKey(), node); err != nil {
		return err
	}
	return nil
}

func (c *clusterCenter) isDeadNode(node string) (bool, error) {
	return c.store.Sismember(c.getDeadNodeKey(), node)
}

func (c *clusterCenter) getModelKey(key string) string {
	return fmt.Sprintf(ModelNodeKey, c.appName, key)
}
func (c *clusterCenter) getNodeListKey(key string) string {
	return fmt.Sprintf(NodeListKey, c.appName, key)
}
func (c *clusterCenter) getNodeAddrKey(direction, key string) string {
	return fmt.Sprintf(NodeAddrKey, c.appName, direction, key)
}
func (c *clusterCenter) getDeadNodeKey() string {
	return fmt.Sprintf(DeadNodeKey, c.appName)
}
