// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-ops/service/internal/dao/internal"
)

// internalVmDao is internal type for wrapping internal DAO implements.
type internalVmDao = *internal.VmDao

// vmDao is the data access object for table vm.
// You can define custom methods on it to extend its functionality as you wish.
type vmDao struct {
	internalVmDao
}

var (
	// Vm is globally public accessible object for table vm operations.
	Vm = vmDao{
		internal.NewVmDao(),
	}
)

// Fill with you ideas below.
