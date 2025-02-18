package MusicGetLogic

import (
	"NepcatGoApiReq/MessageHandle/RandomMusic/MusicDBModel"
	"NepcatGoApiReq/MessageModel"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

type MusicManageHandle struct {
	Handler map[string]func(message MessageModel.Message)
}

func (n *MusicManageHandle) HandlerInit() {
	// 关键词映射到处理函数
	var groupManagekeywordHandlers = map[string]func(MessageModel.Message){
		"随机音乐推荐": n.RandomMusic,
	}
	n.Handler = groupManagekeywordHandlers
}

// **获取按长度排序的关键词**
func (n *MusicManageHandle) getSortedKeywords() []string {
	keys := make([]string, 0, len(n.Handler))
	for key := range n.Handler {
		keys = append(keys, key)
	}
	// 按字符串长度从长到短排序，保证 "解除全体禁言" 先匹配，而不是 "禁言"
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	return keys
}

// 统一处理消息
func (n *MusicManageHandle) HandleMusicManageMessage(message MessageModel.Message) bool {
	sortedKeywords := n.getSortedKeywords() // 获取按长度排序的关键词
	for _, keyword := range sortedKeywords {
		if strings.HasPrefix(message.RawMessage, keyword) || strings.Contains(message.RawMessage, keyword) {
			handler := n.Handler[keyword]
			handler(message)
			return true // 处理完一个就返回，避免重复触发
		}
	}
	return false
}

func (n *MusicManageHandle) RandomMusic(message MessageModel.Message) {
	GetMusicList(3426007951)
}

func GetMusicList(MusicListId int) *[]MusicDBModel.Song {
	test()
	// 设置 Edge 浏览器路径（Windows 默认路径）
	edgePath := `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`

	// 设定 chromedp 运行时的浏览器配置
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.ExecPath(edgePath))

	// 创建浏览器上下文
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 目标 URL（修改为你的歌单 ID）
	url := "https://y.qq.com/musicmac/v6/playlist/detail.html?id=7628192870"

	// 运行任务
	var rawSongs []string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`
			(() => {
				let songs = Array.from(document.querySelectorAll(".songlist__item")).map(song => {
					let songName = song.querySelector(".mod_songname__name")?.title || "Unknown Song";
					let singerName = song.querySelector(".singer_name")?.title || "Unknown Singer";
					let albumName = song.querySelector(".album_name")?.title || "Unknown Album";
					let songTime = song.querySelector(".songlist__time")?.textContent.trim() || "Unknown Time";
					let songId = song.getAttribute(".mod_songname__id") || "Unknown ID";
					return songName + "|" + singerName + "|" + albumName + "|" + songTime + "|" + songId;
				});
				return songs;
			})()
		`, &rawSongs),
	)

	if err != nil {
		fmt.Println("爬取失败:", err)
		return nil
	}

	// 解析数据到结构体列表
	var songList []MusicDBModel.Song
	for _, data := range rawSongs {
		parts := strings.Split(data, "|")
		if len(parts) == 5 {
			songList = append(songList, MusicDBModel.Song{
				Name:   parts[0],
				Singer: parts[1],
				Album:  parts[2],
				Time:   parts[3],
				ID:     parts[4],
			})
		}
	}

	return &songList

}

func getRandomSong(songs []MusicDBModel.Song) MusicDBModel.Song {
	rand.Seed(time.Now().UnixNano()) // 随机种子
	return songs[rand.Intn(len(songs))]
}

func test() {
	// 请求的 URL
	url := "https://c.y.qq.com/splcloud/fcgi-bin/fcg_musiclist_getmyfav.fcg?dirid=201&dirinfo=1&g_tk_new_20200303=194616139&g_tk=194616139&jsonpCallback=MusicJsonCallback254241073633829&loginUin=735439479&hostUin=0&format=jsonp&inCharset=GB2312&outCharset=utf-8&notice=0&platform=mac&needNewCode=0"

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("请求创建失败:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0")
	req.Header.Set("Referer", "https://y.qq.com/")
	req.Header.Set("Cookie", "pgv_pvid=6864890665; fqm_pvqid=85b6a837-ffe1-45e8-9316-0d0d00086cda; ts_refer=cn.bing.com/; ts_uid=4601151800; RK=Wc208KQueQ; ptcz=668e72878d0c6e110c9c7f829e68fe7bd568e958a273280d5857c4f1a812b2bb; fqm_sessionid=6f6c2570-5bec-40da-b1f4-46c2a962fa5d; pgv_info=ssid=s959611415; _qpsvr_localtk=0.44739241331807467; tmeLoginType=2; music_ignore_pskey=202306271436Hn@vBj; ptui_loginuin=735439479; login_type=1; wxunionid=; wxopenid=; euin=7iok7eoq7eSq; qqmusic_key=Q_H_L_63k3NcmLClb-mYilpJA1mjkLYp4YdDtVMvzQLGq6ii5fpG7oMAViG3DeNjU8fvyLn89CCSn5R73xP3d5mH2l5vg; uin=735439479; psrf_musickey_createtime=1739848193; psrf_access_token_expiresAt=1740452993; qm_keyst=Q_H_L_63k3NcmLClb-mYilpJA1mjkLYp4YdDtVMvzQLGq6ii5fpG7oMAViG3DeNjU8fvyLn89CCSn5R73xP3d5mH2l5vg; wxrefresh_token=; psrf_qqopenid=6139E80D0EC45DD70E5F6CB879BAE379; psrf_qqaccess_token=D8985D98B4A7770909A703B3AE071AF4; psrf_qqrefresh_token=E50E9485D4757E1D220FB400925341E1; psrf_qqunionid=41684BB895B27E605D8DD757AD2AAB15; ts_last=y.qq.com/musicmac/v6/playlist/detail.html")

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 检查是否是 gzip 格式
	var reader *gzip.Reader
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("创建 GZIP 解码器失败:", err)
			return
		}
		defer reader.Close()

		// 读取解压后的内容
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println("读取解压后的内容失败:", err)
			return
		}

		// 输出解压后的数据
		handleResponse(string(body))
	} else {
		// 如果不是 gzip 格式，直接读取内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取响应失败:", err)
			return
		}
		handleResponse(string(body))
	}
}

func handleResponse(jsonp string) {
	// 去掉回调函数名部分
	start := strings.Index(jsonp, "(") + 1
	end := strings.LastIndex(jsonp, ")")
	if start > 0 && end > start {
		jsonData := jsonp[start:end]
		fmt.Println("有效的 JSON 数据:", jsonData)
	} else {
		fmt.Println("无法提取有效的 JSON 数据")
	}
}
