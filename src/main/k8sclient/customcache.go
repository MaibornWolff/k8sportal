package k8sclient

import (
    "sync"
    "k8sportal/model"
)

type serviceCustomCache struct {
    servicesMap map[string]*model.Service
    mu sync.RWMutex
}

func (s *serviceCustomCache) GetService(name string) (*model.Service, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    service, ok := s.servicesMap[name]
    return service, ok

}

func (s *serviceCustomCache) ToList() []*model.Service {
    ret := make([]*model.Service, 0, len(s.servicesMap))
    s.mu.RLock()
    defer s.mu.RUnlock()

    for _, service := range s.servicesMap {
        ret = append(ret, service)
    }
    return ret
}

func (s *serviceCustomCache) AddService(service *model.Service) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.servicesMap[service.ServiceName] = service
}

func (s *serviceCustomCache) DeleteService(serviceName string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.servicesMap, serviceName)
}
