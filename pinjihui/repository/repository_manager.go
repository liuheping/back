package repository

import (
    "errors"
)

var (
    Admin *Manager
)

func init()  {
    Admin = NewManager()
}

type Manager struct {
    repositories map[string]interface{}
}

func NewManager() *Manager {
    return &Manager{
        repositories: make(map[string]interface{}),
    }
}

func L(name string) interface{} {
    return Admin.LoadRepository(name)
}

/*func (m *Manager) Request(repositoryName, method string, ctx context.Context, a ...interface{}) ([]reflect.Value, error) {
    repository := m.LoadRepository(ctx, repositoryName)
    if !repository.CheckAuth(ctx) {
        return nil, errors.New(gcontext.CredentialsError)
    }

    args := make([]reflect.Value, len(a))
    for i, v := range a {
        args[i] = reflect.ValueOf(v)
    }
    reflectValues := reflect.ValueOf(repository).MethodByName(method).Call(args)
    return reflectValues, nil
}

func (m *Manager) RequestAsVE(repositoryName, method string, ctx context.Context, a ...interface{}) (interface{}, error) {
    values, err := m.Request(repositoryName, method, ctx, a...)
    if err != nil {
        return nil, err
    }
    l := len(values)
    if l == 0 {
        return nil, nil
    }
    if err, ok := values[l - 1].Interface().(error); ok {
        if l == 1  {
            return nil, err
        } else {
            return values[0].Interface(), err
        }
    } else {
        return values[0].Interface(), nil
    }
}*/

func (m *Manager) LoadRepository(name string) (interface{}) {
    return m.repositories[name]
}

func (m *Manager) Register(name string, repository interface{}) error {
    if _, ok := m.repositories[name]; ok {
        return errors.New("already exist repository: " + name)
    }
    m.repositories[name] = repository
    return nil
}
