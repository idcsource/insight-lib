// Copyright 2016-2017
// CoderG the 2016 project
// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ]
// Normal Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// Use of this source code is governed by GNU LGPL v3 license

package drcm

import (
	"sync"

	"github.com/idcsource/insight00-lib/cpool"
	"github.com/idcsource/insight00-lib/hardstore"
	"github.com/idcsource/insight00-lib/ilogs"
	"github.com/idcsource/insight00-lib/nst"
	"github.com/idcsource/insight00-lib/roles"
)

// 锆存储。
// 这是一个带缓存的存储体系，基本功能方面可以看作是hardstore与rcontrol的合并（虽然还有很大不同），而增强方面它支持分布式存储。
type ZrStorage struct {
	/* 下面这部分是存储相关的 */

	// 配置信息
	config *cpool.Block
	// 本地存储
	local_store *hardstore.HardStore

	/* 下面这部分是缓存相关的 */

	// 角色缓存
	rolesCache map[string]*oneRoleCache
	// 最大缓存角色数
	cacheMax int64
	// 缓存数量
	rolesCount int64
	// 缓存满的触发
	cacheIsFull chan bool
	// 删除缓存
	deleteCache []string
	// 检查缓存数量中
	checkCacheNumOn bool

	/* 下面是分布式服务相关的 */

	// 分布式服务的模式，来自于常量DMODE_*
	dmode uint8
	// 自身的身份码，做服务的时候使用
	code string
	// 请求slave执行或返回数据的连接，string为slave对应的管理第一个值的首字母，而那个切片则是做镜像的
	slaves map[string][]*slaveIn
	// 监听的实例
	listen *nst.TcpServer
	// slave的连接池，从这里分配给slaveIn
	slavepool map[string]*nst.TcpClient
	// slave的slaveIn连接池
	slavecpool map[string]*slaveIn

	// 日志
	logs *ilogs.Logs
	// 全局锁
	lock *sync.RWMutex
}

// 一个角色的缓存，提供了锁
type oneRoleCache struct {
	// 锁
	lock *sync.RWMutex
	// 读写锁状态，CACHE_ROLE_LOCK_*
	lockstatus uint8
	// 角色自身
	role roles.Roleer
}

// 一台从机的信息
type slaveIn struct {
	name    string
	code    string
	tcpconn *nst.TcpClient
}

// 前缀状态，每次向slave发信息都要先把这个状态发出去
type Net_PrefixStat struct {
	// 操作类型，从OPERATE_*
	Operate int
	// 身份验证码
	Code string
}

// slave回执，slave收到PrefixStat之后的第一步返回信息
type Net_SlaveReceipt struct {
	// 数据状态，来自DATA_*
	DataStat uint8
	// 返回的错误
	Error string
}

// slave回执带数据体
type Net_SlaveReceipt_Data struct {
	// 数据状态，来自DATA_*
	DataStat uint8
	// 返回的错误
	Error string
	// 数据体
	Data []byte
}

// 角色的接收与发送格式
type Net_RoleSendAndReceive struct {
	// 角色的身体
	RoleBody []byte
	// 角色的关系
	RoleRela []byte
	// 角色的版本
	RoleVer []byte
}

// 角色的father修改的数据格式
type Net_RoleFatherChange struct {
	Id     string
	Father string
}

// 角色的所有子角色
type Net_RoleAndChildren struct {
	Id       string
	Children []string
}

// 角色的单个子角色关系的网络数据格式
type Net_RoleAndChild struct {
	Id    string
	Child string
}

// 角色的所有朋友
type Net_RoleAndFriends struct {
	Id      string
	Friends map[string]roles.Status
}

// 角色的单个朋友角色关系的网络数据格式
type Net_RoleAndFriend struct {
	Id     string
	Friend string
	Bind   int64
	Status roles.Status
	// 单一的绑定属性修改，1为int，2为float，3为complex
	Single uint8
	// 单一的绑定修改所对应的位置，也就是0到9
	Bit int
	// 单一修改的Int
	Int int64
	// 单一修改的Float
	Float float64
	// 单一修改的Complex
	Complex complex128
}

// 角色的单个上下文关系的网络数据格式
type Net_RoleAndContext struct {
	Id string
	// 上下文的名字
	Context string
	// 这是roles包中的CONTEXT_UP或CONTEXT_DOWN
	UpOrDown uint8
	// 要操作的绑定角色的ID
	BindRole string
}

// 角色的全部上下文
type Net_RoleAndContexts struct {
	Id       string
	Contexts map[string]roles.Context
}

// 角色的单个上下文关系数据的网络数据格式
type Net_RoleAndContext_Data struct {
	Id string
	// 上下文的名字
	Context string
	// 这是roles包中的CONTEXT_UP或CONTEXT_DOWN
	UpOrDown uint8
	// 要操作的绑定角色的ID
	BindRole string
	// 一个的状态位结构
	Status roles.Status
	// 上下文的结构
	ContextBody roles.Context
	// 单一的绑定属性修改，1为int，2为float，3为complex
	Single uint8
	// 单一的绑定修改所对应的位置，也就是0到9
	Bit int
	// 单一修改的Int
	Int int64
	// 单一修改的Float
	Float float64
	// 单一修改的Complex
	Complex complex128
}

// 角色的单个数据的数据体的网络格式
type Net_RoleData_Data struct {
	Id string
	// 数据点的名字
	Name string
	// 数据的字节流
	Data []byte
}
