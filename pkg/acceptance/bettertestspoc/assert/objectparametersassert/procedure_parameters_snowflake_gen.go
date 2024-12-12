// Code generated by assertions generator; DO NOT EDIT.

package objectparametersassert

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ProcedureParametersAssert struct {
	*assert.SnowflakeParametersAssert[sdk.SchemaObjectIdentifierWithArguments]
}

func ProcedureParameters(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) *ProcedureParametersAssert {
	t.Helper()
	return &ProcedureParametersAssert{
		assert.NewSnowflakeParametersAssertWithProvider(id, sdk.ObjectTypeProcedure, acc.TestClient().Parameter.ShowProcedureParameters),
	}
}

func ProcedureParametersPrefetched(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments, parameters []*sdk.Parameter) *ProcedureParametersAssert {
	t.Helper()
	return &ProcedureParametersAssert{
		assert.NewSnowflakeParametersAssertWithParameters(id, sdk.ObjectTypeProcedure, parameters),
	}
}

//////////////////////////////
// Generic parameter checks //
//////////////////////////////

func (p *ProcedureParametersAssert) HasBoolParameterValue(parameterName sdk.ProcedureParameter, expected bool) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterBoolValueSet(parameterName, expected))
	return p
}

func (p *ProcedureParametersAssert) HasIntParameterValue(parameterName sdk.ProcedureParameter, expected int) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterIntValueSet(parameterName, expected))
	return p
}

func (p *ProcedureParametersAssert) HasStringParameterValue(parameterName sdk.ProcedureParameter, expected string) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterValueSet(parameterName, expected))
	return p
}

func (p *ProcedureParametersAssert) HasDefaultParameterValue(parameterName sdk.ProcedureParameter) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterDefaultValueSet(parameterName))
	return p
}

func (p *ProcedureParametersAssert) HasDefaultParameterValueOnLevel(parameterName sdk.ProcedureParameter, parameterType sdk.ParameterType) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterDefaultValueOnLevelSet(parameterName, parameterType))
	return p
}

///////////////////////////////
// Aggregated generic checks //
///////////////////////////////

// HasAllDefaults checks if all the parameters:
// - have a default value by comparing current value of the sdk.Parameter with its default
// - have an expected level
func (p *ProcedureParametersAssert) HasAllDefaults() *ProcedureParametersAssert {
	return p.
		HasDefaultParameterValueOnLevel(sdk.ProcedureParameterAutoEventLogging, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.ProcedureParameterEnableConsoleOutput, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.ProcedureParameterLogLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.ProcedureParameterMetricLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.ProcedureParameterTraceLevel, sdk.ParameterTypeSnowflakeDefault)
}

func (p *ProcedureParametersAssert) HasAllDefaultsExplicit() *ProcedureParametersAssert {
	return p.
		HasDefaultAutoEventLoggingValueExplicit().
		HasDefaultEnableConsoleOutputValueExplicit().
		HasDefaultLogLevelValueExplicit().
		HasDefaultMetricLevelValueExplicit().
		HasDefaultTraceLevelValueExplicit()
}

////////////////////////////
// Parameter value checks //
////////////////////////////

func (p *ProcedureParametersAssert) HasAutoEventLogging(expected sdk.AutoEventLogging) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterStringUnderlyingValueSet(sdk.ProcedureParameterAutoEventLogging, expected))
	return p
}

func (p *ProcedureParametersAssert) HasEnableConsoleOutput(expected bool) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterBoolValueSet(sdk.ProcedureParameterEnableConsoleOutput, expected))
	return p
}

func (p *ProcedureParametersAssert) HasLogLevel(expected sdk.LogLevel) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterStringUnderlyingValueSet(sdk.ProcedureParameterLogLevel, expected))
	return p
}

func (p *ProcedureParametersAssert) HasMetricLevel(expected sdk.MetricLevel) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterStringUnderlyingValueSet(sdk.ProcedureParameterMetricLevel, expected))
	return p
}

func (p *ProcedureParametersAssert) HasTraceLevel(expected sdk.TraceLevel) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterStringUnderlyingValueSet(sdk.ProcedureParameterTraceLevel, expected))
	return p
}

////////////////////////////
// Parameter level checks //
////////////////////////////

func (p *ProcedureParametersAssert) HasAutoEventLoggingLevel(expected sdk.ParameterType) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterLevelSet(sdk.ProcedureParameterAutoEventLogging, expected))
	return p
}

func (p *ProcedureParametersAssert) HasEnableConsoleOutputLevel(expected sdk.ParameterType) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterLevelSet(sdk.ProcedureParameterEnableConsoleOutput, expected))
	return p
}

func (p *ProcedureParametersAssert) HasLogLevelLevel(expected sdk.ParameterType) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterLevelSet(sdk.ProcedureParameterLogLevel, expected))
	return p
}

func (p *ProcedureParametersAssert) HasMetricLevelLevel(expected sdk.ParameterType) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterLevelSet(sdk.ProcedureParameterMetricLevel, expected))
	return p
}

func (p *ProcedureParametersAssert) HasTraceLevelLevel(expected sdk.ParameterType) *ProcedureParametersAssert {
	p.AddAssertion(assert.SnowflakeParameterLevelSet(sdk.ProcedureParameterTraceLevel, expected))
	return p
}

////////////////////////////////////
// Parameter default value checks //
////////////////////////////////////

func (p *ProcedureParametersAssert) HasDefaultAutoEventLoggingValue() *ProcedureParametersAssert {
	return p.HasDefaultParameterValue(sdk.ProcedureParameterAutoEventLogging)
}

func (p *ProcedureParametersAssert) HasDefaultEnableConsoleOutputValue() *ProcedureParametersAssert {
	return p.HasDefaultParameterValue(sdk.ProcedureParameterEnableConsoleOutput)
}

func (p *ProcedureParametersAssert) HasDefaultLogLevelValue() *ProcedureParametersAssert {
	return p.HasDefaultParameterValue(sdk.ProcedureParameterLogLevel)
}

func (p *ProcedureParametersAssert) HasDefaultMetricLevelValue() *ProcedureParametersAssert {
	return p.HasDefaultParameterValue(sdk.ProcedureParameterMetricLevel)
}

func (p *ProcedureParametersAssert) HasDefaultTraceLevelValue() *ProcedureParametersAssert {
	return p.HasDefaultParameterValue(sdk.ProcedureParameterTraceLevel)
}

/////////////////////////////////////////////
// Parameter explicit default value checks //
/////////////////////////////////////////////

func (p *ProcedureParametersAssert) HasDefaultAutoEventLoggingValueExplicit() *ProcedureParametersAssert {
	return p.HasAutoEventLogging(sdk.AutoEventLoggingOff)
}

func (p *ProcedureParametersAssert) HasDefaultEnableConsoleOutputValueExplicit() *ProcedureParametersAssert {
	return p.HasEnableConsoleOutput(false)
}

func (p *ProcedureParametersAssert) HasDefaultLogLevelValueExplicit() *ProcedureParametersAssert {
	return p.HasLogLevel(sdk.LogLevelOff)
}

func (p *ProcedureParametersAssert) HasDefaultMetricLevelValueExplicit() *ProcedureParametersAssert {
	return p.HasMetricLevel(sdk.MetricLevelNone)
}

func (p *ProcedureParametersAssert) HasDefaultTraceLevelValueExplicit() *ProcedureParametersAssert {
	return p.HasTraceLevel(sdk.TraceLevelOff)
}