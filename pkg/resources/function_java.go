package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FunctionJava() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.FunctionJava, CreateContextFunctionJava),
		ReadContext:   TrackingReadWrapper(resources.FunctionJava, ReadContextFunctionJava),
		UpdateContext: TrackingUpdateWrapper(resources.FunctionJava, UpdateContextFunctionJava),
		DeleteContext: TrackingDeleteWrapper(resources.FunctionJava, DeleteFunction),
		Description:   "Resource used to manage java function objects. For more information, check [function documentation](https://docs.snowflake.com/en/sql-reference/sql/create-function).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FunctionJava, customdiff.All(
			// TODO [SNOW-1348103]: ComputedIfAnyAttributeChanged(javaFunctionSchema, ShowOutputAttributeName, ...),
			ComputedIfAnyAttributeChanged(javaFunctionSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(functionParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllFunctionParameters), strings.ToLower)...),
			functionParametersCustomDiff,
			// The language check is more for the future.
			// Currently, almost all attributes are marked as forceNew.
			// When language changes, these attributes also change, causing the object to recreate either way.
			// The only potential option is java staged <-> scala staged (however scala need runtime_version which may interfere).
			RecreateWhenResourceStringFieldChangedExternally("function_language", "JAVA"),
		)),

		Schema: collections.MergeMaps(javaFunctionSchema, functionParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateContextFunctionJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	database := d.Get("database").(string)
	sc := d.Get("schema").(string)
	name := d.Get("name").(string)

	argumentRequests, err := parseFunctionArgumentsCommon(d)
	if err != nil {
		return diag.FromErr(err)
	}
	returns, err := parseFunctionReturnsCommon(d)
	if err != nil {
		return diag.FromErr(err)
	}
	handler := d.Get("handler").(string)

	argumentDataTypes := collections.Map(argumentRequests, func(r sdk.FunctionArgumentRequest) datatypes.DataType { return r.ArgDataType })
	id := sdk.NewSchemaObjectIdentifierWithArgumentsNormalized(database, sc, name, argumentDataTypes...)
	request := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
		WithArguments(argumentRequests)

	errs := errors.Join(
		booleanStringAttributeCreateBuilder(d, "is_secure", request.WithSecure),
		attributeMappedValueCreateBuilder[string](d, "null_input_behavior", request.WithNullInputBehavior, sdk.ToNullInputBehavior),
		attributeMappedValueCreateBuilder[string](d, "return_results_behavior", request.WithReturnResultsBehavior, sdk.ToReturnResultsBehavior),
		stringAttributeCreateBuilder(d, "runtime_version", request.WithRuntimeVersion),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		setFunctionImportsInBuilder(d, request.WithImports),
		setFunctionPackagesInBuilder(d, request.WithPackages),
		setExternalAccessIntegrationsInBuilder(d, request.WithExternalAccessIntegrations),
		setSecretsInBuilder(d, request.WithSecrets),
		setFunctionTargetPathInBuilder(d, request.WithTargetPath),
		stringAttributeCreateBuilder(d, "function_definition", request.WithFunctionDefinitionWrapped),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.Functions.CreateForJava(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	// parameters do not work in create function (query does not fail but parameters stay unchanged)
	setRequest := sdk.NewFunctionSetRequest()
	if parametersCreateDiags := handleFunctionParametersCreate(d, setRequest); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}
	if !reflect.DeepEqual(*setRequest, *sdk.NewFunctionSetRequest()) {
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSet(*setRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextFunctionJava(ctx, d, meta)
}

func ReadContextFunctionJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	allFunctionDetails, diags := queryAllFunctionDetailsCommon(ctx, d, client, id)
	if diags != nil {
		return diags
	}

	// TODO [SNOW-1348103]: handle external changes marking
	// TODO [SNOW-1348103]: handle setting state to value from config

	errs := errors.Join(
		// not reading is_secure on purpose (handled as external change to show output)
		readFunctionOrProcedureArguments(d, allFunctionDetails.functionDetails.NormalizedArguments),
		d.Set("return_type", allFunctionDetails.functionDetails.ReturnDataType.ToSql()),
		// not reading null_input_behavior on purpose (handled as external change to show output)
		// not reading return_results_behavior on purpose (handled as external change to show output)
		setOptionalFromStringPtr(d, "runtime_version", allFunctionDetails.functionDetails.RuntimeVersion),
		d.Set("comment", allFunctionDetails.function.Description),
		readFunctionOrProcedureImports(d, allFunctionDetails.functionDetails.NormalizedImports),
		d.Set("packages", allFunctionDetails.functionDetails.NormalizedPackages),
		setRequiredFromStringPtr(d, "handler", allFunctionDetails.functionDetails.Handler),
		readFunctionOrProcedureExternalAccessIntegrations(d, allFunctionDetails.functionDetails.NormalizedExternalAccessIntegrations),
		readFunctionOrProcedureSecrets(d, allFunctionDetails.functionDetails.NormalizedSecrets),
		readFunctionOrProcedureTargetPath(d, allFunctionDetails.functionDetails.NormalizedTargetPath),
		setOptionalFromStringPtr(d, "function_definition", allFunctionDetails.functionDetails.Body),
		d.Set("function_language", allFunctionDetails.functionDetails.Language),

		handleFunctionParameterRead(d, allFunctionDetails.functionParameters),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.FunctionToSchema(allFunctionDetails.function)}),
		d.Set(ParametersAttributeName, []map[string]any{schemas.FunctionParametersToSchema(allFunctionDetails.functionParameters)}),
	)
	if errs != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContextFunctionJava(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierWithArgumentsInSchema(id.SchemaId(), d.Get("name").(string), id.ArgumentDataTypes()...)

		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithRenameTo(newId.SchemaObjectId()))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming function %v err = %w", d.Id(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	// Batch SET operations and UNSET operations
	setRequest := sdk.NewFunctionSetRequest()
	unsetRequest := sdk.NewFunctionUnsetRequest()

	err = errors.Join(
		stringAttributeUpdate(d, "comment", &setRequest.Comment, &unsetRequest.Comment),
		func() error {
			if d.HasChange("secrets") {
				return setSecretsInBuilder(d, func(references []sdk.SecretReference) *sdk.FunctionSetRequest {
					return setRequest.WithSecretsList(sdk.SecretsListRequest{SecretsList: references})
				})
			}
			return nil
		}(),
		func() error {
			if d.HasChange("external_access_integrations") {
				return setExternalAccessIntegrationsInBuilder(d, func(references []sdk.AccountObjectIdentifier) any {
					if len(references) == 0 {
						return unsetRequest.WithExternalAccessIntegrations(true)
					} else {
						return setRequest.WithExternalAccessIntegrations(references)
					}
				})
			}
			return nil
		}(),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if updateParamDiags := handleFunctionParametersUpdate(d, setRequest, unsetRequest); len(updateParamDiags) > 0 {
		return updateParamDiags
	}

	// Apply SET and UNSET changes
	if !reflect.DeepEqual(*setRequest, *sdk.NewFunctionSetRequest()) {
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSet(*setRequest))
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unsetRequest, *sdk.NewFunctionUnsetRequest()) {
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnset(*unsetRequest))
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}

	// has to be handled separately
	if d.HasChange("is_secure") {
		if v := d.Get("is_secure").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			err = client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetSecure(parsed))
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		} else {
			err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetSecure(true))
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		}
	}

	return ReadContextFunctionJava(ctx, d, meta)
}