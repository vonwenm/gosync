package datastore

import (
	"gosync/utils"
)

var (
	dbstore Datastore
)

type Datastore interface {
	Insert(table string, item utils.FsItem) bool
    Remove(table string, item utils.FsItem) bool
	CheckEmpty(table string) bool
	FetchAll(table string) []utils.DataTable
	CheckIn(listener string) ([]utils.DataTable, error)
    GetOne(listener, path string) (utils.DataTable, error)
    UpdateHost(table, path string)
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
	}
}

func Insert(table string, item utils.FsItem) bool {
	setdbstoreEngine()
    dbstore.initDB()
    defer dbstore.Close()
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
    dbstore.initDB()
    defer dbstore.Close()
	return dbstore.FetchAll(table)
}

func UpdateHost(table, path string){
    setdbstoreEngine()
    dbstore.initDB()
    defer dbstore.Close()
    dbstore.UpdateHost(table, path)
}

func CheckIn(listener string) ([]utils.DataTable, error) {
    setdbstoreEngine()
    dbstore.initDB()
    defer dbstore.Close()
    utils.WriteLn("Starting db checking background script for: " + listener)
	data,err := dbstore.CheckIn(listener)
    return data, err

}

func GetOne(basepath, path string) (utils.DataTable, error){
    setdbstoreEngine()
    dbstore.initDB()
    defer dbstore.Close()
    listener := utils.GetListenerFromDir(basepath)
    dbitem, err := dbstore.GetOne(listener, path)
    return dbitem, err
}

func Remove(table string, item utils.FsItem) bool {
    setdbstoreEngine()
    dbstore.initDB()
    defer dbstore.Close()
    return dbstore.Remove(table, item)
}

func CreateDB() {
    setdbstoreEngine()
    dbstore.initDB()
    defer dbstore.Close()
	dbstore.CreateDB()
}
