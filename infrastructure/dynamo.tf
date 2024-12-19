resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name           = "AlpacaAssets"
  billing_mode   = "PROVISIONED"
  read_capacity  = 25
  write_capacity = 25
  hash_key       = "Symbol"
  range_key      = "OpenPrice"

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

  local_secondary_index {
    name               = "local"
    non_key_attributes = ["OpenPrice"]
    projection_type    = "KEYS_ONLY"
    range_key          = "LastUpdated"
  }

  tags = {
    terraform = "true"
  }
}
