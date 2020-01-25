package detection

import (
	"fmt"

	"github.com/stackrox/rox/generated/storage"
)

type setImpl struct {
	compiler           PolicyCompiler
	policyIDToCompiled StringCompiledPolicyFastRMap
}

func (p *setImpl) Compiler() PolicyCompiler {
	return p.compiler
}

func (p *setImpl) ForEach(pt PolicyExecutor) error {
	m := p.policyIDToCompiled.GetMap()

	for _, compiled := range m {
		if err := pt.Execute(compiled); err != nil {
			return err
		}
	}
	return nil
}

func (p *setImpl) ForOne(pID string, pt PolicyExecutor) error {
	compiled, exists := p.policyIDToCompiled.Get(pID)
	if exists {
		return pt.Execute(compiled)
	}
	return fmt.Errorf("policy with ID not found in set: %s", pID)
}

// UpsertPolicy adds or updates a policy in the set.
func (p *setImpl) UpsertPolicy(policy *storage.Policy) error {
	compiled, err := p.compiler.CompilePolicy(policy)
	if err != nil {
		log.Errorf("unable to compile policy: %s", err)
		return err
	}

	p.policyIDToCompiled.Set(compiled.Policy().GetId(), compiled)
	return nil
}

// RemovePolicy removes a policy from the set.
func (p *setImpl) RemovePolicy(policyID string) error {
	p.policyIDToCompiled.Delete(policyID)
	return nil
}

// GetCompiledPolicies returns all of the compiled policies
func (p *setImpl) GetCompiledPolicies() map[string]CompiledPolicy {
	return p.policyIDToCompiled.GetMap()
}

func (p *setImpl) Recompile(policyID string) error {
	olcCompiled, exists := p.policyIDToCompiled.Get(policyID)
	if !exists {
		return fmt.Errorf("policy %s does not exist to recompile", policyID)
	}

	newCompiled, err := p.compiler.CompilePolicy(olcCompiled.Policy())
	if err != nil {
		log.Errorf("unable to compile policy: %s", err)
		return err
	}

	p.policyIDToCompiled.Set(newCompiled.Policy().GetId(), newCompiled)
	return nil
}
