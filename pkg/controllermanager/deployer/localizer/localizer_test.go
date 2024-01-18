/*
Copyright 2024 The Clusternet Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package localizer

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsapi "github.com/clusternet/clusternet/pkg/apis/apps/v1alpha1"
)

func createTempLocalization(
	name string,
	creationTimestamp metav1.Time,
	priority int32,
	privileged bool,
) *appsapi.Localization {
	return &appsapi.Localization{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         "demo",
			CreationTimestamp: creationTimestamp,
		},
		Spec: appsapi.LocalizationSpec{
			Priority:   priority,
			Privileged: privileged,
		},
	}
}

func TestCompareLocalizationPriority(t *testing.T) {
	demoTime := metav1.Now()

	tests := []struct {
		name string
		i    *appsapi.Localization
		j    *appsapi.Localization
		want bool
	}{
		{
			name: "newer localization with higher priority - 1",
			i:    createTempLocalization("foo", demoTime, 500, false),
			j:    createTempLocalization("bar", metav1.NewTime(demoTime.Add(time.Second)), 600, false),
			want: true,
		},
		{
			name: "newer localization with higher priority -2 ",
			i:    createTempLocalization("foo", metav1.NewTime(demoTime.Add(time.Second)), 600, false),
			j:    createTempLocalization("bar", demoTime, 500, false),
			want: false,
		},
		{
			name: "newer localization with lower priority - 1",
			i:    createTempLocalization("foo", demoTime, 600, false),
			j:    createTempLocalization("bar", metav1.NewTime(demoTime.Add(time.Second)), 500, false),
			want: false,
		},
		{
			name: "newer localization with lower priority - 2",
			i:    createTempLocalization("foo", metav1.NewTime(demoTime.Add(time.Second)), 500, false),
			j:    createTempLocalization("bar", demoTime, 600, false),
			want: true,
		},

		{
			name: "localization with same priority, no privileged - 1",
			i:    createTempLocalization("foo", demoTime, 600, false),
			j:    createTempLocalization("bar", metav1.NewTime(demoTime.Add(time.Second)), 600, false),
			want: true,
		},
		{
			name: "localization with same priority, no privileged - 2",
			i:    createTempLocalization("foo", metav1.NewTime(demoTime.Add(time.Second)), 600, false),
			j:    createTempLocalization("bar", demoTime, 600, false),
			want: false,
		},

		{
			name: "newer localization with lower priority, but privileged",
			i:    createTempLocalization("foo", demoTime, 600, false),
			j:    createTempLocalization("bar", metav1.NewTime(demoTime.Add(time.Second)), 500, true),
			want: false,
		},
		{
			name: "newer localization with higher priority, but privileged",
			i:    createTempLocalization("foo", metav1.NewTime(demoTime.Add(time.Second)), 600, true),
			j:    createTempLocalization("bar", demoTime, 500, false),
			want: false,
		},

		{
			name: "newer localization with same priority, but privileged - 1",
			i:    createTempLocalization("foo", demoTime, 600, false),
			j:    createTempLocalization("bar", metav1.NewTime(demoTime.Add(time.Second)), 600, true),
			want: true,
		},
		{
			name: "newer localization with same priority, but privileged - 2",
			i:    createTempLocalization("foo", metav1.NewTime(demoTime.Add(time.Second)), 600, true),
			j:    createTempLocalization("bar", demoTime, 600, false),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareLocalizationPriority(tt.i, tt.j); got != tt.want {
				t.Errorf("compareLocalizationPriority() = %v, want %v", got, tt.want)
			}
		})
	}
}
