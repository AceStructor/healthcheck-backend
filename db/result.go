package db

import (
    "github.com/go-openapi/errors"
)

func WriteResult(result Result) error {
    if err := DB.Create(&result).Error; err != nil {
        return errors.New(500, "failed to create result entry: %v", err)
    }
    
    return nil
}

func ReadResult() ([]Result, error) {

}