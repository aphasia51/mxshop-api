package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` // 自定validator
	Type   uint   `form:"type"   json:"type"   binding:"required,oneof=1 2"`
}
