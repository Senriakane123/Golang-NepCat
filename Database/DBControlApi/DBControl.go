package DBControlApi

import "gorm.io/gorm"

// 用户服务结构体
type DBcontrol struct {
	DB *gorm.DB
}

func (bcontrol *DBcontrol) Where(model interface{}, query interface{}, args ...interface{}) ([]interface{}, error) {
	// 使用 GORM 的 Where 方法执行查询
	result := bcontrol.DB.Where(query, args...).Find(model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model.([]interface{}), nil
}

// Get 查询数据库中所有记录
func (bcontrol *DBcontrol) Get(Model interface{}) (interface{}, error) {
	if result := bcontrol.DB.Find(&Model); result.Error != nil {
		return nil, result.Error
	} else {
		return Model, nil
	}
}

// Create 插入记录
func (bcontrol *DBcontrol) Create(model interface{}) error {
	if result := bcontrol.DB.Create(model); result.Error != nil {
		return result.Error
	}
	return nil
}

// Update 更新记录
func (bcontrol *DBcontrol) Update(model interface{}, updates map[string]interface{}) error {
	if result := bcontrol.DB.Model(model).Updates(updates); result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete 删除记录
func (bcontrol *DBcontrol) Delete(model interface{}) error {
	if result := bcontrol.DB.Delete(model); result.Error != nil {
		return result.Error
	}
	return nil
}
