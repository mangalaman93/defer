package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/coreos/etcd/client"
)

const (
	cCinderConfFilePath = "priv/cinder/nimble.conf"
	cEtcdKey            = "/cinder/storage/nimble_onprem"

	cCinderConfForNimble = "{\"san_ip\":\"10.23.0.1\",\"san_login\":\"nimble\",\"san_password\":\"nimble123\",\"volume_backend_name\":\"nimble_onprem\"}"
)

var (
	ErrRestartingCinder = errors.New("failed to restart cinder")
	ErrWritingToFile    = errors.New("failed to write to cinder.conf")
	randSource          = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func restartCinderService() error {
	//if randSource.Intn(100) < 100 {
	//	return ErrRestartingCinder
	//}

	log.Println("cinder service restarted")
	return nil
}

func writeToFile(filepath string, data []byte) error {
	//if randSource.Intn(100) < 101 {
	//	return ErrWritingToFile
	//}

	err := ioutil.WriteFile(filepath, data, os.ModePerm)
	if err != nil {
		log.Println("failed to write to file")
		return err
	}

	log.Println("writing to file complete")
	return nil
}

func deleteFile(filepath string) error {
	if err := os.Remove(filepath); err != nil {
		log.Println("error in deleting file")
		return err
	}

	log.Println("deleting file complete")
	return nil
}

func step1(api client.KeysAPI) (*client.Response, error) {
	resp, err := api.Set(context.Background(), cEtcdKey, cCinderConfForNimble, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Set is done. Metadata is %q\n", resp)
	return resp, nil
}

func undoStep1(api client.KeysAPI, resp *client.Response) {
	if resp.PrevNode == nil {
		_, err := api.Delete(context.Background(), cEtcdKey, nil)
		if err != nil {
			log.Printf("unable to delete key:%v :: %v\n", cEtcdKey, err)
		}
		log.Printf("deleted key: %v\n", cEtcdKey)
	} else {
		_, err := api.Set(context.Background(), cEtcdKey, resp.PrevNode.Value, nil)
		if err != nil {
			log.Printf("unable to undo set for key:%v, value:%v :: %v\n", cEtcdKey, resp.PrevNode.Value, err)
		}
		log.Printf("undo set value for key: %v\n", cEtcdKey)
	}
}

func performTransaction() error {
	// BEGIN transaction

	// ################## STEP 1 ##################
	// setup client to etcd
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	}
	c, err := client.New(cfg)
	if err != nil {
		return err
	}
	keysAPI := client.NewKeysAPI(c)
	resp, err := step1(keysAPI)
	if err != nil {
		return err
	}

	// ################## STEP 2 ##################
	if err := writeToFile(cCinderConfFilePath, []byte(cCinderConfForNimble)); err != nil {
		// undo step 1
		undoStep1(keysAPI, resp)
		return err
	}

	// ################## STEP 3 ##################
	if err := restartCinderService(); err != nil {
		// undo step 1
		undoStep1(keysAPI, resp)

		// undo step 2
		deleteFile(cCinderConfFilePath)

		return err
	}

	return nil
	// END transaction
}

func main() {
	if err := performTransaction(); err != nil {
		log.Println("transaction failed ::", err)
	} else {
		log.Println("transaction complete")
	}
}
