package njalla

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/Sighery/gonjalla"
)

func TestAccRecordANAME_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordANAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordANAMECreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordANAMEExists(
						"njalla_record_aname.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_create",
						"name",
						"testacc1-aname-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_create",
						"content",
						"testacc1-aname-create-content.com",
					),
				),
			},
		},
	})
}

func TestAccRecordANAME_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordANAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordANAMEUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordANAMEExists(
						"njalla_record_aname.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_update",
						"name",
						"testacc2-aname-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_update",
						"content",
						"testacc2-aname-update-content1.com",
					),
				),
			},
			{
				Config: testAccCheckRecordANAMEUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordANAMEExists(
						"njalla_record_aname.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_update",
						"name",
						"testacc2-aname-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aname.test_update",
						"content",
						"testacc2-aname-update-content2.com",
					),
				),
			},
		},
	})
}

func TestAccRecordANAME_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordANAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordANAMEImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordANAMEExists(
						"njalla_record_ANAME.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_ANAME.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordANAME_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordANAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordANAMEEmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordANAMEExists(
						"njalla_record_ANAME.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ANAME.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ANAME.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ANAME.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ANAME.test_empty_name",
						"content",
						"testacc4-ANAME-emptyname-content.com",
					),
				),
			},
		},
	})
}

func TestAccRecordANAME_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordANAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordANAMEInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordANAMEDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_ANAME" {
			continue
		}

		records, err := gonjalla.ListRecords(config.Token, domain)
		if err != nil {
			return fmt.Errorf(
				"Error fetching the records data for domain %s: %s",
				domain, err,
			)
		}

		for _, record := range records {
			if record.ID == rs.Primary.ID {
				return fmt.Errorf(
					"Record %s still exists in domain %s",
					rs.Primary.ID, domain,
				)
			}
		}
	}

	return nil
}

func testAccCheckRecordANAMEExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
		records, err := gonjalla.ListRecords(config.Token, domain)
		if err != nil {
			return fmt.Errorf(
				"Error fetching the records data for domain %s: %s",
				domain, err,
			)
		}

		for _, record := range records {
			if record.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf(
			"Record %s doesn't exist for domain %s", rs.Primary.ID, domain,
		)
	}
}

func testAccCheckRecordANAMECreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ANAME test_create {
  domain = %q
  name = "testacc1-ANAME-create-name"
  ttl = 10800
  content = "testacc1-ANAME-create-content.com"
}
`, domain)
}

func testAccCheckRecordANAMEUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ANAME test_update {
  domain = %q
  name = "testacc2-ANAME-update-name1"
  ttl = 10800
  content = "testacc2-ANAME-update-content1.com"
}
`, domain)
}

func testAccCheckRecordANAMEUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ANAME test_update {
  domain = %q
  name = "testacc2-ANAME-update-name2"
  ttl = 3600
  content = "testacc2-ANAME-update-content2.com"
}
`, domain)
}

func testAccCheckRecordANAMEImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ANAME test_import {
  domain = %q
  name = "testacc3-ANAME-import-name"
  ttl = 10800
  content = "testacc3-ANAME-import-content.com"
}
`, domain)
}

func testAccCheckRecordANAMEEmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ANAME test_empty_name {
  domain = %q
  ttl = 10800
  content = "testacc4-ANAME-emptyname-content.com"
}
`, domain)
}

func testAccCheckRecordANAMEInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ANAME test_invalid_ttl {
  domain = %q
  name = "testacc5-ANAME-invalidttl-name"
  ttl = 999
  content = "testacc5-ANAME-invalidttl-content.com"
}
`, domain)
}
