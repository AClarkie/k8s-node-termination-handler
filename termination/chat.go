// Copyright 2017 Google Inc. All Rights Reserved.
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

package termination

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/golang/glog"
)

func sendChat() error {
	url := os.Getenv("CHAT_WEBHOOK_URL")
	if url == "" {
		return nil
	}

	glog.V(4).Infof("Posting to Google Chat webhook")

	instanceName, err := metadata.InstanceName()
	if err != nil {
		return err
	}
	zone, err := metadata.Zone()
	if err != nil {
		return err
	}
	projectID, err := metadata.ProjectID()
	if err != nil {
		return err
	}
	machineType, err := metadata.Get(machineTypeSuffix)
	if err != nil {
		return err
	}

	payload := []byte(fmt.Sprintf(`{"text":"
*InstanceName*:	%s
*MachineType*:	%s
*Zone*: 			%s
*ProjectID*: 		%s
	"}`, instanceName, zone, projectID, machineType))

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
