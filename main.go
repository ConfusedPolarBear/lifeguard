package main

import (
	"log"
	"fmt"

	"github.com/ConfusedPolarBear/zfs-manager/zpool"
)

func main() {
	log.Println("Parsing all pools on system")

	pools := zpool.ParseAllPools()

	for index, pool := range pools {
		fmt.Printf("Pool %d: %#v", index, pool)
	}
}
