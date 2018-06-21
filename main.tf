# Configure the OpenStack Provider
provider "opentelekomcloud" {
  user_name   = "c2c-6"
  domain_name = "OTC00000000001000010501"
  password    = "Newuser@123"
  auth_url    = "https://iam.eu-de.otc.t-systems.com/v3"
  region      = "eu-de"
  tenant_id   = "17fbda95add24720a4038ba4b1c705ed"
}

resource "opentelekomcloud_sfs_file_sharing_v2" "sfs1" {
  size=1
  name="sfs-c2c-dev14"
  availability_zone="eu-de-01"
  metadata={
    "share_used"="0"
    "key1"="value1"
    "key2"="value2"
  }
  vpc_id="5232f396-d6cc-4a81-8de3-afd7a7ecdfd8" //vpcid
  access_level="rw"
  description="sfs_c2c_test-file"
  size=32
}

output "access_id" {
  //  value = "${data.opentelekomcloud_sfs_stack_v1.stacks.template_body["key_name"]}"
  value = "${opentelekomcloud_sfs_file_sharing_v2.sfs1.access_id}"
}
output "access_state" {
  //  value = "${data.opentelekomcloud_sfs_stack_v1.stacks.template_body["key_name"]}"
  value = "${opentelekomcloud_sfs_file_sharing_v2.sfs1.access_state}"
}
output "vpc_id" {
  //  value = "${data.opentelekomcloud_sfs_stack_v1.stacks.template_body["key_name"]}"
  value = "${opentelekomcloud_sfs_file_sharing_v2.sfs1.vpc_id}"
}
output "access_type" {
  //  value = "${data.opentelekomcloud_sfs_stack_v1.stacks.template_body["key_name"]}"
  value = "${opentelekomcloud_sfs_file_sharing_v2.sfs1.access_type}"
}
output "access_level" {
  //  value = "${data.opentelekomcloud_sfs_stack_v1.stacks.template_body["key_name"]}"
  value = "${opentelekomcloud_sfs_file_sharing_v2.sfs1.access_level}"
}
output "size" {
  //  value = "${data.opentelekomcloud_sfs_stack_v1.stacks.template_body["key_name"]}"
  value = "${opentelekomcloud_sfs_file_sharing_v2.sfs1.size}"
}