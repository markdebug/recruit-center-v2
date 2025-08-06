package dao

import (
	"gorm.io/gorm"
	"org.thinkinai.com/recruit-center/internal/model"
)

// DictDAO 字典数据访问对象
type DictDAO struct {
	db *gorm.DB
}

// NewDictDAO 创建字典DAO实例
func NewDictDAO(db *gorm.DB) *DictDAO {
	return &DictDAO{db: db}
}

// Create 创建字典
func (d *DictDAO) Create(dict *model.Dict) error {
	return d.db.Create(dict).Error
}

// Update 更新字典
func (d *DictDAO) Update(dict *model.Dict) error {
	return d.db.Save(dict).Error
}

// Delete 删除字典
func (d *DictDAO) Delete(id uint) error {
	return d.db.Delete(&model.Dict{}, id).Error
}

// GetByID 根据ID获取字典
func (d *DictDAO) GetByID(id uint) (*model.Dict, error) {
	var dict model.Dict
	err := d.db.First(&dict, id).Error
	if err != nil {
		return nil, err
	}
	return &dict, nil
}

// List 获取字典列表
func (d *DictDAO) List(page, size int) ([]model.Dict, int64, error) {
	var total int64
	var dicts []model.Dict

	err := d.db.Model(&model.Dict{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && size > 0 {
		offset := (page - 1) * size
		err = d.db.Offset(offset).Limit(size).Find(&dicts).Error
	} else {
		err = d.db.Find(&dicts).Error
	}

	return dicts, total, err
}

// ListByCategory 根据分类获取字典列表
func (d *DictDAO) ListByCategory(category string) ([]model.Dict, error) {
	var dicts []model.Dict
	err := d.db.Where("category = ?", category).Find(&dicts).Error
	return dicts, err
}

// ListTree 获取字典树形结构
func (d *DictDAO) ListTree(category string) ([]model.Dict, error) {
	var dicts []model.Dict

	// 先获取所有数据
	err := d.db.Where("category = ?", category).Order("sort").Find(&dicts).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	return buildDictTree(dicts), nil
}

// buildDictTree 构建字典树形结构
func buildDictTree(dicts []model.Dict) []model.Dict {
	dictMap := make(map[uint]*model.Dict)
	var result []model.Dict

	// 构建字典map
	for i := range dicts {
		dict := dicts[i]
		dictMap[dict.ID] = &dict
	}

	// 构建树形结构
	for _, dict := range dictMap {
		if dict.ParentID == 0 {
			result = append(result, *dict)
		} else {
			if parent, ok := dictMap[dict.ParentID]; ok {
				parent.Children = append(parent.Children, *dict)
			}
		}
	}

	return result
}
