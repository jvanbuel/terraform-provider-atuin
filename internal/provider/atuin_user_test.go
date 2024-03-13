package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleAtuinUserResource(t *testing.T) {
	random_user := fmt.Sprintf("testUser%d", rand.Intn(1000))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleAtuinUserResourceConfig(random_user, "pa$$word", "test1234@yahoo.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("atuin_user.test", "password", "pa$$word"),
				),
			},
			// ImportState testing
			{
				ResourceName:  "atuin_user.test",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s,pa$$word,indoor dish desk flag debris potato excuse depart ticket judge file exit", random_user),
			},
			// Update and Read testing
			// {
			// 	Config: testAccExampleAtuinUserResourceConfig("testUser523018", "pa$$word2", "test1234$yahoo.com"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("atuin_user.test", "password", "pa$$word2"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleAtuinUserResourceConfig(username, password, email string) string {
	return fmt.Sprintf(`
resource "atuin_user" "test" {
  username = %[1]q
  password = %[2]q
  email    = %[3]q
}
`, username, password, email)
}
