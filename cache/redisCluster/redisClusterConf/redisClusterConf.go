// Copyright (c) 2018-2023 yejiantai
// 基于config的配置文件读取情况
// 配置文件都存储在conf文件夹下
package redisClusterConf

import (
	"fmt"
	"strings"
	"time"

	"github.com/liudanyejiantai/gosdk/datatype"

	"github.com/go-redis/redis"
	"github.com/liudanyejiantai/gosdk/file"
	"github.com/liudanyejiantai/gosdk/string_func"
	"github.com/robfig/config"
)

type RedisClusterConf struct {
	Conf *config.Config
	// 配置文件名称
	name string
	// 配置文件路径
	Dir string
	// 初始化
	bInit bool
}

func (c *RedisClusterConf) notInitError() error {
	return fmt.Errorf("conf is not create by NewConfig")
}
func (c *RedisClusterConf) ReadAllConf() ([]redis.ClusterSlot, error) {
	if !c.bInit {
		return nil, c.notInitError()
	}

	// 读取到全部的节点名称,每个节点对应一个配置文件，这样操作修改起来很简单
	arraySections, err := c.Conf.SectionOptions("conf")
	if err != nil {
		return nil, fmt.Errorf("get SectionOptions error,%s", err.Error())
	}
	var arrayClusterSlot []redis.ClusterSlot
	for i := 0; i < len(arraySections); i++ {
		d, err := readSecitonConf(arraySections[i])
		if err == nil {
			arrayClusterSlot = append(arrayClusterSlot, *d)
		} else {
			fmt.Printf("readSecitonConf(arraySections[%d]) faild, %s\n", i, err.Error())
		}
	}

	return arrayClusterSlot, nil
}

// 读取section的全部配置信息
func readSecitonConf(name string) (*redis.ClusterSlot, error) {
	confName := name + ".conf"
	if !file.FileExists(file.GetCurDir() + "redisClusterConf/" + confName) {
		return nil, fmt.Errorf(" config file [%s] is not exists", file.GetCurDir()+"redisClusterConf/"+confName)
	}
	c, err := NewConfig(confName)
	if err != nil {
		return nil, err
	}

	section := "conf"
	if !c.Conf.HasSection(section) {
		return nil, fmt.Errorf(" config file [%s] section [%s] not exists", confName, section)
	}

	var s redis.ClusterSlot
	var temp string

	// slot 范围值类似0-2047
	temp, err = c.Conf.String(section, "master_slots_range")
	if err != nil {
		return nil, fmt.Errorf("config file [%s] section[%s] key[master_slots_range] faild, %s",
			confName, section, err.Error())
	}
	temp = string_func.Trim(temp)
	arraySlots := strings.Split(temp, "-")
	if len(arraySlots) != 2 {
		return nil, fmt.Errorf("config file [%s] section[%s] key[master_slots_range] faild, %s",
			confName, section, err.Error())
	}
	if s.Start, err = datatype.StringToInt(arraySlots[0]); err != nil {
		return nil, fmt.Errorf("config file [%s] section[%s] key[master_slots_range] faild, %s",
			confName, section, err.Error())
	}
	if s.End, err = datatype.StringToInt(arraySlots[1]); err != nil {
		return nil, fmt.Errorf("config file [%s] section[%s] key[master_slots_range] faild, %s",
			confName, section, err.Error())
	}

	var masterNode redis.ClusterNode
	masterNode.Addr, err = c.Conf.String(section, "master_endpoint")
	if err != nil {
		return nil, fmt.Errorf("config file [%s] section[%s] key[master_endpoint] faild, %s",
			confName, section, err.Error())
	}
	if masterNode.Addr == "" {
		return nil, fmt.Errorf("config file [%s] section[%s] key[master_endpoint] is empty",
			confName, section)
	}

	s.Nodes = append(s.Nodes, masterNode)

	// 如果有配置从机器节点的
	temp, err = c.Conf.String(section, "slaver_endpoints")
	if err == nil {
		temp = string_func.Trim(temp)
		if temp != "" {
			arrayTemp := strings.Split(temp, ";")
			for i := 0; i < len(arrayTemp); i++ {
				var node redis.ClusterNode = redis.ClusterNode{Addr: arrayTemp[i]}
				//fmt.Println("node.Addr", node.Addr)
				s.Nodes = append(s.Nodes, node)
			}
		}

	}

	return &s, nil
}

func NewConfig(fileName string) (*RedisClusterConf, error) {
	var (
		conf RedisClusterConf
		err  error
	)
	conf.name, conf.Dir = fileName, file.GetCurDir()+"redisClusterConf/"
	if err = file.CreateDirTree(conf.Dir); err != nil {
		return nil, err
	}
	c := config.NewDefault()

	if !file.FileExists(conf.Dir + conf.name) {
		defaultStr := fmt.Sprintf("config info for eAgent Auto create at %s", time.Now().Format("2006-01-02 15:04:05"))
		if err = c.WriteFile(conf.Dir+conf.name, 0644, defaultStr); err != nil {
			return nil, fmt.Errorf("create config file [%s] faild, %s", conf.Dir+conf.name, err.Error())
		}
	} else {
		if c, err = config.ReadDefault(conf.Dir + conf.name); err != nil {
			return nil, fmt.Errorf("open config file [%s] faild, %s", conf.Dir+conf.name, err.Error())
		}
	}

	conf.Conf, conf.bInit = c, true

	return &conf, nil
}
