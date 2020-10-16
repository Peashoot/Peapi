package download

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/peashoot/peapi/config"
	uuid "github.com/satori/go.uuid"
)

// Aria2Rpc 创建AriaRpc请求
// 文档地址 https://aria2.github.io/manual/en/html/aria2c.html#rpc-interface
func Aria2Rpc(method string, firstParam ...interface{}) (map[string]interface{}, error) {
	reqMap := make(map[string]interface{})
	reqMap["jsonrpc"] = "2.0"
	reqMap["id"] = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	reqMap["method"] = method
	paramArray := make([]interface{}, 1)
	if config.Config.DownloadConfig.Aria2Secret != "" {
		firstParam = append([]interface{}{fmt.Sprintf("token:%s", config.Config.DownloadConfig.Aria2Secret)}, firstParam)
	}
	paramArray[0] = firstParam
	var respMap map[string]interface{}
	var jsonBytes []byte
	var err error
	if jsonBytes, err = json.Marshal(reqMap); err != nil {
		return respMap, err
	}
	client := http.DefaultClient
	if req, err := http.NewRequest("GET", config.Config.DownloadConfig.Aria2RpcURL, strings.NewReader(string(jsonBytes))); err != nil {
		return respMap, err
	} else if resp, err := client.Do(req); err != nil {
		return respMap, err
	} else if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return respMap, err
	} else {
		if err := json.Unmarshal(bodyBytes, &respMap); err != nil {
			return respMap, err
		}
		return respMap, nil
	}
}

// AddURI 创建下载任务
func AddURI(link string) (string, error) {
	respMap, err := Aria2Rpc("aria2.addUri", []string{link})
	if err != nil {
		return "", err
	}
	var gid string
	var ok bool
	if gid, ok = respMap["result"].(string); !ok {
		return "", ErrNoMatchedResult
	}
	return gid, nil
}

// AddTorrent 添加种子任务
func AddTorrent(torrent string) (string, error) {
	fileBytes, err := ioutil.ReadFile(torrent)
	if err != nil {
		return "", err
	}
	fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)
	var respMap map[string]interface{}
	if respMap, err = Aria2Rpc("aria2.addTorrent", fileBase64); err != nil {
		return "", err
	}
	var gid string
	var ok bool
	if gid, ok = respMap["result"].(string); !ok {
		return "", ErrNoMatchedResult
	}
	return gid, nil
}

// AddMetalink 添加metalink文件
func AddMetalink(metalink string) (string, error) {
	fileBytes, err := ioutil.ReadFile(metalink)
	if err != nil {
		return "", err
	}
	fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)
	var respMap map[string]interface{}
	if respMap, err = Aria2Rpc("aria2.addMetalink", fileBase64); err != nil {
		return "", err
	}
	var gid string
	var ok bool
	if gid, ok = respMap["result"].(string); !ok {
		return "", ErrNoMatchedResult
	}
	return gid, nil
}

// operateWithGid aria有关gid的操作
func operateWithGid(method, gid string) (string, error) {
	respMap, err := Aria2Rpc(method, gid)
	if err != nil {
		return "", err
	}
	var ok bool
	if gid, ok = respMap["result"].(string); !ok {
		return "", ErrNoMatchedResult
	}
	return gid, nil
}

// operateWithNoParams aria不需要参数的操作
func operateWithNoParams(method string) (string, error) {
	respMap, err := Aria2Rpc(method)
	if err != nil {
		return "", err
	}
	var result string
	var ok bool
	if result, ok = respMap["result"].(string); !ok {
		return "", ErrNoMatchedResult
	}
	return result, nil
}

// Remove 根据gid移除任务
func Remove(gid string) (string, error) {
	return operateWithGid("aria2.remove", gid)
}

// ForceRemove 强制删除
func ForceRemove(gid string) (string, error) {
	return operateWithGid("aria2.forceRemove", gid)
}

// Pause 暂停任务
func Pause(gid string) (string, error) {
	return operateWithGid("aria2.pause", gid)
}

// PauseAll 暂停所有的任务
func PauseAll() (string, error) {
	return operateWithNoParams("aria2.pauseAll")
}

// ForcePause 强制暂停任务
func ForcePause(gid string) (string, error) {
	return operateWithGid("aria2.forcePause", gid)
}

// ForcePauseAll 强制暂停所有任务
func ForcePauseAll() (string, error) {
	return operateWithNoParams("aria2.forcePauseAll")
}

// Unpause 继续任务
func Unpause(gid string) (string, error) {
	return operateWithGid("aria2.unpause", gid)
}

// UnpauseAll 继续所有任务
func UnpauseAll() (string, error) {
	return operateWithNoParams("aria2.unpauseAll")
}

const (
	// ActiveStatus 当前下载/播种下载
	ActiveStatus = "active"
	// WaitingStatus 用于队列中的下载；下载未开始
	WaitingStatus = "waiting"
	// PausedStatus 暂停下载
	PausedStatus = "paused"
	// ErrorStatus 对于由于错误而停止的下载
	ErrorStatus = "error"
	// CompleteStatus 停止和完成下载
	CompleteStatus = "complete"
	// RemovedStatus 用户删除的下载
	RemovedStatus = "removed"
)

var (
	// ErrNoMatchedResult 未找到匹配结果错误
	ErrNoMatchedResult = errors.New("can not find correct result from response")
)

// MissionURI 任务URI
type MissionURI struct {
	URI    string `json:"uri"`    // URI
	Status string `json:"status"` // 如果已使用URI，则为“used”。如果URI仍在队列中，则“waiting”。
}

// MissionFile 任务文件信息
type MissionFile struct {
	Index           string       `json:"index"`           // 文件的索引（从1开始），与文件在多文件torrent中的显示顺序相同
	Path            string       `json:"path"`            // 文件路径
	Length          string       `json:"length"`          // 文件大小（以字节为单位）
	CompletedLength string       `json:"completedLength"` // 此文件的完整长度（以字节为单位）。请注意，的总和可能completedLength小于 方法completedLength返回的总和aria2.tellStatus()。这是因为，completedLength在 aria2.getFiles() 仅包括完成作品。在另一方面，completedLength 在aria2.tellStatus()还包括部分完成的块。
	Selected        string       `json:"selected"`        // true如果通过--select-file选项选择了此文件。如果 --select-file未指定，或者这是单文件洪流或根本不是洪流下载，则此值始终为true。否则 false。
	Uris            []MissionURI `json:"uris"`            // 返回此文件的URI列表
}

// MissionTorrent 任务种子信息
type MissionTorrent struct {
	AnnounceList []string           `json:"announceList"` // 公告URI列表列表
	Comment      string             `json:"comment"`      // 种子的评论
	CreationDate string             `json:"creationDate"` // 种子的创建时间。该值是自纪元以来的整数，以秒为单位
	Mode         string             `json:"mode"`         // torrent的文件模式。值为single或multi
	Info         MissionTorrentInfo `json:"info"`         // 包含信息字典中数据的结构
}

// MissionTorrentInfo 任务种子详细信息
type MissionTorrentInfo struct {
	Name string `json:"name"` // 信息字典中的名称
}

// MissionStatus 任务状态
type MissionStatus struct {
	GID                    string         `json:"gid"`                    // 下载的GID
	Status                 string         `json:"status"`                 // 下载状态
	TotalLength            string         `json:"totalLength"`            // 下载的总长度（以字节为单位）
	CompletedLength        string         `json:"completedLength"`        // 下载完成的长度（以字节为单位）
	UploadLength           string         `json:"uploadLength"`           // 上载长度（以字节为单位）
	BitField               string         `json:"bitfield"`               // 下载进度的十六进制表示, 当尚未开始下载时，此密钥将不包括在响应中
	DownloadSpeed          string         `json:"downloadSpeed"`          // 此下载的下载速度以字节/秒为单位
	UploadSpeed            string         `json:"uploadSpeed"`            // 此下载的上传速度（以字节/秒为单位）
	InfoHash               string         `json:"infoHash"`               // InfoHash(仅限BitTorrent)
	NumSeeders             string         `json:"numSeeders"`             // aria2所连接的播种机数量(仅限BitTorrent)
	Seeder                 string         `json:"seeder"`                 // true如果本地端点是播种者。否则false(仅限BitTorrent)
	PieceLength            string         `json:"pieceLength"`            // 片段长度（以字节为单位）
	NumPieces              string         `json:"numPieces"`              // 片段数量
	Connections            string         `json:"connections"`            // aria2已连接的对等/服务器数
	ErrorCode              string         `json:"errorCode"`              // 此项最后错误的代码（如果有）
	ErrorMessage           string         `json:"errorMessage"`           // 错误消息
	FollowedBy             []string       `json:"followedBy"`             // 下载结果生成的GID列表
	Following              string         `json:"following"`              // 和followedBy相对应, followedBy中包含该gid的任务gid
	BelongsTo              string         `json:"belongsTo"`              // 父级下载的GID, 如果此下载没有父项，则此密钥将不包含在响应中
	Dir                    string         `json:"dir"`                    // 保存文件的目录
	Files                  []MissionFile  `json:"files"`                  // 返回文件列表
	Bittorrent             MissionTorrent `json:"bittorrent"`             // 包含从.torrent（文件）中检索到的信息的结构(仅限BitTorrent)
	VerifiedLength         string         `json:"verifiedLength"`         // 在对文件进行哈希检查时，已验证的字节数。仅当对该下载进行哈希检查时，该密钥才存在
	VerifyIntegrityPending string         `json:"verifyIntegrityPending"` // true如果此下载正在等待队列中的哈希检查。仅当此下载在队列中时，此密钥才存在
}

// TellStatus 查询任务状态
func TellStatus(gid string, keys ...string) (status MissionStatus, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.tellStatus", gid, keys)
	if err != nil {
		return
	}
	result, ok := respMap["result"].(map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &status)
	return
}

// GetURIs 返回由gid（字符串）表示的下载中使用的URI
func GetURIs(gid string) (uriArray []MissionURI, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.getUris", gid)
	if err != nil {
		return
	}
	result, ok := respMap["result"].([]map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &uriArray)
	return
}

// GetFiles 返回以gid（字符串）表示的下载文件列表
func GetFiles(gid string) (files []MissionFile, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.getFiles", gid)
	if err != nil {
		return
	}
	result, ok := respMap["result"].([]map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &files)
	return
}

// MissionPeer gid（字符串）表示的下载列表的同级对象
type MissionPeer struct {
	PeerID        string `json:"peerId"`        // 百分比编码的对等ID
	IP            string `json:"ip"`            // 对端的IP地址
	Port          string `json:"port"`          // 对等体的端口号
	BitField      string `json:"bitfield"`      // 对等体下载进度的十六进制表示形式。最高位对应于索引为0的作品。置位表示作品可用，未置位表示作品缺失。最后的任何备用位都设置为零
	AmChoking     string `json:"amChoking"`     // 为true时表示主动屏蔽
	PeerChoking   string `json:"peerChoking"`   // 为true时表示对方屏蔽
	DownloadSpeed string `json:"downloadSpeed"` // 此客户端从对等方获得的下载速度（字节/秒）
	UploadSpeed   string `json:"uploadSpeed"`   // 此客户端上载到对等方的上载速度（字节/秒）
	Seeder        string `json:"seeder"`        // 为true时，此对等方是播种机
}

// GetPeers 返回由gid（字符串）表示的下载列表的同级对象
func GetPeers(gid string) (peers []MissionPeer, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.getPeers", gid)
	if err != nil {
		return
	}
	result, ok := respMap["result"].([]map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &peers)
	return
}

// MissionServer 连接的服务器信息
type MissionServer struct {
	Index   string              `json:"index"`   // 从1开始的文件索引，与文件在多文件metalink中出现的顺序相同
	Servers []MissionServerInfo `json:"servers"` // 服务器具体信息
}

// MissionServerInfo 下载服务器具体信息
type MissionServerInfo struct {
	URI           string `json:"uri"`           // 原始URI
	CurrentURI    string `json:"currentUri"`    // 这是当前用于下载的URI。如果涉及重定向，则currentUri和uri可能会不同
	DownloadSpeed string `json:"downloadSpeed"` // 下载速度（字节/秒）
}

// GetServers 返回以gid（字符串）表示的下载的当前连接的HTTP（S）/ FTP / SFTP服务器
func GetServers(gid string) (servers []MissionServer, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.getServers", gid)
	if err != nil {
		return
	}
	result, ok := respMap["result"].([]map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &servers)
	return
}

// TellActive 返回活动下载列表
func TellActive(keys ...string) (statusArray []MissionStatus, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.tellActive", keys)
	if err != nil {
		return
	}
	result, ok := respMap["result"].([]map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &statusArray)
	return
}

// TellWaiting 返回等待下载的列表，包括暂停的下载
func TellWaiting(offset, num int, keys ...string) (statusArray []MissionStatus, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.tellWaiting", offset, num, keys)
	if err != nil {
		return
	}
	result, ok := respMap["result"].([]map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &statusArray)
	return
}

// TellStopped 返回已停止下载的列表
func TellStopped(offset, num int, keys ...string) (statusArray []MissionStatus, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.tellStopped", offset, num, keys)
	if err != nil {
		return
	}
	result, ok := respMap["result"].([]map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &statusArray)
	return
}

// ChangePosition 更改队列中gid表示的下载位置
func ChangePosition(gid string, pos int, how string) (int, error) {
	respMap, err := Aria2Rpc("", gid, pos, how)
	if err != nil {
		return 0, err
	}
	result, ok := respMap["result"].(int)
	if !ok {
		return 0, ErrNoMatchedResult
	}
	return result, nil
}

// ChangeURI 修改下载链接
func ChangeURI(gid string, fileIndex int, delUris, addUris []string, position int) ([]int, error) {
	respMap, err := Aria2Rpc("aria2.changeUri", gid, fileIndex, delUris, addUris, position)
	if err != nil {
		return []int{}, err
	}
	result, ok := respMap["result"].([]int)
	if !ok {
		return []int{}, ErrNoMatchedResult
	}
	return result, nil
}

// MissionOption 任务配置项
type MissionOption struct {
	AllowOverwrite         bool `json:"allow-overwrite"`           // 是否允许重写
	AllowPieceLengthChange bool `json:"allow-piece-length-change"` // 允许分片长度改变
	AlwaysResume           bool `json:"always-resume"`             // 总是继续下载
	AsyncDNS               bool `json:"async-dns"`                 // 同步DNS
}

// GetOption 获取任务配置信息
func GetOption(gid string) (option MissionOption, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.getOption", gid)
	if err != nil {
		return
	}
	result, ok := respMap["result"].(map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &option)
	return
}

// ChangeOption 更改配置项
func ChangeOption(gid string, option MissionOption) (string, error) {
	respMap, err := Aria2Rpc("aria2.changeOption", gid, option)
	if err != nil {
		return "", err
	}
	result, ok := respMap["result"].(string)
	if !ok {
		return "", ErrNoMatchedResult
	}
	return result, nil
}

// GetGlobalOption 获取全局配置项
func GetGlobalOption(keys ...string) (option MissionOption, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.getGlobalOption", keys)
	if err != nil {
		return
	}
	result, ok := respMap["result"].(map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &option)
	return
}

// ChangeGlobalOption 修改全局配置
func ChangeGlobalOption(keys []string, option MissionOption) (string, error) {
	respMap, err := Aria2Rpc("aria2.changeGlobalOption", keys, option)
	if err != nil {
		return "", err
	}
	result, ok := respMap["result"].(string)
	if !ok {
		return "", ErrNoMatchedResult
	}
	return result, nil
}

// GlobalStat 全局状态
type GlobalStat struct {
	DownloadSpeed   string `json:"downloadSpeed"`   // 总体下载速度（字节/秒）
	UploadSpeed     string `json:"uploadSpeed"`     // 总体上传速度（字节/秒）
	NumActive       string `json:"numActive"`       // 活动下载数
	NumWaiting      string `json:"numWaiting"`      // 活动下载的次数
	NumStopped      string `json:"numStopped"`      // 当前会话中已停止的下载数。该值由--max-download-result选项限制
	NumStoppedTotal string `json:"numStoppedTotal"` // 当前会话中已停止的下载数量，但不受 该--max-download-result选项的限制

}

// GetGlobalStat 获取全局状态
func GetGlobalStat() (stat GlobalStat, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.getGlobalStat")
	if err != nil {
		return
	}
	result, ok := respMap["result"].(map[string]interface{})
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &stat)
	return
}

// PurgeDownloadResult 将完成/错误/删除的下载清除到可用内存中
func PurgeDownloadResult() (string, error) {
	return operateWithNoParams("aria2.purgeDownloadResult")
}

// RemoveDownloadResult 从内存中删除由gid表示的完成/错误/删除的下载
func RemoveDownloadResult(gid string) (string, error) {
	return operateWithGid("aria2.removeDownloadResult", gid)
}

// Aria2Version aria2的版本信息和启用的功能列表
type Aria2Version struct {
	Version        string   `json:"version"`         // aria2的版本号作为字符串
	EnableFeatures []string `json:"enabledFeatures"` // 启用的功能列表
}

// GetVerrsion 获取aria2的版本和已启用功能的列表
func GetVerrsion() (version Aria2Version, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.getVersion")
	if err != nil {
		return
	}
	result, ok := respMap["result"]
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &version)
	return
}

// Aria2Session 会话信息
type Aria2Session struct {
	SessionID string `json:"sessionId"` // 每次调用aria2时生成的会话ID
}

// GetSessionInfo 返回会话信息
func GetSessionInfo() (session Aria2Session, err error) {
	var respMap map[string]interface{}
	respMap, err = Aria2Rpc("aria2.getSessionInfo")
	if err != nil {
		return
	}
	result, ok := respMap["result"]
	if !ok {
		err = ErrNoMatchedResult
		return
	}
	err = mapstructure.Decode(result, &session)
	return
}

// ShutDown 关闭aria2
func ShutDown() (string, error) {
	return operateWithNoParams("aria2.shutdown")
}

// ForceShutdown 强制关闭aria2
func ForceShutdown() (string, error) {
	return operateWithNoParams("aria2.forceShutdown")
}

// SaveSession 将当前会话保存到该--save-session选项指定的文件中
func SaveSession() (string, error) {
	return operateWithNoParams("aria2.saveSession")
}

// MethodCallInfo 方法调用信息
type MethodCallInfo struct {
	MethodName string        `json:"methodName"` // 参数名
	Params     []interface{} `json:"params"`     // 方法参数
}

// MultiCall 将多个方法调用封装在单个请求中
func MultiCall(calls ...MethodCallInfo) ([]map[string]interface{}, error) {
	respMap, err := Aria2Rpc("system.multicall", calls)
	if err != nil {
		return nil, err
	}
	result, ok := respMap["result"].([]map[string]interface{})
	if !ok {
		return nil, ErrNoMatchedResult
	}
	return result, nil
}

// ListMethods 以字符串数组形式返回所有可用的RPC方法
func ListMethods() ([]string, error) {
	respMap, err := Aria2Rpc("system.listMethods")
	if err != nil {
		return nil, err
	}
	result, ok := respMap["result"].([]string)
	if !ok {
		return nil, ErrNoMatchedResult
	}
	return result, nil
}

// ListNotifications 以字符串数组形式返回所有可用的RPC通知
func ListNotifications() ([]string, error) {
	respMap, err := Aria2Rpc("system.listNotifications")
	if err != nil {
		return nil, err
	}
	result, ok := respMap["result"].([]string)
	if !ok {
		return nil, ErrNoMatchedResult
	}
	return result, nil
}
