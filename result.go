package nebulaorm

import (
	"fmt"
	"github.com/haysons/nebulaorm/internal/utils"
	"github.com/haysons/nebulaorm/resolver"
	"github.com/haysons/nebulaorm/statement"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"reflect"
)

// NGQL get this generated statement does not actually execute the statement
func (db *DB) NGQL() (string, error) {
	tx := db.getInstance()
	return tx.Statement.NGQL()
}

// RawResult exec the statement and return the result of nebula-go directly
func (db *DB) RawResult() (*nebula.ResultSet, error) {
	nGQL, err := db.NGQL()
	if err != nil {
		return nil, err
	}
	return db.sessionPool.Execute(nGQL)
}

// Exec the statement, but don't care about the result as long as it is used for insert, update, delete operations
func (db *DB) Exec() error {
	nGQL, err := db.NGQL()
	if err != nil {
		return err
	}
	res, err := db.sessionPool.Execute(nGQL)
	if err != nil {
		return err
	}
	if !res.IsSucceed() {
		return fmt.Errorf("nebulaorm: result is not succeed, err code: %d, msg: %s", res.GetErrorCode(), res.GetErrorMsg())
	}
	return nil
}

// Find exec the statement and assign the returned result to the dest variable
func (db *DB) Find(dest interface{}) error {
	rawRes, err := db.RawResult()
	if err != nil {
		return err
	}
	return Scan(rawRes, dest)
}

// FindCol parse one column of the result, it is used to easily get the value of a field
func (db *DB) FindCol(col string, dest interface{}) error {
	rawRes, err := db.RawResult()
	if err != nil {
		return err
	}
	return Pluck(rawRes, col, dest)
}

// Take get a single test result, if no limit is specified, limit 1 will be added automatically,
// if the final return value is empty, will return  ErrRecordNotFound
func (db *DB) Take(dest interface{}) error {
	tx := db.getInstance()
	lastPart := tx.Statement.LastPart()
	if lastPart.GetType() != statement.PartTypeLimit {
		tx.Statement.Limit(1)
	}
	nGQL, err := tx.Statement.NGQL()
	if err != nil {
		return err
	}
	rawRes, err := db.sessionPool.Execute(nGQL)
	if err != nil {
		return err
	}
	return scan(rawRes, dest, true)
}

// TakeCol parse one column of the result, it is used to easily get the value of a field
// if the final return value is empty, will return  ErrRecordNotFound
func (db *DB) TakeCol(col string, dest interface{}) error {
	tx := db.getInstance()
	lastPart := tx.Statement.LastPart()
	if lastPart.GetType() != statement.PartTypeLimit {
		tx.Statement.Limit(1)
	}
	nGQL, err := tx.Statement.NGQL()
	if err != nil {
		return err
	}
	rawRes, err := db.sessionPool.Execute(nGQL)
	if err != nil {
		return err
	}
	return pluck(rawRes, col, dest, true)
}

// Scan assign the results to the target variable
func Scan(rawRes *nebula.ResultSet, dest interface{}) error {
	return scan(rawRes, dest, false)
}

func scan(rawRes *nebula.ResultSet, dest interface{}, raiseNotFound bool) error {
	if !rawRes.IsSucceed() {
		return fmt.Errorf("nebulaorm: result is not succeed, err code: %d, msg: %s", rawRes.GetErrorCode(), rawRes.GetErrorMsg())
	}
	if rawRes.GetRowSize() == 0 {
		if raiseNotFound {
			return ErrRecordNotFound
		} else {
			return nil
		}
	}
	switch v := dest.(type) {
	case *map[string]interface{}:
		if *v == nil {
			*v = make(map[string]interface{})
		}
		record, err := rawRes.GetRowValuesByIndex(0)
		if err != nil {
			return err
		}
		return scanIntoMap(record, rawRes.GetColNames(), *v)
	case map[string]interface{}:
		record, err := rawRes.GetRowValuesByIndex(0)
		if err != nil {
			return err
		}
		return scanIntoMap(record, rawRes.GetColNames(), v)
	case *[]map[string]interface{}:
		for i := 0; i < rawRes.GetRowSize(); i++ {
			record, err := rawRes.GetRowValuesByIndex(i)
			if err != nil {
				return err
			}
			value := make(map[string]interface{}, len(rawRes.GetColNames()))
			if err = scanIntoMap(record, rawRes.GetColNames(), value); err != nil {
				return err
			}
			*v = append(*v, value)
		}
		return nil
	default:
		destValue := reflect.ValueOf(dest)
		if destValue.Kind() != reflect.Ptr {
			return fmt.Errorf("nebulaorm: %w, scan dest should be pointer to struct, slice or array", ErrInvalidValue)
		}
		destValue = utils.PtrValue(destValue)
		if !destValue.IsValid() {
			return fmt.Errorf("nebulaorm: %w, scan dest should be pointer to struct, slice or array", ErrInvalidValue)
		}
		rv := resolver.NewResolver()
		switch destValue.Kind() {
		case reflect.Slice, reflect.Array:
			return utils.SliceSetElem(destValue, rawRes.GetRowSize(), func(i int, elem reflect.Value) (bool, error) {
				if i >= rawRes.GetRowSize() {
					return false, nil
				}
				record, _ := rawRes.GetRowValuesByIndex(i)
				if err := rv.ScanRecord(record, rawRes.GetColNames(), elem); err != nil {
					return false, err
				}
				return true, nil
			})
		case reflect.Struct:
			colNames := rawRes.GetColNames()
			record, err := rawRes.GetRowValuesByIndex(0)
			if err != nil {
				return err
			}
			return rv.ScanRecord(record, colNames, destValue)
		default:
			return fmt.Errorf("nebulaorm: %w, scan dest should be pointer to struct, slice or array", ErrInvalidValue)
		}
	}
}

// scanIntoMap scan a row into map
func scanIntoMap(record *nebula.Record, colNames []string, dest map[string]interface{}) error {
	for _, colName := range colNames {
		colValue, err := record.GetValueByColName(colName)
		if err != nil {
			return err
		}
		dest[colName], err = resolver.GetValueIface(colValue)
		if err != nil {
			return err
		}
	}
	return nil
}

// Pluck assign one of the fields of the return value into dest
func Pluck(rawRes *nebula.ResultSet, col string, dest interface{}) error {
	return pluck(rawRes, col, dest, false)
}

func pluck(rawRes *nebula.ResultSet, col string, dest interface{}, raiseNotFound bool) error {
	if !rawRes.IsSucceed() {
		return fmt.Errorf("nebulaorm: result is not succeed, err code: %d, msg: %s", rawRes.GetErrorCode(), rawRes.GetErrorMsg())
	}
	if rawRes.GetRowSize() == 0 {
		if raiseNotFound {
			return ErrRecordNotFound
		} else {
			return nil
		}
	}
	values, err := rawRes.GetValuesByColName(col)
	if err != nil {
		return fmt.Errorf("nebulaorm: get values by col name failed: %w", err)
	}
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("nebulaorm: %w, scan dest should be pointer to struct, slice or array", ErrInvalidValue)
	}
	destValue = utils.PtrValue(destValue)
	if !destValue.IsValid() {
		return fmt.Errorf("nebulaorm: %w, scan dest should be pointer to struct, slice or array", ErrInvalidValue)
	}
	rv := resolver.NewResolver()
	switch destValue.Kind() {
	case reflect.Slice, reflect.Array:
		return utils.SliceSetElem(destValue, rawRes.GetRowSize(), func(i int, elem reflect.Value) (bool, error) {
			if i >= len(values) {
				return false, nil
			}
			if err := rv.ScanValue(values[i], elem); err != nil {
				return false, err
			}
			return true, nil
		})
	default:
		return rv.ScanValue(values[0], destValue)
	}
}
