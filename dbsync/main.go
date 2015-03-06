package dbsync

import (
	"gosync/utils"
)

var (
	dbstore Datastore
)

type Datastore interface {
	Insert(table string, item utils.FsItem) bool
    Remove(table string, relPath string) bool
	CheckEmpty(table string) bool
	FetchAll(table string) []utils.DataTable
	CheckIn(listener string) ([]utils.DataTable, error)
    GetOne(listener, path string) (utils.DataTable, error)
	CreateDB()
	Close() error // call this method when you want to close the connection
	initDB()
}

func setdbstoreEngine() {
	cfg := utils.GetConfig()
	var engine = cfg.Database.Type
	switch engine {
	case "mysql":
		dbstore = &MySQLDB{config: cfg}
		dbstore.initDB()
	}
}

func Insert(table string, item utils.FsItem) bool {
	setdbstoreEngine()
	return dbstore.Insert(table, item)
}

func CheckEmpty(table string) bool {
	setdbstoreEngine()
	empty := dbstore.CheckEmpty(table)
	if empty {
        utils.WriteLn("Database is EMPTY, starting creation")
	} else {
        utils.WriteLn("Using existing table: " + table)
	}
	return empty
}

func FetchAll(table string) []utils.DataTable {
	setdbstoreEngine()
	return dbstore.FetchAll(table)
}

func CheckIn(listener string) ([]utils.DataTable, error) {
    utils.WriteLn("Starting db checking background script for: " + listener)
	data,err := dbstore.CheckIn(listener)
    return data, err

}

func GetOne(basepath, path string) (utils.DataTable, error){
    setdbstoreEngine()
    listener := utils.GetListenerFromDir(basepath)
    dbitem, err := dbstore.GetOne(listener, path)
    return dbitem, err
}

func Remove(basepath, relPath string) bool {
    setdbstoreEngine()
    listener := utils.GetListenerFromDir(basepath)
    return dbstore.Remove(listener, relPath)
}

func CreateDB() {
	setdbstoreEngine()
	dbstore.CreateDB()
}
