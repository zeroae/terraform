resource "template_file" "test" {
    count = 1
    template = "content"
}

output "bad" {
    value = "${template_file.test.*}"
}
