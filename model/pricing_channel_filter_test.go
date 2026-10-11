package model

import (
	"done-hub/common/logger"
	"testing"
)

func TestFilterPricesByChannelModels(t *testing.T) {
	logger.SetupLogger()
	origRule := ChannelGroup.Rule
	defer func() { ChannelGroup.Rule = origRule }()

	ChannelGroup.Rule = map[string]map[string][][]int{
		"default": {"gpt-5": nil, "claude-sonnet-4-5": nil},
		"vip":     {"gpt-5": nil, "deepseek-v3.2": nil},
	}

	got := ChannelGroup.GetAllModels()
	want := []string{"claude-sonnet-4-5", "deepseek-v3.2", "gpt-5"}
	if len(got) != len(want) {
		t.Fatalf("GetAllModels = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("GetAllModels = %v, want %v", got, want)
		}
	}

	prices := []*Price{
		{Model: "gpt-5"},
		{Model: "GPT-5"}, // 大小写不同视为不同模型，不匹配
		{Model: "babbage-002"},
		{Model: "claude-sonnet-4-5"},
		{Model: "claude-sonnet-4-5*"},
		nil,
	}
	filtered := FilterPricesByChannelModels(prices)
	if len(filtered) != 2 || filtered[0].Model != "gpt-5" || filtered[1].Model != "claude-sonnet-4-5" {
		t.Fatalf("unexpected filtered result: %+v", filtered)
	}
}
