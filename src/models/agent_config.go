package models

type AgentConfig struct {
	Id       int64  `json:"id" gorm:"primaryKey"`
	TargetId int64  `json:"target_id"`
	GroupId  int64  `json:"group_id"` // busi group id
	Version  string `json:"version"`
	Instance string `json:"instance"` //instance
	Ident    string `json:"ident"`
	Lables   string `json:"lables"`
	Format   string `json:"format"`
	Config   string `json:"config"`
}
