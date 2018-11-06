package tedd

import (
	"context"
	"errors"
	"time"

	"github.com/chain/txvm/protocol/bc"
)

type testReserver struct {
	utxos []UTXO
}

func (r *testReserver) Reserve(_ context.Context, amount int64, assetID bc.Hash, exp time.Time) (Reservation, error) {
	res := &testReservation{reserver: r}
	for amount > 0 {
		if len(r.utxos) == 0 {
			return nil, errors.New("insufficient funds")
		}
		amount -= r.utxos[0].Amount()
		res.utxos = append(res.utxos, r.utxos[0])
		r.utxos = r.utxos[1:]
	}
	res.change = -amount
	return res, nil
}

type testReservation struct {
	reserver *testReserver
	utxos    []UTXO
	change   int64
}

func (r *testReservation) UTXOs() []UTXO {
	return r.utxos
}

func (r *testReservation) Change() int64 {
	return r.change
}

func (r *testReservation) Cancel(context.Context) error {
	r.reserver.utxos = append(r.reserver.utxos, r.utxos...)
	return nil
}

type testUTXO struct {
	amount  int64
	assetID bc.Hash
	anchor  []byte
}

func (u *testUTXO) Amount() int64 {
	return u.amount
}

func (u *testUTXO) AssetID() bc.Hash {
	return u.assetID
}

func (u *testUTXO) Anchor() []byte {
	return u.anchor
}