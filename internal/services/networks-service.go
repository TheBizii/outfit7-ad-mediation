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

func GetDashboardPriorityLists() ([]models.PriorityListSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// select all priority lists and their ad networks and scores
	const getPriorityListsStmt = `
		SELECT list.country_code, list.ad_type, list.last_updated, network.network_name, network.score
		FROM priority_lists list
		LEFT JOIN priority_networks network ON list.id = network.priority_list_id
		ORDER BY list.country_code ASC, list.ad_type ASC, network.score DESC;`

	rows, err := db.Conn.QueryContext(ctx, getPriorityListsStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// used to group results together by (country, adType) pair
	type key struct {
		countryCode string
		adType      string
	}
	groups := make(map[key]*models.PriorityListSummary)

	// process all rows and group
	for rows.Next() {
		var countryCode, adType, networkName string
		var lastUpdated sql.NullTime
		var score sql.NullFloat64

		if err := rows.Scan(&countryCode, &adType, &lastUpdated, &networkName, &score); err != nil {
			return nil, err
		}

		k := key{countryCode: countryCode, adType: adType}
		group, exists := groups[k]
		if !exists {
			group = &models.PriorityListSummary{
				CountryCode: countryCode,
				AdType:      adType,
				LastUpdated: lastUpdated.Time,
				Networks:    []models.NetworkScore{},
			}
			groups[k] = group
		}

		networkScore := models.NetworkScore{
			NetworkName: networkName,
			Score:       float32(score.Float64),
		}
		group.Networks = append(group.Networks, networkScore)
	}

	groupsSlice := make([]models.PriorityListSummary, 0, len(groups))
	for _, group := range groups {
		if group != nil {
			groupsSlice = append(groupsSlice, *group)
		}
	}

	return groupsSlice, nil
}

func GetAdNetworks(req models.GetNetworksRequest) ([]string, error) {
	if req.CountryCode == "" || req.AdType == "" {
		return nil, fmt.Errorf("You must supply the countryCode and adType parameters.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// to keep the code clean and simple, I apply contextual filtering separately
	const getNetworksStmt = `
		SELECT network.network_name
		FROM priority_networks network
		JOIN priority_lists list ON list.id = network.priority_list_id
		WHERE list.country_code = $1 AND list.ad_type = $2
		ORDER BY network.score DESC;`

	rows, err := db.Conn.QueryContext(ctx, getNetworksStmt, req.CountryCode, req.AdType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process rows returned by the SELECT statement
	var networks []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		networks = append(networks, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// apply contextual filtering
	filtered := applyContextFilters(networks, req)
	return filtered, nil
}

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

// helper functions
func applyContextFilters(networks []string, req models.GetNetworksRequest) []string {
	filtered := make([]string, 0)
	osMajorVersion := extractMajorVersion(req.OSVersion)
	hasAdMob := false
	hasAdMobOptOut := false

	for _, network := range networks {
		if network == "AdMob" {
			if strings.EqualFold(req.Platform, "android") && osMajorVersion == "9" {
				// AdMob does not work on Android OS if major version is 9
				continue
			}
			hasAdMob = true
		} else if network == "AdMob-OptOut" {
			// AdMob-OptOut will be included if it was registered as a valid ad network for this country and
			// AdMob was present in the original list of ad networks
			hasAdMobOptOut = true
			continue
		} else if network == "Facebook" {
			if strings.EqualFold(req.CountryCode, "CN") {
				// Facebook should not be served in CN
				continue
			}
		}

		filtered = append(filtered, network)
	}

	// AdMob-OptOut should be present in the list only if there is no AdMob in the list
	if !hasAdMob && hasAdMobOptOut {
		filtered = append(filtered, "AdMob-OptOut")
	}

	return filtered
}

func extractMajorVersion(version string) string {
	if i := strings.Index(version, "."); i != -1 {
		return version[:i]
	}
	return version
}
