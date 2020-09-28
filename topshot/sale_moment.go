package topshot

import (
	"context"
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

func GetSaleMomentFromOwnerAtBlock(flowClient *client.Client, blockHeight uint64, ownerAddress flow.Address, momentFlowID uint64) (*SaleMoment, error) {
	getSaleMomentScript := `
		import TopShot from 0x0b2a3299cc857e29
        import Market from 0xc1e4f4f4c4257510
		pub fun main(owner:Address, momentID:UInt64): [UInt32] {
			let acct = getAccount(owner)
            let collectionRef = acct.getCapability(/public/topshotSaleCollection)!.borrow<&{Market.SalePublic}>() ?? panic("Could not borrow capability from public collection")
			let saleMoment = collectionRef.borrowMoment(id: momentID)!
			return [saleMoment.data.playID, saleMoment.data.setID, saleMoment.data.serialNumber]
		}
`
	res, err := flowClient.ExecuteScriptAtBlockHeight(context.Background(), blockHeight, []byte(getSaleMomentScript), []cadence.Value{
		cadence.BytesToAddress(ownerAddress.Bytes()),
		cadence.UInt64(momentFlowID),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching sale moment from flow: %w", err)
	}
	saleMoment := SaleMoment(res.(cadence.Array))
	return &saleMoment, nil
}

type SaleMoment cadence.Array

func (s SaleMoment) PlayID() uint32 {
	return uint32(s.Values[0].(cadence.UInt32))

}

func (s SaleMoment) SetID() uint32 {
	return uint32(s.Values[1].(cadence.UInt32))
}

func (s SaleMoment) SerialNumber() uint32 {
	return uint32(s.Values[2].(cadence.UInt32))
}

func (s SaleMoment) String() string {
	return fmt.Sprintf("saleMoment: serialNumber: %d, setID: %d, playID: %d",
		s.SerialNumber(), s.SetID(), s.PlayID())
}
