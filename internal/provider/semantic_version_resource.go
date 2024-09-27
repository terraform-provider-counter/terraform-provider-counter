package provider

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SemanticVersionResource{}
var _ resource.ResourceWithModifyPlan = &SemanticVersionResource{}

func NewSemanticVersionResource() resource.Resource {
	return &SemanticVersionResource{}
}

type SemanticVersionResource struct {
}

func (s SemanticVersionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_semantic_version"
}

func (s SemanticVersionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A semantic version number whose components increment according to the configured triggers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"major_value": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The current major version number.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"minor_value": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The current minor version number.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"patch_value": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The current patch version number.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"value": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The semantic version number as a string in `<major>.<minor>.<patch>` form.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_history": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(1000),
				MarkdownDescription: "Maximum number of versions this resource should store in the `history` attribute.",
			},
			"history": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "A list of semantic versions that this resource has produced.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Computed: true,
						},
						"major_value": schema.Int64Attribute{
							Computed: true,
						},
						"minor_value": schema.Int64Attribute{
							Computed: true,
						},
						"patch_value": schema.Int64Attribute{
							Computed: true,
						},
						"major_triggers": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"minor_triggers": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"patch_triggers": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"major_initial_value": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "The initial major version value.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"minor_initial_value": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "The initial minor version value.",
				Default:             int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"patch_initial_value": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "The initial patch version value.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"major_triggers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A map of strings that will cause the major version number to increment when any of the values change.",
			},
			"minor_triggers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A map of strings that will cause the minor version number to increment when any of the values change.",
			},
			"patch_triggers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A map of strings that will cause the patch version number to increment when any of the values change.",
			},
		},
	}
}

func (s SemanticVersionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data semanticVersionModelV1
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = types.StringValue(uuid.New().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (s SemanticVersionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.State = req.State
}

func (s SemanticVersionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data semanticVersionModelV1
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (s SemanticVersionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (s SemanticVersionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the entire plan is null, the resource is planned for destruction.
	if req.Plan.Raw.IsNull() {
		return
	}

	var majorValue types.Int64
	var minorValue types.Int64
	var patchValue types.Int64
	var maxHistory types.Int64
	var majorTriggers types.Map
	var minorTriggers types.Map
	var patchTriggers types.Map
	creation := req.State.Raw.IsNull()

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("major_triggers"), &majorTriggers)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("minor_triggers"), &minorTriggers)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("patch_triggers"), &patchTriggers)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("max_history"), &maxHistory)...)

	if creation {
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("major_initial_value"), &majorValue)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("minor_initial_value"), &minorValue)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("patch_initial_value"), &patchValue)...)

		value := types.StringValue(fmt.Sprintf("%d.%d.%d", majorValue.ValueInt64(), minorValue.ValueInt64(), patchValue.ValueInt64()))
		history := appendAndTruncate([]basetypes.ObjectValue{}, s.createHistoryEntry(value, majorValue, majorTriggers, minorValue, minorTriggers, patchValue, patchTriggers), maxHistory.ValueInt64())

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("value"), value)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("major_value"), majorValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("minor_value"), minorValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("patch_value"), patchValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("history"), history)...)
		return
	}

	if !mapAttributeIsEqual(ctx, req, resp, "major_triggers") {
		var history []basetypes.ObjectValue
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("major_value"), &majorValue)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("history"), &history)...)

		majorValue = types.Int64Value(majorValue.ValueInt64() + 1)
		minorValue = types.Int64Value(0)
		patchValue = types.Int64Value(0)
		value := types.StringValue(fmt.Sprintf("%d.%d.%d", majorValue.ValueInt64(), minorValue.ValueInt64(), patchValue.ValueInt64()))
		history = appendAndTruncate(history, s.createHistoryEntry(value, majorValue, majorTriggers, minorValue, minorTriggers, patchValue, patchTriggers), maxHistory.ValueInt64())

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("value"), value)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("major_value"), majorValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("minor_value"), minorValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("patch_value"), patchValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("history"), history)...)
		return
	}

	if !mapAttributeIsEqual(ctx, req, resp, "minor_triggers") {
		var history []basetypes.ObjectValue
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("major_value"), &majorValue)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("minor_value"), &minorValue)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("history"), &history)...)

		minorValue = types.Int64Value(minorValue.ValueInt64() + 1)
		patchValue = types.Int64Value(0)
		value := types.StringValue(fmt.Sprintf("%d.%d.%d", majorValue.ValueInt64(), minorValue.ValueInt64(), patchValue.ValueInt64()))
		history = appendAndTruncate(history, s.createHistoryEntry(value, majorValue, majorTriggers, minorValue, minorTriggers, patchValue, patchTriggers), maxHistory.ValueInt64())

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("value"), value)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("major_value"), majorValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("minor_value"), minorValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("patch_value"), patchValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("history"), history)...)
		return
	}

	if !mapAttributeIsEqual(ctx, req, resp, "patch_triggers") {
		var history []basetypes.ObjectValue
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("major_value"), &majorValue)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("minor_value"), &minorValue)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("patch_value"), &patchValue)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("history"), &history)...)

		patchValue = types.Int64Value(patchValue.ValueInt64() + 1)
		value := types.StringValue(fmt.Sprintf("%d.%d.%d", majorValue.ValueInt64(), minorValue.ValueInt64(), patchValue.ValueInt64()))
		history = appendAndTruncate(history, s.createHistoryEntry(value, majorValue, majorTriggers, minorValue, minorTriggers, patchValue, patchTriggers), maxHistory.ValueInt64())

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("value"), value)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("major_value"), majorValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("minor_value"), minorValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("patch_value"), patchValue)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("history"), history)...)
	}
}

func (s SemanticVersionResource) createHistoryEntry(value types.String, majorValue types.Int64, majorTriggers types.Map, minorValue types.Int64, minorTriggers types.Map, patchValue types.Int64, patchTriggers types.Map) basetypes.ObjectValue {
	return types.ObjectValueMust(
		map[string]attr.Type{
			"value":          types.StringType,
			"major_value":    types.Int64Type,
			"minor_value":    types.Int64Type,
			"patch_value":    types.Int64Type,
			"major_triggers": types.MapType{ElemType: types.StringType},
			"minor_triggers": types.MapType{ElemType: types.StringType},
			"patch_triggers": types.MapType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"value":          value,
			"major_value":    majorValue,
			"minor_value":    minorValue,
			"patch_value":    patchValue,
			"major_triggers": majorTriggers,
			"minor_triggers": minorTriggers,
			"patch_triggers": patchTriggers,
		},
	)
}

type semanticVersionModelV1 struct {
	Id                types.String            `tfsdk:"id"`
	MajorValue        types.Int64             `tfsdk:"major_value"`
	MinorValue        types.Int64             `tfsdk:"minor_value"`
	PatchValue        types.Int64             `tfsdk:"patch_value"`
	Value             types.String            `tfsdk:"value"`
	MaxHistory        types.Int64             `tfsdk:"max_history"`
	History           []basetypes.ObjectValue `tfsdk:"history"`
	MajorInitialValue types.Int64             `tfsdk:"major_initial_value"`
	MinorInitialValue types.Int64             `tfsdk:"minor_initial_value"`
	PatchInitialValue types.Int64             `tfsdk:"patch_initial_value"`
	MajorTriggers     types.Map               `tfsdk:"major_triggers"`
	MinorTriggers     types.Map               `tfsdk:"minor_triggers"`
	PatchTriggers     types.Map               `tfsdk:"patch_triggers"`
}
