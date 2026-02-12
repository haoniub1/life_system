package logic

import (
	"context"
	"fmt"

	"life-system-backend/internal/model"
	"life-system-backend/internal/realm"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

type ShopLogic struct {
	svcCtx *svc.ServiceContext
}

func NewShopLogic(svcCtx *svc.ServiceContext) *ShopLogic {
	return &ShopLogic{
		svcCtx: svcCtx,
	}
}

func (l *ShopLogic) GetShopItems(ctx context.Context, userID int64) (*types.ShopItemListResp, error) {
	items, err := l.svcCtx.ShopModel.GetItemsByUserID(userID)
	if err != nil {
		return nil, err
	}

	resp := &types.ShopItemListResp{
		Items: make([]types.ShopItemResp, 0),
	}

	for _, item := range items {
		resp.Items = append(resp.Items, types.ShopItemResp{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			SellPrice:   item.SellPrice,
			ItemType:    item.ItemType,
			Icon:        item.Icon,
			Image:       item.Image,
			Stock:       item.Stock,
		})
	}

	return resp, nil
}

func (l *ShopLogic) CreateShopItem(ctx context.Context, userID int64, req *types.CreateShopItemReq) (*types.ShopItemResp, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("商品名称不能为空")
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("价格不能为负数")
	}

	itemType := req.ItemType
	if itemType == "" {
		itemType = "consumable"
	}

	item := &model.ShopItem{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SellPrice:   req.SellPrice,
		ItemType:    itemType,
		Icon:        req.Icon,
		Image:       req.Image,
		Stock:       req.Stock,
	}

	id, err := l.svcCtx.ShopModel.CreateItem(item)
	if err != nil {
		return nil, err
	}

	return &types.ShopItemResp{
		ID:          id,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		SellPrice:   item.SellPrice,
		ItemType:    item.ItemType,
		Icon:        item.Icon,
		Image:       item.Image,
		Stock:       item.Stock,
	}, nil
}

func (l *ShopLogic) UpdateShopItem(ctx context.Context, userID int64, itemID int64, req *types.UpdateShopItemReq) (*types.ShopItemResp, error) {
	existing, err := l.svcCtx.ShopModel.GetItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("商品不存在")
	}
	if existing.UserID != userID {
		return nil, fmt.Errorf("无权修改此商品")
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.SellPrice != nil {
		existing.SellPrice = *req.SellPrice
	}
	if req.ItemType != nil {
		existing.ItemType = *req.ItemType
	}
	if req.Icon != nil {
		existing.Icon = *req.Icon
	}
	if req.Image != nil {
		existing.Image = *req.Image
	}
	if req.Stock != nil {
		existing.Stock = *req.Stock
	}

	if err := l.svcCtx.ShopModel.UpdateItem(existing); err != nil {
		return nil, err
	}

	return &types.ShopItemResp{
		ID:          existing.ID,
		Name:        existing.Name,
		Description: existing.Description,
		Price:       existing.Price,
		SellPrice:   existing.SellPrice,
		ItemType:    existing.ItemType,
		Icon:        existing.Icon,
		Image:       existing.Image,
		Stock:       existing.Stock,
	}, nil
}

func (l *ShopLogic) DeleteShopItem(ctx context.Context, userID int64, itemID int64) error {
	existing, err := l.svcCtx.ShopModel.GetItemByID(itemID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("商品不存在")
	}
	if existing.UserID != userID {
		return fmt.Errorf("无权删除此商品")
	}

	return l.svcCtx.ShopModel.DeleteItem(itemID, userID)
}

func (l *ShopLogic) PurchaseItem(ctx context.Context, userID int64, req *types.PurchaseItemReq) (*types.PurchaseResult, error) {
	item, err := l.svcCtx.ShopModel.GetItemByID(req.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("物品不存在")
	}

	if item.Stock != -1 && item.Stock < req.Quantity {
		return nil, fmt.Errorf("库存不足")
	}

	totalPrice := item.Price * req.Quantity

	stats, err := l.svcCtx.CharacterModel.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	// Check if user has enough spirit stones
	if stats.SpiritStones < totalPrice {
		return nil, fmt.Errorf("灵石不足！需要 %d 灵石，当前只有 %d 灵石", totalPrice, stats.SpiritStones)
	}

	// Deduct spirit stones
	stats.SpiritStones -= totalPrice

	if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
		return nil, err
	}

	if item.Stock != -1 {
		if err := l.svcCtx.ShopModel.UpdateItemStock(item.ID, req.Quantity); err != nil {
			return nil, err
		}
	}

	if err := l.svcCtx.ShopModel.AddToInventory(userID, item.ID, req.Quantity); err != nil {
		return nil, err
	}

	if err := l.svcCtx.ShopModel.RecordPurchase(userID, item.ID, item.Name, req.Quantity, totalPrice); err != nil {
		return nil, err
	}

	return &types.PurchaseResult{
		Success:              true,
		Message:              fmt.Sprintf("成功购买 %d 个「%s」", req.Quantity, item.Name),
		RemainingSpiritStones: stats.SpiritStones,
	}, nil
}

func (l *ShopLogic) GetInventory(ctx context.Context, userID int64) (*types.InventoryListResp, error) {
	inventoryItems, err := l.svcCtx.ShopModel.GetUserInventory(userID)
	if err != nil {
		return nil, err
	}

	resp := &types.InventoryListResp{
		Items: make([]types.InventoryItemResp, 0),
	}

	for _, invItem := range inventoryItems {
		item, err := l.svcCtx.ShopModel.GetItemByID(invItem.ItemID)
		if err != nil {
			continue
		}
		if item == nil {
			continue
		}

		resp.Items = append(resp.Items, types.InventoryItemResp{
			ID:          invItem.ID,
			ItemID:      invItem.ItemID,
			Name:        item.Name,
			Description: item.Description,
			ItemType:    item.ItemType,
			SellPrice:   item.SellPrice,
			Icon:        item.Icon,
			Image:       item.Image,
			Quantity:    invItem.Quantity,
		})
	}

	return resp, nil
}

func (l *ShopLogic) UseItem(ctx context.Context, userID int64, req *types.UseItemReq) (*types.UseItemResult, error) {
	invItem, err := l.svcCtx.ShopModel.GetInventoryItemByItemID(userID, req.ItemID)
	if err != nil {
		return nil, err
	}
	if invItem == nil || invItem.Quantity < req.Quantity {
		return nil, fmt.Errorf("物品不足")
	}

	item, err := l.svcCtx.ShopModel.GetItemByID(req.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("物品不存在")
	}

	stats, err := l.svcCtx.CharacterModel.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	// Load attributes
	attrs, err := l.svcCtx.CharacterModel.FindAttributesByUserID(userID)
	if err != nil {
		return nil, err
	}
	attrMap := make(map[string]*model.CharacterAttribute)
	for _, a := range attrs {
		attrMap[a.AttrKey] = a
	}

	message := ""
	switch item.Effect {
	case "fatigue_restore":
		stats.Fatigue -= item.EffectValue * req.Quantity
		if stats.Fatigue < 0 {
			stats.Fatigue = 0
		}
		message = fmt.Sprintf("恢复了 %d 点精力（降低疲劳）", item.EffectValue*req.Quantity)

	case "physique_boost":
		if attr, ok := attrMap["physique"]; ok {
			gain := float64(item.EffectValue * req.Quantity)
			result := realm.ProcessAttrGain(attr.Value, gain, attr.Realm, attr.RealmExp, attr.IsBottleneck, attr.AccumulationPool)
			attr.Value = result.NewValue
			attr.AccumulationPool = result.NewAccPool
			attr.RealmExp = result.NewRealmExp
			attr.IsBottleneck = result.NewIsBottleneck
			if err := l.svcCtx.CharacterModel.UpdateAttribute(attr); err != nil {
				return nil, err
			}
		}
		message = fmt.Sprintf("体魄提升了 %d 点", item.EffectValue*req.Quantity)

	case "willpower_boost":
		if attr, ok := attrMap["willpower"]; ok {
			gain := float64(item.EffectValue * req.Quantity)
			result := realm.ProcessAttrGain(attr.Value, gain, attr.Realm, attr.RealmExp, attr.IsBottleneck, attr.AccumulationPool)
			attr.Value = result.NewValue
			attr.AccumulationPool = result.NewAccPool
			attr.RealmExp = result.NewRealmExp
			attr.IsBottleneck = result.NewIsBottleneck
			if err := l.svcCtx.CharacterModel.UpdateAttribute(attr); err != nil {
				return nil, err
			}
		}
		message = fmt.Sprintf("意志提升了 %d 点", item.EffectValue*req.Quantity)

	case "intelligence_boost":
		if attr, ok := attrMap["intelligence"]; ok {
			gain := float64(item.EffectValue * req.Quantity)
			result := realm.ProcessAttrGain(attr.Value, gain, attr.Realm, attr.RealmExp, attr.IsBottleneck, attr.AccumulationPool)
			attr.Value = result.NewValue
			attr.AccumulationPool = result.NewAccPool
			attr.RealmExp = result.NewRealmExp
			attr.IsBottleneck = result.NewIsBottleneck
			if err := l.svcCtx.CharacterModel.UpdateAttribute(attr); err != nil {
				return nil, err
			}
		}
		message = fmt.Sprintf("智力提升了 %d 点", item.EffectValue*req.Quantity)

	case "perception_boost":
		if attr, ok := attrMap["perception"]; ok {
			gain := float64(item.EffectValue * req.Quantity)
			result := realm.ProcessAttrGain(attr.Value, gain, attr.Realm, attr.RealmExp, attr.IsBottleneck, attr.AccumulationPool)
			attr.Value = result.NewValue
			attr.AccumulationPool = result.NewAccPool
			attr.RealmExp = result.NewRealmExp
			attr.IsBottleneck = result.NewIsBottleneck
			if err := l.svcCtx.CharacterModel.UpdateAttribute(attr); err != nil {
				return nil, err
			}
		}
		message = fmt.Sprintf("感知提升了 %d 点", item.EffectValue*req.Quantity)

	case "charisma_boost":
		if attr, ok := attrMap["charisma"]; ok {
			gain := float64(item.EffectValue * req.Quantity)
			result := realm.ProcessAttrGain(attr.Value, gain, attr.Realm, attr.RealmExp, attr.IsBottleneck, attr.AccumulationPool)
			attr.Value = result.NewValue
			attr.AccumulationPool = result.NewAccPool
			attr.RealmExp = result.NewRealmExp
			attr.IsBottleneck = result.NewIsBottleneck
			if err := l.svcCtx.CharacterModel.UpdateAttribute(attr); err != nil {
				return nil, err
			}
		}
		message = fmt.Sprintf("魅力提升了 %d 点", item.EffectValue*req.Quantity)

	case "agility_boost":
		if attr, ok := attrMap["agility"]; ok {
			gain := float64(item.EffectValue * req.Quantity)
			result := realm.ProcessAttrGain(attr.Value, gain, attr.Realm, attr.RealmExp, attr.IsBottleneck, attr.AccumulationPool)
			attr.Value = result.NewValue
			attr.AccumulationPool = result.NewAccPool
			attr.RealmExp = result.NewRealmExp
			attr.IsBottleneck = result.NewIsBottleneck
			if err := l.svcCtx.CharacterModel.UpdateAttribute(attr); err != nil {
				return nil, err
			}
		}
		message = fmt.Sprintf("敏捷提升了 %d 点", item.EffectValue*req.Quantity)

	case "spirit_stone_gain":
		stats.SpiritStones += item.EffectValue * req.Quantity
		message = fmt.Sprintf("获得了 %d 灵石", item.EffectValue*req.Quantity)

	case "", "none":
		message = fmt.Sprintf("已使用「%s」", item.Name)

	default:
		return nil, fmt.Errorf("未知的物品效果")
	}

	// Update character stats
	if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
		return nil, err
	}

	// Remove from inventory (only for consumables)
	if item.ItemType == "consumable" {
		if err := l.svcCtx.ShopModel.RemoveFromInventory(userID, req.ItemID, req.Quantity); err != nil {
			return nil, err
		}
	}

	// Reload attributes for response
	attrs, err = l.svcCtx.CharacterModel.FindAttributesByUserID(userID)
	if err != nil {
		return nil, err
	}

	charLogic := NewCharacterLogic(l.svcCtx)
	charResp := charLogic.statsToResp(stats, attrs)

	return &types.UseItemResult{
		Success:   true,
		Message:   message,
		Character: *charResp,
	}, nil
}

func (l *ShopLogic) SellItem(ctx context.Context, userID int64, req *types.SellItemReq) (*types.SellItemResult, error) {
	invItem, err := l.svcCtx.ShopModel.GetInventoryItemByItemID(userID, req.ItemID)
	if err != nil {
		return nil, err
	}
	if invItem == nil || invItem.Quantity < req.Quantity {
		return nil, fmt.Errorf("物品不足")
	}

	item, err := l.svcCtx.ShopModel.GetItemByID(req.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("物品不存在")
	}

	if item.SellPrice <= 0 {
		return nil, fmt.Errorf("该物品不可出售")
	}

	totalGain := item.SellPrice * req.Quantity

	stats, err := l.svcCtx.CharacterModel.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	stats.SpiritStones += totalGain
	if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
		return nil, err
	}

	if err := l.svcCtx.ShopModel.RemoveFromInventory(userID, req.ItemID, req.Quantity); err != nil {
		return nil, err
	}

	return &types.SellItemResult{
		Success:              true,
		Message:              fmt.Sprintf("成功出售 %d 个「%s」，获得 %d 灵石", req.Quantity, item.Name, totalGain),
		RemainingSpiritStones: stats.SpiritStones,
	}, nil
}

func (l *ShopLogic) GetPurchaseHistory(ctx context.Context, userID int64) (*types.PurchaseHistoryResp, error) {
	history, err := l.svcCtx.ShopModel.GetPurchaseHistory(userID, 50)
	if err != nil {
		return nil, err
	}

	resp := &types.PurchaseHistoryResp{
		History: make([]types.PurchaseRecordResp, 0),
	}

	for _, record := range history {
		resp.History = append(resp.History, types.PurchaseRecordResp{
			ID:         record.ID,
			ItemName:   record.ItemName,
			Quantity:   record.Quantity,
			TotalPrice: record.TotalPrice,
			CreatedAt:  record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}
