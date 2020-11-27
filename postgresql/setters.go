package postgresql

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/models"
)

func AddOrUpdateGoEngine(token string, machineID string, ipAddress string, port int) error {
	newData := &models.EngineModel{
		Token:     token,
		MachineID: machineID,
		IPAddress: ipAddress,
		Port:      port,
		Enabled:   true,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	engine := GetGoEngine(machineID)
	if engine.ID != "" {
		newData.CreatedAt = engine.CreatedAt
		_, err := db.Model(newData).Where("id = ?", engine.ID).Update()
		if err != nil {
			return err
		}
		return nil
	}
	newID, _ := uuid.NewUUID()
	newData.ID = newID.String()
	newData.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Model(newData).Insert()
	if err != nil {
		return err
	}
	return nil
}

func AddorUpdateReplication(name string, completed bool, log string) error {
	newData := &models.ReplicationModel{
		MachineID: helpers.MachineID,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Completed: completed,
		Log:       log,
		Key:       name,
	}
	replication := GetReplication(name)
	if replication.ID != "" {
		newData.CreatedAt = replication.CreatedAt
		_, err := db.Model(newData).Where("id = ?", replication.ID).Update()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		return nil
	}
	newID, _ := uuid.NewUUID()
	newData.ID = newID.String()
	newData.Key = name
	newData.Completed = completed
	newData.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Model(newData).Insert()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func StoreEngineData() {
	key, _ := uuid.NewUUID()
	machineID, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		panic(err.Error())
	}
	localIP := helpers.GetLocalIP()
	if localIP == "" {
		panic("Cannot find local IP Address, please add CURRENT_IP to configuration.")
	} else {
		log.Printf("Current IP Address %v\n", localIP)
	}
	helpers.MachineID = strings.TrimSpace(strings.ToUpper(string(machineID)))
	err = AddOrUpdateGoEngine(key.String(), helpers.MachineID, localIP, 5454)
	if err != nil {
		log.Panic(err.Error())
	}
}
