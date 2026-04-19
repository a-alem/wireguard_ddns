# resource "aws_route53_zone" "vpn_zone" {
#   name = "kfupm-masters-aalem.com"
#   comment = "HostedZone created by Route53 Registrar"
# }

# resource "aws_route53_record" "vpn_onprem_server_record" {
#   zone_id = aws_route53_zone.vpn_zone.zone_id
#   name    = "onprem.vpn.kfupm-masters-aalem.com"
#   type    = "A"
#   ttl     = 60
#   records = ["157.175.59.210"]
# }

# resource "aws_route53_record" "vpn_offsite_server_record" {
#   zone_id = aws_route53_zone.vpn_zone.zone_id
#   name    = "offsite.vpn.kfupm-masters-aalem.com"
#   type    = "A"
#   ttl     = 60
#   records = ["157.175.149.198"]
# }