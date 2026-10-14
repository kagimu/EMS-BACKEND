package repositories

import "gorm.io/gorm/clause"

// lockForUpdate returns the row-locking clause used to serialise concurrent
// decisions. MySQL/MariaDB translate it to SELECT ... FOR UPDATE.
func lockForUpdate() clause.Locking {
	return clause.Locking{Strength: "UPDATE"}
}
