package db

import (
	"log"
	"math/big"
	"time"

	"github.com/memoio/backend/contract"
)

func NewCron() {
	db, err := OpenDataBase()
	if err != nil {
		log.Println("open database error")
		return
	}
	defer db.Close()

	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		rows, err := db.Query(`SELECT p.id, p.hash, p.size, a.address
				FROM pkginfo p
				JOIN address a ON p.address_id = a.id
				WHERE p.is_updated = false
			`)
		if err != nil {
			log.Println(err)
			return
		}

		var updateid []int64

		for rows.Next() {
			var (
				id      int64
				address string
				hash    string
				size    int64
			)
			if err := rows.Scan(&id, &hash, &size, &address); err != nil {
				log.Printf("Error scanning pkginfo: %v", err)
				continue
			}

			if contract.StoreOrderPkg(address, hash, big.NewInt(size)) {
				log.Println(id)
				updateid = append(updateid, id)
			}
		}

		rows.Close()
		log.Println("updateid ", updateid)

		for _, id := range updateid {
			log.Println(id)
			_, err = db.Exec(`
                    UPDATE pkginfo
                    SET is_updated = true, time = ?
                    WHERE id = ?`, time.Now(), id)
			if err != nil {
				log.Printf("Error updating pkginfo: %v", err)
				return
			}
		}
	}
}
