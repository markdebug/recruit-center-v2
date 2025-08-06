package model

import "time"

// Dict 字典表
type Dict struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	ParentID  uint       `gorm:"default:0;index" json:"parentId"`        // 父级ID，0表示顶级
	Category  string     `gorm:"size:50;not null;index" json:"category"` // 字典分类
	Code      string     `gorm:"size:50;not null" json:"code"`           // 字典编码
	Name      string     `gorm:"size:100;not null" json:"name"`          // 字典名称
	Value     string     `gorm:"size:255" json:"value"`                  // 字典值
	Sort      int        `gorm:"default:0" json:"sort"`                  // 排序
	Status    int        `gorm:"default:1" json:"status"`                // 状态 1-启用 0-禁用
	Remarks   string     `gorm:"size:255" json:"remarks"`                // 备注
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"`

	Children []Dict `gorm:"-" json:"children,omitempty"` // 子项列表
}

// TableName 指定表名
func (Dict) TableName() string {
	return "t_rc_dict"
}
