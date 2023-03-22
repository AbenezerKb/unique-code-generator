package main

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 26257
	user     = "root"
	password = ""
	dbname   = "defaultdb"
)

func main() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	//  postgres://root@localhost:26257/defaultdb?sslmode=disable
	// open database

	// 	testConnURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	// 	viper.GetString("test_database.user"),
	// 	viper.GetString("test_database.password"),
	// 	viper.GetString("test_database.host"),
	// 	viper.GetString("test_database.port"),
	// 	dbName)

	// testConn := initiator.InitDB(testConnURL, logg)
	// logg.Info(context.Background(), "test database initialized")
	// PersistDB := persistencedb.New(testConn, logg)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")
	list := make(map[string]int)
	status := false
	count := 0
	var wg sync.WaitGroup
	var someMapMutex sync.RWMutex
	for i := 0; i<5; i++{		
		wg.Add(1)
	go	generate(count, status,  list, db,  &wg, &someMapMutex, i )
	}
	
	for k, v := range list {
		if v > 0 {
			fmt.Printf("%s exists", k)
		}
	}
	fmt.Println("done")
	CheckError(err)
	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()
	fmt.Println("Done!")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
func generate(count int, status bool,  list map[string]int, db *sql.DB,  wg *sync.WaitGroup, someMapMutex *sync.RWMutex, routine int) {
	defer wg.Done()
	someMapMutex.Lock()
	for i := 0; i < 100000; i++ {
		insertDynStmt := ` SELECT 
		concat(
		  substr(
			'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', 
			floor(random() * 62)::int + 1, 
			1
		  ), 
		  substr(
			'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', 
			floor(random() * 62)::int + 1, 
			1
		  ), 
		  substr(
			'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', 
			floor(random() * 62)::int + 1, 
			1
		  ), 
		  substr(
			'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', 
			floor(random() * 62)::int + 1, 
			1
		  ), 
		  substr(
			'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', 
			floor(random() * 62)::int + 1, 
			1
		  ), 
		  substr(
			'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', 
			floor(random() * 62)::int + 1, 
			1
		  ), 
		  substr(
			'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', 
			floor(random() * 62)::int + 1, 
			1
		  ), 
		  substr(
			'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', 
			floor(random() * 62)::int + 1, 
			1
		  )
		)`
		// result, err := db.Exec(insertDynStmt, "Jane", 2)
		rows, err := db.Query(insertDynStmt)
		defer rows.Close()
		for rows.Next() {
			var name string

			err = rows.Scan(&name)
			CheckError(err)
			fmt.Println("the key: ", status, " count: ", count, " index: ", i, "routine: ", routine)
			_, exists := list[name]
			if exists {
				status = true
				count++
				list[name]++
				// time.Sleep(time.Second*10)
			}
			list[name] = 0
		}		
	}
	someMapMutex.Unlock()
}
