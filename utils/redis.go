package utils

/*
	redis 工具类
*/
import (
	//"fmt"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
	"time"
)

/*
	初始化redis结构体
*/
type RedisHander struct {
	Pool *redis.Pool
}

//声明redis map
var redisHandleList map[string]*RedisHander

//创建redis锁
var createRedisSync *sync.Mutex = &sync.Mutex{}

/**
 * GetRedisInstance
 * 对外单例获取redis对象
 * @param string redisName 配置名称
 * @param int isLoad 是否需要取余加载 1不取余  2取余
 * @param int calculate 取余变量
 * @return RAddAppNews 操作结果
 */
func GetRedisInstance(configGroup string, isLoad int, calculate int) *RedisHander {
	//获取redis配置
	hostConfig := GetRedisConfig(configGroup, isLoad, calculate)
	//获取redis对象
	thisRedisHandle, ok := redisHandleList[hostConfig]
	//如果不存在 或 已失效 重新获取设置
	if ok == false || thisRedisHandle.Pool == nil {
		//加锁
		createRedisSync.Lock()
		//执行结束解锁
		defer createRedisSync.Unlock()
		//双重锁
		if ok == false || thisRedisHandle.Pool == nil {
			//创建redis句柄
			thisRedisHandle = CreateRedis(hostConfig)
			fmt.Println(thisRedisHandle, "创建redis")
			//thisRedisHandle = RedisHander{newRedisHandle}
		}
		//赋值
		redisHandleList[hostConfig] = thisRedisHandle
	}
	//返回redis句柄
	return thisRedisHandle
}

/**
 * GetRedisConfig
 * 获取redis配置
 * @param string redisName 配置名称
 * @param int isLoad 是否需要取余加载 1不取余  2取余
 * @param int calculate 取余变量
 * @return RAddAppNews 操作结果
 */
func GetRedisConfig(redisName string, isLoad int, calculate int) string {
	//获取配置信息
	config, err := LoadConfig("redis", redisName)
	if err != nil {
		return ""
	}

	//判断是否取余
	if isLoad == 1 {
		return config["host"] + ":" + config["port"]
	}
	//取模计算
	hostNum, _ := strconv.Atoi(config["hostNum"])
	remainder := calculate % hostNum

	//拼串rediskey
	hostName := "host" + strconv.Itoa(remainder+1)
	portName := "port" + strconv.Itoa(remainder+1)
	return config[hostName] + ":" + config[portName]
}

/*
	初始化redis引擎
*/
func CreateRedis(configGroup string) *RedisHander {
	config, err := LoadConfig("redis", configGroup)
	if err != nil {
		return nil
	}
	configInfo := config["host"] + ":" + config["port"]
	return &RedisHander{
		Pool: newPool(configInfo),
	}
}

/*
	创建redis connection pool
*/
func newPool(configGroup string) *redis.Pool {
	if configGroup == "" {
		WriteLog("log_redis", "redis配置获取失败！", configGroup)
		return nil
	}
	//连接redis
	return &redis.Pool{
		MaxIdle:     80,
		MaxActive:   500, // max number of connections
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", configGroup)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

/*
	获取redis值
*/
func (cacheRedis *RedisHander) GetInt(key string) (int, error, bool) {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()

	v, err := thisRedis.Do("GET", key)
	if v == nil || err != nil {
		return 0, err, false
	}
	v, err = redis.Int(v, err)

	return v.(int), err, true
}

/*
	获取redis值
*/
func (cacheRedis *RedisHander) GetString(key string) (string, error, bool) {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()
	v, err := thisRedis.Do("GET", key)
	if v == nil || err != nil {
		return "", err, false
	}
	v, err = redis.String(v, err)

	return v.(string), err, true
}

/*
	获取redis值
*/
func (cacheRedis *RedisHander) Lpop(key string) (interface{}, error) {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()
	v, err := thisRedis.Do("LPOP", key)
	if v == nil || err != nil {
		return nil, err
	}
	return redis.String(v, err)

}

/*
	获取redis值
*/
func (cacheRedis *RedisHander) HMGetToInterface(key string, arr interface{}) []interface{} {

	//获取redis Pool
	thisRedis := cacheRedis.Pool.Get()
	//最后关闭
	defer thisRedis.Close()
	//组装key为interface{} 集合
	keylist := ToSlice(key, arr)
	//获取数据
	reply, err := redis.Values(thisRedis.Do("HMGET", keylist...))
	//验证是否获取到数据
	if err != nil {
		WriteLog("log_redis", "HMGet获取数据失败,KEY：", keylist, "  错误原因：", err)
		return nil
	}
	//初始化返回数据
	result := make([]interface{}, len(reply))
	//循环转义数据
	for rkey, rvalue := range reply {
		//判断rke是否为nil
		if rvalue == nil {
			result[rkey] = nil
			continue
		}
		//初始化接参数据
		var onlyResult string
		//初始化当前需转义数据
		newValue := make([]interface{}, 1)
		//赋值当前需转义数据
		newValue[0] = rvalue
		//转义数据 成功赋值
		if _, err := redis.Scan(newValue, &onlyResult); err != nil {
			continue
		}
		result[rkey] = onlyResult
	}
	//返回数据
	return result
}

/*
	获取redis值
*/
func (cacheRedis *RedisHander) HMGet(key string, arr interface{}) []string {

	//获取redis Pool
	thisRedis := cacheRedis.Pool.Get()
	//最后关闭
	defer thisRedis.Close()
	//组装key为interface{} 集合
	keylist := ToSlice(key, arr)
	//获取数据
	reply, err := redis.Values(thisRedis.Do("HMGET", keylist...))
	//验证是否获取到数据
	if err != nil {
		WriteLog("log_redis", "HMGet获取数据失败,KEY：", keylist, "  错误原因：", err)
		return nil
	}
	//初始化返回数据
	result := make([]string, len(reply))
	//循环转义数据
	for rkey, rvalue := range reply {
		//判断rkey是否为nil
		if rvalue == nil {
			continue
		}
		//初始化接参数据
		var onlyResult string
		//初始化当前需转义数据
		newValue := make([]interface{}, 1)
		//赋值当前需转义数据
		newValue[0] = rvalue
		//转义数据 成功赋值
		if _, err := redis.Scan(newValue, &onlyResult); err != nil {
			continue
		}
		result[rkey] = onlyResult
	}
	//返回数据
	return result
}

/*
	获取redis值
*/
func (cacheRedis *RedisHander) HVals(key string) []string {
	//获取redis Pool
	thisRedis := cacheRedis.Pool.Get()
	//最后关闭
	defer thisRedis.Close()
	//获取数据
	reply, err := redis.Values(thisRedis.Do("HVALS", key))
	//验证是否获取到数据
	if err != nil {
		WriteLog("log_redis", "HVals获取数据失败,KEY：", key, "  错误原因：", err)
		return nil
	}
	//初始化返回数据
	result := make([]string, len(reply))
	//循环转义数据
	for rkey, rvalue := range reply {
		//初始化接参数据
		var onlyResult string
		//初始化当前需转义数据
		newValue := make([]interface{}, 1)
		//赋值当前需转义数据
		newValue[0] = rvalue
		//转义数据 成功赋值
		if _, err := redis.Scan(newValue, &onlyResult); err != nil {
			continue
		}
		result[rkey] = onlyResult
	}
	//返回数据
	return result
}

/*
	获取redis总数
*/
func (cacheRedis *RedisHander) Llen(key string) (interface{}, error) {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()
	return redis.Int(thisRedis.Do("LLEN", key))

}

/*
	压入redis中
*/
func (cacheRedis *RedisHander) Lpush(key string, value string) (interface{}, error) {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()
	return redis.Int(thisRedis.Do("LPUSH", key, value))
}

/*
	自增计数器redis中
	返回int型
*/
func (cacheRedis *RedisHander) Incr(key string) (interface{}, error) {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()
	v, err := redis.Int(thisRedis.Do("INCR", key))
	return v, err
}

/*
	设置超时时间
*/
func (cacheRedis *RedisHander) Expire(key string, timeOut int) interface{} {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()
	v, err := thisRedis.Do("EXPIRE", key, timeOut)
	if err != nil {
		return err
	}
	return v
}

/*
	压入redis中
	返回的是bool型
*/
func (cacheRedis *RedisHander) Set(key string, value int64) (interface{}, error) {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()

	v, err := redis.Bool(thisRedis.Do("SET", key, value))

	return v, err
}

/*
	压入redis中
	返回的是bool型
*/
func (cacheRedis *RedisHander) SetNx(key string, value int64) (interface{}, error) {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()

	v, err := redis.Bool(thisRedis.Do("SETNX", key, value))

	return v, err
}

/*
删除redis
*/
func (cacheRedis *RedisHander) Delete(key string) (interface{}, error) {
	thisRedis := cacheRedis.Pool.Get()
	defer thisRedis.Close()

	v, err := redis.Bool(thisRedis.Do("DEL", key))

	return v, err
}

/*
	关闭redis
*/
func (cacheRedis *RedisHander) Close() {
	if cacheRedis.Pool != nil {
		cacheRedis.Pool.Close()
	}
}
