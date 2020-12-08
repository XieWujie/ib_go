package db

const COV_TYPE_SINGLE = 0
const COV_TYPE_ROOM = 1

type Conversation struct {
	ConversationId   int          `json:"conversationId" xorm:"pk autoincr"`
	Members          []MemberInfo `json:"members" xorm:"json"`
	ConversationType int          `json:"conversation_type"`
}

func (c *Conversation) Save() error {
	_, err := engine.InsertOne(c)
	return err
}

func (c *Conversation) Update() error {
	_, err := engine.ID(c.ConversationId).Update(c)
	return err
}
func (c *Conversation) Get() error {
	_, err := engine.Get(c)
	return err
}
