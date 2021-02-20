# awls
Lookup EC2 information in the terminal

## Usage

```bash
Usage: awls <search>

Arguments:
  <search>    EC2 Instance Name search term

Flags:
  -h, --help                      Show context-sensitive help.
  -i, --ip-only                   Output only Private IPs
  -d, --delimiter=" "             IP delimiter
  -f, --filter-type="tag:Name"    EC2 Filter Type (https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html)
```

### Examples

```bash
# basic tag:Name search
awls searchterm

# fuzzy tag:Name search
awls 'search*term'

# fuzzy tag:cluster search
awls 'dev*ec2*' -f tag:cluster

# Search by instance-id and output private IP delimited by ','
awls 'i-0bfxxxxxxx' -f instance-id -i -d,
```
