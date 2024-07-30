package cache

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
	lineNumber := 0 // 初始化行号计数器

	for scanner.Scan() {
		lineNumber++ // 每次循环迭代时递增行号
		line := scanner.Text()
		parts := strings.Split(line, " ")

		if len(parts) > 0 && parts[0] == key { // 确保分割后的数组有足够的元素
			loadedValue = parts[1]
			// 将真正访问的 key 对应的数据放入老年区（假设已经在 Access 中处理）
			l.Old.Add(parts[0], parts[1], l.Young)

			// 打印行号，如果需要
			fmt.Printf("找到所需的数据在第 %d 行\n", lineNumber)

			break // 找到所需的数据，不需要继续扫描
		}
	}

	if loadedValue == "" {
		return nil
	}

	// 根据空间局部性原则，加载邻近的键
	nearbyKeys, err := l.nearbyKeysAndValues(lineNumber)

	if err != nil {

		return loadedValue
	}

	for _, meta := range nearbyKeys {
		// 这里需要调用 loadFromDisk 或其他方法来加载邻近键的数据
		// 并将它们放入青年区
		l.Young.Add(meta.Key, meta.Value)

	}

	return loadedValue
}

// nearbyKeysAndValues 获取指定行号附近的键值对，不包括目标行本身
func (l *LRUCache) nearbyKeysAndValues(targetLineNum int) ([]ItemCache, error) {
	filePath := "/home/zwm/go_projects/bash_algorithm/LRU/test/createFile.txt"

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var items []ItemCache
	lineNum := 0
	fmt.Println(targetLineNum, targetLineNum-5, targetLineNum+5)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(lineNum)
		if lineNum >= targetLineNum-5 && lineNum <= targetLineNum+5 {
			if lineNum == targetLineNum {
				lineNum++
				continue
			}
			parts := strings.SplitN(line, " ", 2)
			items = append(items, ItemCache{Key: parts[0], Value: parts[1]})
		}
		if lineNum > targetLineNum+5 {
			break
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// 检查是否成功找到数据
	if len(items) == 0 {
		return nil, errors.New("data not found")
	}

	return items, nil
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
