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

/*
Run etcd
./etcd -name etcd0 -data-dir data/0/ -advertise-client-urls http://127.0.0.1:2379 \
 -listen-client-urls http://127.0.0.1:2379 \
 -initial-advertise-peer-urls http://127.0.0.1:2380 \
 -listen-peer-urls http://127.0.0.1:2380 \
 -initial-cluster-token etcd-cluster \
 -initial-cluster etcd0=http://127.0.0.1:2380 \
 -initial-cluster-state new
*/

const (
	cCinderConfForNimble = "{\"san_ip\":\"10.23.0.1\",\"san_login\":\"nimble\",\"san_password\":\"nimble123\",\"volume_backend_name\":\"nimble_onprem\"}"
)

var (
	ErrRestartingCinder = errors.New("failed to restart cinder")
	randSource          = rand.New(rand.NewSource(time.Now().UnixNano()))
)

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

func restartCinderService() error {
	//if randSource.Intn(100) < 50 {
	//	return ErrRestartingCinder
	//}

	log.Println("cinder service restarted")
	return nil
}

func main() {
	// BEGIN transaction

	// ################## STEP 1 ##################
	// setup client to etcd
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	}
	c, err := client.New(cfg)
	if err != nil {
		panic(err)
	}
	keysAPI := client.NewKeysAPI(c)

	// set the key
	resp, err := keysAPI.Set(context.Background(), "/cinder/storage/nimble_onprem",
		cCinderConfForNimble, nil)
	if err != nil {
		panic(err)
	}
	log.Printf("Set is done. Metadata is %q\n", resp)

	// ################## STEP 2 ##################
	if err := writeToFile("priv/cinder/cinder.conf",
		[]byte(cCinderConfForNimble)); err != nil {

		panic(err)
	}

	// ################## STEP 3 ##################
	if err := restartCinderService(); err != nil {
		panic(err)
	}

	// END transaction
	log.Println("transaction complete")
}
