package dirs

import (
	"os"
	"path"
	"sync"

	"github.com/tomogoma/go-commons/errors"
)

type Helper struct {
	wdLock      sync.Mutex
	workingDirs []string
}

func NewHelper() *Helper {
	return &Helper{wdLock: sync.Mutex{}, workingDirs: make([]string, 0)}
}

func (p *Helper) PushD(dir string) error {
	p.wdLock.Lock()
	defer p.wdLock.Unlock()
	execF, err := os.Executable()
	if err != nil {
		return errors.Newf("error getting executable location: %v", err)
	}
	p.workingDirs = append(p.workingDirs, path.Dir(execF))
	if err := os.Chdir(dir); err != nil {
		return errors.Newf("pushd error: %v", err)
	}
	return nil
}

func (p *Helper) PopD() error {
	p.wdLock.Lock()
	defer p.wdLock.Unlock()
	size := len(p.workingDirs)
	if size == 0 {
		return nil
	}
	dir := p.workingDirs[size-1]
	p.workingDirs = p.workingDirs[:size-1]
	if err := os.Chdir(dir); err != nil {
		return errors.Newf("popd error: %v", err)
	}
	return nil
}
