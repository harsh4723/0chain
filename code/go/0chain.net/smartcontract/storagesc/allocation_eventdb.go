package storagesc

import (
	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/state"
	"0chain.net/core/common"
	"0chain.net/smartcontract/dbs/event"
	"0chain.net/smartcontract/stakepool"
	"encoding/json"
	"fmt"
	"time"
)

func allocationTableToStorageAllocation(alloc *event.Allocation, eventDb *event.EventDb) (*StorageAllocation, error) {
	storageNodes := make([]*StorageNode, 0)
	blobberDetails := make([]*BlobberAllocation, 0)
	blobberIDs := make([]string, 0)
	blobberIDTermMapping := make(map[string]struct {
		AllocationID string
		Terms
	})
	blobberMap := make(map[string]*BlobberAllocation)

	curators, err := eventDb.GetCuratorsByAllocationID(alloc.AllocationID)
	if err != nil {
		return nil, fmt.Errorf("error finding curators: %v", err)
	}

	var allocTerms []event.AllocationTerm
	err = json.Unmarshal([]byte(alloc.Terms), &allocTerms)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling allocation terms: %v", err)
	}

	for _, t := range allocTerms {
		blobberIDs = append(blobberIDs, t.BlobberID)
		blobberIDTermMapping[t.BlobberID] = struct {
			AllocationID string
			Terms
		}{
			AllocationID: t.AllocationID,
			Terms: Terms{
				ReadPrice:               t.ReadPrice,
				WritePrice:              t.WritePrice,
				MinLockDemand:           t.MinLockDemand,
				MaxOfferDuration:        t.MaxOfferDuration,
				ChallengeCompletionTime: t.ChallengeCompletionTime,
			}}
	}

	blobbers, err := eventDb.GetBlobbersFromIDs(blobberIDs)
	if err != nil {
		return nil, fmt.Errorf("error retrieving blobbers from db: %v", err)
	}

	var dpsSze = alloc.DataShards + alloc.ParityShards
	var gbSize = sizeInGB((alloc.Size + int64(dpsSze-1)) / int64(dpsSze))
	var rdtu = float64(time.Second*time.Duration(alloc.Expiration-alloc.StartTime)) / float64(alloc.TimeUnit)

	for _, b := range blobbers {
		storageNodes = append(storageNodes, &StorageNode{
			ID:      b.BlobberID,
			BaseURL: b.BaseURL,
			Geolocation: StorageNodeGeolocation{
				Latitude:  b.Latitude,
				Longitude: b.Longitude,
			},
			Terms:           blobberIDTermMapping[b.BlobberID].Terms,
			Capacity:        b.Capacity,
			Used:            b.Used,
			SavedData:       b.SavedData,
			LastHealthCheck: common.Timestamp(b.LastHealthCheck),
			StakePoolSettings: stakepool.StakePoolSettings{
				DelegateWallet:  b.DelegateWallet,
				MinStake:        state.Balance(b.MinStake),
				MaxStake:        state.Balance(b.MaxStake),
				MaxNumDelegates: b.NumDelegates,
				ServiceCharge:   b.ServiceCharge,
			},
		})

		terms := blobberIDTermMapping[b.BlobberID].Terms
		tempBlobberAllocation := &BlobberAllocation{
			BlobberID:     b.BlobberID,
			AllocationID:  blobberIDTermMapping[b.BlobberID].AllocationID,
			Size:          b.Used,
			Terms:         terms,
			MinLockDemand: state.Balance(float64(terms.WritePrice) * gbSize * terms.MinLockDemand * rdtu),
		}
		blobberDetails = append(blobberDetails, tempBlobberAllocation)
		blobberMap[b.BlobberID] = tempBlobberAllocation
	}

	sa := &StorageAllocation{
		ID:             alloc.AllocationID,
		Tx:             alloc.TransactionID,
		DataShards:     alloc.DataShards,
		ParityShards:   alloc.ParityShards,
		Size:           alloc.Size,
		Expiration:     common.Timestamp(alloc.Expiration),
		Blobbers:       storageNodes,
		Owner:          alloc.Owner,
		OwnerPublicKey: alloc.OwnerPublicKey,
		Stats: &StorageAllocationStats{
			UsedSize:                  alloc.UsedSize,
			NumWrites:                 alloc.NumWrites,
			ReadsSize:                 alloc.ReadSize,
			TotalChallenges:           alloc.TotalChallenges,
			OpenChallenges:            alloc.OpenChallenges,
			SuccessChallenges:         alloc.SuccessfulChallenges,
			FailedChallenges:          alloc.FailedChallenges,
			LastestClosedChallengeTxn: alloc.LatestClosedChallengeTxn,
		},
		BlobberDetails:             blobberDetails,
		BlobberMap:                 blobberMap,
		IsImmutable:                alloc.IsImmutable,
		ReadPriceRange:             PriceRange{alloc.ReadPriceMin, alloc.ReadPriceMax},
		WritePriceRange:            PriceRange{alloc.WritePriceMin, alloc.WritePriceMax},
		MaxChallengeCompletionTime: time.Duration(alloc.MaxChallengeCompletionTime),
		// todo: to be added with WritePool : select user_id from WritePools where allocation_id = ?
		// WritePoolOwners:            nil,
		ChallengeCompletionTime: time.Duration(alloc.ChallengeCompletionTime),
		StartTime:               common.Timestamp(alloc.StartTime),
		Finalized:               alloc.Finalized,
		Canceled:                alloc.Cancelled,
		UsedSize:                alloc.UsedSize,
		MovedToChallenge:        alloc.MovedToChallenge,
		MovedBack:               alloc.MovedBack,
		MovedToValidators:       alloc.MovedToValidators,
		TimeUnit:                time.Duration(alloc.TimeUnit),
		Curators:                curators,
	}

	return sa, nil
}

func storageAllocationToAllocationTable(sa *StorageAllocation) (*event.Allocation, error) {
	allocationTerms := make([]event.AllocationTerm, 0)
	for _, b := range sa.BlobberDetails {
		allocationTerms = append(allocationTerms, event.AllocationTerm{
			BlobberID:               b.BlobberID,
			AllocationID:            b.AllocationID,
			ReadPrice:               b.Terms.ReadPrice,
			WritePrice:              b.Terms.WritePrice,
			MinLockDemand:           b.Terms.MinLockDemand,
			MaxOfferDuration:        b.Terms.MaxOfferDuration,
			ChallengeCompletionTime: b.Terms.ChallengeCompletionTime,
		})
	}

	termsByte, err := json.Marshal(allocationTerms)
	if err != nil {
		return nil, fmt.Errorf("error marshalling terms: %v", err)
	}

	alloc := &event.Allocation{
		AllocationID:               sa.ID,
		TransactionID:              sa.Tx,
		DataShards:                 sa.DataShards,
		ParityShards:               sa.ParityShards,
		Size:                       sa.Size,
		Expiration:                 int64(sa.Expiration),
		Terms:                      string(termsByte),
		Owner:                      sa.Owner,
		OwnerPublicKey:             sa.OwnerPublicKey,
		IsImmutable:                sa.IsImmutable,
		ReadPriceMin:               sa.ReadPriceRange.Min,
		ReadPriceMax:               sa.ReadPriceRange.Max,
		WritePriceMin:              sa.WritePriceRange.Min,
		WritePriceMax:              sa.WritePriceRange.Max,
		MaxChallengeCompletionTime: int64(sa.MaxChallengeCompletionTime),
		ChallengeCompletionTime:    int64(sa.ChallengeCompletionTime),
		StartTime:                  int64(sa.StartTime),
		Finalized:                  sa.Finalized,
		Cancelled:                  sa.Canceled,
		UsedSize:                   sa.UsedSize,
		MovedToChallenge:           sa.MovedToChallenge,
		MovedBack:                  sa.MovedBack,
		MovedToValidators:          sa.MovedToValidators,
		TimeUnit:                   int64(sa.TimeUnit),
	}

	if sa.Stats != nil {
		alloc.NumWrites = sa.Stats.NumWrites
		alloc.ReadSize = sa.Stats.ReadsSize
		alloc.TotalChallenges = sa.Stats.TotalChallenges
		alloc.OpenChallenges = sa.Stats.OpenChallenges
		alloc.SuccessfulChallenges = sa.Stats.SuccessChallenges
		alloc.FailedChallenges = sa.Stats.FailedChallenges
		alloc.LatestClosedChallengeTxn = sa.Stats.LastestClosedChallengeTxn
	}

	return alloc, nil
}

func emitAddOrOverwriteAllocation(sa *StorageAllocation, balances cstate.StateContextI) error {
	alloc, err := storageAllocationToAllocationTable(sa)
	if err != nil {
		return err
	}

	data, err := json.Marshal(alloc)
	if err != nil {
		return fmt.Errorf("error marshalling allocation: %v", err)
	}

	balances.EmitEvent(event.TypeStats, event.TagAddOrOverwriteAllocation, alloc.AllocationID, string(data))

	return nil
}

func getStorageAllocationFromDb(id string, eventDb *event.EventDb) (*StorageAllocation, error) {
	alloc, err := eventDb.GetAllocation(id)
	if err != nil {
		return nil, err
	}

	sa, err := allocationTableToStorageAllocation(alloc, eventDb)
	if err != nil {
		return nil, err
	}

	return sa, nil
}

func getClientAllocationsFromDb(clientID string, eventDb *event.EventDb) ([]*StorageAllocation, error) {

	sas := make([]*StorageAllocation, 0)

	allocs, err := eventDb.GetClientsAllocation(clientID)
	if err != nil {
		return nil, err
	}

	for _, alloc := range allocs {
		sa, err := allocationTableToStorageAllocation(&alloc, eventDb)
		if err != nil {
			return nil, err
		}

		sas = append(sas, sa)
	}

	return sas, nil
}