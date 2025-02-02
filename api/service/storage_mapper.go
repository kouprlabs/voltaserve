// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

type StorageMapper interface {
	mapStorageUsage(byteCount int64, maxBytes int64) *StorageUsage
}

type storageMapper struct{}

func newStorageMapper() StorageMapper {
	return &storageMapper{}
}

func (mp *storageMapper) mapStorageUsage(byteCount int64, maxBytes int64) *StorageUsage {
	res := StorageUsage{
		Bytes:    byteCount,
		MaxBytes: maxBytes,
	}
	if maxBytes != 0 {
		res.Percentage = int(byteCount * 100 / maxBytes)
	}
	return &res
}
