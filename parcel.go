package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// Выполняем запрос INSERT
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return -1, err
	}
	// Возвращаем индификатор последней добавленной записи при помощи функции LastInsertId()
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// Выполняем SELECT запрос по номеру посылки (только одна запись)
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	// Заполняем структуру полученными занными
	p := Parcel{}

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	// Возвращаем структуру и ошибку
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// Выполняем SELECT запрос по номеру клиента (выводятся все записи с данным клиентом)
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close() // закрываем курсор по завершению
	// заполняем массив структурами с данными из полученных записей
	var res []Parcel

	for rows.Next() {
		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}
	// возвращаем массив со структурами и ошибку
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// Выолняем UPDATE запрос с обновлением статуса
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// Выполняем UPDATE запрос с обновлением адресса (дополнительное условие на статус)

	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// Выполняем DELETE запрос с удалением записи (дополнительное условие на статус)

	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		return err
	}
	return nil
}
