package main

import (
	"testing"

	"github.com/ConfusedPolarBear/zfs-manager/zpool"
)

func areEqual(name string, expected string, actual string, test *testing.T) {
	if expected != actual {
		test.Errorf("Error testing %s - expected %s, was %s.", name, expected, actual)
	}
}

// TODO: test containers in all pools. property is containers []*ContainerStatus

func TestHealthy(t *testing.T) {
	var output = `  pool: test
 state: ONLINE
  scan: scrub repaired 0B in 0 days 00:01:32 with 0 errors on Wed May 27 01:25:56 2020
config:

        NAME                               STATE     READ WRITE CKSUM
        test                               ONLINE       0     0     0
          mirror-0                         ONLINE       0     0     0
            pci-0000:03:00.0-scsi-0:0:1:0  ONLINE       0     0     0
            pci-0000:03:00.0-scsi-0:0:2:0  ONLINE       0     0     0

errors: No known data errors`

	var parsed = zpool.ParseZpoolStatus(output)

	areEqual("healthy pool name", "test", parsed.Name, t)
	areEqual("healthy pool state", "ONLINE", parsed.State, t)
	areEqual("healthy pool status", "OK", parsed.Status, t)
	areEqual("healthy pool action", "No action needed", parsed.Action, t)
	areEqual("healthy pool see", "", parsed.See, t)
	areEqual("healthy pool scan", "scrub repaired 0B in 0 days 00:01:32 with 0 errors on Wed May 27 01:25:56 2020", parsed.Scan, t)
	areEqual("healthy pool errors", "No known data errors", parsed.Errors, t)
	areEqual("healthy pool raw", output, parsed.Raw, t)
}

func TestChecksumErrors(t *testing.T) {
	var output = `  pool: test
 state: ONLINE
status: One or more devices has experienced an unrecoverable error.  An
        attempt was made to correct the error.  Applications are unaffected.
action: Determine if the device needs to be replaced, and clear the errors
        using 'zpool clear' or replace the device with 'zpool replace'.
   see: http://zfsonlinux.org/msg/ZFS-8000-9P
  scan: scrub repaired 177M in 0 days 00:01:25 with 0 errors on Tue May 19 00:06:01 2020
config:

        NAME                               STATE     READ WRITE CKSUM
        test                               ONLINE       0     0     0
          mirror-0                         ONLINE       0     0     0
            pci-0000:03:00.0-scsi-0:0:1:0  ONLINE       0     0     0
            pci-0000:03:00.0-scsi-0:0:2:0  ONLINE       0     0 11.7K

errors: No known data errors`

	var parsed = zpool.ParseZpoolStatus(output)

	areEqual("checksum pool name", "test", parsed.Name, t)
	areEqual("checksum pool state", "ONLINE", parsed.State, t)
	areEqual("checksum pool status", "One or more devices has experienced an unrecoverable error.  An\nattempt was made to correct the error.  Applications are unaffected.", parsed.Status, t)
	areEqual("checksum pool action", "Determine if the device needs to be replaced, and clear the errors\nusing 'zpool clear' or replace the device with 'zpool replace'.", parsed.Action, t)
	areEqual("checksum pool see", "http://zfsonlinux.org/msg/ZFS-8000-9P", parsed.See, t)
	areEqual("checksum pool scan", "scrub repaired 177M in 0 days 00:01:25 with 0 errors on Tue May 19 00:06:01 2020", parsed.Scan, t)
	areEqual("checksum pool errors", "No known data errors", parsed.Errors, t)
	areEqual("checksum pool raw", output, parsed.Raw, t)
}

func TestResilvering(t *testing.T) {
	var output = `  pool: test
 state: DEGRADED
status: One or more devices is currently being resilvered.  The pool will
        continue to function, possibly in a degraded state.
action: Wait for the resilver to complete.
  scan: resilver in progress since Wed May 27 01:32:44 2020
        23.5G scanned at 2.61G/s, 228M issued at 25.3M/s, 23.5G total
        236M resilvered, 0.95% done, 0 days 00:15:40 to go
config:

        NAME                                 STATE     READ WRITE CKSUM
        test                                 DEGRADED     0     0     0
          mirror-0                           DEGRADED     0     0     0
            pci-0000:03:00.0-scsi-0:0:1:0    ONLINE       0     0     0
            replacing-1                      DEGRADED     0     0     0
              old                            UNAVAIL      4     1     0
              pci-0000:03:00.0-scsi-0:0:2:0  ONLINE       0     0     0  (resilvering)

errors: No known data errors`

	var parsed = zpool.ParseZpoolStatus(output)

	areEqual("resilvering pool name", "test", parsed.Name, t)
	areEqual("resilvering pool state", "DEGRADED", parsed.State, t)
	areEqual("resilvering pool status", "One or more devices is currently being resilvered.  The pool will\ncontinue to function, possibly in a degraded state.", parsed.Status, t)
	areEqual("resilvering pool action", "Wait for the resilver to complete.", parsed.Action, t)
	areEqual("resilvering pool see", "", parsed.See, t)
	areEqual("resilvering pool scan", "resilver in progress since Wed May 27 01:32:44 2020\n23.5G scanned at 2.61G/s, 228M issued at 25.3M/s, 23.5G total\n236M resilvered, 0.95% done, 0 days 00:15:40 to go", parsed.Scan, t)
	areEqual("resilvering pool errors", "No known data errors", parsed.Errors, t)
	areEqual("resilvering pool raw", output, parsed.Raw, t)
}

func TestPermanentError(t *testing.T) {
	// zpool status output from https://serverfault.com/questions/800628/what-does-a-permanent-zfs-error-indicate
	var output = `  pool: seagate3tb
 state: ONLINE
status: One or more devices has experienced an error resulting in data
        corruption.  Applications may be affected.
action: Restore the file in question if possible.  Otherwise restore the
        entire pool from backup.
   see: http://zfsonlinux.org/msg/ZFS-8000-8A
  scan: none requested
config:

        NAME        STATE     READ WRITE CKSUM
        seagate3tb  ONLINE       0     0    28
          sda       ONLINE       0     0    56

errors: Permanent errors have been detected in the following files:

        /mnt/seagate3tb/Install.iso
        /mnt/seagate3tb/some-other-file1.txt
        /mnt/seagate3tb/some-other-file2.txt`

	var parsed = zpool.ParseZpoolStatus(output)

	areEqual("resilvering pool name", "seagate3tb", parsed.Name, t)
	areEqual("resilvering pool state", "ONLINE", parsed.State, t)
	areEqual("resilvering pool status", "One or more devices has experienced an error resulting in data\ncorruption.  Applications may be affected.", parsed.Status, t)
	areEqual("resilvering pool action", "Restore the file in question if possible.  Otherwise restore the\nentire pool from backup.", parsed.Action, t)
	areEqual("resilvering pool see", "http://zfsonlinux.org/msg/ZFS-8000-8A", parsed.See, t)
	areEqual("resilvering pool scan", "none requested", parsed.Scan, t)
	areEqual("resilvering pool errors", "Permanent errors have been detected in the following files:\n/mnt/seagate3tb/Install.iso\n/mnt/seagate3tb/some-other-file1.txt\n/mnt/seagate3tb/some-other-file2.txt", parsed.Errors, t)
	areEqual("resilvering pool raw", output, parsed.Raw, t)
}

func TestSanitizing(t *testing.T) {
	var evil = "; cat /etc/shadow"

	areEqual("argument sanitization", "catetcshadow", zpool.Sanitize(evil), t)
}

