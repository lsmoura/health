package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/lsmoura/health/pkg/dbfieldvalues"
	"github.com/lsmoura/health/pkg/health"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func connect(ctx context.Context, dbURL string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("pgx.Connect: %w", err)
	}

	return conn, nil
}

func insertType(ctx context.Context, conn *pgx.Conn, tableName string, data []any) error {
	if len(data) == 0 {
		fmt.Println("No data to insert in", tableName)
		return nil
	}
	columns := dbfieldvalues.Fields(data[0], "id")

	var workoutData [][]any
	for i, workout := range data {
		values, err := dbfieldvalues.Values(workout, "id")
		if err != nil {
			return fmt.Errorf("dbfieldvalues.Values: %w", err)
		}

		workoutData = append(workoutData, values)

		fmt.Printf("inserting fields in %s %d/%d\t\r", tableName, i, len(data))
	}

	copyCount, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{tableName},
		columns,
		pgx.CopyFromRows(workoutData),
	)

	if err != nil {
		return fmt.Errorf("db.CopyFrom: %w", err)
	}

	fmt.Printf("Copied %d rows\n", copyCount)

	return nil
}

func toInterfaceArray[T any](data []T) []any {
	result := make([]any, len(data))
	for i, item := range data {
		result[i] = item
	}

	return result
}

func insert(ctx context.Context, db *pgx.Conn, health *health.HealthData) error {
	fmt.Printf("Inserting %d records\n", len(health.Records))

	insertList := []struct {
		tableName string
		indexName string
		data      []any
	}{
		{"records", "records_id_seq", toInterfaceArray(health.Records)},
		{"correlations", "", toInterfaceArray(health.Correlations)},
		{"workouts", "workouts_id_seq", toInterfaceArray(health.Workouts)},
		{"activity_summaries", "", toInterfaceArray(health.ActivitySummary)},
		{"clinical_records", "", toInterfaceArray(health.ClinicalRecord)},
		{"audiograms", "", toInterfaceArray(health.Audiogram)},
		{"vision_prescriptions", "", toInterfaceArray(health.VisionPrescription)},
	}

	for _, insert := range insertList {
		if _, err := db.Exec(ctx, "DELETE FROM "+insert.tableName); err != nil {
			return fmt.Errorf("DELETE FROM %s: %w", insert.tableName, err)
		}
		if insert.indexName != "" {
			if _, err := db.Exec(ctx, "ALTER SEQUENCE "+insert.indexName+" RESTART WITH 1"); err != nil {
				return fmt.Errorf("ALTER SEQUENCE %s RESTART WITH 1: %w", insert.indexName, err)
			}
		}
		if err := insertType(ctx, db, insert.tableName, insert.data); err != nil {
			return fmt.Errorf("insertRecords: %w", err)
		}
	}

	return nil
}

type Options struct {
	Help bool

	Input string // defaults to export.xml

	DBHost   string // defaults to localhost
	DBUser   string // defaults to postgres
	DBPort   int    // defaults to 5432
	DBPass   string // defaults to no password
	DBSSL    bool   // defaults to false
	Database string // defaults to health

	Version     bool
	ApplySchema bool
}

func (o Options) DBURL() string {
	sb := strings.Builder{}
	sb.WriteString("postgres://")
	sb.WriteString(o.DBUser)
	if o.DBPass != "" {
		sb.WriteString(":")
		sb.WriteString(o.DBPass)
	}
	sb.WriteString("@")
	sb.WriteString(o.DBHost)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(o.DBPort))
	sb.WriteString("/")
	sb.WriteString(o.Database)
	if !o.DBSSL {
		sb.WriteString("?sslmode=disable")
	}

	return sb.String()
}

func getOptions() Options {
	var options Options

	flag.BoolVar(&options.Help, "help", false, "show help")
	flag.BoolVar(&options.Version, "version", false, "show version and exit")

	flag.BoolVar(&options.ApplySchema, "apply-schema", false, "apply schema (this will recreate all tables. It should be used on the first run.)")
	flag.StringVar(&options.Input, "input", "export.xml", "input file")
	flag.StringVar(&options.DBHost, "dbhost", "localhost", "database host")
	flag.StringVar(&options.DBUser, "dbuser", "postgres", "database user")
	flag.IntVar(&options.DBPort, "dbport", 5432, "database port")
	flag.StringVar(&options.DBPass, "dbpass", "", "database password")
	flag.BoolVar(&options.DBSSL, "dbssl", false, "database ssl")
	flag.StringVar(&options.Database, "database", "health", "database name")

	flag.Usage = usage

	flag.Parse()

	return options
}

func printVersion() {
	fmt.Printf("health converter version %s (%s) built on %s\n", version, commit, date)
}

func usage() {
	printVersion()
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	options := getOptions()

	if options.Version {
		printVersion()
		return
	}

	if options.Help {
		flag.Usage()
		return
	}

	file, err := os.Open("export.xml")
	if err != nil {
		log.Panicf("open: %v\n", err)
	}
	decoder := xml.NewDecoder(file)

	ctx := context.Background()

	db, err := connect(ctx, options.DBURL())
	if err != nil {
		log.Panicf("connect: %v\n", err)
	}
	defer db.Close(ctx)

	if err := db.Ping(ctx); err != nil {
		log.Panicf("ping: %v\n", err)
	}

	if options.ApplySchema {
		fmt.Println("applying schema...")
		schema, err := health.Schema()
		if err != nil {
			log.Panicf("cannot read schema: %v\n", err)
		}
		if _, err := db.Exec(ctx, schema); err != nil {
			log.Panicf("cannot apply schema: %v\n", err)
		}
	}

	var data health.HealthData
	fmt.Println("decoding data...")
	if err := decoder.Decode(&data); err != nil {
		log.Panicf("decode error: %v\n", err)
	}

	fmt.Println("inserting data...")
	if err := insert(ctx, db, &data); err != nil {
		log.Panicf("insert error: %v\n", err)
	}
}
