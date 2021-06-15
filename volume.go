package main

import (
	"strconv"
)

type VolumeStatus struct {
	DiskStatuses []struct {
		Dir string `json:"dir"`
	} `json:"DiskStatuses"`
	Volumes []Vol `json:"Volumes"`
}
type Vol struct {
	Id         int    `json:"Id"`
	Collection string `json:"Collection"`
}

func VolumeCollectionAndDir(volumeId string) (string, string) {
	var volumeCollection string
	var volumeServerDir string

	for _, vs := range CONFIG.MVS {
		for _, vol := range vs.Volumes {
			if strconv.FormatInt(int64(vol.Id), 10) == volumeId {
				volumeCollection = vol.Collection
				volumeServerDir = vs.DiskStatuses[0].Dir
			}
		}
	}
	return volumeCollection, volumeServerDir
}
