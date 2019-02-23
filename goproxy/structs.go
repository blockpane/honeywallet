package goproxy

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	WalletUrl          string
	WalletAddr         string
	WalletStatus       = `Unlocked`
	LogChan            = make(chan LogEntry)
	GethUrl            = `http://10.99.172.2:8545/`
	IPRateLimit        = make(map[string]*RateLimit)
	IPRateLimitChannel = make(chan IpUpdate)
	MaxRequests        = 20
	WaitTime           = 3600
)

func randomAccount() string {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	list := []string{
		// just a bunch of high-value accounts grabbed from etherscan.io:
		`3f82128150b804d3b4bc75fc76f322fc2fd6da0e`,
		`483f5d9a15c1f9222ba65e5e2e93760355c795f7`,
		`a7e4fecddc20d83f36971b67e13f1abc98dfcfa6`,
		`eec606a66edb6f497662ea31b5eb1610da87ab5f`,
		`cea2b9186ece677f9b8ff38dc8ff792e9a9e7f8a`,
		`850c0224f37f67c471e860375ac8e39fea61e8b0`,
		`fbb1b73c4f0bda4f67dca266ce6ef42f520fbb98`,
		`3bfc20f0b9afcace800d73d2191166ff16540258`,
		`3bf86ed8a3153ec933786a02ac090301855e576b`,
		`69c6dcc8f83b196605fa1076897af0e7e2b6b044`,
		`0a4c79ce84202b03e95b7a692e5d728d83c44c76`,
		`2b6ed29a95753c3ad948348e3e7b1a251080ffb9`,
		`db6fd484cfa46eeeb73c71edee823e4812f9e2e1`,
		`35a85499ccf7c8505a88e23017b745c671cc5aaf`,
		`0b2708b52b7f248c69c5d43002b243aa249a4aac`,
		`6cc5f688a315f3dc28a7781717a9a798a59fda7b`,
		`75fddbf107ce88ba66ed7fd8182ee6c3a39e9420`,
		`f4a2eff88a408ff4c4550148151c33c93442619e`,
		`22650fcf7e175ffe008ea18a90486d7ba0f51e41`,
		`e3ecccd6c67da25871fc5ff9a32a6f5c379167a6`,
		`73bceb1cd57c711feac4224d062b0f6ff338501e`,
		`a0e239b0abf4582366adaff486ee268c848c4409`,
		`1e2fcfd26d36183f1a5d90f0e6296915b02bcb40`,
		`8cf23cd535a240eb0ab8667d24eedbd9eccd5cba`,
		`d69b0089d9ca950640f5dc9931a41a5965f00303`,
		`036b96eea235880a9e82fb128e5f6c107dfe8f57`,
		`851b7f3ab81bd8df354f0d7640efcd7288553419`,
		`bf3aeb96e164ae67e763d9e050ff124e7c3fdd28`,
		`af10cc6c50defff901b535691550d7af208939c5`,
		`07ee55aa48bb72dcc6e9d78256648910de513eca`,
		`1c4b70a3968436b9a0a9cf5205c787eb81bb558c`,
		`cafe1a77e84698c83ca8931f54a755176ef75f2c`,
		`04f2894d12662f2728d02b74ea10056c11467dba`,
		`a646e29877d52b9e2de457eca09c724ff16d0a2b`,
		`5b5b69f4e0add2df5d2176d7dbd20b4897bc7ec4`,
		`6812f391fd38375316f6613ee1b46b77ad846c52`,
		`ce127836a13e235ff191767c290b584417149e71`,
		`9295022fb35b28c65f7efaf678254823d3acbe53`,
		`167a9333bf582556f35bd4d16a7e80e191aa6476`,
		`3cceb0d443ca4b1320ae4fa60a053eac163ca512`,
		`c78310231aa53bd3d0fea2f8c705c67730929d8f`,
		`4baf012726cb5ec7dda57bc2770798a38100c44d`,
		`c4cf565a5d25ee2803c9b8e91fc3d7c62e79fe69`,
		`e04cf52e9fafa3d9bf14c407afff94165ef835f7`,
		`4b4a011c420b91260a272afd91e54accdafdfc1d`,
		`77afe94859163abf0b90725d69e904ea91446c7b`,
		`19d599012788b991ff542f31208bab21ea38403e`,
		`ca582d9655a50e6512045740deb0de3a7ee5281f`,
		`d05e6bf1a00b5b4c9df909309f19e29af792422b`,
		`0f00294c6e4c30d9ffc0557fec6c586e6f8c3935`,
		`eb2b00042ce4522ce2d1aacee6f312d26c4eb9d6`,
		`7ae92148e79d60a0749fd6de374c8e81dfddf792`,
		`554f4476825293d4ad20e02b54aca13956acc40a`,
		`9cf36e93a8e2b1eaaa779d9965f46c90b820048c`,
		`4756eeebf378046f8dd3cb6fa908d93bfd45f139`,
		`091933ee1088cdf5daace8baec0997a4e93f0dd6`,
		`a8dcc0373822b94d7f57326be24ca67bafcaad6b`,
		`a0efb63be0db8fc11681a598bf351a42a6ff50e0`,
		`8b83b9c4683aa4ec897c569010f09c6d04608163`,
		`550cd530bc893fc6d2b4df6bea587f17142ab64e`,
		`828103b231b39fffce028562412b3c04a4640e64`,
		`367989c660881e1ca693730f7126fe0ffc0963fb`,
		`e35b0ef92452c353dbb93775e0df97cedf873c72`,
		`844ada2ed8ecd77a1a9c72912df0fcb8b8c495a7`,
		`0518f5bb058f6215a0ff5f4df54dae832d734e04`,
		`0e86733eab26cfcc04bb1752a62ec88e910b4cf5`,
		`0ff64c53d295533a37f913bb78be9e2adc78f5fe`,
		`b8b6fe7f357adeab950ac6c270ce340a46989ce1`,
		`eddf8eb4984cc27407a568cae1c78a1ddb0c2c1b`,
		`7145cfedcf479bede20c0a4ba1868c93507d5786`,
		`3ba25081d3935fcc6788e6220abcace39d58d95d`,
		`90a9e09501b70570f9b11df2a6d4f047f8630d6d`,
		`7712bdab7c9559ec64a1f7097f36bc805f51ff1a`,
		`d65bd7f995bcc3bdb2ea2f8ff7554a61d1bf6e53`,
		`fc39f0dc7a1c5d5cd1cdf3b460d5fa99a56abf65`,
		`1ffedd7837bcbc53f91ad4004263deb8e9107540`,
		`024861e9f89d44d00a7ada4aa89fe03cab9387cd`,
		`d44023d2710dd7bef797a074ecec4fc74fdd52b2`,
		`1a71b118ac6c9086f43bcf2bb6ada3393be82a5c`,
		`657e46adad8be23d569ba3105d7a02124e8def97`,
		`73263803def2ac8b1f8a42fac6539f5841f4e673`,
		`6047a74d635262fb73ebce6c12bb6b14b3da70b4`,
		`78b96178e7ae1ff9adc5d8609e000811657993c8`,
		`81153b940932f49f42e5719fe9d1ec04e0e5c119`,
		`40f0d6fb7c9ddd9cbc1c02a208380c08cf77189b`,
		`3f5ce5fbfe3e9af3971dd833d26ba9b5c936f0be`,
		`bfc868b0c0af3885389a2242a2afdb841b78812f`,
		`2fa9f9efc767650aace0422668444c3ff63e1f8d`,
		`d57479b8287666b44978255f1677e412d454d4f0`,
		`07c62a47ebe0fa853bb83375e488896ce71266df`,
		`840760aed6bbd878c46c5850d3af0a61afcd09c8`,
		`42ada615203749550a51a0678b8e7d5f853c6a03`,
	}
	return list[r.Intn(len(list))]
}

func UpdateRand() {
	for {
		randAccount := randomAccount()
		WalletUrl = fmt.Sprintf(`keystore:///home/ethereum/.ethereum/keystore/UTC--%s--%s`,
			strings.Replace(time.Now().AddDate(0, -4, 0).String(), ` `, `-`, -1),
			randAccount)
		WalletAddr = `0x` + randAccount
		time.Sleep(time.Hour * 3)
	}
}

type FakeEthAccount struct {
	JsonRpc string   `json:"jsonrpc"`
	Id      int      `json:"id"`
	Result  []string `json:"result"`
}

func NewFakeAccount(id int) FakeEthAccount {
	return FakeEthAccount{
		JsonRpc: "2.0",
		Id:      id,
		Result:  []string{WalletAddr},
	}
}

type FakeWallet struct {
	JsonRpc string               `json:"jsonrpc"`
	Id      int                  `json:"id"`
	Result  []FakePersonalWallet `json:"result"`
}

type FakePersonalWallet struct {
	Url      string `json:"url"`
	Status   string `json:"status"`
	Accounts []FakeWalletAccounts
}

type FakeWalletAccounts struct {
	Address string `json:"address"`
	Url     string `json:"url"`
}

func NewFakeWallet(id int) FakeWallet {
	return FakeWallet{
		JsonRpc: "2.0",
		Id:      id,
		Result: []FakePersonalWallet{
			{
				Url:    WalletAddr,
				Status: WalletStatus,
				Accounts: []FakeWalletAccounts{
					{
						Address: WalletAddr,
						Url:     WalletUrl,
					},
				},
			},
		},
	}
}

type LogEntry struct {
	sync.RWMutex
	Ip       string                 `json:"ip"`
	Port     int64                  `json:"src_port"`
	UnixTime int64                  `json:"unixtime"`
	Request  map[string]interface{} `json:"request"`
	ReqBytes []byte                 `json:"-"`
}

func (le *LogEntry) Marshall() (j []byte, e error) {
	le.RLock()
	j, e = json.Marshal(le)
	le.RUnlock()
	return
}

func (le *LogEntry) MarshallRequest() (j []byte, e error) {
	le.RLock()
	j, e = json.Marshal(le.Request)
	le.RUnlock()
	return
}

type RateLimit struct {
	Counter    int
	AllowAfter int64
}

type Method struct {
	Ips      map[string]int `json:"ips"`
	Accounts map[string]int `json:"accounts"`
	LastSeen string         `json:"last_seen"`
}

type IpAddress struct {
	Methods  map[string]int `json:"methods"`
	Accounts map[string]int `json:"accounts"`
	LastSeen string         `json:"last_seen"`
}

type Address struct {
	Ips      map[string]int `json:"ips"`
	Methods  map[string]int `json:"methods"`
	LastSeen string         `json:"last_seen"`
}

type resultsStub struct {
	Method string        `json:"method"`
	Params []*paramsStub `json:"params"`
}

type paramsStub struct {
	To string `json:"to"`
}

type Stats struct {
	Ips       map[string]*IpAddress `json:"ips"`
	Methods   map[string]*Method    `json:"methods"`
	Addresses map[string]*Address
}

type IpUpdate struct {
	Ip     string
	Action string
}
