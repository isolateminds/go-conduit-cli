package conduit

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/isolateminds/go-conduit-cli/internal/compose/composeopt"
)

type variableMap map[string]string
type templateFormatter struct {
	formatter   composeopt.TemplateFormatter
	VariableMap variableMap
}

func (tf *templateFormatter) Format(in []byte) (out io.Reader, err error) {
	env := string(in)
	for k, v := range tf.VariableMap {
		env = strings.ReplaceAll(env, fmt.Sprintf("{{%s}}", k), v)
	}
	return bytes.NewReader([]byte(env)), nil
}

func newTemplateFormatter(variableMap variableMap) *templateFormatter {
	return &templateFormatter{
		VariableMap: variableMap,
	}
}
