// Package annotate provides a lightweight store for attaching user-defined
// annotations to Vault secret paths within a given environment.
//
// Annotations include a free-text note, an owner field, and a creation
// timestamp. They are persisted as JSON files on disk, one file per
// environment, and can be queried by path.
//
// Typical usage:
//
//	store := annotate.NewStore("/var/vaultwatch/annotations")
//	store.Save("prod", annotate.Annotation{
//		Path:  "secret/db/password",
//		Note:  "rotated quarterly",
//		Owner: "dba-team",
//	})
//	ann, err := store.Get("prod", "secret/db/password")
package annotate
