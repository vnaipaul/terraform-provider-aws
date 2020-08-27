package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/tfawsresource"
)

func testAccAwsGuardDutyFilter_basic(t *testing.T) {
	var v1, v2 guardduty.GetFilterOutput
	resourceName := "aws_guardduty_filter.test"
	detectorResourceName := "aws_guardduty_detector.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsGuardDutyFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardDutyFilterConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsGuardDutyFilterExists(resourceName, &v1),
					resource.TestCheckResourceAttrPair(resourceName, "detector_id", detectorResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "test-filter"),
					resource.TestCheckResourceAttr(resourceName, "action", "ARCHIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "rank", "1"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "guardduty", regexp.MustCompile("detector/[a-z0-9]{32}/filter/test-filter$")),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "finding_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "finding_criteria.0.criterion.#", "4"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "finding_criteria.0.criterion.*", map[string]string{
						"field":     "region",
						"values.#":  "1",
						"values.0":  "eu-west-1",
						"condition": "equals",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "finding_criteria.0.criterion.*", map[string]string{
						"field":     "service.additionalInfo.threatListName",
						"values.#":  "2",
						"values.0":  "some-threat",
						"values.1":  "another-threat",
						"condition": "not_equals",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "finding_criteria.0.criterion.*", map[string]string{
						"field":     "updatedAt",
						"values.#":  "1",
						"values.0":  "1570744740000",
						"condition": "less_than",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "finding_criteria.0.criterion.*", map[string]string{
						"field":     "updatedAt",
						"values.#":  "1",
						"values.0":  "1570744240000",
						"condition": "greater_than",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGuardDutyFilterConfigNoop_full(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsGuardDutyFilterExists(resourceName, &v2),
					resource.TestCheckResourceAttrPair(resourceName, "detector_id", detectorResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "test-filter"),
					resource.TestCheckResourceAttr(resourceName, "action", "NOOP"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a NOOP"),
					resource.TestCheckResourceAttr(resourceName, "rank", "1"),
					resource.TestCheckResourceAttr(resourceName, "finding_criteria.#", "1"),
				),
			},
		},
	})
}

func testAccAwsGuardDutyFilter_update(t *testing.T) {
	var v1, v2 guardduty.GetFilterOutput
	resourceName := "aws_guardduty_filter.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsGuardDutyFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardDutyFilterConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsGuardDutyFilterExists(resourceName, &v1),
					resource.TestCheckResourceAttr(resourceName, "finding_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "finding_criteria.0.criterion.#", "4"),
				),
			},
			{
				Config: testAccGuardDutyFilterConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsGuardDutyFilterExists(resourceName, &v2),
					resource.TestCheckResourceAttr(resourceName, "finding_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "finding_criteria.0.criterion.#", "2"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "finding_criteria.0.criterion.*", map[string]string{
						"field":     "region",
						"values.#":  "1",
						"values.0":  "us-west-2",
						"condition": "equals",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "finding_criteria.0.criterion.*", map[string]string{
						"field":     "service.additionalInfo.threatListName",
						"values.#":  "2",
						"values.0":  "some-threat",
						"values.1":  "yet-another-threat",
						"condition": "not_equals",
					}),
				),
			},
		},
	})
}

func testAccAwsGuardDutyFilter_tags(t *testing.T) {
	var v1, v2, v3 guardduty.GetFilterOutput
	resourceName := "aws_guardduty_filter.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsGuardDutyFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardDutyFilterConfig_multipleTags(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsGuardDutyFilterExists(resourceName, &v1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "test-filter"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key", "Value"),
				),
			},
			{
				Config: testAccGuardDutyFilterConfig_updateTags(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsGuardDutyFilterExists(resourceName, &v2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key", "Updated"),
				),
			},
			{
				Config: testAccGuardDutyFilterConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsGuardDutyFilterExists(resourceName, &v3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccAwsGuardDutyFilter_disappears(t *testing.T) {
	var v guardduty.GetFilterOutput
	resourceName := "aws_guardduty_filter.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsAcmpcaCertificateAuthorityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardDutyFilterConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsGuardDutyFilterExists(resourceName, &v),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsGuardDutyFilter(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAwsGuardDutyFilterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).guarddutyconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_guardduty_filter" {
			continue
		}

		detectorID, filterName, err := parseImportedId(rs.Primary.ID)
		if err != nil {
			return err
		}

		input := &guardduty.GetFilterInput{
			DetectorId: aws.String(detectorID),
			FilterName: aws.String(filterName),
		}

		_, err = conn.GetFilter(input)
		if err != nil {
			if isAWSErr(err, guardduty.ErrCodeBadRequestException, "The request is rejected because the input detectorId is not owned by the current account.") {
				return nil
			}
			return err
		}

		return fmt.Errorf("Expected GuardDuty Filter to be destroyed, %s found", rs.Primary.Attributes["filter_name"])
	}

	return nil
}

func testAccCheckAwsGuardDutyFilterExists(name string, filter *guardduty.GetFilterOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No GuardDuty filter is set")
		}

		detectorID, name, err := parseImportedId(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := testAccProvider.Meta().(*AWSClient).guarddutyconn
		input := guardduty.GetFilterInput{
			DetectorId: aws.String(detectorID),
			FilterName: aws.String(name),
		}
		filter, err = conn.GetFilter(&input)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccGuardDutyFilterConfig_full() string {
	return `
resource "aws_guardduty_filter" "test" {
  detector_id = "${aws_guardduty_detector.test.id}"
	name        = "test-filter"
	action      = "ARCHIVE"
	rank        = 1

  finding_criteria {
    criterion {
      field     = "region"
      values    = ["eu-west-1"]
      condition = "equals"
    }

    criterion {
      field     = "service.additionalInfo.threatListName"
      values    = ["some-threat", "another-threat"]
      condition = "not_equals"
    }

    criterion {
      field     = "updatedAt"
      values    = ["1570744740000"]
      condition = "less_than"
    }

    criterion {
      field     = "updatedAt"
      values    = ["1570744240000"]
      condition = "greater_than"
    }
  }
}

resource "aws_guardduty_detector" "test" {
  enable = true
}`
}

func testAccGuardDutyFilterConfigNoop_full() string {
	return `
resource "aws_guardduty_filter" "test" {
  detector_id = "${aws_guardduty_detector.test.id}"
	name        = "test-filter"
	action      = "NOOP"
	description = "This is a NOOP"
	rank        = 1

  finding_criteria {
    criterion {
      field     = "region"
      values    = ["eu-west-1"]
      condition = "equals"
    }

    criterion {
      field     = "service.additionalInfo.threatListName"
      values    = ["some-threat", "another-threat"]
      condition = "not_equals"
    }

    criterion {
      field     = "updatedAt"
      values    = ["1570744740000"]
      condition = "less_than"
    }

    criterion {
      field     = "updatedAt"
      values    = ["1570744240000"]
      condition = "greater_than"
    }
  }
}

resource "aws_guardduty_detector" "test" {
  enable = true
}`
}

func testAccGuardDutyFilterConfig_multipleTags() string {
	return `
resource "aws_guardduty_filter" "test" {
  detector_id = "${aws_guardduty_detector.test.id}"
	name        = "test-filter"
	action      = "ARCHIVE"
	rank        = 1

  finding_criteria {
    criterion {
		field     = "region"
		values    = ["us-west-2"]
		condition = "equals"
	  }
	}

  tags = {
	  Name= "test-filter"
	  Key = "Value"
  }
}

resource "aws_guardduty_detector" "test" {
  enable = true
}`
}

func testAccGuardDutyFilterConfig_update() string {
	return `
resource "aws_guardduty_filter" "test" {
  detector_id = "${aws_guardduty_detector.test.id}"
	name        = "test-filter"
	action      = "ARCHIVE"
	rank        = 1

  finding_criteria {
    criterion {
      field     = "region"
      values    = ["us-west-2"]
      condition = "equals"
    }

    criterion {
      field     = "service.additionalInfo.threatListName"
      values    = ["some-threat", "yet-another-threat"]
      condition = "not_equals"
    }
  }
}

resource "aws_guardduty_detector" "test" {
  enable = true
}`
}

func testAccGuardDutyFilterConfig_updateTags() string {
	return `
resource "aws_guardduty_filter" "test" {
  detector_id = "${aws_guardduty_detector.test.id}"
	name        = "test-filter"
	action      = "ARCHIVE"
	rank        = 1

  finding_criteria {
    criterion {
		field     = "region"
		values    = ["us-west-2"]
		condition = "equals"
	  }
	}

  tags = {
	  Key = "Updated"
  }
}

resource "aws_guardduty_detector" "test" {
  enable = true
}`
}
