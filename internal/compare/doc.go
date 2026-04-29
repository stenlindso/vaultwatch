// Package compare implements multi-environment secret path comparison.
//
// It builds a pairwise diff matrix from a set of environment snapshots,
// identifying paths unique to each environment and paths shared across them.
// This is useful for auditing consistency across dev, staging, and production.
//
// # Overview
//
// The primary entry point is [Matrix], which accepts a map of environment
// names to their respective sets of secret paths and returns a [Result]
// describing which paths are exclusive to each environment and which are
// present in all of them.
//
// # Example
//
//	envs := map[string][]string{
//		"dev":  {"secret/db", "secret/api"},
//		"prod": {"secret/db", "secret/tls"},
//	}
//	result, err := compare.Matrix(envs)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result.Common)   // paths present in all environments
//	fmt.Println(result.Unique)   // map of env name -> paths only in that env
package compare
