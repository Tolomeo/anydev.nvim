package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getMetadata(symbolOrigin origin.Origin) (*symbol.Metadata, error) {
	meta := symbol.NewMetadata()
	meta.SetStatic(tr.context.Target().IsChild())

	switch symbolOriginType := symbolOrigin.(type) {
	case *origin.AliasOrigin:
		return meta, nil

	case *origin.AliasEnumeratorOrigin:
		return meta, nil

	case *origin.ClassOrigin:
		return meta, nil

	case *origin.FieldAnnotationOrigin:
		meta.SetPrivate(symbolOriginType.Private())
		meta.SetProtected(symbolOriginType.Protected())
		meta.SetPackage(symbolOriginType.Package())
		return meta, nil

	case *origin.FunctionOrigin:
		meta.SetStatic(symbolOriginType.Static())

	case *origin.FunctionCallOrigin:
		switch symbolOriginType.FunctionName() {
		case "vim._defer_deprecated_module":
			meta.SetDeprecated(true)
			return meta, nil
		}
	}

	buffer, err := tr.context.Nvim().OpenScratchBuffer()

	if err != nil {
		return meta, err
	}

	defer buffer.Close()

	err = buffer.SetLines(symbolOrigin.Annotations())

	if err != nil {
		return meta, nil
	}

	atPrivateMatch, err := buffer.TsQueryOne(annotation.AtPrivateAnnotationQuery)

	if err != nil {
		return meta, nil
	}

	atProtectedMatch, err := buffer.TsQueryOne(annotation.AtProtectedAnnotationQuery)

	if err != nil {
		return meta, nil
	}

	atPackageMatch, err := buffer.TsQueryOne(annotation.AtPackageAnnotationQuery)

	if err != nil {
		return meta, nil
	}

	atDeprecatedMatch, err := buffer.TsQueryOne(annotation.AtDeprectedAnnotationQuery)

	if err != nil {
		return meta, nil
	}

	meta.SetPrivate(atPrivateMatch != nil)
	meta.SetProtected(atProtectedMatch != nil)
	meta.SetPackage(atPackageMatch != nil)
	meta.SetDeprecated(atDeprecatedMatch != nil)

	return meta, nil
}
