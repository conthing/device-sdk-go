package cache

import (
	"fmt"
	"strings"
	"sync"

	"github.com/conthing/device-sdk-go/internal/common"

	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

var (
	pc *profileCache
)

type ProfileCache interface {
	All() []contract.DeviceProfile
	ForName(name string) (contract.DeviceProfile, bool)
	ForId(id string) (contract.DeviceProfile, bool)
	Add(profile contract.DeviceProfile) error
	Update(profile contract.DeviceProfile) error
	Remove(id string) error
	RemoveByName(name string) error
	CommandExists(profileName string, cmd string) (bool, error)
	DeviceResource(profileName string, resourceName string) (contract.DeviceResource, bool)
	ResourceOperation(profileName string, object string, method string) (contract.ResourceOperation, error)
	ResourceOperations(profileName string, cmd string, method string) ([]contract.ResourceOperation, error)
}

type profileCache struct {
	dpMap    sync.Map
	nameMap  sync.Map
	drMap    sync.Map
	dcMap    sync.Map
	setOpMap sync.Map
	ccMap    sync.Map
}

func (p *profileCache) All() []contract.DeviceProfile {
	var ps []contract.DeviceProfile
	f := func(k, v interface{}) bool {
		ps = append(ps, v.(contract.DeviceProfile))
		return true
	}
	p.dpMap.Range(f)
	return ps
}

func (p *profileCache) ForName(name string) (contract.DeviceProfile, bool) {
	dp, ok := p.dpMap.Load(name)
	return dp.(contract.DeviceProfile), ok
}

func (p *profileCache) ForId(id string) (contract.DeviceProfile, bool) {
	name, ok := p.nameMap.Load(id)
	if !ok {
		return contract.DeviceProfile{}, ok
	}

	dp, ok := p.dpMap.Load(name.(string))
	return dp.(contract.DeviceProfile), ok
}

func (p *profileCache) Add(profile contract.DeviceProfile) error {
	if _, ok := p.dpMap.Load(profile.Name); ok {
		return fmt.Errorf("device profile %s has already existed in cache", profile.Name)
	}

	p.dpMap.Store(profile.Name, profile)
	p.nameMap.Store(profile.Id, profile.Name)
	p.drMap.Store(profile.Name, deviceResourceSliceToMap(profile.DeviceResources))
	get, set := profileResourceSliceToMaps(profile.DeviceCommands)
	p.dcMap.Store(profile.Name, get)
	p.setOpMap.Store(profile.Name, set)
	p.ccMap.Store(profile.Name, commandSliceToMap(profile.CoreCommands))
	return nil
}

func (p *profileCache) Update(profile contract.DeviceProfile) error {
	if err := p.Remove(profile.Id); err != nil {
		return err
	}

	return p.Add(profile)
}

func (p *profileCache) Remove(id string) error {
	name, ok := p.nameMap.Load(id)
	if !ok {
		return fmt.Errorf("device profile %s does not exist in cache", id)
	}

	return p.RemoveByName(name.(string))
}

func (p *profileCache) RemoveByName(name string) error {
	profile, ok := p.dpMap.Load(name)
	if !ok {
		return fmt.Errorf("device profile %s does not exist in cache", name)
	}

	p.dpMap.Delete(name)
	pf := profile.(contract.DeviceProfile)
	p.nameMap.Delete(pf.Id)
	p.drMap.Delete(name)
	p.dcMap.Delete(name)
	p.setOpMap.Delete(name)
	p.ccMap.Delete(name)
	return nil
}

func (p *profileCache) CommandExists(profileName string, cmd string) (bool, error) {
	commands, ok := p.ccMap.Load(profileName)
	if !ok {
		err := fmt.Errorf("specified profile: %s not found", profileName)
		return false, err
	}
	cmds := commands.(map[string]contract.Command)
	if _, ok := cmds[cmd]; !ok {
		return false, nil
	}

	return true, nil
}

func (p *profileCache) DeviceResource(profileName string, resourceName string) (contract.DeviceResource, bool) {
	drsm, ok := p.drMap.Load(profileName)
	if !ok {
		return contract.DeviceResource{}, ok
	}
	drs := drsm.(map[string]contract.DeviceResource)
	dr, ok := drs[resourceName]
	return dr, ok
}

func (p *profileCache) ResourceOperations(profileName string, cmd string, method string) ([]contract.ResourceOperation, error) {
	var resOps []contract.ResourceOperation
	var rosMap map[string][]contract.ResourceOperation
	var ok bool
	if strings.ToLower(method) == common.GetCmdMethod {
		rosMapget, ok := p.dcMap.Load(profileName)
		if !ok {
			return nil, fmt.Errorf("specified profile:%s not found", profileName)
		}
		rosMap = rosMapget.(map[string][]contract.ResourceOperation)
	} else if strings.ToLower(method) == common.SetCmdMethod {
		rosMapset, ok := p.setOpMap.Load(profileName)
		if !ok {
			return nil, fmt.Errorf("specified profile: %s not found", profileName)
		}
		rosMap = rosMapset.(map[string][]contract.ResourceOperation)
	}
	if resOps, ok = rosMap[cmd]; !ok {
		return nil, fmt.Errorf("specified cmd: %s not found", cmd)
	}
	return resOps, nil
}

func (p *profileCache) ResourceOperation(profileName string, object string, method string) (contract.ResourceOperation, error) {
	var ro contract.ResourceOperation
	var rosMap map[string][]contract.ResourceOperation
	var ok bool
	if strings.ToLower(method) == common.GetCmdMethod {
		rosMapget, ok := p.dcMap.Load(profileName)
		if !ok {
			return ro, fmt.Errorf("specified profile:%s not found", profileName)
		}
		rosMap = rosMapget.(map[string][]contract.ResourceOperation)
	} else if strings.ToLower(method) == common.SetCmdMethod {
		rosMapset, ok := p.setOpMap.Load(profileName)
		if !ok {
			return ro, fmt.Errorf("specified profile:%s not found", profileName)
		}
		rosMap = rosMapset.(map[string][]contract.ResourceOperation)
	}

	if ro, ok = retrieveFirstRObyObject(rosMap, object); !ok {
		return ro, fmt.Errorf("specified ResourceOperation by object %s not found", object)
	}
	return ro, nil
}

func retrieveFirstRObyObject(rosMap map[string][]contract.ResourceOperation, object string) (contract.ResourceOperation, bool) {
	for _, ros := range rosMap {
		for _, ro := range ros {
			if ro.Object == object {
				return ro, true
			}
		}
	}
	return contract.ResourceOperation{}, false
}

func deviceResourceSliceToMap(deviceResources []contract.DeviceResource) map[string]contract.DeviceResource {
	result := make(map[string]contract.DeviceResource, len(deviceResources))
	for _, dr := range deviceResources {
		result[dr.Name] = dr
	}
	return result
}

func profileResourceSliceToMaps(profileResources []contract.ProfileResource) (map[string][]contract.ResourceOperation, map[string][]contract.ResourceOperation) {
	getResult := make(map[string][]contract.ResourceOperation, len(profileResources))
	setResult := make(map[string][]contract.ResourceOperation, len(profileResources))
	for _, pr := range profileResources {
		if len(pr.Get) > 0 {
			getResult[pr.Name] = pr.Get
		}
		if len(pr.Set) > 0 {
			setResult[pr.Name] = pr.Set
		}
	}
	return getResult, setResult
}

func commandSliceToMap(commands []contract.Command) map[string]contract.Command {
	result := make(map[string]contract.Command, len(commands))
	for _, cmd := range commands {
		result[cmd.Name] = cmd
	}
	return result
}

func newProfileCache(profiles []contract.DeviceProfile) ProfileCache {
	var dpMap sync.Map
	var nameMap sync.Map
	var drMap sync.Map
	var getOpMap sync.Map
	var setOpMap sync.Map
	var cmdMap sync.Map
	for _, profile := range profiles {
		dpMap.Store(profile.Name, profile)
		nameMap.Store(profile.Id, profile.Name)
		drMap.Store(profile.Name, deviceResourceSliceToMap(profile.DeviceResources))
		getResult, setResult := profileResourceSliceToMaps(profile.DeviceCommands)
		getOpMap.Store(profile.Name, getResult)
		setOpMap.Store(profile.Name, setResult)
		cmdMap.Store(profile.Name, commandSliceToMap(profile.CoreCommands))
	}

	return &profileCache{dpMap: dpMap, nameMap: nameMap, drMap: drMap, dcMap: getOpMap, setOpMap: setOpMap, ccMap: cmdMap}
}

func Profiles() ProfileCache {
	if pc == nil {
		InitCache()
	}
	return pc
}
