package fswatcher

import (
	"gopkg.in/fsnotify.v1"
	"gosync/datastore"
	"gosync/storage"
	"gosync/utils"
)

func SysPathWatcher(path string) {
	utils.LogWriteF("Starting new watcher for %s:", path)
	watcher, err := fsnotify.NewWatcher()
	if !utils.ErrorCheckF(err, 400, "Cannot create new watcher for %s ", path) {
		defer watcher.Close()
		done := make(chan bool)
		go func() {
			listener := utils.GetListenerFromDir(path)
			rel_path := utils.GetRelativePath(listener, path)
			for {
				select {
				case event := <-watcher.Events:
					//logs.WriteLn("event:", event)
					if event.Op&fsnotify.Chmod == fsnotify.Chmod {
						utils.LogWriteF("Chmod occurred on:", event.Name)
						if checkItemInDB(path, rel_path, event.Name) {
							runFileUpdate(path, event.Name, "chmod")
						}
					}
					if event.Op&fsnotify.Rename == fsnotify.Rename {
						utils.LogWriteF("Rename occurred on:", event.Name)
						if checkItemInDB(path, rel_path, event.Name) {
							runFileUpdate(path, event.Name, "rename")
						}
					}
					if event.Op&fsnotify.Create == fsnotify.Create {
						utils.LogWriteF("New File:", event.Name)
						if checkItemInDB(path, rel_path, event.Name) {
							runFileUpdate(path, event.Name, "create")
						}
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						utils.LogWriteF("modified file:", event.Name)
						if checkItemInDB(path, rel_path, event.Name) {
							runFileUpdate(path, event.Name, "write")
						}
					}
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						utils.LogWriteF("Removed File: ", event.Name)
						if checkItemInDB(path, rel_path, event.Name) {
							runFileUpdate(path, event.Name, "remove")
						}
					}
				case err := <-watcher.Errors:
					utils.LogWriteF("error:", err)
				}
			}

		}()
		err = watcher.Add(path)
		utils.ErrorCheckF(err, 500, "Cannot add watcher to %s ", path)
		<-done
	}

}

func runFileUpdate(base_path, path, operation string) bool {
	listener := utils.GetListenerFromDir(base_path)
	rel_path := utils.GetRelativePath(listener, path)
	fsItem, err := utils.GetFileInfo(path)

	if err != nil {
		utils.LogWriteF("Error getting file details for %s: %+v", path, err)
	}

	dbItem, err := datastore.GetOne(base_path, rel_path)

	if !utils.ErrorCheckF(err, 400, "Error getting file row (%s)", rel_path) {
		switch operation {
		/*case "chmod:":
		  if dbItem.Perms != 0664{
		      logs.WriteLn("Perms don't match")
		  }*/
		case "create":
			if dbItem.Checksum != fsItem.Checksum {
				utils.LogWriteF("Creating:->")
				datastore.Insert(listener, fsItem)
				utils.LogWriteF("Putting in storage:->")
				storage.PutFile(path, listener)
			}
		case "write":
			if dbItem.Checksum != fsItem.Checksum {
				utils.LogWriteF("Writing:->")
				datastore.Insert(listener, fsItem)
				utils.LogWriteF("Putting in storage:->")
				storage.PutFile(path, listener)
			}
		case "remove":
			// Item was removed so we just wipe it from DB and storage
			datastore.Remove(listener, fsItem)
		}
	}

	return false
}

func checkItemInDB(base_path, rel_path, abspath string) bool {
	dbItem, err := datastore.GetOne(base_path, rel_path)
	if err != nil {
		if dbItem.Checksum != utils.GetMd5Checksum(abspath) {
			utils.WriteLn("Checksums do not match running update:")
			return true
		} else {
			utils.WriteLn("File already matches the DB")
			return false
		}
	}
	return false
}
