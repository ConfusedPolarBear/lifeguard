package main

import (
	"log"
	"fmt"

	"github.com/ConfusedPolarBear/zfs-manager/zpool"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pools := string(zpool.GetPools())

	fmt.Println(pools)
}
