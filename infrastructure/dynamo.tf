resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name           = "AlpacaAssets"
  billing_mode   = "PROVISIONED"
  read_capacity  = 25
  write_capacity = 25
  hash_key       = "Symbol"

  attribute {
    name = "Symbol"
    type = "S"
  }

  attribute {
    name = "OpenPrice"
    type = "S"
  }

  attribute {
    name = "LastUpdated"
    type = "S"
  }

  tags = {
    terraform = "true"
  }
}
