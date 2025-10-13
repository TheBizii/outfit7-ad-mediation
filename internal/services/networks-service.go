package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/TheBizii/outfit7-ad-mediation/internal/db"
	"github.com/TheBizii/outfit7-ad-mediation/internal/models"
)

func UpsertPriorityList(countryCode string, adType string, req models.UpdateNetworksRequest) error {
	if len(req.Networks) == 0 {
		return fmt.Errorf("Request body is missing networks and their performance scores.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// begin the transaction
	tx, err := db.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// insert a new priority list or update the last_updated if it already exists
	const upsertListStmt = `
		INSERT INTO priority_lists (country_code, ad_type, last_updated)
		VALUES ($1, $2, NOW())
		ON CONFLICT (country_code, ad_type)
		DO UPDATE SET last_updated = EXCLUDED.last_updated
		RETURNING id;`

	var listId int
	if err = tx.QueryRowContext(ctx, upsertListStmt, countryCode, adType).Scan(&listId); err != nil {
		return err
	}

	// networks tied to this priority list also have to be upserted
	const upsertNetworkStmt = `
		INSERT INTO priority_networks (priority_list_id, network_name, score)
		VALUES ($1, $2, $3)
		ON CONFLICT (priority_list_id, network_name)
		DO UPDATE SET score = EXCLUDED.score;`

	var networkNames []string
	for _, network := range req.Networks {
		if _, err = tx.ExecContext(ctx, upsertNetworkStmt, listId, network.NetworkName, network.Score); err != nil {
			return err
		}

		networkNames = append(networkNames, network.NetworkName)
	}

	// delete stale networks (those that aren't present in this request)
	args := make([]interface{}, 0, len(networkNames)+1)
	args = append(args, listId)
	placeholders := make([]string, len(networkNames))
	for i := range networkNames {
		placeholders[i] = fmt.Sprintf("$%d", i+2) // $1 is listId, we therefore must start at $2
		args = append(args, networkNames[i])
	}
	deleteStaleNetworksStmt := fmt.Sprintf(`
		DELETE FROM priority_networks
		WHERE priority_list_id = $1 AND network_name NOT IN (%s);
	`, strings.Join(placeholders, ", "))
	if _, err = tx.ExecContext(ctx, deleteStaleNetworksStmt, args...); err != nil {
		return err
	}

	// nothing went wrong, commit the transaction
	return tx.Commit()
}
