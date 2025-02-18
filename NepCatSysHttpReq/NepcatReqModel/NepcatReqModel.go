package NepcatReqModel

// 定义 JSON 结构体
type ConfigResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Config `json:"data"`
}

type Config struct {
	Network             Network `json:"network"`
	MusicSignUrl        string  `json:"musicSignUrl"`
	EnableLocalFile2Url bool    `json:"enableLocalFile2Url"`
	ParseMultMsg        bool    `json:"parseMultMsg"`
}

type Network struct {
	HTTPServers      []HTTPServer      `json:"httpServers"`
	HTTPSseServers   []interface{}     `json:"httpSseServers"` // 空数组
	HTTPClients      []HTTPClient      `json:"httpClients"`
	WebsocketServers []WebsocketServer `json:"websocketServers"`
	WebsocketClients []interface{}     `json:"websocketClients"` // 空数组
	Plugins          []interface{}     `json:"plugins"`          // 空数组
}

type HTTPServer struct {
	Enable            bool   `json:"enable"`
	Name              string `json:"name"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	EnableCors        bool   `json:"enableCors"`
	EnableWebsocket   bool   `json:"enableWebsocket"`
	MessagePostFormat string `json:"messagePostFormat"`
	Token             string `json:"token"`
	Debug             bool   `json:"debug"`
}

type HTTPClient struct {
	Enable            bool   `json:"enable"`
	Name              string `json:"name"`
	URL               string `json:"url"`
	ReportSelfMessage bool   `json:"reportSelfMessage"`
	MessagePostFormat string `json:"messagePostFormat"`
	Token             string `json:"token"`
	Debug             bool   `json:"debug"`
}

type WebsocketServer struct {
	Enable               bool   `json:"enable"`
	Name                 string `json:"name"`
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	ReportSelfMessage    bool   `json:"reportSelfMessage"`
	EnableForcePushEvent bool   `json:"enableForcePushEvent"`
	MessagePostFormat    string `json:"messagePostFormat"`
	Token                string `json:"token"`
	Debug                bool   `json:"debug"`
	HeartInterval        int    `json:"heartInterval"`
}

// 定义结构体
type AuthResponse struct {
	Code    int      `json:"code"`
	Data    AuthData `json:"data"`
	Message string   `json:"message"`
}

type AuthData struct {
	Credential string `json:"Credential"`
}
