package git

import (
	"giks/test/gittest"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const testGitDir = "../test/output/git-dir"

func TestApplyMixins(t *testing.T) {
	r := gittest.NewTestRepository(testGitDir)
	defer r.Clean()
	r.WriteFile("README", "please read me")
	r.AddAll()
	r.WriteFile("README", "please read me2")
	vars := map[string]string{}
	ApplyMixins(r.AbsDir(), vars)
	assert.Contains(t, strings.Split(vars["GIKS_MIXIN_STAGED_FILES"], " "), "README", "expected affected files to have added file")
	assert.Contains(t, strings.Split(vars["GIKS_MIXIN_MODIFIED_FILES"], " "), "README", "expected affected files to have added file")
}
