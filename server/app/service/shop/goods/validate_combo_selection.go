package goods

import (
	"fmt"
)

// ValidateComboService 验证套餐选择服务
type ValidateComboService struct{}

// Execute 执行套餐选择验证
func (s *ValidateComboService) Execute(comboId int64, selections []*ComboSelection) error {
	// 1. 获取套餐配置
	configService := &GetComboConfigService{}
	config, err := configService.Execute(comboId)
	if err != nil {
		return err
	}

	// 2. 按组分类选择
	groupSelections := make(map[int64][]*ComboSelection)
	for _, selection := range selections {
		groupSelections[selection.GroupId] = append(groupSelections[selection.GroupId], selection)
	}

	// 3. 验证每个组的选择
	for _, group := range config.Groups {
		selections := groupSelections[group.Id]

		// 检查必选组
		if group.IsRequired && len(selections) == 0 {
			return fmt.Errorf("套餐组 %s 是必选的", group.Name)
		}

		// 检查选择数量
		if len(selections) < group.MinSelect {
			return fmt.Errorf("套餐组 %s 至少需要选择 %d 个选项", group.Name, group.MinSelect)
		}

		if len(selections) > group.MaxSelect {
			return fmt.Errorf("套餐组 %s 最多只能选择 %d 个选项", group.Name, group.MaxSelect)
		}

		// 检查选项有效性
		for _, selection := range selections {
			found := false
			for _, option := range group.Options {
				if option.Id == selection.SpecValueId {
					found = true
					// 检查库存
					if option.Stock < selection.Quantity {
						return fmt.Errorf("选项 %s 库存不足", option.Name)
					}
					break
				}
			}
			if !found {
				return fmt.Errorf("无效的套餐选项")
			}
		}
	}

	return nil
}
