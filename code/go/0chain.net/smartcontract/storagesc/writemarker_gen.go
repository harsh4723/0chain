package storagesc

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *writeMarkerV1) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 9
	// string "AllocationRoot"
	o = append(o, 0x89, 0xae, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x6f, 0x74)
	o = msgp.AppendString(o, z.AllocationRoot)
	// string "PreviousAllocationRoot"
	o = append(o, 0xb6, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x6f, 0x74)
	o = msgp.AppendString(o, z.PreviousAllocationRoot)
	// string "FileMetaRoot"
	o = append(o, 0xac, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x6f, 0x6f, 0x74)
	o = msgp.AppendString(o, z.FileMetaRoot)
	// string "AllocationID"
	o = append(o, 0xac, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44)
	o = msgp.AppendString(o, z.AllocationID)
	// string "Size"
	o = append(o, 0xa4, 0x53, 0x69, 0x7a, 0x65)
	o = msgp.AppendInt64(o, z.Size)
	// string "BlobberID"
	o = append(o, 0xa9, 0x42, 0x6c, 0x6f, 0x62, 0x62, 0x65, 0x72, 0x49, 0x44)
	o = msgp.AppendString(o, z.BlobberID)
	// string "Timestamp"
	o = append(o, 0xa9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	o, err = z.Timestamp.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Timestamp")
		return
	}
	// string "ClientID"
	o = append(o, 0xa8, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44)
	o = msgp.AppendString(o, z.ClientID)
	// string "Signature"
	o = append(o, 0xa9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65)
	o = msgp.AppendString(o, z.Signature)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *writeMarkerV1) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "AllocationRoot":
			z.AllocationRoot, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "AllocationRoot")
				return
			}
		case "PreviousAllocationRoot":
			z.PreviousAllocationRoot, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PreviousAllocationRoot")
				return
			}
		case "FileMetaRoot":
			z.FileMetaRoot, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "FileMetaRoot")
				return
			}
		case "AllocationID":
			z.AllocationID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "AllocationID")
				return
			}
		case "Size":
			z.Size, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Size")
				return
			}
		case "BlobberID":
			z.BlobberID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BlobberID")
				return
			}
		case "Timestamp":
			bts, err = z.Timestamp.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Timestamp")
				return
			}
		case "ClientID":
			z.ClientID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ClientID")
				return
			}
		case "Signature":
			z.Signature, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Signature")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *writeMarkerV1) Msgsize() (s int) {
	s = 1 + 15 + msgp.StringPrefixSize + len(z.AllocationRoot) + 23 + msgp.StringPrefixSize + len(z.PreviousAllocationRoot) + 13 + msgp.StringPrefixSize + len(z.FileMetaRoot) + 13 + msgp.StringPrefixSize + len(z.AllocationID) + 5 + msgp.Int64Size + 10 + msgp.StringPrefixSize + len(z.BlobberID) + 10 + z.Timestamp.Msgsize() + 9 + msgp.StringPrefixSize + len(z.ClientID) + 10 + msgp.StringPrefixSize + len(z.Signature)
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *writeMarkerV2) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 12
	// string "version"
	o = append(o, 0x8c, 0xa7, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Version)
	// string "AllocationRoot"
	o = append(o, 0xae, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x6f, 0x74)
	o = msgp.AppendString(o, z.AllocationRoot)
	// string "PreviousAllocationRoot"
	o = append(o, 0xb6, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x6f, 0x74)
	o = msgp.AppendString(o, z.PreviousAllocationRoot)
	// string "FileMetaRoot"
	o = append(o, 0xac, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x6f, 0x6f, 0x74)
	o = msgp.AppendString(o, z.FileMetaRoot)
	// string "AllocationID"
	o = append(o, 0xac, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44)
	o = msgp.AppendString(o, z.AllocationID)
	// string "Size"
	o = append(o, 0xa4, 0x53, 0x69, 0x7a, 0x65)
	o = msgp.AppendInt64(o, z.Size)
	// string "ChainSize"
	o = append(o, 0xa9, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x53, 0x69, 0x7a, 0x65)
	o = msgp.AppendInt64(o, z.ChainSize)
	// string "ChainHash"
	o = append(o, 0xa9, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x48, 0x61, 0x73, 0x68)
	o = msgp.AppendString(o, z.ChainHash)
	// string "BlobberID"
	o = append(o, 0xa9, 0x42, 0x6c, 0x6f, 0x62, 0x62, 0x65, 0x72, 0x49, 0x44)
	o = msgp.AppendString(o, z.BlobberID)
	// string "Timestamp"
	o = append(o, 0xa9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	o, err = z.Timestamp.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Timestamp")
		return
	}
	// string "ClientID"
	o = append(o, 0xa8, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44)
	o = msgp.AppendString(o, z.ClientID)
	// string "Signature"
	o = append(o, 0xa9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65)
	o = msgp.AppendString(o, z.Signature)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *writeMarkerV2) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "version":
			z.Version, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		case "AllocationRoot":
			z.AllocationRoot, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "AllocationRoot")
				return
			}
		case "PreviousAllocationRoot":
			z.PreviousAllocationRoot, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PreviousAllocationRoot")
				return
			}
		case "FileMetaRoot":
			z.FileMetaRoot, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "FileMetaRoot")
				return
			}
		case "AllocationID":
			z.AllocationID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "AllocationID")
				return
			}
		case "Size":
			z.Size, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Size")
				return
			}
		case "ChainSize":
			z.ChainSize, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ChainSize")
				return
			}
		case "ChainHash":
			z.ChainHash, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ChainHash")
				return
			}
		case "BlobberID":
			z.BlobberID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BlobberID")
				return
			}
		case "Timestamp":
			bts, err = z.Timestamp.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Timestamp")
				return
			}
		case "ClientID":
			z.ClientID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ClientID")
				return
			}
		case "Signature":
			z.Signature, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Signature")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *writeMarkerV2) Msgsize() (s int) {
	s = 1 + 8 + msgp.StringPrefixSize + len(z.Version) + 15 + msgp.StringPrefixSize + len(z.AllocationRoot) + 23 + msgp.StringPrefixSize + len(z.PreviousAllocationRoot) + 13 + msgp.StringPrefixSize + len(z.FileMetaRoot) + 13 + msgp.StringPrefixSize + len(z.AllocationID) + 5 + msgp.Int64Size + 10 + msgp.Int64Size + 10 + msgp.StringPrefixSize + len(z.ChainHash) + 10 + msgp.StringPrefixSize + len(z.BlobberID) + 10 + z.Timestamp.Msgsize() + 9 + msgp.StringPrefixSize + len(z.ClientID) + 10 + msgp.StringPrefixSize + len(z.Signature)
	return
}