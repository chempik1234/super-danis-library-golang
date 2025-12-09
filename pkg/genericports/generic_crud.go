package genericports

import "context"

// ObjectWithIdentifier describes an object that has an ID.
//
// ID must be comparable in order to find objects.
type ObjectWithIdentifier[I comparable] interface {
	// GetUniqueIdentifier retrieves a value that uniquely identifies object
	//
	// Example:
	//   func (u *User) GetUniqueIdentifier() types.UUID {
	//       return u.ID
	//   }
	GetUniqueIdentifier() I
}

// GenericStoragePort describes a permanent storage for objects
//
// Supposed to be working along with GenericCachePort
type GenericStoragePort[I comparable, T ObjectWithIdentifier[I]] interface {
	// GetObjects gets all objects list
	GetObjects(ctx context.Context) ([]*T, error)
	// GetObjectByID retrieves 1 object by given ID, if exists, else err
	GetObjectByID(ctx context.Context, id I) (*T, error)
	// CreateObject creates a new object, returns created object
	//
	// Given ID is used, so pre-generate it!
	//
	// Error on conflict
	CreateObject(ctx context.Context, fullyReadyObject *T) (*T, error)
	// UpdateObject fully updates object by given ID, if exists, else err
	//
	// ID is in the model
	UpdateObject(ctx context.Context, fullyReadyObject *T) (*T, error)
	// DeleteObject deletes a object by given ID, if exists, else err
	DeleteObject(ctx context.Context, id I) error
}

// GenericCachePort describes a temporary KV storage for object objects
//
// Supposed to be working along with GenericStoragePort
type GenericCachePort[I comparable, T ObjectWithIdentifier[I]] interface {
	// GetObjectByID retrieves 1 object by given ID, if exists, else err
	GetObjectByID(ctx context.Context, id I) (*T, error)
	// SaveObject creates a new object, returns created object
	//
	// Given ID is used, so pre-generate it!
	//
	// Error on conflict
	SaveObject(ctx context.Context, fullyReadyObject *T) (*T, error)
	// DeleteObject deletes a object by given ID, if exists, else err
	DeleteObject(ctx context.Context, id I) error
}
