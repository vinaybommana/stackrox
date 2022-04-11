// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"reflect"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

const (
	baseTable  = "node_cves"
	countStmt  = "SELECT COUNT(*) FROM node_cves"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM node_cves WHERE Id = $1 AND OperatingSystem = $2)"

	getStmt    = "SELECT serialized FROM node_cves WHERE Id = $1 AND OperatingSystem = $2"
	deleteStmt = "DELETE FROM node_cves WHERE Id = $1 AND OperatingSystem = $2"
	walkStmt   = "SELECT serialized FROM node_cves"

	batchAfter = 100

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	schema = walker.Walk(reflect.TypeOf((*storage.CVE)(nil)), baseTable)
	log    = logging.LoggerForModule()
)

func init() {
	globaldb.RegisterTable(schema)
}

type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string, operatingSystem string) (bool, error)
	Get(ctx context.Context, id string, operatingSystem string) (*storage.CVE, bool, error)
	Upsert(ctx context.Context, obj *storage.CVE) error
	UpsertMany(ctx context.Context, objs []*storage.CVE) error
	Delete(ctx context.Context, id string, operatingSystem string) error

	Walk(ctx context.Context, fn func(obj *storage.CVE) error) error

	AckKeysIndexed(ctx context.Context, keys ...string) error
	GetKeysToIndex(ctx context.Context) ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

func createTableNodeCves(ctx context.Context, db *pgxpool.Pool) {
	table := `
create table if not exists node_cves (
    Id varchar,
    Cve varchar,
    OperatingSystem varchar,
    Cvss numeric,
    ImpactScore numeric,
    PublishedOn timestamp,
    CreatedAt timestamp,
    Suppressed bool,
    SuppressExpiry timestamp,
    Severity integer,
    serialized bytea,
    PRIMARY KEY(Id, OperatingSystem)
)
`

	_, err := db.Exec(ctx, table)
	if err != nil {
		log.Panicf("Error creating table %s: %v", table, err)
	}

	indexes := []string{}
	for _, index := range indexes {
		if _, err := db.Exec(ctx, index); err != nil {
			log.Panicf("Error creating index %s: %v", index, err)
		}
	}

}

func insertIntoNodeCves(ctx context.Context, tx pgx.Tx, obj *storage.CVE) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		obj.GetId(),
		obj.GetCve(),
		obj.GetOperatingSystem(),
		obj.GetCvss(),
		obj.GetImpactScore(),
		pgutils.NilOrTime(obj.GetPublishedOn()),
		pgutils.NilOrTime(obj.GetCreatedAt()),
		obj.GetSuppressed(),
		pgutils.NilOrTime(obj.GetSuppressExpiry()),
		obj.GetSeverity(),
		serialized,
	}

	finalStr := "INSERT INTO node_cves (Id, Cve, OperatingSystem, Cvss, ImpactScore, PublishedOn, CreatedAt, Suppressed, SuppressExpiry, Severity, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT(Id, OperatingSystem) DO UPDATE SET Id = EXCLUDED.Id, Cve = EXCLUDED.Cve, OperatingSystem = EXCLUDED.OperatingSystem, Cvss = EXCLUDED.Cvss, ImpactScore = EXCLUDED.ImpactScore, PublishedOn = EXCLUDED.PublishedOn, CreatedAt = EXCLUDED.CreatedAt, Suppressed = EXCLUDED.Suppressed, SuppressExpiry = EXCLUDED.SuppressExpiry, Severity = EXCLUDED.Severity, serialized = EXCLUDED.serialized"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func (s *storeImpl) copyFromNodeCves(ctx context.Context, tx pgx.Tx, objs ...*storage.CVE) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"id",

		"cve",

		"operatingsystem",

		"cvss",

		"impactscore",

		"publishedon",

		"createdat",

		"suppressed",

		"suppressexpiry",

		"severity",

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

			obj.GetId(),

			obj.GetCve(),

			obj.GetOperatingSystem(),

			obj.GetCvss(),

			obj.GetImpactScore(),

			pgutils.NilOrTime(obj.GetPublishedOn()),

			pgutils.NilOrTime(obj.GetCreatedAt()),

			obj.GetSuppressed(),

			pgutils.NilOrTime(obj.GetSuppressExpiry()),

			obj.GetSeverity(),

			serialized,
		})

		if _, err := tx.Exec(ctx, deleteStmt, obj.GetId(), obj.GetOperatingSystem()); err != nil {
			return err
		}

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"node_cves"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

// New returns a new Store instance using the provided sql instance.
func New(ctx context.Context, db *pgxpool.Pool) Store {
	createTableNodeCves(ctx, db)

	return &storeImpl{
		db: db,
	}
}

func (s *storeImpl) copyFrom(ctx context.Context, objs ...*storage.CVE) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "CVE")
	if err != nil {
		return err
	}
	defer release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	if err := s.copyFromNodeCves(ctx, tx, objs...); err != nil {
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

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.CVE) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "CVE")
	if err != nil {
		return err
	}
	defer release()

	for _, obj := range objs {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}

		if err := insertIntoNodeCves(ctx, tx, obj); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return err
			}
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *storeImpl) Upsert(ctx context.Context, obj *storage.CVE) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "CVE")

	return s.upsert(ctx, obj)
}

func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.CVE) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.UpdateMany, "CVE")

	if len(objs) < batchAfter {
		return s.upsert(ctx, objs...)
	} else {
		return s.copyFrom(ctx, objs...)
	}
}

// Count returns the number of objects in the store
func (s *storeImpl) Count(ctx context.Context) (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "CVE")

	row := s.db.QueryRow(ctx, countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(ctx context.Context, id string, operatingSystem string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "CVE")

	row := s.db.QueryRow(ctx, existsStmt, id, operatingSystem)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, pgutils.ErrNilIfNoRows(err)
	}
	return exists, nil
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(ctx context.Context, id string, operatingSystem string) (*storage.CVE, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "CVE")

	conn, release, err := s.acquireConn(ctx, ops.Get, "CVE")
	if err != nil {
		return nil, false, err
	}
	defer release()

	row := conn.QueryRow(ctx, getStmt, id, operatingSystem)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	var msg storage.CVE
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}

func (s *storeImpl) acquireConn(ctx context.Context, op ops.Op, typ string) (*pgxpool.Conn, func(), error) {
	defer metrics.SetAcquireDBConnDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.Release, nil
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(ctx context.Context, id string, operatingSystem string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "CVE")

	conn, release, err := s.acquireConn(ctx, ops.Remove, "CVE")
	if err != nil {
		return err
	}
	defer release()

	if _, err := conn.Exec(ctx, deleteStmt, id, operatingSystem); err != nil {
		return err
	}
	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.CVE) error) error {
	rows, err := s.db.Query(ctx, walkStmt)
	if err != nil {
		return pgutils.ErrNilIfNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		var msg storage.CVE
		if err := proto.Unmarshal(data, &msg); err != nil {
			return err
		}
		if err := fn(&msg); err != nil {
			return err
		}
	}
	return nil
}

//// Used for testing

func dropTableNodeCves(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS node_cves CASCADE")

}

func Destroy(ctx context.Context, db *pgxpool.Pool) {
	dropTableNodeCves(ctx, db)
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
