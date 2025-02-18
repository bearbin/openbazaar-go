package core

import (
	peer "gx/ipfs/QmRBqJF7hb8ZSpRcMwUt8hNhydWcxGEhtk81HKq6oUwKvs/go-libp2p-peer"
	crypto "gx/ipfs/QmUWER4r4qMvaCnX5zREcfyiWN7cXN9g3a7fkRqNz8qWPP/go-libp2p-crypto"

	"github.com/OpenBazaar/openbazaar-go/pb"
	"github.com/btcsuite/btcd/btcec"
	"github.com/golang/protobuf/proto"
)

type noSigError struct{}
type invalidSigError struct{}
type matchKeyError struct{}

func (noSigError) Error() string {
	return "Contract does not contain a signature for the given section"
}

func (invalidSigError) Error() string {
	return "Invalid signature"
}

func (matchKeyError) Error() string {
	return "Signature does not match public key"
}

// verifyMessageSignature accepts message, public key butes, list of signatures,
// signature section to be looked up in this list, and GUID string. Returns an
// error, with special cases:
//     - noSigError (signature not present in the list);
//     - invalidSigError (signature is invalid);
//     - matchKeyError (GUID does not match public key).
// Example usage can be seen in verifySignaturesOnOrder() in 'core/order.go'.
func verifyMessageSignature(msg proto.Message, pk []byte, sigs []*pb.Signature, sigType pb.Signature_Section, guid string) error {
	ser, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	guidPubkey, err := crypto.UnmarshalPublicKey(pk)
	if err != nil {
		return err
	}
	sig, err := selectSignature(sigs, sigType)
	if err != nil {
		return err
	}
	valid, err := guidPubkey.Verify(ser, sig.SignatureBytes)
	if err != nil {
		return err
	}
	if !valid {
		return invalidSigError{}
	}
	pid, err := peer.IDB58Decode(guid)
	if err != nil {
		return err
	}
	if !pid.MatchesPublicKey(guidPubkey) {
		return matchKeyError{}
	}
	return nil
}

func verifyBitcoinSignature(pubkeyBytes, sigBytes []byte, guid string) error {
	bitcoinPubkey, err := btcec.ParsePubKey(pubkeyBytes, btcec.S256())
	if err != nil {
		return err
	}
	bitcoinSig, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		return err
	}
	if !bitcoinSig.Verify([]byte(guid), bitcoinPubkey) {
		return invalidSigError{}
	}
	return nil
}

func selectSignature(signatures []*pb.Signature, sigType pb.Signature_Section) (*pb.Signature, error) {
	var sig *pb.Signature
	var err error = noSigError{}
	for _, s := range signatures {
		if s.Section == sigType {
			sig, err = s, nil
		}
	}
	return sig, err
}
