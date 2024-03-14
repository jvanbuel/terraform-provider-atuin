package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleAtuinUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleAtuinUserResourceConfig("rincewind", "pa$$word", "test1234@yahoo.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("atuin_user.test", "password", "pa$$word"),
				),
			},
			// ImportState testing
			{
				ResourceName:  "atuin_user.test",
				ImportState:   true,
				ImportStateId: "rincewind,pa$$word,indoor dish desk flag debris potato excuse depart ticket judge file exit",
			},
			// Update and Read testing
			{
				Config: testAccExampleAtuinUserResourceConfig("twoflower", "pa$$word2", "test1234$yahoo.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("atuin_user.test", "password", "pa$$word2"),
				),
			},
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
