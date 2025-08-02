package db

import (
    "github.com/go-openapi/errors"
)

func WriteResult(result Result, configId uint) error {
	writeableResult := Result{ConfigID: configId, Status: result.Status, Text: result.Text, ResponseTime: result.ResponseTime}
	
    if err := DB.Create(&writeableResult).Error; err != nil {
        return errors.New(500, "failed to create result entry: %v", err)
    }
    
    return nil
}

func ReadResult(query string, limit int32, offset int32, orderBy string, orderDirection string) ([]JoinedResult, error) {
	var results []JoinedResult
	
	if err := DB.Table("results")
			    .Select("results.name, results.type, results.address, configs.status, configs.text")
			    .Joins("INNER JOIN configs ON results.config_id = configs.id")
			    .Where("configs.name ILIKE ? OR configs.type ILIKE ? OR configs.address ILIKE ? OR results.text ILIKE ?", query, query, query, query)
			    .Limit(int(limit))
			    .Offset(int(offset))
			    .Order(orderBy + " " + orderDirection)
			    .Scan(&results); err != nil {
		return nil, errors.New(500, "failed to search: %v", err)
	}
	
	return results, nil
}
