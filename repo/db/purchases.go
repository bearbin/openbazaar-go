package db

import (
	"database/sql"
	"encoding/json"
	"github.com/OpenBazaar/jsonpb"
	"github.com/OpenBazaar/openbazaar-go/pb"
	"github.com/OpenBazaar/spvwallet"
	btc "github.com/btcsuite/btcutil"
	"strings"
	"sync"
)

type PurchasesDB struct {
	db   *sql.DB
	lock *sync.Mutex
}

func (p *PurchasesDB) Put(orderID string, contract pb.RicardianContract, state pb.OrderState, read bool) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	readInt := 0
	if read {
		readInt = 1
	}
	m := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "    ",
		OrigName:     false,
	}
	out, err := m.MarshalToString(&contract)

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into purchases(orderID, contract, state, read, date, total, thumbnail, vendorID, vendorBlockchainID, title, shippingName, shippingAddress, paymentAddr, funded, transactions) values(?,?,?,?,?,?,?,?,?,?,?,?,?,(select funded from purchases where orderID="` + orderID + `"),(select transactions from purchases where orderID="` + orderID + `"))`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		return err
	}
	blockchainID := contract.VendorListings[0].VendorID.BlockchainID
	shippingName := ""
	shippingAddress := ""
	if contract.BuyerOrder.Shipping != nil {
		shippingName = contract.BuyerOrder.Shipping.ShipTo
		shippingAddress = contract.BuyerOrder.Shipping.Address
	}
	var paymentAddr string
	if contract.BuyerOrder.Payment.Method == pb.Order_Payment_DIRECT || contract.BuyerOrder.Payment.Method == pb.Order_Payment_MODERATED {
		paymentAddr = contract.BuyerOrder.Payment.Address
	} else if contract.BuyerOrder.Payment.Method == pb.Order_Payment_ADDRESS_REQUEST {
		paymentAddr = contract.VendorOrderConfirmation.PaymentAddress
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		orderID,
		out,
		int(state),
		readInt,
		int(contract.BuyerOrder.Timestamp.Seconds),
		int(contract.BuyerOrder.Payment.Amount),
		contract.VendorListings[0].Item.Images[0].Tiny,
		contract.VendorListings[0].VendorID.Guid,
		blockchainID,
		strings.ToLower(contract.VendorListings[0].Item.Title),
		shippingName,
		shippingAddress,
		paymentAddr,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (p *PurchasesDB) MarkAsRead(orderID string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, err := p.db.Exec("update purchases set read=? where orderID=?", 1, orderID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PurchasesDB) UpdateFunding(orderId string, funded bool, records []*spvwallet.TransactionRecord) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	fundedInt := 0
	if funded {
		fundedInt = 1
	}
	serializedTransactions, err := json.Marshal(records)
	if err != nil {
		return err
	}
	_, err = p.db.Exec("update purchases set funded=?, transactions=? where orderID=?", fundedInt, string(serializedTransactions), orderId)
	if err != nil {
		return err
	}
	return nil
}

func (p *PurchasesDB) Delete(orderID string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, err := p.db.Exec("delete from purchases where orderID=?", orderID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PurchasesDB) GetAll() ([]string, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	stm := "select orderID from purchases"
	rows, err := p.db.Query(stm)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var ret []string
	for rows.Next() {
		var orderID string
		if err := rows.Scan(&orderID); err != nil {
			return ret, err
		}
		ret = append(ret, orderID)
	}
	return ret, nil
}

func (p *PurchasesDB) GetByPaymentAddress(addr btc.Address) (*pb.RicardianContract, pb.OrderState, bool, []*spvwallet.TransactionRecord, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	stmt, err := p.db.Prepare("select contract, state, funded, transactions from purchases where paymentAddr=?")
	defer stmt.Close()
	var contract []byte
	var stateInt int
	var fundedInt *int
	var serializedTransactions []byte
	err = stmt.QueryRow(addr.EncodeAddress()).Scan(&contract, &stateInt, &fundedInt, &serializedTransactions)
	if err != nil {
		return nil, pb.OrderState(0), false, nil, err
	}
	rc := new(pb.RicardianContract)
	err = jsonpb.UnmarshalString(string(contract), rc)
	if err != nil {
		return nil, pb.OrderState(0), false, nil, err
	}
	funded := false
	if fundedInt != nil && *fundedInt == 1 {
		funded = true
	}
	var records []*spvwallet.TransactionRecord
	json.Unmarshal(serializedTransactions, &records)
	return rc, pb.OrderState(stateInt), funded, records, nil
}

func (p *PurchasesDB) GetByOrderId(orderId string) (*pb.RicardianContract, pb.OrderState, bool, []*spvwallet.TransactionRecord, bool, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	stmt, err := p.db.Prepare("select contract, state, funded, transactions, read from purchases where orderID=?")
	defer stmt.Close()
	var contract []byte
	var stateInt int
	var fundedInt *int
	var readInt *int
	var serializedTransactions []byte
	err = stmt.QueryRow(orderId).Scan(&contract, &stateInt, &fundedInt, &serializedTransactions, &readInt)
	if err != nil {
		return nil, pb.OrderState(0), false, nil, false, err
	}
	rc := new(pb.RicardianContract)
	err = jsonpb.UnmarshalString(string(contract), rc)
	if err != nil {
		return nil, pb.OrderState(0), false, nil, false, err
	}
	funded := false
	if fundedInt != nil && *fundedInt == 1 {
		funded = true
	}
	read := false
	if readInt != nil && *readInt == 1 {
		read = true
	}
	var records []*spvwallet.TransactionRecord
	json.Unmarshal(serializedTransactions, &records)
	return rc, pb.OrderState(stateInt), funded, records, read, nil
}
