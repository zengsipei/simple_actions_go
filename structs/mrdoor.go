package structs

type Mrdoor struct {
	Uuid            string      `json:"uuid"`
	Password        string      `json:"password"`
	New             int         `json:"new"`
	Id              int         `json:"id"`
	Type            int         `json:"type"`
	OrderId         int         `json:"orderId"`
	UserId          int         `json:"userId"`
	Server          []int       `json:"server"`
	Port            int         `json:"port"`
	Key             interface{} `json:"key"`
	Data            string      `json:"data"`
	Status          int         `json:"status"`
	AutoRemove      int         `json:"autoRemove"`
	AutoRemoveDelay int         `json:"autoRemoveDelay"`
	MultiServerFlow int         `json:"multiServerFlow"`
	Active          int         `json:"active"`
	User            interface{} `json:"user"`
	Expire          int64       `json:"expire"`
	From            int64       `json:"from"`
	To              int64       `json:"to"`
	Flow            int         `json:"flow"`
	FlowPack        int         `json:"flowPack"`
	FlowOfServer    int         `json:"flowOfServer"`
	PublicKey       string      `json:"publicKey"`
	PrivateKey      string      `json:"privateKey"`
	Username        string      `json:"username"`
}
