package cache

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// ItemCache 用于在列表和map中存储缓存项
type ItemCache struct {
	Key   interface{}
	Value interface{}
}

// LRUCache 包含老年区和青年区
type LRUCache struct {
	Old   *OldCache
	Young *YoungCache
	sync.RWMutex
}

func NewLRUCache() *LRUCache {
	return &LRUCache{
		Old:   NewOldCache(),
		Young: NewYoungCache(),
	}
}

// Access 访问缓存中的项，根据空间局部性决定放在老年区还是青年区
func (l *LRUCache) Access(key interface{}) (interface{}, error) {
	keyStr, ok := key.(string)
	if !ok {
		return nil, fmt.Errorf("invalid key type")
	}
	l.RLock()
	inOld, oldOk := l.Old.Items.Load(key)
	inYoung, youngOk := l.Young.Items.Load(key)
	l.RUnlock()

	if oldOk {
		// 如果数据在老年区，直接返回值
		return inOld, nil
	} else if youngOk {
		// 如果数据在青年区，可以考虑晋升到老年区
		// 假设 value 不需要，因为可以从青年区项中获取
		l.Young.PromoteToOld(keyStr, l.Old)
		return inYoung, nil
	} else {
		// 数据不在缓存中，需要从磁盘加载
		// 这里假设 loadFromDisk 是从磁盘加载数据的函数
		value := l.loadFromDisk(key.(string))

		if value != nil {
			return value, nil
		} else {
			return nil, errors.New("don`t find data")
		}

	}
}

// loadFromDisk 从磁盘加载数据到 LRU 缓存中，并考虑空间局部性
func (l *LRUCache) loadFromDisk(key string) interface{} {
	filePath := "/home/zwm/go_projects/bash_algorithm/LRU/test/createFile.txt" // 使用正确的文件路径
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var loadedValue string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if parts[0] == key { // extractKey 需要你根据文件格式定义一个解析键的函数
			loadedValue = parts[1]
			break // 找到所需的数据，不需要继续扫描
		}
	}
	if loadedValue == "" {
		return nil
	}

	// 将真正访问的 key 对应的数据放入老年区（假设已经在 Access 中处理）
	// 根据空间局部性原则，加载邻近的键
	nearbyKeys := l.nearbyKeys(key)
	for _, nearbyKey := range nearbyKeys {
		// 这里需要调用 loadFromDisk 或其他方法来加载邻近键的数据
		// 并将它们放入青年区
		nearbyValue := l.loadFromDisk(nearbyKey)
		if nearbyValue != nil {
			l.Young.Add(nearbyKey, nearbyValue)
		}
	}

	return loadedValue
}

// nearbyKeys 确定并返回给定键的邻近键，基于文件中的数据格式
func (l *LRUCache) nearbyKeys(baseKey string) []string {
	// 将 baseKey 转换为整数
	_, err := strconv.Atoi(baseKey)
	if err != nil {
		// 如果转换失败，返回空切片
		return []string{}
	}

	// 假设我们有一个函数获取所有键的列表
	allKeys, err := l.allKeys()
	if err != nil {
		// 如果获取所有键失败，返回空切片
		return []string{}
	}

	// 找到 baseKey 在 allKeys 中的索引
	baseIndex := sort.Search(len(allKeys), func(i int) bool {
		return allKeys[i] >= "12123"
	}) - 1 // 使用 -1 因为 Search 找到的是第一个大于或等于 baseNum 的索引

	// 定义一个函数来获取邻近键
	getNearbyKeys := func(index int, count int, reverse bool) []string {
		if index < 0 || index >= len(allKeys) {
			return []string{}
		}
		start := index
		end := index + count
		if reverse {
			start, end = end-1, start+1
		}
		return allKeys[start:end]
	}

	// 获取前面五个和后面五个邻近键
	nearbyBefore := getNearbyKeys(baseIndex, 5, true)
	nearbyAfter := getNearbyKeys(baseIndex+1, 5, false)

	// 合并前后邻近键，避免重复
	nearbyKeys := append(nearbyBefore, nearbyAfter...)
	// 去除 baseKey 本身
	if contains(nearbyKeys, baseKey) {
		nearbyKeys = removeElement(nearbyKeys, baseKey)
	}

	return nearbyKeys
}

func (l *LRUCache) allKeys() ([]string, error) {
	return nil, nil
}

// contains 检查切片中是否包含特定的元素
func contains(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// removeElement 从切片中移除特定的元素
func removeElement(slice []string, element string) []string {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// 解析文件中的每一行，提取键和值
func parseLine(line string) (key string, value string, err error) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		err = errors.New("invalid line format")
		return
	}

	key = parts[0]
	value, err = strconv.Unquote(`"` + parts[1] + `"`) // 假设值是被双引号包围的
	return
}

// 解析文件的起始行，获取所有键
func parseAllKeys(line string) map[string]struct{} {
	allKeys := make(map[string]struct{})
	// 假设第一行包含所有行号，用空格分隔
	for _, key := range strings.Fields(line) {
		allKeys[key] = struct{}{}
	}
	return allKeys
}
