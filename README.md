# LifeGuard
![GitHub](https://img.shields.io/github/license/ConfusedPolarBear/lifeguard) [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ConfusedPolarBear_lifeguard&metric=security_rating)](https://sonarcloud.io/dashboard?id=ConfusedPolarBear_lifeguard)

LifeGuard is a ZFS pool management web interface which supports pool status, snapshots, and current data usage.

## Installation
LifeGuard can be installed by source using make. First download the source.

`git clone https://github.com/ConfusedPolarBear/lifeguard.git && cd lifeguard`

Within the config folder, create a file named `browser.ini` following the template of `example_browser.ini`. This is used to control whether the included file browser is enabled and what folders it can access.

Once created, install using make.

`make && make install`

Run lifeguard. By default, the application is accessible at port 5120.

`./lifeguard`

## FAQ

### How is this licensed?
All files (except one) are AGPL version 3 licensed. The file ``pkg/structs/pool.go`` is licensed under the GPL version 3 and only contains the structs used for holding pool layout.

### Where did the name lifeguard come from?
Since this project monitors and reports on ZFS pool(s), the name lifeguard seems fitting.
