package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func truncate(list []basetypes.ObjectValue, maximum int64) []basetypes.ObjectValue {
	if int64(len(list)) > maximum {
		return list[(int64(len(list)) - maximum):]
	}
	return list
}

func appendAndTruncate(list []basetypes.ObjectValue, item basetypes.ObjectValue, maximum int64) []basetypes.ObjectValue {
	return truncate(append(list, item), maximum)
}

func mapAttributeIsEqual(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, attribute string) bool {
	var state types.Map
	var plan types.Map
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root(attribute), &plan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(attribute), &state)...)

	return state.Equal(plan)
}
