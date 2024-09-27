package provider

import (
	"context"
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
var _ resource.Resource = &MonotonicResource{}
var _ resource.ResourceWithModifyPlan = &MonotonicResource{}

func NewMonotonicResource() resource.Resource {
	return &MonotonicResource{}
}

type MonotonicResource struct {
}

func (m MonotonicResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monotonic"
}

func (m MonotonicResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A monotonic counter which increments according to the configured triggers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The current value of the counter.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"step": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "The amount used to increment / decrement the counter on each revision.",
			},
			"max_history": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(1000),
				MarkdownDescription: "Maximum number of versions this resource should store in the `history` attribute.",
			},
			"history": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "A list of counter values that this resource has produced.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.Int64Attribute{
							Computed: true,
						},
						"triggers": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"initial_value": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "The initial value of the counter.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"triggers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A map of strings that will cause a change to the counter when any of the values change.",
			},
		},
	}
}

func (m MonotonicResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data monotonicModelV1
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = types.StringValue(uuid.New().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m MonotonicResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.State = req.State
}

func (m MonotonicResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data monotonicModelV1
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m MonotonicResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (m MonotonicResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the entire plan is null, the resource is planned for destruction.
	if req.Plan.Raw.IsNull() {
		return
	}

	var value types.Int64
	var maxHistory types.Int64
	var triggers types.Map
	creation := req.State.Raw.IsNull()

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("triggers"), &triggers)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("max_history"), &maxHistory)...)

	if creation {
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("initial_value"), &value)...)
		history := appendAndTruncate([]basetypes.ObjectValue{}, m.createHistoryEntry(value, triggers), maxHistory.ValueInt64())

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("value"), value)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("history"), history)...)
		return
	}

	if !mapAttributeIsEqual(ctx, req, resp, "triggers") {
		var step types.Int64
		var history []basetypes.ObjectValue
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("value"), &value)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("step"), &step)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("history"), &history)...)

		value = types.Int64Value(value.ValueInt64() + step.ValueInt64())
		history = appendAndTruncate(history, m.createHistoryEntry(value, triggers), maxHistory.ValueInt64())
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("value"), value)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("history"), history)...)
	}
}

func (m MonotonicResource) createHistoryEntry(value types.Int64, triggers types.Map) basetypes.ObjectValue {
	return types.ObjectValueMust(
		map[string]attr.Type{
			"value":    types.Int64Type,
			"triggers": types.MapType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"value":    value,
			"triggers": triggers,
		},
	)
}

type monotonicModelV1 struct {
	Id           types.String            `tfsdk:"id"`
	Value        types.Int64             `tfsdk:"value"`
	Step         types.Int64             `tfsdk:"step"`
	MaxHistory   types.Int64             `tfsdk:"max_history"`
	History      []basetypes.ObjectValue `tfsdk:"history"`
	InitialValue types.Int64             `tfsdk:"initial_value"`
	Triggers     types.Map               `tfsdk:"triggers"`
}
