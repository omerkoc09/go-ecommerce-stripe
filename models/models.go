package models

import (
	"context"
	"database/sql"
	"time"
)

// DBModel is the type for database connection
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for the all models
type Models struct {
	DB DBModel
}

// newmodel returns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Mac struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Price          int       `json:"price"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Image          string    `json:"image"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Order struct {
	ID            int       `json:"id"`
	MacID         int       `json:"mac_id"`
	TransactionID int       `json:"transaction_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *DBModel) GetMac(id int) (Mac, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var mac Mac

	row := m.DB.QueryRowContext(ctx,
		`SELECT 
			id, name, price, description, inventory_level, coalesce(image, ''), created_at, updated_at
	 FROM 
	        macs 
	 WHERE id = ?`, id)

	err := row.Scan(
		&mac.ID,
		&mac.Name,
		&mac.Price,
		&mac.Description,
		&mac.InventoryLevel,
		&mac.Image,
		&mac.CreatedAt,
		&mac.UpdatedAt,
	)
	if err != nil {
		return mac, err
	}
	return mac, nil
}

// InsertTransaction inserts a new txn and returns the id of the txn
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO transactions 
	(amount, currency, last_four, bank_return_code, transaction_status_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.TransactionStatusID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertOrder inserts a new order and returns the id of the order
func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO orders 
	(mac_id, transaction_id, status_id, quantity, amount, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, stmt,
		order.MacID,
		order.TransactionID,
		order.StatusID,
		order.Quantity,
		order.Amount,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
