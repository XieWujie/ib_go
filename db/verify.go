package db

type verifyState int

const Agree verifyState = 0
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

func FindVerifyList(userId int) []Verify {
	var result []Verify
	_ = engine.Where("user_to=?", userId).Find(&result)
	return result
}
