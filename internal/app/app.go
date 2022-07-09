package app

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/stan.go"
	"github.com/valen0k/wb_l0/internal/model"
	"log"
	"reflect"
	"sync"
)

type App struct {
	DB  *pgx.Conn
	Buf struct {
		sync.RWMutex
		memory map[string]model.Model
	}
}

func (a *App) MemoryRecovery() error {
	query := `SELECT id, order_info FROM test`
	rows, err := a.DB.Query(context.Background(), query)
	if err != nil {
		return err
	}

	var order model.Model
	var id string
	a.Buf.memory = make(map[string]model.Model)

	for rows.Next() {
		err = rows.Scan(&id, &order)
		if err != nil {
			return err
		}
		a.Set(id, order)
	}

	log.Println("record recovery completed")
	return nil
}

func (a *App) Set(key string, val model.Model) {
	a.Buf.Lock()
	defer a.Buf.Unlock()

	log.Println("recorded in memory")
	a.Buf.memory[key] = val
}

func (a *App) Get(key string) (model.Model, bool) {
	a.Buf.RLock()
	defer a.Buf.RUnlock()

	value, ok := a.Buf.memory[key]
	if ok {
		return value, ok
	}
	return model.Model{}, ok
}

func (a *App) Rec(msg *stan.Msg) {
	var order model.Model

	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Println(err)
		return
	}

	var actual model.Model
	if reflect.DeepEqual(order, actual) {
		log.Println("type model not equal")
		return
	}

	var id string
	query := `INSERT INTO test (order_info) VALUES ($1) RETURNING id`

	err = a.DB.QueryRow(context.Background(), query, order).Scan(&id)
	if err != nil {
		log.Println(err)
		return
	}

	a.Set(id, order)

	log.Println("data recorded")
}
