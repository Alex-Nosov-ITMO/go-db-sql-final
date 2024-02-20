package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1001,
		Status:    ParcelStatusRegistered,
		Address:   "test_new",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавление новой записи в бд
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// get
	// проверка на корректность добавления
	parcel.Number = id
	p, err := store.Get(id)
	assert.Equal(t, parcel, p)
	require.NoError(t, err)

	// delete
	// удаление записи из бд
	err = store.Delete(id)
	require.NoError(t, err)

	// chek
	// проверка на корректность удаления
	p, err = store.Get(id)
	assert.Error(t, err)
	assert.Empty(t, p)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавление новой записи в бд
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// set address
	// обновление адреса
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	// проверка корректности обновления
	p, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, p.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавление новой записи в бд
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// set status
	// обновление статуса
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	// check
	// проверка корректности обновления
	p, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, p.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	// добавление новых записей в бд
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		assert.NotEmpty(t, id)
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get
	// получаем все строки с данным клиентом
	storedParcels, err := store.GetByClient(client)
	// проверка на корректность выполнения функции
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		// проверка, что все полученные строки есть в нашей мапе
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
