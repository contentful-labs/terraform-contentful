provider "contentful" {
  cma_token = ""
  organization_id = ""
}
resource "contentful_space" "test" {
  name = "my-update-space-name"
}

resource "contentful_webhook" "test" {
  name = "b3"
  space_id = "${contentful_space.test.id}"
  url = "http://yeah-updated.com"
  headers = {
    "wwwwwwwwwwwwking" = "update just now"
  }
  topics = ["*.*"]
}

resource "contentful_apikey" "test" {
  space_id = "${contentful_space.test.id}"
  name = "foobar3"
}

output "space_id" {
  value = "${contentful_space.test.id}"
}

output "space_name" {
  value = "${contentful_space.test.name}"
}

output "space_version" {
  value = "${contentful_space.test.version}"
}


output "webhook_id" {
  value = "${contentful_webhook.test.id}"
}

output "webhook_name" {
  value = "${contentful_webhook.test.name}"
}

# needs special plan
# resource "contentful_locale" "test" {
  # space_id = "${contentful_space.test.id}"
  # name = "foo"
  # code = "ru"
# }
