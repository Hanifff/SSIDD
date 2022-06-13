package ssidd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

// OffchainDB is an instance of database
type OffchainDB struct {
	Config mysql.Config
	Conn   *sql.DB
}

// Linkdata is an instance of database query response
type Linkdata struct {
	ID   string
	data string
}

// NewOffchainDB configures and start a database connection.
// This function returnes a pointer to OffchainDB
// user: Username of mysql database adminstrator
// passw: Passowod of mysql database adminstrator
// net: TCP address of mysql server
// dbname: Offchain database
func NewOffchainDB(user, passw, net, addr, dbname string) (*OffchainDB, error) {
	cfg := mysql.Config{
		User:                 user,
		Passwd:               passw,
		Net:                  net,
		Addr:                 addr,
		DBName:               dbname,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return &OffchainDB{}, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &OffchainDB{
		Config: cfg,
		Conn:   db,
	}, nil
}

// InsertData inserts data to the Offchain database.
// ID: Transaction ID which links data to blockchain
// Data: The data that correspond to the Transaction ID
func (oc *OffchainDB) InsertData(ID string, offdata string) error {
	values := make([]interface{}, 2)
	values[0] = ID
	values[1] = offdata
	insert, err := oc.Conn.Query("INSERT INTO linkdata(ID, Offdata) VALUES(?,?)", values...)
	if err != nil {
		fmt.Printf("Something went wrong while trying to insert data to the database.\n%v\n", err)
		return err
	}
	defer insert.Close()
	return nil
}

// ReadData reads data from the Offchain database.
// ID: Transaction ID which links data to blockchain
func (oc *OffchainDB) ReadData(ID string) error {
	results, err := oc.Conn.Query("SELECT Offdata FROM linkdata where ID=?", ID)
	if err != nil {
		fmt.Printf("Something went wrong while trying to select from database.\n%v", err)
		return err
	}
	defer results.Close()
	for results.Next() {
		var links Linkdata
		// for each row, scan the result into our tag composite object
		err = results.Scan(&links.data)
		if err != nil {
			fmt.Printf("Something went wrong while casting data to the composite object.\n%v", err)
			return err
		}
		log.Printf("DATA: %s\n", links.data)
	}
	return nil
}

// ReadAllData reads all data stored in the Offchain database.
func (oc *OffchainDB) ReadAllData() error {
	fmt.Printf("Reading all data from the linkdata table in OffchainDB database..\n")
	results, err := oc.Conn.Query("SELECT ID, Offdata FROM linkdata")
	if err != nil {
		fmt.Printf("Something went wrong while trying to select from database.\n%v", err)
		return err
	}
	defer results.Close()
	for results.Next() {
		var links Linkdata
		// for each row, scan the result into our tag composite object
		err = results.Scan(&links.ID, &links.data)
		if err != nil {
			fmt.Printf("Something went wrong while casting data to the composite object.\n%v", err)
			return err
		}
		log.Printf("ID: %s\nDATA: %s\n", links.ID, links.data)
	}
	fmt.Println("Done!")
	return nil
}

// DeleteData deletes data from the database given its hash.
func (oc *OffchainDB) DeleteData(ID string) error {
	results, err := oc.Conn.Query("DELETE Offdata FROM linkdata where ID=?", ID)
	if err != nil {
		fmt.Printf("Something went wrong while trying to select from database.\n%v", err)
		return err
	}
	defer results.Close()
	return nil
}
