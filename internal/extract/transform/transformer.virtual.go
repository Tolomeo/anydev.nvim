package transform

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

type virtualOriginAtAnnotations struct {
	AtType   *AtTypeAnnotation
	AtModule *annotation.AtModule
}

func (tr *Transformer) getVirtualOriginAnnotations(docblock []string) (virtualOriginAtAnnotations, error) {
	annotations := virtualOriginAtAnnotations{
		AtType:   nil,
		AtModule: nil,
	}

	buffer, err := tr.context.Nvim().NewBuffer()

	if err != nil {
		return annotations, err
	}

	defer buffer.Close()

	err = buffer.SetLines(docblock)

	if err != nil {
		return annotations, err
	}

	atType, err := tr.getAtTypeAnnotations(buffer)

	if err != nil {
		return annotations, fmt.Errorf("Error lexing type annotation: %w", err)
	}

	atModuleMatch, err := buffer.TsQueryOne(annotation.AtModuleQuery)

	if err != nil {
		return annotations, err
	}

	if atModuleMatch != nil {
		annotations.AtModule = annotation.NewAtModule(*atModuleMatch)
	}

	annotations.AtType = atType
	return annotations, nil
}

func (tr *Transformer) getVirtualOriginType(virtualOrigin *origin.VirtualOrigin) (symbol.Type, error) {
	annotations, err := tr.getVirtualOriginAnnotations(virtualOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	if annotations.AtType != nil {
		// NB: we don't check for the presence of multiple types here
		lexedType, err := tr.getType(annotations.AtType.Types[0])

		if err != nil {
			return nil, err
		}

		return lexedType, nil
	}

	if annotations.AtModule != nil {
		moduleName := strings.Trim(annotations.AtModule.Name().Text, "'\"")
		err := tr.context.Extract("module", moduleName)

		if err != nil {
			return nil, err
		}

		return symbol.NewReference(moduleName), nil
	}

	tr.context.Logger().Warn(fmt.Sprintf("Unknown meta type '%s' received", tr.context.Target().Name()))
	return symbol.NewUnknown(), nil
}
