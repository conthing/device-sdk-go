package cache

import (
	"fmt"
	"sync"

	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

var (
	vdc *valueDescriptorCache
)

type ValueDescriptorCache interface {
	ForName(name string) (contract.ValueDescriptor, bool)
	All() []contract.ValueDescriptor
	Add(descriptor contract.ValueDescriptor) error
	Update(descriptor contract.ValueDescriptor) error
	Remove(id string) error
	RemoveByName(name string) error
}

type valueDescriptorCache struct {
	vdMap   sync.Map
	nameMap sync.Map
}

func (v *valueDescriptorCache) ForName(name string) (contract.ValueDescriptor, bool) {
	vd, ok := v.vdMap.Load(name)
	return vd.(contract.ValueDescriptor), ok
}

func (v *valueDescriptorCache) All() []contract.ValueDescriptor {
	var vds []contract.ValueDescriptor
	f := func(k, v interface{}) bool {
		//这个函数的入参、出参的类型都已经固定，不能修改
		//可以在函数体内编写自己的代码，调用map中的k,v

		vds = append(vds, v.(contract.ValueDescriptor))
		return true
	}
	v.vdMap.Range(f)
	return vds
}

func (v *valueDescriptorCache) Add(descriptor contract.ValueDescriptor) error {
	_, ok := v.vdMap.Load(descriptor.Name)
	if ok {
		return fmt.Errorf("value descriptor %s has already existed in cache", descriptor.Name)
	}
	v.vdMap.Store(descriptor.Name, descriptor)
	v.nameMap.Store(descriptor.Id, descriptor.Name)
	return nil
}

func (v *valueDescriptorCache) Update(descriptor contract.ValueDescriptor) error {
	if err := v.Remove(descriptor.Id); err != nil {
		return err
	}
	return v.Add(descriptor)
}

func (v *valueDescriptorCache) Remove(id string) error {
	name, ok := v.nameMap.Load(id)
	if !ok {
		return fmt.Errorf("value descriptor %s does not exist in cache", id)
	}

	return v.RemoveByName(name.(string))
}

func (v *valueDescriptorCache) RemoveByName(name string) error {
	valueDescriptor, ok := v.vdMap.Load(name)
	if !ok {
		return fmt.Errorf("value descriptor %s does not exist in cache", name)
	}
	vd := valueDescriptor.(contract.ValueDescriptor)
	v.vdMap.Delete(vd.Id)
	v.nameMap.Delete(name)
	return nil
}

func newValueDescriptorCache(descriptors []contract.ValueDescriptor) ValueDescriptorCache {
	var descriptorMap sync.Map
	var nameMap sync.Map
	for _, descriptor := range descriptors {
		descriptorMap.Store(descriptor.Name, descriptor)
		nameMap.Store(descriptor.Id, descriptor.Name)
	}
	dec, err := descriptorMap.Load("mode")
	fmt.Sprintf(dec.(contract.ValueDescriptor).String(), err)
	vdc = &valueDescriptorCache{
		vdMap:   descriptorMap,
		nameMap: nameMap,
	}
	return vdc
}

func ValueDescriptors() ValueDescriptorCache {
	if vdc == nil {
		InitCache()
	}
	return vdc
}
