package parser

import (
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

type prsRwSet struct {
	kvRWSet kvrwset.KVRWSet
}

type prsAction struct {
	payload *peer.ChaincodeActionPayload
}

// proposalResponsePayload returns a pointer to peer.ProposalResponsePayload.
// proposalResponsePayload is the payload of a proposal response.
// This message is the "bridge" between the client's request and the endorser's action in response to that request.
// Concretely, for chaincodes, it contains a hashed representation of the proposal (proposalHash) and a representation of the chaincode state changes and events inside the extension field.
func (a *prsAction) proposalResponsePayload() (*peer.ProposalResponsePayload, error) {
	proposalResponsePayload := &peer.ProposalResponsePayload{}
	if err := proto.Unmarshal(a.payload.Action.ProposalResponsePayload, proposalResponsePayload); err != nil {
		return nil, errors.Wrap(err, "unmarshal proposal response payload error")
	}
	return proposalResponsePayload, nil
}

// chaincodeAction returns a pointer to peer.ChaincodeAction that contains and actions the events generated by the execution of the chaincode.
func (a *prsAction) chaincodeAction() (*peer.ChaincodeAction, error) {
	proposalResponsePayload, err := a.proposalResponsePayload()
	if err != nil {
		return nil, err
	}
	chaincodeAction := &peer.ChaincodeAction{}
	if err := proto.Unmarshal(proposalResponsePayload.Extension, chaincodeAction); err != nil {
		return nil, errors.Wrap(err, "unmarshal chaincode action error")
	}
	return chaincodeAction, nil
}

// chaincodeEvent returns a pointer to peer.ChaincodeEvent that contains events info.
func (a *prsAction) chaincodeEvent() (*peer.ChaincodeEvent, error) {
	chaincodeAction, err := a.chaincodeAction()
	if err != nil {
		return nil, err
	}
	chaincodeEvent := &peer.ChaincodeEvent{}
	if err := proto.Unmarshal(chaincodeAction.Events, chaincodeEvent); err != nil {
		return nil, errors.Wrap(err, "unmarshal chaincode event error")
	}
	return chaincodeEvent, nil
}

// rwSets returns to a read-write sets slice.
func (a *prsAction) rwSets() ([]prsRwSet, error) {
	chaincodeAction, err := a.chaincodeAction()
	if err != nil {
		return nil, err
	}

	txReadWriteSet := &rwset.TxReadWriteSet{}
	if err := proto.Unmarshal(chaincodeAction.Results, txReadWriteSet); err != nil {
		return nil, errors.Wrap(err, "unmarshal tx read-write set error")
	}

	result := make([]prsRwSet, len(txReadWriteSet.NsRwset))
	for i, rwSet := range txReadWriteSet.NsRwset {
		if err := proto.Unmarshal(rwSet.Rwset, &result[i].kvRWSet); err != nil {
			return nil, errors.Wrap(err, "unmarshal kv read-write set error")
		}
	}
	return result, err
}
