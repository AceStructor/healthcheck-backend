package db

import (
    "fmt"
    "log"
    
    "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func WriteResult(res Result, configId uint, WarningLog *log.Logger, InfoLog *log.Logger) error {
	InfoLog.Printf("Writing Result for congig %v \n", configId)
	writeableResult := Result{ConfigID: configId, Status: res.Status, Text: res.Text, ResponseTime: res.ResponseTime}
	
    if err := DB.Create(&writeableResult).Error; err != nil {
        return fmt.Errorf("failed to create result entry: %w", err)
    }
    
    return nil
}

func ReadResult(query string, limit int32, offset int32, orderBy string, orderDirection string, WarningLog *log.Logger, InfoLog *log.Logger) ([]JoinedResult, error) {
	InfoLog.Println("Reading Results...")
	var ress []JoinedResult
	
	InfoLog.Println("Query: %v, Limit: %v, Offset: %v, Order By: %v, Order Direction: %v", query, limit, offset, orderBy, orderDirection)
	if err := DB.Table("results")
			    .Select("results.name, results.type, results.address, configs.status, configs.text")
			    .Joins("INNER JOIN configs ON results.config_id = configs.id")
			    .Where("configs.name ILIKE ? OR configs.type ILIKE ? OR configs.address ILIKE ? OR results.text ILIKE ?", query, query, query, query)
			    .Limit(int(limit))
			    .Offset(int(offset))
			    .Order(orderBy + " " + orderDirection)
			    .Scan(&ress); err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	
	return ress, nil
}
