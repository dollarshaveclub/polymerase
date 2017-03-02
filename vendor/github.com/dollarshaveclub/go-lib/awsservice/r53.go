package awsservice

import (
	"github.com/aws/aws-sdk-go/service/route53"
)

type Route53RecordDefinition struct {
	ZoneID string
	Name   string
	Value  string
	Type   string
	TTL    int64
}

func (aws *RealAWSService) executeR53Action(a string, rd *Route53RecordDefinition) error {
	param := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: &a,
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: &rd.Name,
						Type: &rd.Type,
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: &rd.Value,
							},
						},
						TTL: &rd.TTL,
					},
				},
			},
		},
		HostedZoneId: &rd.ZoneID,
	}
	_, err := aws.r53c.ChangeResourceRecordSets(param)
	return err
}

func (aws *RealAWSService) CreateDNSRecord(rd *Route53RecordDefinition) error {
	return aws.executeR53Action("CREATE", rd)
}

func (aws *RealAWSService) DeleteDNSRecord(rd *Route53RecordDefinition) error {
	return aws.executeR53Action("DELETE", rd)
}

// Testing mocks

func (aws *TestingAWSService) CreateDNSRecord(rd *Route53RecordDefinition) error {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "CreateDNSRecord",
		NotableParams: map[string]string{
			"name":  rd.Name,
			"value": rd.Value,
		},
	})
	return nil
}

func (aws *TestingAWSService) DeleteDNSRecord(rd *Route53RecordDefinition) error {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "CreateDNSRecord",
		NotableParams: map[string]string{
			"name": rd.Name,
		},
	})
	return nil
}
