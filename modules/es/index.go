// Copyright © 2020 Uzhinskiy Boris
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package es

import (
	"context"
	"fmt"
	"regexp"

	"github.com/olivere/elastic/v7"
)

type indexTemplate struct {
	IndexPatterns []string `json:"index_patterns"`
	Settings      struct {
		Shards                      int    `json:"number_of_shards"`
		Replicas                    int    `json:"number_of_replicas"`
		IndexLifecycleName          string `json:"index.lifecycle.name"`
		IndexLifecycleRolloverAlias string `json:"index.lifecycle.rollover_alias"`
	} `json:"settings"`
	Mappings struct {
		Properties struct {
			Value struct {
				Type string `json:"type"`
			} `json:"value"`
			Timestamp struct {
				Type string `json:"type"`
			} `json:"timestamp"`
			Datetime struct {
				Type string `json:"type"`
			} `json:"datetime"`
		} `json:"properties"`
	} `json:"mappings"`
}

func indextemplate(index string, shards, replicas int) indexTemplate {
	something := indexTemplate{}

	something.IndexPatterns = []string{fmt.Sprintf("%s-*", index)}
	something.Settings.Shards = shards
	something.Settings.Replicas = replicas
	something.Settings.IndexLifecycleName = fmt.Sprintf("%s-ilm-policy", index)
	something.Settings.IndexLifecycleRolloverAlias = index
	something.Mappings.Properties.Value.Type = "long"
	something.Mappings.Properties.Timestamp.Type = "long"
	something.Mappings.Properties.Datetime.Type = "date"

	return something
}

type ilmPolicy struct {
	Policy struct {
		Phases struct {
			Hot struct {
				MinAge  string `json:"min_age"`
				Actions struct {
					Rollover struct {
						MaxAge string `json:"max_age"`
					} `json:"rollover"`
				} `json:"actions"`
			} `json:"hot"`
			Cold struct {
				MinAge  string `json:"min_age"`
				Actions struct {
					Allocate struct {
						NumberOfReplicas int `json:"number_of_replicas"`
						Include          struct {
						} `json:"include"`
						Exclude struct {
						} `json:"exclude"`
						Require struct {
							BoxType string `json:"box_type"`
						} `json:"require"`
					} `json:"allocate"`
				} `json:"actions"`
			} `json:"cold"`
			Delete struct {
				MinAge  string `json:"min_age"`
				Actions struct {
					Delete struct {
					} `json:"delete"`
				} `json:"actions"`
			} `json:"delete"`
		} `json:"phases"`
	} `json:"policy"`
}

func ilmpolicy(index, hot, warm, cold string) ilmPolicy {
	something := ilmPolicy{}
	something.Policy.Phases.Hot.MinAge = "0ms"
	something.Policy.Phases.Hot.Actions.Rollover.MaxAge = hot
	something.Policy.Phases.Cold.MinAge = hot
	something.Policy.Phases.Cold.Actions.Allocate.NumberOfReplicas = 0
	something.Policy.Phases.Cold.Actions.Allocate.Require.BoxType = "cold"
	something.Policy.Phases.Delete.MinAge = cold
	return something
}

func (esc *ESClient) SetupIndex(c esconf) error {

	indices, err := esc.ec.IndexNames()
	if err != nil {
		return err
	}
	exists := grepIndexName(indices, c.Index)
	if !exists {
		fmt.Println("First start!")
		ilm := ilmpolicy(c.Index, c.Ilm.Hot, c.Ilm.Warm, c.Ilm.Cold)

		ilmservice := elastic.NewXPackIlmPutLifecycleService(esc.ec)

		policy_create, err := ilmservice.Policy(fmt.Sprintf("%s-ilm-policy", c.Index)).
			BodyJson(ilm).
			Do(context.Background())
		if err != nil {
			return err
		}
		if policy_create.Acknowledged {
			nit := indextemplate(c.Index, c.Shards, c.Replicas)
			templservice := elastic.NewIndicesPutTemplateService(esc.ec)
			templ_create, err := templservice.Name(fmt.Sprintf("%s-template", c.Index)).
				BodyJson(nit).
				Do(context.Background())
			if err != nil {
				return err
			}
			if templ_create.Acknowledged {

				mapping := `{ "aliases": { "` + c.Index + `": { "is_write_index": true } } }`
				_, err := esc.ec.CreateIndex(fmt.Sprintf("%s-000001", c.Index)).BodyString(mapping).Do(context.Background())
				if err != nil {
					return err
				}

			}

		}
	}

	return nil
}

func grepIndexName(indices []string, index string) bool {
	ret := false
	re := regexp.MustCompile(fmt.Sprintf(`^((%s)|(%s-{1}\w+)|(%s\w+))$`, index, index, index))
	for _, name := range indices {
		if len(re.FindStringIndex(name)) > 0 {
			ret = true
			break
		} else {
			ret = false
		}
	}
	return ret
}
