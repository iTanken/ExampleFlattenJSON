package test

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	data map[string]any
	want []map[string]any
)

func init() {
	bytesData, err := os.ReadFile("../data.json")
	if err != nil {
		log.Fatalln(`os.Open("data.json") got error: `, err)
	}
	if err = json.Unmarshal(bytesData, &data); err != nil {
		log.Fatalln("json.Unmarshal(bytesData, &data) got error: ", err)
	}

	bytesWant, err := os.ReadFile("../want.json")
	if err != nil {
		log.Fatalln(`os.Open("want.json") got error: `, err)
	}
	if err = json.Unmarshal(bytesWant, &want); err != nil {
		log.Fatalln("json.Unmarshal(bytesWant, &want) got error: ", err)
	}
}

func TestFlattenJSON(t *testing.T) {
	result := unwind(data)
	if result == nil {
		t.Fatalf("unwind() got nil")
	}

	if !reflect.DeepEqual(result, want) {
	    jsonBytes, _ := json.Marshal(result)
		t.Fatalf("unwind() got result:\n%s", jsonBytes)
	}
}

// 通过深度优先搜索遍历给定的 map 结构，并将其展开为一个包含多个层级 map 的切片。
// 这个函数主要处理复杂的数据结构，其中包含嵌套的 map 和切片。
//
//   - root 入参：需要展开的嵌套 map 结构
//   - result 返回值：包含展开后的多个 map 的切片
func unwind(root map[string]any) (result []map[string]any) {
	// 用于追踪当前遍历状态的切片，包括路径和中间结果。初始化为包含一个空的 map
	states := []map[string]any{make(map[string]any)}
	// 用于记录当前遍历的路径的切片，初始化为根路径 $
	paths := []string{"$"}

	// dfs 是一个递归函数，用于深度优先搜索遍历 map 结构
	//
	// node：当前遍历的 map 节点
	// inherit：是否继承父状态
	var dfs func(node map[string]any, inherit bool)
	dfs = func(node map[string]any, inherit bool) {
		// 获取当前状态 map
		state := states[len(states)-1]
		if !inherit {
			// 如果不继承父状态，则创建一个新的状态 map 并将其添加到 states 切片中
			newState := make(map[string]any)
			for k, v := range state {
				newState[k] = v
			}
			states = append(states, newState)
			state = newState
		}
		// 遍历当前节点的所有键值对
		for k, v := range node {
			if _, ok := v.([]any); ok {
				// 忽略任何值为切片的键值对
				continue
			}
			// 将当前键添加到路径中
			paths = append(paths, k)
			switch v := v.(type) {
			case map[string]any:
				// 对于值为 map 的情况，递归调用 dfs
				dfs(v, true)
			default:
				// 否则将当前键值对添加到状态 map 中
				state[strings.Join(paths, ".")] = v
			}
			// 回溯，移除当前键路径，返回上一级节点
			paths = paths[:len(paths)-1]
		}
		// 检查当前节点是否包含任何切片值，如果包含，isParent 设置为 true
		isParent := false
		for _, v := range node {
			if _, ok := v.([]any); ok {
				isParent = true
				break
			}
		}
		// 如果存在切片值，处理切片中的每个元素
		if isParent {
			for k, v := range node {
				if list, ok := v.([]any); ok {
					for _, vv := range list {
						// 为切片元素添加路径标记
						paths = append(paths, k+"[*]")
						if vvMap, ok := vv.(map[string]any); ok {
							// 递归调用 dfs 处理嵌套 map
							dfs(vvMap, false)
						}
						// 回溯，移除切片元素路径标记
						paths = paths[:len(paths)-1]
					}
				}
			}
		} else if !inherit {
			// 如果当前节点不包含切片值且不继承父状态，将当前状态 map 添加到结果切片 result 中
			newState := make(map[string]any)
			for k, v := range state {
				newState[k] = v
			}
			result = append(result, newState)
		}
		// 如果不继承父状态，回溯并移除当前状态 map
		if !inherit {
			states = states[:len(states)-1]
		}
	}
	// 从根节点开始递归调用 dfs
	dfs(root, true)
	return result
}
