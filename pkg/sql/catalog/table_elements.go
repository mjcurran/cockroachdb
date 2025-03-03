// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package catalog

import (
	"fmt"
	"time"

	"github.com/cockroachdb/cockroach/pkg/geo/geoindex"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/iterutil"
)

// TableElementMaybeMutation is an interface used as a subtype for the various
// table descriptor elements which may be present in a mutation.
type TableElementMaybeMutation interface {
	// IsMutation returns true iff this table element is in a mutation.
	IsMutation() bool

	// IsRollback returns true iff the table element is in a rollback mutation.
	IsRollback() bool

	// MutationID returns the table element's mutationID if applicable,
	// descpb.InvalidMutationID otherwise.
	MutationID() descpb.MutationID

	// WriteAndDeleteOnly returns true iff the table element is in a mutation in
	// the delete-and-write-only state.
	WriteAndDeleteOnly() bool

	// DeleteOnly returns true iff the table element is in a mutation in the
	// delete-only state.
	DeleteOnly() bool

	// Backfilling returns true iff the table element is in a
	// mutation in the backfilling state.
	Backfilling() bool

	// Mergin returns true iff the table element is in a
	// mutation in the merging state.
	Merging() bool

	// Adding returns true iff the table element is in an add mutation.
	Adding() bool

	// Dropped returns true iff the table element is in a drop mutation.
	Dropped() bool
}

// ConstraintProvider is an interface for something which might unwrap
// a constraint.
type ConstraintProvider interface {

	// AsCheck returns the corresponding CheckConstraint if there is one,
	// nil otherwise.
	AsCheck() CheckConstraint

	// AsForeignKey returns the corresponding ForeignKeyConstraint if
	// there is one, nil otherwise.
	AsForeignKey() ForeignKeyConstraint

	// AsUniqueWithIndex returns the corresponding UniqueWithIndexConstraint if
	// there is one, nil otherwise.
	AsUniqueWithIndex() UniqueWithIndexConstraint

	// AsUniqueWithoutIndex returns the corresponding
	// UniqueWithoutIndexConstraint if there is one, nil otherwise.
	AsUniqueWithoutIndex() UniqueWithoutIndexConstraint
}

// Mutation is an interface around a table descriptor mutation.
type Mutation interface {
	TableElementMaybeMutation
	ConstraintProvider

	// AsColumn returns the corresponding Column if the mutation is on a column,
	// nil otherwise.
	AsColumn() Column

	// AsIndex returns the corresponding Index if the mutation is on an index,
	// nil otherwise.
	AsIndex() Index

	// AsConstraintWithoutIndex returns the corresponding WithoutIndexConstraint
	// if the mutation is on a check constraint or on a foreign key constraint or
	// on a non-index-backed unique constraint, nil otherwise.
	AsConstraintWithoutIndex() WithoutIndexConstraint

	// AsPrimaryKeySwap returns the corresponding PrimaryKeySwap if the mutation
	// is a primary key swap, nil otherwise.
	AsPrimaryKeySwap() PrimaryKeySwap

	// AsComputedColumnSwap returns the corresponding ComputedColumnSwap if the
	// mutation is a computed column swap, nil otherwise.
	AsComputedColumnSwap() ComputedColumnSwap

	// AsMaterializedViewRefresh returns the corresponding MaterializedViewRefresh
	// if the mutation is a materialized view refresh, nil otherwise.
	AsMaterializedViewRefresh() MaterializedViewRefresh

	// AsModifyRowLevelTTL returns the corresponding ModifyRowLevelTTL
	// if the mutation is a row-level TTL alter, nil otherwise.
	AsModifyRowLevelTTL() ModifyRowLevelTTL

	// NOTE: When adding new types of mutations to this interface, be sure to
	// audit the code which unpacks and introspects mutations to be sure to add
	// cases for the new type.

	// MutationOrdinal returns the ordinal of the mutation in the underlying table
	// descriptor's Mutations slice.
	MutationOrdinal() int
}

// Index is an interface around the index descriptor types.
type Index interface {
	TableElementMaybeMutation
	ConstraintProvider

	// IndexDesc returns the underlying protobuf descriptor.
	// Ideally, this method should be called as rarely as possible.
	IndexDesc() *descpb.IndexDescriptor

	// IndexDescDeepCopy returns a deep copy of the underlying proto.
	IndexDescDeepCopy() descpb.IndexDescriptor

	// Ordinal returns the ordinal of the index in its parent table descriptor.
	//
	// The ordinal of an index in a `tableDesc descpb.TableDescriptor` is
	// defined as follows:
	// - 0 is the ordinal of the primary index,
	// - [1:1+len(tableDesc.Indexes)] is the range of public non-primary indexes,
	// - [1+len(tableDesc.Indexes):] is the range of non-public indexes.
	//
	// In terms of a `table catalog.TableDescriptor` interface, it is defined
	// as the catalog.Index object's position in the table.AllIndexes() slice.
	Ordinal() int

	// Primary returns true iff the index is the primary index for the table
	// descriptor.
	Primary() bool

	// Public returns true iff the index is active, i.e. readable, in the table
	// descriptor.
	Public() bool

	// The remaining methods operate on the underlying descpb.IndexDescriptor object.

	GetID() descpb.IndexID
	GetConstraintID() descpb.ConstraintID
	GetName() string
	IsPartial() bool
	IsUnique() bool
	IsDisabled() bool
	IsSharded() bool
	IsNotVisible() bool
	IsCreatedExplicitly() bool
	GetPredicate() string
	GetType() descpb.IndexDescriptor_Type
	GetGeoConfig() geoindex.Config
	GetVersion() descpb.IndexDescriptorVersion
	GetEncodingType() descpb.IndexDescriptorEncodingType

	GetSharded() catpb.ShardedDescriptor
	GetShardColumnName() string

	// IsValidOriginIndex returns whether the index can serve as an origin index
	// for a foreign key constraint.
	IsValidOriginIndex(fk ForeignKeyConstraint) bool

	GetPartitioning() Partitioning
	PartitioningColumnCount() int
	ImplicitPartitioningColumnCount() int

	ExplicitColumnStartIdx() int

	NumKeyColumns() int
	GetKeyColumnID(columnOrdinal int) descpb.ColumnID
	GetKeyColumnName(columnOrdinal int) string
	GetKeyColumnDirection(columnOrdinal int) catpb.IndexColumn_Direction

	CollectKeyColumnIDs() TableColSet
	CollectKeySuffixColumnIDs() TableColSet
	CollectPrimaryStoredColumnIDs() TableColSet
	CollectSecondaryStoredColumnIDs() TableColSet
	CollectCompositeColumnIDs() TableColSet

	// InvertedColumnID returns the ColumnID of the inverted column of the
	// inverted index.
	//
	// Panics if the index is not inverted.
	InvertedColumnID() descpb.ColumnID

	// InvertedColumnName returns the name of the inverted column of the inverted
	// index.
	//
	// Panics if the index is not inverted.
	InvertedColumnName() string

	// InvertedColumnKeyType returns the type of the data element that is encoded
	// as the inverted index key. This is currently always EncodedKey.
	//
	// Panics if the index is not inverted.
	InvertedColumnKeyType() *types.T

	// InvertedColumnKind returns the kind of the inverted column of the inverted
	// index.
	InvertedColumnKind() catpb.InvertedIndexColumnKind

	NumPrimaryStoredColumns() int
	NumSecondaryStoredColumns() int
	GetStoredColumnID(storedColumnOrdinal int) descpb.ColumnID
	GetStoredColumnName(storedColumnOrdinal int) string
	HasOldStoredColumns() bool

	NumKeySuffixColumns() int
	GetKeySuffixColumnID(extraColumnOrdinal int) descpb.ColumnID

	NumCompositeColumns() int
	GetCompositeColumnID(compositeColumnOrdinal int) descpb.ColumnID
	UseDeletePreservingEncoding() bool
	// ForcePut forces all writes to use Put rather than CPut or InitPut.
	//
	// Users of this options should take great care as it
	// effectively mean unique constraints are not respected.
	//
	// Currently (2022-07-15) there are three users:
	//   * delete preserving indexes
	//   * merging indexes
	//   * dropping primary indexes
	//
	// Delete preserving encoding indexes are used only as a log of
	// index writes during backfill, thus we can blindly put values into
	// them.
	//
	// New indexes may miss updates during the backfilling process
	// that would lead to CPut failures until the missed updates
	// are merged into the index. Uniqueness for such indexes is
	// checked by the schema changer before they are brought back
	// online.
	//
	// In the case of dropping primary indexes, we always ensure that
	// there's a replacement primary index which has become public.
	// The reason we must not use cput is that the new primary index
	// may not store all the columns stored in this index.
	ForcePut() bool

	// CreatedAt is an approximate timestamp at which the index was created.
	// It is derived from the statement time at which the relevant statement
	// was issued.
	CreatedAt() time.Time

	// IsTemporaryIndexForBackfill returns true iff the index is
	// an index being used as the temporary index being used by an
	// in-progress index backfill.
	IsTemporaryIndexForBackfill() bool
}

// Column is an interface around the column descriptor types.
type Column interface {
	TableElementMaybeMutation

	// ColumnDesc returns the underlying protobuf descriptor.
	// Ideally, this method should be called as rarely as possible.
	ColumnDesc() *descpb.ColumnDescriptor

	// ColumnDescDeepCopy returns a deep copy of the underlying proto.
	ColumnDescDeepCopy() descpb.ColumnDescriptor

	// DeepCopy returns a deep copy of the receiver.
	DeepCopy() Column

	// Ordinal returns the ordinal of the column in its parent table descriptor.
	//
	// The ordinal of a column in a `tableDesc descpb.TableDescriptor` is
	// defined as follows:
	// - [0:len(tableDesc.Columns)] is the range of public columns,
	// - [len(tableDesc.Columns):] is the range of non-public columns.
	//
	// In terms of a `table catalog.TableDescriptor` interface, it is defined
	// as the catalog.Column object's position in the table.AllColumns() slice.
	Ordinal() int

	// Public returns true iff the column is active, i.e. readable, in the table
	// descriptor.
	Public() bool

	// GetID returns the column ID.
	GetID() descpb.ColumnID

	// GetName returns the column name as a string.
	GetName() string

	// ColName returns the column name as a tree.Name.
	ColName() tree.Name

	// HasType returns true iff the column type is set.
	HasType() bool

	// GetType returns the column type.
	GetType() *types.T

	// IsNullable returns true iff the column allows NULL values.
	IsNullable() bool

	// HasDefault returns true iff the column has a default expression set.
	HasDefault() bool

	// HasNullDefault returns true if the column has a default expression and
	// that expression is NULL.
	HasNullDefault() bool

	// GetDefaultExpr returns the column default expression if it exists,
	// empty string otherwise.
	GetDefaultExpr() string

	// HasOnUpdate returns true iff the column has an on update expression set.
	HasOnUpdate() bool

	// GetOnUpdateExpr returns the column on update expression if it exists,
	// empty string otherwise.
	GetOnUpdateExpr() string

	// IsComputed returns true iff the column is a computed column.
	IsComputed() bool

	// GetComputeExpr returns the column computed expression if it exists,
	// empty string otherwise.
	GetComputeExpr() string

	// IsHidden returns true iff the column is not visible.
	IsHidden() bool

	// IsInaccessible returns true iff the column is inaccessible.
	IsInaccessible() bool

	// IsExpressionIndexColumn returns true iff the column is an an inaccessible
	// virtual computed column that represents an expression in an expression
	// index.
	IsExpressionIndexColumn() bool

	// NumUsesSequences returns the number of sequences used by this column.
	NumUsesSequences() int

	// GetUsesSequenceID returns the ID of a sequence used by this column.
	GetUsesSequenceID(usesSequenceOrdinal int) descpb.ID

	// NumOwnsSequences returns the number of sequences owned by this column.
	NumOwnsSequences() int

	// GetOwnsSequenceID returns the ID of a sequence owned by this column.
	GetOwnsSequenceID(ownsSequenceOrdinal int) descpb.ID

	// IsVirtual returns true iff the column is a virtual column.
	IsVirtual() bool

	// CheckCanBeInboundFKRef returns whether the given column can be on the
	// referenced (target) side of a foreign key relation.
	CheckCanBeInboundFKRef() error

	// CheckCanBeOutboundFKRef returns whether the given column can be on the
	// referencing (origin) side of a foreign key relation.
	CheckCanBeOutboundFKRef() error

	// GetPGAttributeNum returns the PGAttributeNum of the column descriptor
	// if the PGAttributeNum is set (non-zero). Returns the ID of the
	// column descriptor if the PGAttributeNum is not set.
	GetPGAttributeNum() descpb.PGAttributeNum

	// IsSystemColumn returns true iff the column is a system column.
	IsSystemColumn() bool

	// IsGeneratedAsIdentity returns true iff the column is created
	// with GENERATED {ALWAYS | BY DEFAULT} AS IDENTITY syntax.
	IsGeneratedAsIdentity() bool

	// IsGeneratedAlwaysAsIdentity returns true iff the column is created
	// with GENERATED ALWAYS AS IDENTITY syntax.
	IsGeneratedAlwaysAsIdentity() bool

	// IsGeneratedByDefaultAsIdentity returns true iff the column is created
	// with GENERATED BY DEFAULT AS IDENTITY syntax.
	IsGeneratedByDefaultAsIdentity() bool

	// GetGeneratedAsIdentityType returns the type of how the column was
	// created as an IDENTITY column.
	// If the column is created with `GENERATED ALWAYS AS IDENTITY` syntax,
	// it will return descpb.GeneratedAsIdentityType_GENERATED_ALWAYS;
	// if the column is created with `GENERATED BY DEFAULT AS IDENTITY` syntax,
	// it will return descpb.GeneratedAsIdentityType_GENERATED_BY_DEFAULT;
	// otherwise, returns descpb.GeneratedAsIdentityType_NOT_IDENTITY_COLUMN.
	GetGeneratedAsIdentityType() catpb.GeneratedAsIdentityType

	// HasGeneratedAsIdentitySequenceOption returns true if there is a
	// customized sequence option when this column is created as a
	// `GENERATED AS IDENTITY` column.
	HasGeneratedAsIdentitySequenceOption() bool

	// GetGeneratedAsIdentitySequenceOptionStr returns the string representation
	// of the column's `GENERATED AS IDENTITY` sequence option if it exists, empty
	// string otherwise.
	GetGeneratedAsIdentitySequenceOptionStr() string

	// GetGeneratedAsIdentitySequenceOption returns the column's `GENERATED AS
	// IDENTITY` sequence option if it exists, and possible error.
	// If the column is not an identity column, return nil for both sequence option
	// and the error.
	// Note it doesn't return the sequence owner info.
	GetGeneratedAsIdentitySequenceOption(defaultIntSize int32) (*descpb.TableDescriptor_SequenceOpts, error)
}

// Constraint is an interface around a constraint.
type Constraint interface {
	TableElementMaybeMutation
	ConstraintProvider
	fmt.Stringer

	// GetConstraintID returns the ID for the constraint.
	GetConstraintID() descpb.ConstraintID

	// GetConstraintValidity returns the validity of this constraint.
	GetConstraintValidity() descpb.ConstraintValidity

	// IsEnforced returns true iff the constraint is enforced for all writes to
	// the parent table.
	IsEnforced() bool

	// GetName returns the name of this constraint update mutation.
	GetName() string

	// IsConstraintValidated returns true iff the constraint is enforced for
	// all writes to the table data and has also been validated on the table's
	// existing data prior to the addition of the constraint, if there was any.
	IsConstraintValidated() bool

	// IsConstraintUnvalidated returns true iff the constraint is enforced for
	// all writes to the table data but which is explicitly NOT to be validated
	// on any data in the table prior to the addition of the constraint.
	IsConstraintUnvalidated() bool
}

// UniqueConstraint is an interface for a unique constraint.
// These are either backed by an index, or not, see respectively
// UniqueWithIndexConstraint and UniqueWithoutIndexConstraint.
type UniqueConstraint interface {
	Constraint

	// IsValidReferencedUniqueConstraint returns whether the unique constraint can
	// serve as a referenced unique constraint for a foreign key constraint.
	IsValidReferencedUniqueConstraint(fk ForeignKeyConstraint) bool

	// NumKeyColumns returns the number of columns in this unique constraint.
	NumKeyColumns() int

	// GetKeyColumnID returns the ID of the column in the unique constraint at
	// ordinal `columnOrdinal`.
	GetKeyColumnID(columnOrdinal int) descpb.ColumnID

	// CollectKeyColumnIDs returns the columns in the unique constraint in a new
	// TableColSet.
	CollectKeyColumnIDs() TableColSet

	// IsPartial returns true iff this is a partial uniqueness constraint.
	IsPartial() bool

	// GetPredicate returns the partial predicate if there is one, "" otherwise.
	GetPredicate() string
}

// UniqueWithIndexConstraint is an interface around a unique constraint
// which backed by an index.
type UniqueWithIndexConstraint interface {
	UniqueConstraint
	Index
}

// WithoutIndexConstraint is the supertype of all constraint subtypes which are
// not backed by an index.
type WithoutIndexConstraint interface {
	Constraint
}

// CheckConstraint is an interface around a check constraint.
type CheckConstraint interface {
	WithoutIndexConstraint

	// CheckDesc returns the underlying descriptor protobuf.
	CheckDesc() *descpb.TableDescriptor_CheckConstraint

	// GetExpr returns the check expression as a string.
	GetExpr() string

	// NumReferencedColumns returns the number of column references in the check
	// expression. Note that a column may be referenced multiple times in an
	// expression; the number returned here is the number of references, not the
	// number of distinct columns.
	NumReferencedColumns() int

	// GetReferencedColumnID returns the ID of the column referenced in the check
	// expression at ordinal `columnOrdinal`.
	GetReferencedColumnID(columnOrdinal int) descpb.ColumnID

	// CollectReferencedColumnIDs returns the columns referenced in the check
	// constraint expression in a new TableColSet.
	CollectReferencedColumnIDs() TableColSet

	// IsNotNullColumnConstraint returns true iff this check constraint is a
	// NOT NULL on a column.
	IsNotNullColumnConstraint() bool

	// IsHashShardingConstraint returns true iff this check constraint is
	// associated with a hash-sharding column in this table.
	IsHashShardingConstraint() bool
}

// ForeignKeyConstraint is an interface around a check constraint.
type ForeignKeyConstraint interface {
	WithoutIndexConstraint

	// ForeignKeyDesc returns the underlying descriptor protobuf.
	ForeignKeyDesc() *descpb.ForeignKeyConstraint

	// GetOriginTableID returns the ID of the table at the origin of this foreign
	// key.
	GetOriginTableID() descpb.ID

	// NumOriginColumns returns the number of origin columns in the foreign key.
	NumOriginColumns() int

	// GetOriginColumnID returns the ID of the origin column in the foreign
	// key at ordinal `columnOrdinal`.
	GetOriginColumnID(columnOrdinal int) descpb.ColumnID

	// CollectOriginColumnIDs returns the origin columns in the foreign key
	// in a new TableColSet.
	CollectOriginColumnIDs() TableColSet

	// GetReferencedTableID returns the ID of the table referenced by this
	// foreign key.
	GetReferencedTableID() descpb.ID

	// NumReferencedColumns returns the number of columns referenced by this
	// foreign key.
	NumReferencedColumns() int

	// GetReferencedColumnID returns the ID of the column referenced by the
	// foreign key at ordinal `columnOrdinal`.
	GetReferencedColumnID(columnOrdinal int) descpb.ColumnID

	// CollectReferencedColumnIDs returns the columns referenced by the foreign
	// key in a new TableColSet.
	CollectReferencedColumnIDs() TableColSet

	// OnDelete returns the action to take ON DELETE.
	OnDelete() catpb.ForeignKeyAction

	// OnUpdate returns the action to take ON UPDATE.
	OnUpdate() catpb.ForeignKeyAction

	// Match returns the type of algorithm used to match composite keys.
	Match() descpb.ForeignKeyReference_Match
}

// UniqueWithoutIndexConstraint is an interface around a unique constraint
// which is not backed by an index.
type UniqueWithoutIndexConstraint interface {
	WithoutIndexConstraint
	UniqueConstraint

	// UniqueWithoutIndexDesc returns the underlying descriptor protobuf.
	UniqueWithoutIndexDesc() *descpb.UniqueWithoutIndexConstraint

	// ParentTableID returns the ID of the table this constraint applies to.
	ParentTableID() descpb.ID
}

// PrimaryKeySwap is an interface around a primary key swap mutation.
type PrimaryKeySwap interface {
	TableElementMaybeMutation

	// PrimaryKeySwapDesc returns the underlying protobuf descriptor.
	PrimaryKeySwapDesc() *descpb.PrimaryKeySwap

	// NumOldIndexes returns the number of old active indexes to swap out.
	NumOldIndexes() int

	// ForEachOldIndexIDs iterates through each of the old index IDs.
	// iterutil.StopIteration is supported.
	ForEachOldIndexIDs(fn func(id descpb.IndexID) error) error

	// NumNewIndexes returns the number of new active indexes to swap in.
	NumNewIndexes() int

	// ForEachNewIndexIDs iterates through each of the new index IDs.
	// iterutil.StopIteration is supported.
	ForEachNewIndexIDs(fn func(id descpb.IndexID) error) error

	// HasLocalityConfig returns true iff the locality config is swapped also.
	HasLocalityConfig() bool

	// LocalityConfigSwap returns the locality config swap, if there is one.
	LocalityConfigSwap() descpb.PrimaryKeySwap_LocalityConfigSwap
}

// ComputedColumnSwap is an interface around a computed column swap mutation.
type ComputedColumnSwap interface {
	TableElementMaybeMutation

	// ComputedColumnSwapDesc returns the underlying protobuf descriptor.
	ComputedColumnSwapDesc() *descpb.ComputedColumnSwap
}

// MaterializedViewRefresh is an interface around a materialized view refresh
// mutation.
type MaterializedViewRefresh interface {
	TableElementMaybeMutation

	// MaterializedViewRefreshDesc returns the underlying protobuf descriptor.
	MaterializedViewRefreshDesc() *descpb.MaterializedViewRefresh

	// ShouldBackfill returns true iff the query should be backfilled into the
	// indexes.
	ShouldBackfill() bool

	// AsOf returns the timestamp at which the query should be run.
	AsOf() hlc.Timestamp

	// ForEachIndexID iterates through each of the index IDs.
	// iterutil.StopIteration is supported.
	ForEachIndexID(func(id descpb.IndexID) error) error

	// TableWithNewIndexes returns a new TableDescriptor based on the old one
	// but with the refreshed indexes put in.
	TableWithNewIndexes(tbl TableDescriptor) TableDescriptor
}

// ModifyRowLevelTTL is an interface around a modify row level TTL mutation.
type ModifyRowLevelTTL interface {
	TableElementMaybeMutation

	// RowLevelTTL returns the row level TTL for the mutation.
	RowLevelTTL() *catpb.RowLevelTTL
}

// Partitioning is an interface around an index partitioning.
type Partitioning interface {

	// PartitioningDesc returns the underlying protobuf descriptor.
	PartitioningDesc() *catpb.PartitioningDescriptor

	// DeepCopy returns a deep copy of the receiver.
	DeepCopy() Partitioning

	// FindPartitionByName recursively searches the partitioning for a partition
	// whose name matches the input and returns it, or nil if no match is found.
	FindPartitionByName(name string) Partitioning

	// ForEachPartitionName applies fn on each of the partition names in this
	// partition and recursively in its subpartitions.
	// Supports iterutil.StopIteration.
	ForEachPartitionName(fn func(name string) error) error

	// ForEachList applies fn on each list element of the wrapped partitioning.
	// Supports iterutil.StopIteration.
	ForEachList(fn func(name string, values [][]byte, subPartitioning Partitioning) error) error

	// ForEachRange applies fn on each range element of the wrapped partitioning.
	// Supports iterutil.StopIteration.
	ForEachRange(fn func(name string, from, to []byte) error) error

	// NumColumns is how large of a prefix of the columns in an index are used in
	// the function mapping column values to partitions. If this is a
	// subpartition, this is offset to start from the end of the parent
	// partition's columns. If NumColumns is 0, then there is no partitioning.
	NumColumns() int

	// NumImplicitColumns specifies the number of columns that implicitly prefix a
	// given index. This occurs if a user specifies a PARTITION BY which is not a
	// prefix of the given index, in which case the ColumnIDs are added in front
	// of the index and this field denotes the number of columns added as a
	// prefix.
	// If NumImplicitColumns is 0, no implicit columns are defined for the index.
	NumImplicitColumns() int

	// NumLists returns the number of list elements in the underlying partitioning
	// descriptor.
	NumLists() int

	// NumRanges returns the number of range elements in the underlying
	// partitioning descriptor.
	NumRanges() int
}

func isIndexInSearchSet(desc TableDescriptor, opts IndexOpts, idx Index) bool {
	if !opts.NonPhysicalPrimaryIndex && idx.Primary() && !desc.IsPhysicalTable() {
		return false
	}
	if !opts.AddMutations && idx.Adding() {
		return false
	}
	if !opts.DropMutations && idx.Dropped() {
		return false
	}
	return true
}

// ForEachIndex runs f over each index in the table descriptor according to
// filter parameters in opts. Indexes are visited in their canonical order,
// see Index.Ordinal(). ForEachIndex supports iterutil.StopIteration().
func ForEachIndex(desc TableDescriptor, opts IndexOpts, f func(idx Index) error) error {
	for _, idx := range desc.AllIndexes() {
		if !isIndexInSearchSet(desc, opts, idx) {
			continue
		}
		if err := f(idx); err != nil {
			return iterutil.Map(err)
		}
	}
	return nil
}

func forEachIndex(slice []Index, f func(idx Index) error) error {
	for _, idx := range slice {
		if err := f(idx); err != nil {
			return iterutil.Map(err)
		}
	}
	return nil
}

// ForEachActiveIndex is like ForEachIndex over ActiveIndexes().
func ForEachActiveIndex(desc TableDescriptor, f func(idx Index) error) error {
	return forEachIndex(desc.ActiveIndexes(), f)
}

// ForEachNonDropIndex is like ForEachIndex over NonDropIndexes().
func ForEachNonDropIndex(desc TableDescriptor, f func(idx Index) error) error {
	return forEachIndex(desc.NonDropIndexes(), f)
}

// ForEachPartialIndex is like ForEachIndex over PartialIndexes().
func ForEachPartialIndex(desc TableDescriptor, f func(idx Index) error) error {
	return forEachIndex(desc.PartialIndexes(), f)
}

// ForEachNonPrimaryIndex is like ForEachIndex over
// NonPrimaryIndexes().
func ForEachNonPrimaryIndex(desc TableDescriptor, f func(idx Index) error) error {
	return forEachIndex(desc.NonPrimaryIndexes(), f)
}

// ForEachPublicNonPrimaryIndex is like ForEachIndex over
// PublicNonPrimaryIndexes().
func ForEachPublicNonPrimaryIndex(desc TableDescriptor, f func(idx Index) error) error {
	return forEachIndex(desc.PublicNonPrimaryIndexes(), f)
}

// ForEachWritableNonPrimaryIndex is like ForEachIndex over
// WritableNonPrimaryIndexes().
func ForEachWritableNonPrimaryIndex(desc TableDescriptor, f func(idx Index) error) error {
	return forEachIndex(desc.WritableNonPrimaryIndexes(), f)
}

// ForEachDeletableNonPrimaryIndex is like ForEachIndex over
// DeletableNonPrimaryIndexes().
func ForEachDeletableNonPrimaryIndex(desc TableDescriptor, f func(idx Index) error) error {
	return forEachIndex(desc.DeletableNonPrimaryIndexes(), f)
}

// ForEachDeleteOnlyNonPrimaryIndex is like ForEachIndex over
// DeleteOnlyNonPrimaryIndexes().
func ForEachDeleteOnlyNonPrimaryIndex(desc TableDescriptor, f func(idx Index) error) error {
	return forEachIndex(desc.DeleteOnlyNonPrimaryIndexes(), f)
}

// FindIndex returns the first index for which test returns true, nil otherwise,
// according to the parameters in opts just like ForEachIndex.
// Indexes are visited in their canonical order, see Index.Ordinal().
func FindIndex(desc TableDescriptor, opts IndexOpts, test func(idx Index) bool) Index {
	for _, idx := range desc.AllIndexes() {
		if !isIndexInSearchSet(desc, opts, idx) {
			continue
		}
		if test(idx) {
			return idx
		}
	}
	return nil
}

func findIndex(slice []Index, test func(idx Index) bool) Index {
	for _, idx := range slice {
		if test(idx) {
			return idx
		}
	}
	return nil
}

// FindActiveIndex returns the first index in ActiveIndex() for which test
// returns true.
func FindActiveIndex(desc TableDescriptor, test func(idx Index) bool) Index {
	return findIndex(desc.ActiveIndexes(), test)
}

// FindNonDropIndex returns the first index in NonDropIndex() for which test
// returns true.
func FindNonDropIndex(desc TableDescriptor, test func(idx Index) bool) Index {
	return findIndex(desc.NonDropIndexes(), test)
}

// FindPartialIndex returns the first index in PartialIndex() for which test
// returns true.
func FindPartialIndex(desc TableDescriptor, test func(idx Index) bool) Index {
	return findIndex(desc.PartialIndexes(), test)
}

// FindPublicNonPrimaryIndex returns the first index in PublicNonPrimaryIndex()
// for which test returns true.
func FindPublicNonPrimaryIndex(desc TableDescriptor, test func(idx Index) bool) Index {
	return findIndex(desc.PublicNonPrimaryIndexes(), test)
}

// FindWritableNonPrimaryIndex returns the first index in
// WritableNonPrimaryIndex() for which test returns true.
func FindWritableNonPrimaryIndex(desc TableDescriptor, test func(idx Index) bool) Index {
	return findIndex(desc.WritableNonPrimaryIndexes(), test)
}

// FindDeletableNonPrimaryIndex returns the first index in
// DeletableNonPrimaryIndex() for which test returns true.
func FindDeletableNonPrimaryIndex(desc TableDescriptor, test func(idx Index) bool) Index {
	return findIndex(desc.DeletableNonPrimaryIndexes(), test)
}

// FindNonPrimaryIndex returns the first index in
// NonPrimaryIndex() for which test returns true.
func FindNonPrimaryIndex(desc TableDescriptor, test func(idx Index) bool) Index {
	return findIndex(desc.NonPrimaryIndexes(), test)
}

// FindDeleteOnlyNonPrimaryIndex returns the first index in
// DeleteOnlyNonPrimaryIndex() for which test returns true.
func FindDeleteOnlyNonPrimaryIndex(desc TableDescriptor, test func(idx Index) bool) Index {
	return findIndex(desc.DeleteOnlyNonPrimaryIndexes(), test)
}

// FindCorrespondingTemporaryIndexByID finds the temporary index that
// corresponds to the currently mutated index identified by ID. It
// assumes that the temporary index for a given index ID exists
// directly after it in the mutations array.
//
// Callers should take care that AllocateIDs() has been called before
// using this function.
func FindCorrespondingTemporaryIndexByID(desc TableDescriptor, id descpb.IndexID) Index {
	mutations := desc.AllMutations()
	var ord int
	for _, m := range mutations {
		idx := m.AsIndex()
		if idx != nil && idx.IndexDesc().ID == id {
			// We want the mutation after this mutation
			// since the temporary index is added directly
			// after.
			ord = m.MutationOrdinal() + 1
		}
	}

	// A temporary index will never be found at index 0 since we
	// always add them _after_ the index they correspond to.
	if ord == 0 {
		return nil
	}

	if len(mutations) >= ord+1 {
		candidateMutation := mutations[ord]
		if idx := candidateMutation.AsIndex(); idx != nil {
			if idx.IsTemporaryIndexForBackfill() {
				return idx
			}
		}
	}
	return nil
}

// IsCorrespondingTemporaryIndex returns true iff idx is a temporary index
// created during a backfill and is the corresponding temporary index for
// otherIdx. It assumes that idx and otherIdx are both indexes from the same
// table.
func IsCorrespondingTemporaryIndex(idx Index, otherIdx Index) bool {
	return idx.IsTemporaryIndexForBackfill() && idx.Ordinal() == otherIdx.Ordinal()+1
}

// UserDefinedTypeColsHaveSameVersion returns whether one table descriptor's
// columns with user defined type metadata have the same versions of metadata
// as in the other descriptor. Note that this function is only valid on two
// descriptors representing the same table at the same version.
func UserDefinedTypeColsHaveSameVersion(desc TableDescriptor, otherDesc TableDescriptor) bool {
	otherCols := otherDesc.UserDefinedTypeColumns()
	for i, thisCol := range desc.UserDefinedTypeColumns() {
		this, other := thisCol.GetType(), otherCols[i].GetType()
		if this.TypeMeta.Version != other.TypeMeta.Version {
			return false
		}
	}
	return true
}

// UserDefinedTypeColsInFamilyHaveSameVersion returns whether one table descriptor's
// columns with user defined type metadata have the same versions of metadata
// as in the other descriptor, for all columns in the same family.
// Note that this function is only valid on two
// descriptors representing the same table at the same version.
func UserDefinedTypeColsInFamilyHaveSameVersion(
	desc TableDescriptor, otherDesc TableDescriptor, familyID descpb.FamilyID,
) (bool, error) {
	family, err := desc.FindFamilyByID(familyID)
	if err != nil {
		return false, err
	}

	familyCols := util.FastIntSet{}
	for _, colID := range family.ColumnIDs {
		familyCols.Add(int(colID))
	}

	otherCols := otherDesc.UserDefinedTypeColumns()
	for i, thisCol := range desc.UserDefinedTypeColumns() {
		this, other := thisCol.GetType(), otherCols[i].GetType()
		if familyCols.Contains(int(thisCol.GetID())) && this.TypeMeta.Version != other.TypeMeta.Version {
			return false, nil
		}
	}
	return true, nil
}

// ColumnIDToOrdinalMap returns a map from Column ID to the ordinal
// position of that column.
func ColumnIDToOrdinalMap(columns []Column) TableColMap {
	var m TableColMap
	for _, col := range columns {
		m.Set(col.GetID(), col.Ordinal())
	}
	return m
}

// ColumnTypes returns the types of the given columns
func ColumnTypes(columns []Column) []*types.T {
	return ColumnTypesWithInvertedCol(columns, nil /* invertedCol */)
}

// ColumnTypesWithInvertedCol returns the types of all given columns,
// If invertedCol is non-nil, substitutes the type of the inverted
// column instead of the column with the same ID.
func ColumnTypesWithInvertedCol(columns []Column, invertedCol Column) []*types.T {
	t := make([]*types.T, len(columns))
	for i, col := range columns {
		t[i] = col.GetType()
		if invertedCol != nil && col.GetID() == invertedCol.GetID() {
			t[i] = invertedCol.GetType()
		}
	}
	return t
}

// ColumnNeedsBackfill returns true if adding or dropping (according to
// the direction) the given column requires backfill.
func ColumnNeedsBackfill(col Column) bool {
	if col.IsVirtual() {
		// Virtual columns are not stored in the primary index, so they do not need
		// backfill.
		return false
	}
	if col.Dropped() {
		// In all other cases, DROP requires backfill.
		return true
	}
	// ADD requires backfill for:
	//  - columns with non-NULL default value
	//  - computed columns
	//  - non-nullable columns (note: if a non-nullable column doesn't have a
	//    default value, the backfill will fail unless the table is empty).
	if col.HasNullDefault() {
		return false
	}
	return col.HasDefault() || !col.IsNullable() || col.IsComputed()
}
