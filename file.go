package main

import (
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/storage"
	"github.com/chrislusf/seaweedfs/weed/storage/idx"
	"github.com/chrislusf/seaweedfs/weed/storage/needle"
	"github.com/chrislusf/seaweedfs/weed/storage/super_block"
	"github.com/chrislusf/seaweedfs/weed/storage/types"
	"github.com/chrislusf/seaweedfs/weed/util"
	"log"
	"os"
	"path"
	"strconv"
)

func FileId(chunk *filer_pb.FileChunk) *needle.FileId {
	VolumeId, _ := needle.NewVolumeId(strconv.FormatInt(int64(chunk.Fid.GetVolumeId()), 10))
	return needle.NewFileId(VolumeId, chunk.Fid.GetFileKey(), chunk.Fid.GetCookie())
}

func DataFromVolumeFile(chunk *filer_pb.FileChunk) ([]byte, error) {
	fileId := FileId(chunk)
	n, err := Needle(fileId)
	if err != nil {
		return nil, err
	}
	data := n.Data
	return getMaybeDecryptData(chunk, data)
}

func Needle(fileId *needle.FileId) (*needle.Needle, error) {
	log.Println("fileId: ", fileId)
	volumeCollection, volumeServerDir := VolumeCollectionAndDir(fileId.GetVolumeId().String())
	// Needle Index
	needleId := fileId.GetNeedleId()
	offset, size, err := NeedleOffsetAndSizeFromIdx(fileId.GetVolumeId().String(), volumeCollection, volumeServerDir, needleId)
	if err != nil {
		log.Printf("read Needle idx error: %v", err)
		return nil, err
	}
	log.Printf("根据fileId %v 找到了 needleId %v offset %v size %s", fileId, needleId, offset, util.BytesToHumanReadable(uint64(size)))

	// Volume
	v, err := storage.NewVolume(volumeServerDir, volumeServerDir, volumeCollection, fileId.GetVolumeId(), storage.NeedleMapInMemory, &super_block.ReplicaPlacement{}, &needle.TTL{}, 0, 0)
	if err != nil {
		log.Printf("volume creation: %v", err)
		return nil, err
	}
	backend := v.DataBackend

	// Needle
	n := new(needle.Needle)
	n.Id = needleId
	// read needle data
	err = n.ReadData(backend, offset.ToActualOffset(), size, v.Version())
	if err != nil {
		log.Printf("read needle data:%v", err)
		return nil, err
	}
	return n, nil
}

func NeedleOffsetAndSizeFromIdx(fixVolumeId, fixVolumeCollection, fixVolumePath string, targetNeedleId types.NeedleId) (types.Offset, types.Size, error) {
	fileName := fixVolumeId
	if fixVolumeCollection != "" {
		fileName = fixVolumeCollection + "_" + fileName
	}
	indexFile, err := os.OpenFile(path.Join(fixVolumePath, fileName+".idx"), os.O_RDONLY, 0644)
	if err != nil {
		glog.Fatalf("Create Volume Index [ERROR] %s\n", err)
	}
	defer indexFile.Close()

	var _offset types.Offset
	var _size types.Size
	err = fmt.Errorf("the target needle is not in this volume")
	idx.WalkIndexFile(indexFile, func(key types.NeedleId, offset types.Offset, size types.Size) error {
		//fmt.Printf("key:%v offset:%v size:%v(%v)\n", key, offset, size, util.BytesToHumanReadable(uint64(size)))
		if key == targetNeedleId {
			_offset = offset
			_size = size
			err = nil
		}
		return nil
	})
	return _offset, _size, err
}
