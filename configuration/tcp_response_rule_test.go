// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models/v2"
)

func TestGetTCPResponseRules(t *testing.T) {
	v, tRules, err := client.GetTCPResponseRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(tRules) != 2 {
		t.Errorf("%v tcp response rules returned, expected 2", len(tRules))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range tRules {
		if *r.Index == 0 {
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "accept" {
				t.Errorf("%v: Action not accept: %v", *r.Index, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not src TRUE: %v", *r.Index, r.CondTest)
			}
		} else if *r.Index == 1 {
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "reject" {
				t.Errorf("%v: Action not reject: %v", *r.Index, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not src FALSE: %v", *r.Index, r.CondTest)
			}
		} else {
			t.Errorf("Expext only tcp-response 0 or 1, %v found", *r.Index)
		}
	}

	_, tRules, err = client.GetTCPResponseRules("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(tRules) > 0 {
		t.Errorf("%v TCP Response Rules returned, expected 0", len(tRules))
	}
}

func TestGetTCPResponseRule(t *testing.T) {
	v, r, err := client.GetTCPResponseRule(0, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if r.Type != "content" {
		t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
	}
	if r.Action != "accept" {
		t.Errorf("%v: Action not accept: %v", *r.Index, r.Action)
	}
	if r.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
	}
	if r.CondTest != "TRUE" {
		t.Errorf("%v: CondTest not src TRUE: %v", *r.Index, r.CondTest)
	}

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetTCPResponseRule(3, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant TCP Response Rule")
	}
}

func TestCreateEditDeleteTCPResponseRule(t *testing.T) {
	id := int64(2)
	tOut := int64(1000)
	// TestCreateTCPResponseRule
	r := &models.TCPResponseRule{
		Index:   &id,
		Type:    "inspect-delay",
		Timeout: &tOut,
	}

	err := client.CreateTCPResponseRule("test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := client.GetTCPResponseRule(2, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created TCP response rule: %v\n", ondiskR)
		fmt.Printf("Given TCP response rule: %v\n", r)
		t.Error("Created TCP response rule not equal to given TCP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditTCPResponseRule
	r = &models.TCPResponseRule{
		Index:    &id,
		Type:     "content",
		Action:   "accept",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = client.EditTCPResponseRule(2, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = client.GetTCPResponseRule(2, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited TCP response rule: %v\n", ondiskR)
		fmt.Printf("Given TCP response rule: %v\n", r)
		t.Error("Edited TCP response rule not equal to given TCP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteTCPResponse
	err = client.DeleteTCPResponseRule(2, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetTCPResponseRule(2, "test", "")
	if err == nil {
		t.Error("DeleteTCPResponseRule failed, TCP Response Rule 3 still exists")
	}

	err = client.DeleteTCPResponseRule(2, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant TCP Response Rule")
		version++
	}
}
