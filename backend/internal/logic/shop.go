package logic

import (
	"context"
	"fmt"

	"life-system-backend/internal/model"
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

	item := &model.ShopItem{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
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
	// Get item
	item, err := l.svcCtx.ShopModel.GetItemByID(req.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("物品不存在")
	}

	// Check stock
	if item.Stock != -1 && item.Stock < req.Quantity {
		return nil, fmt.Errorf("库存不足")
	}

	// Calculate total price
	totalPrice := item.Price * req.Quantity

	// Get character stats
	stats, err := l.svcCtx.CharacterModel.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	// Check if user has enough gold
	if stats.Gold < totalPrice {
		return nil, fmt.Errorf("金币不足！需要 %d 金币，当前只有 %d 金币", totalPrice, stats.Gold)
	}

	// Deduct gold
	stats.Gold -= totalPrice

	// Update character
	if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
		return nil, err
	}

	// Update item stock if not unlimited
	if item.Stock != -1 {
		if err := l.svcCtx.ShopModel.UpdateItemStock(item.ID, req.Quantity); err != nil {
			return nil, err
		}
	}

	// Add to inventory
	if err := l.svcCtx.ShopModel.AddToInventory(userID, item.ID, req.Quantity); err != nil {
		return nil, err
	}

	// Record purchase
	if err := l.svcCtx.ShopModel.RecordPurchase(userID, item.ID, item.Name, req.Quantity, totalPrice); err != nil {
		return nil, err
	}

	return &types.PurchaseResult{
		Success:       true,
		Message:       fmt.Sprintf("成功购买 %d 个「%s」", req.Quantity, item.Name),
		RemainingGold: stats.Gold,
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
		// Get item details
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
			Icon:        item.Icon,
			Image:       item.Image,
			Quantity:    invItem.Quantity,
		})
	}

	return resp, nil
}

func (l *ShopLogic) UseItem(ctx context.Context, userID int64, req *types.UseItemReq) (*types.UseItemResult, error) {
	// Get inventory item
	invItem, err := l.svcCtx.ShopModel.GetInventoryItemByItemID(userID, req.ItemID)
	if err != nil {
		return nil, err
	}
	if invItem == nil || invItem.Quantity < req.Quantity {
		return nil, fmt.Errorf("物品不足")
	}

	// Get item details
	item, err := l.svcCtx.ShopModel.GetItemByID(req.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("物品不存在")
	}

	// Get character stats
	stats, err := l.svcCtx.CharacterModel.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	// Apply item effect
	message := ""
	switch item.Effect {
	case "hp_restore":
		stats.HP += item.EffectValue * req.Quantity
		if stats.HP > stats.MaxHP {
			stats.HP = stats.MaxHP
		}
		message = fmt.Sprintf("恢复了 %d 点生命值", item.EffectValue*req.Quantity)

	case "energy_restore":
		stats.Energy += item.EffectValue * req.Quantity
		if stats.Energy > stats.MaxEnergy {
			stats.Energy = stats.MaxEnergy
		}
		message = fmt.Sprintf("恢复了 %d 点能量", item.EffectValue*req.Quantity)

	case "strength_boost":
		stats.Strength += float64(item.EffectValue * req.Quantity)
		message = fmt.Sprintf("力量永久提升了 %d 点", item.EffectValue*req.Quantity)

	case "intelligence_boost":
		stats.Intelligence += float64(item.EffectValue * req.Quantity)
		message = fmt.Sprintf("智力永久提升了 %d 点", item.EffectValue*req.Quantity)

	case "vitality_boost":
		stats.Vitality += float64(item.EffectValue * req.Quantity)
		message = fmt.Sprintf("体力永久提升了 %d 点", item.EffectValue*req.Quantity)

	case "spirit_boost":
		stats.Spirit += float64(item.EffectValue * req.Quantity)
		message = fmt.Sprintf("精神永久提升了 %d 点", item.EffectValue*req.Quantity)

	case "full_restore":
		stats.HP = stats.MaxHP
		stats.Energy = stats.MaxEnergy
		message = "完全恢复了生命值和能量"

	case "exp_gain":
		stats.Exp += item.EffectValue * req.Quantity
		CheckAndApplyLevelUp(stats)
		message = fmt.Sprintf("获得了 %d 点经验值", item.EffectValue*req.Quantity)

	default:
		return nil, fmt.Errorf("未知的物品效果")
	}

	// Recalculate MaxHP if attributes changed
	stats.MaxHP = 100 + int(stats.Strength*2) + int(stats.Vitality*3)
	if stats.HP > stats.MaxHP {
		stats.HP = stats.MaxHP
	}

	// Update character
	if err := l.svcCtx.CharacterModel.Update(stats); err != nil {
		return nil, err
	}

	// Remove from inventory (only for consumables)
	if item.ItemType == "consumable" {
		if err := l.svcCtx.ShopModel.RemoveFromInventory(userID, req.ItemID, req.Quantity); err != nil {
			return nil, err
		}
	}

	// Get character logic to return updated character
	charLogic := NewCharacterLogic(l.svcCtx)
	charResp := charLogic.statsToResp(stats)

	return &types.UseItemResult{
		Success:   true,
		Message:   message,
		Character: *charResp,
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
