package expr

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/schemabuilder"
	sdkapi "github.com/grafana/grafana-plugin-sdk-go/v0alpha1"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/expr/classic"
	"github.com/grafana/grafana/pkg/expr/mathexp"
)

func TestQueryTypeDefinitions(t *testing.T) {
	builder, err := schemabuilder.NewSchemaBuilder(
		schemabuilder.BuilderOptions{
			PluginID: []string{DatasourceType},
			ScanCode: []schemabuilder.CodePaths{{
				BasePackage: "github.com/grafana/grafana/pkg/expr",
				CodePath:    "./",
			}},
			Enums: []reflect.Type{
				reflect.TypeOf(mathexp.ReducerSum),   // pick an example value (not the root)
				reflect.TypeOf(mathexp.UpsamplerPad), // pick an example value (not the root)
				reflect.TypeOf(ReduceModeDrop),       // pick an example value (not the root)
				reflect.TypeOf(ThresholdIsAbove),
				reflect.TypeOf(classic.ConditionOperatorAnd),
			},
		})
	require.NoError(t, err)
	err = builder.AddQueries(
		schemabuilder.QueryTypeInfo{
			Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeMath),
			GoType:         reflect.TypeOf(&MathQuery{}),
			Examples: []sdkapi.QueryExample{
				{
					Name: "constant addition",
					SaveModel: sdkapi.AsUnstructured(MathQuery{
						Expression: "$A + 10",
					}),
				},
				{
					Name: "math with two queries",
					SaveModel: sdkapi.AsUnstructured(MathQuery{
						Expression: "$A - $B",
					}),
				},
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeReduce),
			GoType:         reflect.TypeOf(&ReduceQuery{}),
			Examples: []sdkapi.QueryExample{
				{
					Name: "get max value",
					SaveModel: sdkapi.AsUnstructured(ReduceQuery{
						Expression: "$A",
						Reducer:    mathexp.ReducerMax,
						Settings: &ReduceSettings{
							Mode: ReduceModeDrop,
						},
					}),
				},
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeResample),
			GoType:         reflect.TypeOf(&ResampleQuery{}),
			Examples: []sdkapi.QueryExample{
				{
					Name: "resample at a every day",
					SaveModel: sdkapi.AsUnstructured(ResampleQuery{
						Expression:  "$A",
						Window:      "1d",
						Downsampler: mathexp.ReducerLast,
						Upsampler:   mathexp.UpsamplerPad,
					}),
				},
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeSQL),
			GoType:         reflect.TypeOf(&SQLExpression{}),
			Examples: []sdkapi.QueryExample{
				{
					Name: "Select the first row from A",
					SaveModel: sdkapi.AsUnstructured(SQLExpression{
						Expression: "SELECT * FROM A limit 1",
					}),
				},
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeClassic),
			GoType:         reflect.TypeOf(&ClassicQuery{}),
			Examples:       []sdkapi.QueryExample{
				// {
				// 	Name: "do classic query (TODO)",
				// 	SaveModel: ClassicQuery{
				// 		// ????
				// 		Conditions: []classic.ConditionJSON{},
				// 	},
				// },
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeThreshold),
			GoType:         reflect.TypeOf(&ThresholdQuery{}),
			Examples:       []sdkapi.QueryExample{
				// {
				// 	Name: "TODO... a threshold query",
				// 	SaveModel: ThresholdQuery{
				// 		Expression: "$A",
				// 	},
				// },
			},
		},
	)

	require.NoError(t, err)
	_ = builder.UpdateQueryDefinition(t, "./")
}