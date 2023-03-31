package staffrepo

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/smart48ru/FaceIDApp/internal/domain"
	"github.com/smart48ru/FaceIDApp/internal/service/staffservice"
)

var _ staffservice.StaffRepo = &Repo{}

type Repo struct {
	mu  sync.Mutex
	m   map[uint64]domain.Employee
	seq uint64
}

func (r *Repo) Serialize() []domain.Employee {
	if len(r.m) == 0 {
		return []domain.Employee{}
	}
	allEmp := make([]domain.Employee, 0, len(r.m))
	for _, employee := range r.m {
		allEmp = append(allEmp, employee)
	}
	sort.Slice(allEmp, func(i, j int) bool {
		return allEmp[i].ID < allEmp[j].ID
	})
	return allEmp
}

// Create implements staffservice.StaffRepo
func (r *Repo) Create(ctx context.Context, u domain.Employee) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	r.seq += 1
	u.ID = r.seq
	r.m[r.seq] = u

	return r.seq, nil
}

// Delete implements staffservice.StaffRepo
func (r *Repo) Delete(ctx context.Context, id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, ok := r.m[id]
	if !ok {
		return fmt.Errorf("employee id=%d not exist", id)
	}
	delete(r.m, id)

	return nil
}

// Read implements staffservice.StaffRepo
func (r *Repo) Read(ctx context.Context, id uint64) (domain.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	select {
	case <-ctx.Done():
		return domain.Employee{}, ctx.Err()
	default:
	}

	if v, ok := r.m[id]; ok {
		return v, nil
	}

	return domain.Employee{}, fmt.Errorf("employee id=%d not found", id)
}

// ReadAll implements staffservice.StaffRepo
func (r *Repo) ReadAll(ctx context.Context) ([]domain.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var employees []domain.Employee
	select {
	case <-ctx.Done():
		return []domain.Employee{}, ctx.Err()
	default:
	}

	employees = r.Serialize()

	return employees, nil
}

// Update implements staffservice.StaffRepo
func (r *Repo) Update(ctx context.Context, id uint64, u domain.Employee) (domain.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	select {
	case <-ctx.Done():
		return domain.Employee{}, ctx.Err()
	default:
	}

	r.m[id] = u

	return u, nil
}

func New() *Repo {
	return &Repo{
		m: make(map[uint64]domain.Employee),
	}
}
