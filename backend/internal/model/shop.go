package model

import (
	"database/sql"
	"time"
)

type ShopItem struct {
	ID          int64
	UserID      int64
	Name        string
	Description string
	Price       int    // Gold cost
	ItemType    string // consumable, permanent
	Effect      string // hp_restore, energy_restore, attribute_boost, etc.
	EffectValue int    // Amount of effect
	Icon        string // Emoji or icon identifier
	Image       string // Image file path
	Stock       int    // -1 for unlimited
	CreatedAt   time.Time
}

type InventoryItem struct {
	ID        int64
	UserID    int64
	ItemID    int64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PurchaseHistory struct {
	ID         int64
	UserID     int64
	ItemID     int64
	ItemName   string
	Quantity   int
	TotalPrice int
	CreatedAt  time.Time
}

type ShopModel struct {
	db *sql.DB
}

func NewShopModel(db *sql.DB) *ShopModel {
	return &ShopModel{db: db}
}

// GetItemsByUserID returns all shop items created by a specific user
func (m *ShopModel) GetItemsByUserID(userID int64) ([]*ShopItem, error) {
	rows, err := m.db.Query(`
		SELECT id, user_id, name, description, price, item_type, effect, effect_value, icon, image, stock, created_at
		FROM shop_items
		WHERE user_id = ? AND stock != 0
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*ShopItem
	for rows.Next() {
		var item ShopItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.Name, &item.Description, &item.Price, &item.ItemType,
			&item.Effect, &item.EffectValue, &item.Icon, &item.Image, &item.Stock, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, rows.Err()
}

// GetItemByID returns a specific shop item
func (m *ShopModel) GetItemByID(id int64) (*ShopItem, error) {
	var item ShopItem
	err := m.db.QueryRow(`
		SELECT id, user_id, name, description, price, item_type, effect, effect_value, icon, image, stock, created_at
		FROM shop_items
		WHERE id = ?
	`, id).Scan(
		&item.ID, &item.UserID, &item.Name, &item.Description, &item.Price, &item.ItemType,
		&item.Effect, &item.EffectValue, &item.Icon, &item.Image, &item.Stock, &item.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

// CreateItem creates a new shop item
func (m *ShopModel) CreateItem(item *ShopItem) (int64, error) {
	result, err := m.db.Exec(`
		INSERT INTO shop_items (user_id, name, description, price, item_type, effect, effect_value, icon, image, stock, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
	`, item.UserID, item.Name, item.Description, item.Price, item.ItemType,
		item.Effect, item.EffectValue, item.Icon, item.Image, item.Stock)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UpdateItem updates an existing shop item
func (m *ShopModel) UpdateItem(item *ShopItem) error {
	_, err := m.db.Exec(`
		UPDATE shop_items
		SET name = ?, description = ?, price = ?, item_type = ?, effect = ?, effect_value = ?, icon = ?, image = ?, stock = ?
		WHERE id = ? AND user_id = ?
	`, item.Name, item.Description, item.Price, item.ItemType,
		item.Effect, item.EffectValue, item.Icon, item.Image, item.Stock,
		item.ID, item.UserID)

	return err
}

// DeleteItem deletes a shop item by id for a specific user
func (m *ShopModel) DeleteItem(id, userID int64) error {
	_, err := m.db.Exec(`
		DELETE FROM shop_items WHERE id = ? AND user_id = ?
	`, id, userID)

	return err
}

// UpdateItemStock updates the stock of an item
func (m *ShopModel) UpdateItemStock(id int64, quantity int) error {
	_, err := m.db.Exec(`
		UPDATE shop_items
		SET stock = stock - ?
		WHERE id = ? AND stock >= ?
	`, quantity, id, quantity)

	return err
}

// GetUserInventory returns all items in user's inventory
func (m *ShopModel) GetUserInventory(userID int64) ([]*InventoryItem, error) {
	rows, err := m.db.Query(`
		SELECT id, user_id, item_id, quantity, created_at, updated_at
		FROM inventory
		WHERE user_id = ? AND quantity > 0
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*InventoryItem
	for rows.Next() {
		var item InventoryItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.ItemID, &item.Quantity,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, rows.Err()
}

// GetInventoryItemByItemID gets a specific inventory item
func (m *ShopModel) GetInventoryItemByItemID(userID, itemID int64) (*InventoryItem, error) {
	var item InventoryItem
	err := m.db.QueryRow(`
		SELECT id, user_id, item_id, quantity, created_at, updated_at
		FROM inventory
		WHERE user_id = ? AND item_id = ?
	`, userID, itemID).Scan(
		&item.ID, &item.UserID, &item.ItemID, &item.Quantity,
		&item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

// AddToInventory adds an item to user's inventory
func (m *ShopModel) AddToInventory(userID, itemID int64, quantity int) error {
	// Check if item already exists in inventory
	existing, err := m.GetInventoryItemByItemID(userID, itemID)
	if err != nil {
		return err
	}

	if existing != nil {
		// Update existing
		_, err = m.db.Exec(`
			UPDATE inventory
			SET quantity = quantity + ?, updated_at = datetime('now')
			WHERE user_id = ? AND item_id = ?
		`, quantity, userID, itemID)
	} else {
		// Insert new
		_, err = m.db.Exec(`
			INSERT INTO inventory (user_id, item_id, quantity, created_at, updated_at)
			VALUES (?, ?, ?, datetime('now'), datetime('now'))
		`, userID, itemID, quantity)
	}

	return err
}

// RemoveFromInventory removes quantity from user's inventory
func (m *ShopModel) RemoveFromInventory(userID, itemID int64, quantity int) error {
	_, err := m.db.Exec(`
		UPDATE inventory
		SET quantity = quantity - ?, updated_at = datetime('now')
		WHERE user_id = ? AND item_id = ? AND quantity >= ?
	`, quantity, userID, itemID, quantity)

	return err
}

// RecordPurchase records a purchase in history
func (m *ShopModel) RecordPurchase(userID, itemID int64, itemName string, quantity, totalPrice int) error {
	_, err := m.db.Exec(`
		INSERT INTO purchase_history (user_id, item_id, item_name, quantity, total_price, created_at)
		VALUES (?, ?, ?, ?, ?, datetime('now'))
	`, userID, itemID, itemName, quantity, totalPrice)

	return err
}

// GetPurchaseHistory returns user's purchase history
func (m *ShopModel) GetPurchaseHistory(userID int64, limit int) ([]*PurchaseHistory, error) {
	query := `
		SELECT id, user_id, item_id, item_name, quantity, total_price, created_at
		FROM purchase_history
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	if limit > 0 {
		query += ` LIMIT ?`
	}

	var rows *sql.Rows
	var err error

	if limit > 0 {
		rows, err = m.db.Query(query, userID, limit)
	} else {
		rows, err = m.db.Query(query, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*PurchaseHistory
	for rows.Next() {
		var record PurchaseHistory
		err := rows.Scan(
			&record.ID, &record.UserID, &record.ItemID, &record.ItemName,
			&record.Quantity, &record.TotalPrice, &record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, &record)
	}

	return history, rows.Err()
}
