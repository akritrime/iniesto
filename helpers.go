package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

type guildDB struct {
	id     string
	name   string
	prefix string
}

func getDB() (db *sql.DB, err error) {
	dbname := os.Getenv("DBNAME")
	if dbname == "" {
		err = errors.New("var DBNAME not set")
		return
	}

	dbuser := os.Getenv("DBUSER")
	if dbuser == "" {
		err = errors.New("var DBUSER not set")
		return
	}

	dbpwd := os.Getenv("DBPWD")
	if dbpwd == "" {
		err = errors.New("var DBPWD not set")
		return
	}

	conf := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=disable", dbuser, dbname, dbpwd)

	db, err = sql.Open("postgres", conf)
	return
}

func isNewGuild(gID string) (flag bool, err error) {
	var name string
	err = db.QueryRow("SELECT name FROM guilds WHERE id = $1", gID).Scan(&name)

	if err == sql.ErrNoRows {
		return true, nil
	}

	return false, err
}

func getGuildsCount() (count int, err error) {
	err = db.QueryRow("SELECT COUNT(*) FROM guilds").Scan(&count)
	return
}

func addNewGuild(gID, gName string) (err error) {
	_, err = db.Exec("INSERT INTO guilds(id, name, prefix) VALUES ($1, $2, $3)", gID, gName, ">")
	if err != nil {
		fmt.Println(err)
		// } else {
		// 	fmt.Println("res: ", res)
	}
	return
}

func getAllGuilds() (l []guildDB, err error) {
	rows, err := db.Query("SELECT name, id FROM guilds")
	if err != nil {
		fmt.Println("err in Query select quild ", err)
		return
	}
	for rows.Next() {
		var g guildDB
		err = rows.Scan(&g.name, &g.id)
		if err != nil {
			fmt.Println("err in iterating ", err)
		}
		l = append(l, g)
	}

	return

}

func removeGuild(gID string) {
	_, err := db.Exec("DELETE FROM guilds WHERE id = $1", gID)
	if err != nil {
		fmt.Println("Err in guild Delete ", err)
	}

}

func setPrefix(gID, prefix string) (err error) {
	_, err = db.Exec("UPDATE guilds SET prefix = $1 WHERE id = $2", prefix, gID)
	if err != nil {
		fmt.Println("Err in guild prefix update ", err)
	}
	return
}

func getPrefix(gID string) (prefix string) {
	err := db.QueryRow("SELECT prefix FROM guilds WHERE id = $1", gID).Scan(&prefix)
	if err != nil {
		fmt.Println("Err in getting prefix", err)
		prefix = ">"
	}
	return
}
