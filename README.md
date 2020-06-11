# lifeguard
![GitHub](https://img.shields.io/github/license/ConfusedPolarBear/lifeguard)

## Short description
lifeguard is a ZFS pool management web interface which supports pool status, snapshots, and current data usage. 

## Miscellaneous
### Why create this instead of installing TrueNAS/OpenMediaVault/etc?
This project was created because I wanted to learn Go and Vue.

### How is this licensed?
All files (except one) are AGPL version 3 licensed. The file ``pkg/zpool/parser-gpl.go`` is licensed under the GPL version 3 and only contains the structs used for holding pool layout.

### Where did the name lifeguard come from?
Since this project monitors and reports on ZFS pool(s), the name lifeguard seems fitting.
