// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	pkgSchema "github.com/stackrox/rox/migrator/migrations/frozenschema/v73"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

const (
	baseTable = "pods"

	batchAfter = 100

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000

	cursorBatchSize = 50
	deleteBatchSize = 5000
)

var (
	log    = logging.LoggerForModule()
	schema = pkgSchema.PodsSchema
)

type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, id string) (*storage.Pod, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storage.Pod, error)
	Upsert(ctx context.Context, obj *storage.Pod) error
	UpsertMany(ctx context.Context, objs []*storage.Pod) error
	Delete(ctx context.Context, id string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	GetIDs(ctx context.Context) ([]string, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.Pod, []int, error)
	DeleteMany(ctx context.Context, ids []string) error

	Walk(ctx context.Context, fn func(obj *storage.Pod) error) error

	AckKeysIndexed(ctx context.Context, keys ...string) error
	GetKeysToIndex(ctx context.Context) ([]string, error)
}

type storeImpl struct {
	db    *pgxpool.Pool
	mutex sync.Mutex
}

// New returns a new Store instance using the provided sql instance.
func New(db *pgxpool.Pool) Store {
	return &storeImpl{
		db: db,
	}
}

func insertIntoPods(ctx context.Context, batch *pgx.Batch, obj *storage.Pod) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(obj.GetId()),
		obj.GetName(),
		pgutils.NilOrUUID(obj.GetDeploymentId()),
		obj.GetNamespace(),
		pgutils.NilOrUUID(obj.GetClusterId()),
		serialized,
	}
	if pgutils.NilOrUUID(obj.GetId()) == nil {
		utils.Should(errors.Errorf("Id is not a valid uuid -- %v", obj))
		return nil
	}

	finalStr := "INSERT INTO pods (Id, Name, DeploymentId, Namespace, ClusterId, serialized) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Name = EXCLUDED.Name, DeploymentId = EXCLUDED.DeploymentId, Namespace = EXCLUDED.Namespace, ClusterId = EXCLUDED.ClusterId, serialized = EXCLUDED.serialized"
	batch.Queue(finalStr, values...)

	var query string

	for childIdx, child := range obj.GetLiveInstances() {
		if err := insertIntoPodsLiveInstances(ctx, batch, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from pods_live_instances where pods_Id = $1 AND idx >= $2"
	batch.Queue(query, pgutils.NilOrUUID(obj.GetId()), len(obj.GetLiveInstances()))
	return nil
}

func insertIntoPodsLiveInstances(ctx context.Context, batch *pgx.Batch, obj *storage.ContainerInstance, pods_Id string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(pods_Id),
		idx,
		obj.GetImageDigest(),
	}
	if pgutils.NilOrUUID(pods_Id) == nil {
		utils.Should(errors.Errorf("pods_Id is not a valid uuid -- %v", obj))
		return nil
	}

	finalStr := "INSERT INTO pods_live_instances (pods_Id, idx, ImageDigest) VALUES($1, $2, $3) ON CONFLICT(pods_Id, idx) DO UPDATE SET pods_Id = EXCLUDED.pods_Id, idx = EXCLUDED.idx, ImageDigest = EXCLUDED.ImageDigest"
	batch.Queue(finalStr, values...)

	return nil
}

func (s *storeImpl) copyFromPods(ctx context.Context, tx pgx.Tx, objs ...*storage.Pod) error {

	inputRows := [][]interface{}{}

	var err error

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	var deletes []string

	copyCols := []string{

		"id",

		"name",

		"deploymentid",

		"namespace",

		"clusterid",

		"serialized",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		serialized, marshalErr := obj.Marshal()
		if marshalErr != nil {
			return marshalErr
		}

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(obj.GetId()),

			obj.GetName(),

			pgutils.NilOrUUID(obj.GetDeploymentId()),

			obj.GetNamespace(),

			pgutils.NilOrUUID(obj.GetClusterId()),

			serialized,
		})
		if pgutils.NilOrUUID(obj.GetId()) == nil {
			utils.Should(errors.Errorf("Id is not a valid uuid -- %v", obj))
			continue
		}

		// Add the id to be deleted.
		deletes = append(deletes, obj.GetId())

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if err := s.DeleteMany(ctx, deletes); err != nil {
				return err
			}
			// clear the inserts and vals for the next batch
			deletes = nil

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"pods"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err = s.copyFromPodsLiveInstances(ctx, tx, obj.GetId(), obj.GetLiveInstances()...); err != nil {
			return err
		}
	}

	return err
}

func (s *storeImpl) copyFromPodsLiveInstances(ctx context.Context, tx pgx.Tx, pods_Id string, objs ...*storage.ContainerInstance) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"pods_id",

		"idx",

		"imagedigest",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(pods_Id),

			idx,

			obj.GetImageDigest(),
		})
		if pgutils.NilOrUUID(pods_Id) == nil {
			utils.Should(errors.Errorf("pods_Id is not a valid uuid -- %v", obj))
			continue
		}

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"pods_live_instances"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) copyFrom(ctx context.Context, objs ...*storage.Pod) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "Pod")
	if err != nil {
		return err
	}
	defer release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	if err := s.copyFromPods(ctx, tx, objs...); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.Pod) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "Pod")
	if err != nil {
		return err
	}
	defer release()

	for _, obj := range objs {
		batch := &pgx.Batch{}
		if err := insertIntoPods(ctx, batch, obj); err != nil {
			return err
		}
		batchResults := conn.SendBatch(ctx, batch)
		var result *multierror.Error
		for i := 0; i < batch.Len(); i++ {
			_, err := batchResults.Exec()
			result = multierror.Append(result, err)
		}
		if err := batchResults.Close(); err != nil {
			return err
		}
		if err := result.ErrorOrNil(); err != nil {
			return err
		}
	}
	return nil
}

func (s *storeImpl) Upsert(ctx context.Context, obj *storage.Pod) error {

	return pgutils.Retry(func() error {
		return s.upsert(ctx, obj)
	})
}

func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.Pod) error {

	return pgutils.Retry(func() error {
		// Lock since copyFrom requires a delete first before being executed.  If multiple processes are updating
		// same subset of rows, both deletes could occur before the copyFrom resulting in unique constraint
		// violations
		s.mutex.Lock()
		defer s.mutex.Unlock()

		if len(objs) < batchAfter {
			return s.upsert(ctx, objs...)
		} else {
			return s.copyFrom(ctx, objs...)
		}
	})
}

// Count returns the number of objects in the store
func (s *storeImpl) Count(ctx context.Context) (int, error) {

	var sacQueryFilter *v1.Query

	return postgres.RunCountRequestForSchema(ctx, schema, sacQueryFilter, s.db)
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(ctx context.Context, id string) (bool, error) {

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	count, err := postgres.RunCountRequestForSchema(ctx, schema, q, s.db)
	// With joins and multiple paths to the scoping resources, it can happen that the Count query for an object identifier
	// returns more than 1, despite the fact that the identifier is unique in the table.
	return count > 0, err
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(ctx context.Context, id string) (*storage.Pod, bool, error) {

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	data, err := postgres.RunGetQueryForSchema[storage.Pod](ctx, schema, q, s.db)
	if err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	return data, true, nil
}

func (s *storeImpl) acquireConn(ctx context.Context, op ops.Op, typ string) (*pgxpool.Conn, func(), error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.Release, nil
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(ctx context.Context, id string) error {

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	return postgres.RunDeleteRequestForSchema(ctx, schema, q, s.db)
}

// DeleteByQuery removes the objects based on the passed query
func (s *storeImpl) DeleteByQuery(ctx context.Context, query *v1.Query) error {

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		query,
	)

	return postgres.RunDeleteRequestForSchema(ctx, schema, q, s.db)
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs(ctx context.Context) ([]string, error) {
	var sacQueryFilter *v1.Query
	result, err := postgres.RunSearchRequestForSchema(ctx, schema, sacQueryFilter, s.db)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(result))
	for _, entry := range result {
		ids = append(ids, entry.ID)
	}

	return ids, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice
func (s *storeImpl) GetMany(ctx context.Context, ids []string) ([]*storage.Pod, []int, error) {

	if len(ids) == 0 {
		return nil, nil, nil
	}

	var sacQueryFilter *v1.Query
	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(ids...).ProtoQuery(),
	)

	rows, err := postgres.RunGetManyQueryForSchema[storage.Pod](ctx, schema, q, s.db)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			missingIndices := make([]int, 0, len(ids))
			for i := range ids {
				missingIndices = append(missingIndices, i)
			}
			return nil, missingIndices, nil
		}
		return nil, nil, err
	}
	resultsByID := make(map[string]*storage.Pod, len(rows))
	for _, msg := range rows {
		resultsByID[msg.GetId()] = msg
	}
	missingIndices := make([]int, 0, len(ids)-len(resultsByID))
	// It is important that the elems are populated in the same order as the input ids
	// slice, since some calling code relies on that to maintain order.
	elems := make([]*storage.Pod, 0, len(resultsByID))
	for i, id := range ids {
		if result, ok := resultsByID[id]; !ok {
			missingIndices = append(missingIndices, i)
		} else {
			elems = append(elems, result)
		}
	}
	return elems, missingIndices, nil
}

// GetByQuery returns the objects matching the query
func (s *storeImpl) GetByQuery(ctx context.Context, query *v1.Query) ([]*storage.Pod, error) {

	var sacQueryFilter *v1.Query
	q := search.ConjunctionQuery(
		sacQueryFilter,
		query,
	)

	rows, err := postgres.RunGetManyQueryForSchema[storage.Pod](ctx, schema, q, s.db)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return rows, nil
}

// Delete removes the specified IDs from the store
func (s *storeImpl) DeleteMany(ctx context.Context, ids []string) error {

	var sacQueryFilter *v1.Query

	// Batch the deletes
	localBatchSize := deleteBatchSize
	numRecordsToDelete := len(ids)
	for {
		if len(ids) == 0 {
			break
		}

		if len(ids) < localBatchSize {
			localBatchSize = len(ids)
		}

		idBatch := ids[:localBatchSize]
		q := search.ConjunctionQuery(
			sacQueryFilter,
			search.NewQueryBuilder().AddDocIDs(idBatch...).ProtoQuery(),
		)

		if err := postgres.RunDeleteRequestForSchema(ctx, schema, q, s.db); err != nil {
			err = errors.Wrapf(err, "unable to delete the records.  Successfully deleted %d out of %d", numRecordsToDelete-len(ids), numRecordsToDelete)
			log.Error(err)
			return err
		}

		// Move the slice forward to start the next batch
		ids = ids[localBatchSize:]
	}

	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.Pod) error) error {
	var sacQueryFilter *v1.Query
	fetcher, closer, err := postgres.RunCursorQueryForSchema[storage.Pod](ctx, schema, sacQueryFilter, s.db)
	if err != nil {
		return err
	}
	defer closer()
	for {
		rows, err := fetcher(cursorBatchSize)
		if err != nil {
			return pgutils.ErrNilIfNoRows(err)
		}
		for _, data := range rows {
			if err := fn(data); err != nil {
				return err
			}
		}
		if len(rows) != cursorBatchSize {
			break
		}
	}
	return nil
}

//// Used for testing

func dropTablePods(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS pods CASCADE")
	dropTablePodsLiveInstances(ctx, db)

}

func dropTablePodsLiveInstances(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS pods_live_instances CASCADE")

}

func Destroy(ctx context.Context, db *pgxpool.Pool) {
	dropTablePods(ctx, db)
}

//// Stubs for satisfying legacy interfaces

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(ctx context.Context, keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex(ctx context.Context) ([]string, error) {
	return nil, nil
}
