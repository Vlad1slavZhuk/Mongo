package postgres

import (
	"Mongo/internal/pkg/constErr"
	"Mongo/internal/pkg/data"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type PostgresStorage struct {
	DbName string
	DB     *sql.DB
}

func NewPostgresStorage() *PostgresStorage {
	connSTR := fmt.Sprintf("user=%v password=%v host=%v port=%v sslmode=disable", "postgres", "postgres_password", "postgres", "8003")
	time.Sleep(15 * time.Second)
	db, err := sql.Open("postgres", connSTR)
	if err != nil {
		log.Fatal(err)
	}
	storage := &PostgresStorage{
		DbName: "postgres",
		DB:     db,
	}
	driver, err := postgres.WithInstance(storage.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migrate, err := migrate.NewWithDatabaseInstance("file://./migrations", storage.DbName, driver)
	if err != nil {
		log.Fatal(err)
	}
	if err = migrate.Up(); err != nil {
		log.Fatal(err)
	}
	return storage
}

func (s *PostgresStorage) Add(ad *data.Ad) error {
	row := s.DB.QueryRow(fmt.Sprintf("INSERT INTO %v (brand, model, color, price) VALUES ($1, $2, $3, $4);", "ads"),
		ad.GetBrand(), ad.GetModel(), ad.GetColor(), ad.GetPrice())
	if row != nil {
		return row.Err()
	}
	log.Println(row)

	return nil
}

func (s *PostgresStorage) Get(id uint) (*data.Ad, error) {
	ad := new(data.Ad)
	row := s.DB.QueryRow(fmt.Sprintf("SELECT * FROM %v WHERE id = $1;", "ads"), id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	scan := row.Scan(&ad.ID, &ad.Brand, &ad.Model, &ad.Color, &ad.Price)
	if scan != nil {
		return nil, scan
	}

	return ad, nil
}

func (s *PostgresStorage) GetAll() ([]*data.Ad, error) {
	length, err := s.Size()
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(fmt.Sprintf("SELECT * FROM %v ORDER BY id ASC;", "ads"))
	if err != nil {
		return nil, err
	}

	ads := make([]*data.Ad, 0, length)

	for rows.Next() {
		ad := new(data.Ad)
		err = rows.Scan(&ad.ID, &ad.Brand, &ad.Model, &ad.Color, &ad.Price)
		if err != nil {
			return nil, err
		}

		ads = append(ads, ad)
	}

	return ads, nil
}

func (s *PostgresStorage) Update(temp *data.Ad, id uint) error {
	if _, err := s.Get(id); err != nil {
		return err
	}

	_, err := s.DB.Exec(fmt.Sprintf("UPDATE %v SET brand=$1, model=$2, color=$3, price=$4 WHERE id=$5", "ads"),
		temp.GetBrand(), temp.GetModel(), temp.GetColor(), temp.GetPrice(), id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) Delete(id uint) error {
	if _, err := s.Get(id); err != nil {
		return err
	}

	_, err := s.DB.Exec(fmt.Sprintf("DELETE FROM %v WHERE id=$1;", "ads"), id)
	if err != nil {
		return err
	}

	size, err := s.Size()
	if err != nil {
		return err
	}

	if size == 0 {
		return nil
	}

	ads, _ := s.GetAll()
	log.Println(ads)
	newID := 1
	for _, ad := range ads {
		_, err := s.DB.Exec(fmt.Sprintf("UPDATE %v SET id=$1 WHERE id=$2;", "ads"), newID, ad.ID)
		if err != nil {
			return err
		}
		newID++
	}

	// newIndex := 1
	// for i := 0; true; i++ {
	// 	_, err := s.DB.Exec(fmt.Sprintf("UPDATE %v SET id=$1 WHERE id=$2;", "ads"), newIndex, i)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	if size == newIndex {
	// 		break
	// 	}
	// 	newIndex++
	// }

	_, err = s.DB.Exec(fmt.Sprintf("ALTER SEQUENCE %v RESTART WITH %v;", "ads_id_seq", size+1))
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) AddAccount(acc *data.Account) error {
	row := s.DB.QueryRow(fmt.Sprintf("INSERT INTO %v (username, password, token) VALUES ($1, $2, $3);", "accounts"),
		acc.GetUserName(), acc.GetPassword(), acc.GetToken())
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (s *PostgresStorage) GetAccounts() ([]*data.Account, error) {
	var size int

	err := s.DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %v;", "accounts")).Scan(&size)
	if err != nil || size == 0 {
		return nil, constErr.AccountBaseIsEmpty
	}

	rows, err := s.DB.Query(fmt.Sprintf("SELECT * FROM %v;", "accounts"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accs := make([]*data.Account, 0, size)

	for rows.Next() {
		acc := new(data.Account)

		err = rows.Scan(&acc.ID, &acc.Username, &acc.Password, &acc.Token)
		if err != nil {
			return nil, err
		}

		accs = append(accs, acc)
	}

	return accs, nil
}

func (s *PostgresStorage) UpdateTokenCurrentAcc(acc *data.Account, token string) error {
	_, err := s.DB.Exec(fmt.Sprintf("UPDATE %v SET token=$1 WHERE username=$2 AND password=$3;", "accounts"),
		token, acc.GetUserName(), acc.GetPassword())
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) Size() (int, error) {
	var size int

	row := s.DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %v;", "ads"))
	err := row.Scan(&size)
	if err != nil {
		return 0, err
	}
	return size, nil
}
