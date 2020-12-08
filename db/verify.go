package db

type verifyState int

const Agree verifyState = 3
const Defy verifyState = 1
const NoAction verifyState = 2

type Verify struct {
	VerifyId   int         `json:"verifyId" xorm:"pk autoincr"`
	State      verifyState `json:"state"`
	UserFrom   int         `json:"userFrom"`
	UserTo     int         `json:"userTo"`
	VerifyInfo string      `json:"verifyInfo"`
	CreateAt   int64       `json:"createAt"`
}

func (verify *Verify) Save() error {
	_, err := engine.InsertOne(verify)
	return err
}

func (verify *Verify) Get() error {
	_, err := engine.Get(verify)
	return err
}

func (verify *Verify) UpdateState() error {
	_, err := engine.ID(verify.VerifyId).Update(verify)
	return err
}

func FindVerifyList(userId int) []map[string]interface{} {
	var result []Verify
	_ = engine.Where("user_to=?", userId).Find(&result)
	var list = make([]map[string]interface{}, len(result))
	for i, v := range result {
		var user = User{UserId: v.UserFrom}
		user.Get()
		var m = make(map[string]interface{})
		m["avatar"] = user.Avatar
		m["userId"] = user.UserId
		m["username"] = user.Username
		m["description"] = user.Description
		m["description"] = user.Description
		var r = make(map[string]interface{})
		r["user"] = m
		r["verify"] = v
		list[i] = r
	}
	return list
}
