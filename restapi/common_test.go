package restapi

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

func testAccCheckRestapiObjectExists(n string, id string, client *api_client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			keys := make([]string, 0, len(s.RootModule().Resources))
			for k := range s.RootModule().Resources {
				keys = append(keys, k)
			}
			return fmt.Errorf("RestAPI object not found in terraform state: %s. Found: %s", n, strings.Join(keys, ", "))
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("RestAPI object id not set in terraform")
		}

		/* Make a throw-away API object to read from the API */
		path := "/api/objects"
		obj, err := NewAPIObject(
			client,
			path+"/{id}",
			path,
			path+"/{id}",
			path+"/{id}",
			id,
			"id",
			"{}",
			true,
		)
		if err != nil {
			return err
		}

		err = obj.read_object()
		if err != nil {
			return err
		}

		return nil
	}
}

func TestGetStringAtKey(t *testing.T) {
	debug := false
	test_obj := make(map[string]interface{})
	err := json.Unmarshal([]byte(`
    {
      "rootFoo": "bar",
      "top": {
        "foo": "bar",
        "number": 1,
        "middle": {
          "bottom": {
            "foo": "bar"
          }
        },
        "list": [
          "bar",
          "baz"
        ]
      }
    }
  `), &test_obj)
	if nil != err {
		t.Fatalf("Error unmarshalling JSON: %s", err)
	}

	var res string

	res, err = GetStringAtKey(test_obj, "rootFoo", debug)
	if err != nil {
		t.Fatalf("Error extracting 'rootFoo' from JSON payload: %s", err)
	} else if "bar" != res {
		t.Fatalf("Error: Expected 'bar', but got %s", res)
	}

	res, err = GetStringAtKey(test_obj, "top/foo", debug)
	if err != nil {
		t.Fatalf("Error extracting 'top/foo' from JSON payload: %s", err)
	} else if "bar" != res {
		t.Fatalf("Error: Expected 'bar', but got %s", res)
	}

	res, err = GetStringAtKey(test_obj, "top/middle/bottom/foo", debug)
	if err != nil {
		t.Fatalf("Error extracting top/foo from JSON payload: %s", err)
	} else if "bar" != res {
		t.Fatalf("Error: Expected 'bar', but got %s", res)
	}

	res, err = GetStringAtKey(test_obj, "top/middle/junk", debug)
	if err == nil {
		t.Fatalf("Error expected when trying to extract 'top/middle/junk' from payload")
	}

	res, err = GetStringAtKey(test_obj, "top/number", debug)
	if err != nil {
		t.Fatalf("Error extracting 'top/number' from JSON payload: %s", err)
	} else if "1" != res {
		t.Fatalf("Error: Expected '1', but got %s", res)
	}
}

func TestGetListStringAtKey(t *testing.T) {
	debug := false
	test_obj := make(map[string]interface{})
	err := json.Unmarshal([]byte(`
    {
      "rootFoo": "bar",
      "items": [
        {
          "foo": "bar",
          "number": 1,
          "list_numbers": [1, 2, 3],
          "test": [{"id": "3333"}, {"id": "1337"}],
          "resource": {
            "id": "123"
          }
        }
      ]
    }
  `), &test_obj)
	if nil != err {
		t.Fatalf("Error unmarshalling JSON: %s", err)
	}

	var res string

	res, err = GetStringAtKey(test_obj, "items/0/resource/id", debug)
	if err != nil {
		t.Fatalf("Error extracting 'resource' from JSON payload: %s", err)
	} else if "123" != res {
		t.Fatalf("Error: Expected '123', but got %s", res)
	}

	res, err = GetStringAtKey(test_obj, "items/0/test/1/id", debug)
	if err != nil {
		t.Fatalf("Error extracting 'resource' from JSON payload: %s", err)
	} else if "1337" != res {
		t.Fatalf("Error: Expected '1337', but got %s", res)
	}

	res, err = GetStringAtKey(test_obj, "items/0/list_numbers/1", debug)
	if err != nil {
		t.Fatalf("Error extracting 'resource' from JSON payload: %s", err)
	} else if "2" != res {
		t.Fatalf("Error: Expected '2', but got %s", res)
	}
}
