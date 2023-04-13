package database

import (
	"log"
	"math/big"
	"time"

	"github.com/memoio/backend/contract"
)

type UpdateAddress struct {
	id      int64
	address string
	stype   uint8
	hash    string
	size    int64
}

func NewCron() {
	db, err := OpenDataBase()
	if err != nil {
		log.Println("open database error")
		return
	}
	defer db.Close()

	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		rows, err := db.Query(`SELECT p.id, p.stype, p.hash, p.size, a.address
				FROM pkginfo p
				JOIN address a ON p.address_id = a.id
				WHERE p.is_updated = false
			`)
		if err != nil {
			log.Println(err)
			return
		}

		var update []UpdateAddress

		for rows.Next() {
			u := UpdateAddress{}
			if err := rows.Scan(&u.id, &u.stype, &u.hash, &u.size, &u.address); err != nil {
				log.Printf("Error scanning pkginfo: %v", err)
				continue
			}

			if contract.StoreOrderPkg(u.address, u.hash, u.stype, big.NewInt(u.size)) {
				update = append(update, u)
			}
		}

		rows.Close()
		log.Println("update ", update)

		for _, up := range update {
			log.Println(up)
			_, err = db.Exec(`
                    UPDATE pkginfo
                    SET is_updated = true, time = ?
                    WHERE id = ?`, time.Now(), up.id)
			if err != nil {
				log.Printf("Error updating pkginfo: %v", err)
				return
			}
			si, err := contract.GetPkgSize(up.stype, up.address)
			if err != nil {
				return
			}
			log.Println("si", si)
			ai := Storage{
				Address:    up.address,
				Buysize:    si.Buysize,
				Free:       si.Free,
				Used:       si.Used,
				Files:      si.Files,
				UpdateTime: time.Now(),
			}

			err = ai.Update()
			if err != nil {
				log.Println(err)
				return
			}
		}

	}
}
