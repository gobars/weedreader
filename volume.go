package main

import (
	"encoding/json"
	"log"
	"strconv"
)

type VolumeStatus struct {
	DiskStatuses []struct {
		Dir         string  `json:"dir"`
		All         int64   `json:"all"`
		Used        int64   `json:"used"`
		Free        int64   `json:"free"`
		PercentFree float64 `json:"percent_free"`
		PercentUsed float64 `json:"percent_used"`
	} `json:"DiskStatuses"`
	Version string `json:"Version"`
	Volumes []Vol  `json:"Volumes"`
}
type Vol struct {
	Id               int `json:"Id"`
	Size             int `json:"Size"`
	ReplicaPlacement struct {
		Node int `json:"node"`
	} `json:"ReplicaPlacement"`
	Ttl struct {
		Count int `json:"Count"`
		Unit  int `json:"Unit"`
	} `json:"Ttl"`
	DiskType          string `json:"DiskType"`
	Collection        string `json:"Collection"`
	Version           int    `json:"Version"`
	FileCount         int    `json:"FileCount"`
	DeleteCount       int    `json:"DeleteCount"`
	DeletedByteCount  int    `json:"DeletedByteCount"`
	ReadOnly          bool   `json:"ReadOnly"`
	CompactRevision   int    `json:"CompactRevision"`
	ModifiedAtSecond  int    `json:"ModifiedAtSecond"`
	RemoteStorageName string `json:"RemoteStorageName"`
	RemoteStorageKey  string `json:"RemoteStorageKey"`
}

var mvs []VolumeStatus

func init() {
	// volume status
	mvs = volStatus()
}

// Such as http://localhost:9431/status?pretty=y http://localhost:9432/status?pretty=y http://localhost:9433/status?pretty=y
func volStatus() []VolumeStatus {
	var mvs []VolumeStatus
	vss := []string{
		`{"DiskStatuses":[{"dir":"/Users/anan/Downloads/weed/data/v1","all":250685575168,"used":100160446464,"free":150525128704,"percent_free":60.045387,"percent_used":39.954613}],"Version":"30GB 2.50 4233ad3","Volumes":[{"Id":1,"Size":12583032,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":3,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""},{"Id":2,"Size":15113070,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":5,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""}]}`,
		`{"DiskStatuses":[{"dir":"/Users/anan/Downloads/weed/data/v2","all":250685575168,"used":100170776576,"free":150514798592,"percent_free":60.041267,"percent_used":39.958733}],"Version":"30GB 2.50 4233ad3","Volumes":[{"Id":1,"Size":12583032,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":3,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""},{"Id":3,"Size":28378,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":2,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""},{"Id":4,"Size":4194344,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":1,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""},{"Id":5,"Size":8388688,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":2,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""},{"Id":6,"Size":12585526,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":5,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""}]}`,
		`{"DiskStatuses":[{"dir":"/Users/anan/Downloads/weed/data/v3","all":250685575168,"used":100153532416,"free":150532042752,"percent_free":60.048145,"percent_used":39.951855}],"Version":"30GB 2.50 4233ad3","Volumes":[{"Id":2,"Size":15113070,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":5,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""},{"Id":3,"Size":28378,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":2,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""},{"Id":4,"Size":4194344,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":1,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""},{"Id":5,"Size":8388688,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":2,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""},{"Id":6,"Size":12585526,"ReplicaPlacement":{"node":1},"Ttl":{"Count":0,"Unit":0},"DiskType":"","Collection":"","Version":3,"FileCount":5,"DeleteCount":0,"DeletedByteCount":0,"ReadOnly":false,"CompactRevision":0,"ModifiedAtSecond":0,"RemoteStorageName":"","RemoteStorageKey":""}]}`,
	}
	for _, j := range vss {
		var vs VolumeStatus
		err := json.Unmarshal([]byte(j), &vs)
		if err != nil {
			log.Printf("parse volStatus error: %v", err)
		}
		mvs = append(mvs, vs)
	}

	return mvs
}

func VolumeCollectionAndDir(volumeId string) (string, string) {
	var volumeCollection string
	var volumeServerDir string
	for _, vs := range mvs {
		for _, vol := range vs.Volumes {
			if strconv.FormatInt(int64(vol.Id), 10) == volumeId {
				volumeCollection = vol.Collection
				volumeServerDir = vs.DiskStatuses[0].Dir
			}
		}
	}
	return volumeCollection, volumeServerDir
}
