package awsservice

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/route53"
)

var awsRegion = "us-west-2"

type AWSLoadBalancerService interface {
	CreateLoadBalancer(*LoadBalancerDefinition) (string, error)
	DeleteLoadBalancer(string) error
	RegisterInstances(string, []string) error
	DeregisterInstances(string, []string) error
	GetLoadBalancerInfo(string) (*LoadBalancerInfo, error)
	GetInstanceHealth(string) (*LBInstanceHealthInfo, error)
	SetHealthCheck(string, *LBHealthCheck) error
}

type AWSRoute53Service interface {
	CreateDNSRecord(*Route53RecordDefinition) error
	DeleteDNSRecord(*Route53RecordDefinition) error
}

type AWSEC2Service interface {
	RunInstances(*InstancesDefinition) ([]string, error)
	StartInstances([]string) error
	StopInstances([]string) error
	FindInstancesByTag(string, string) ([]string, error)
	TagInstances([]string, string, string) error
	DeleteTag([]string, string) error
	GetSubnetInfo(string) (*SubnetInfo, error)
	GetInstancesInfo([]string) ([]InstanceInfo, error)
	TerminateInstances([]string) error
}

type AWSService interface {
	AWSLoadBalancerService
	AWSRoute53Service
	AWSEC2Service
}

type LimitedRoute53API interface {
	ChangeResourceRecordSets(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
}

type LimitedELBAPI interface {
	CreateLoadBalancer(*elb.CreateLoadBalancerInput) (*elb.CreateLoadBalancerOutput, error)
	DeleteLoadBalancer(*elb.DeleteLoadBalancerInput) (*elb.DeleteLoadBalancerOutput, error)
	RegisterInstancesWithLoadBalancer(*elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error)
	DeregisterInstancesFromLoadBalancer(*elb.DeregisterInstancesFromLoadBalancerInput) (*elb.DeregisterInstancesFromLoadBalancerOutput, error)
}

type RealAWSService struct {
	elbc *elb.ELB
	r53c *route53.Route53
	ec2  *ec2.EC2
}

// Testing types
type AWSActionLog struct {
	Action        string
	NotableParams map[string]string
}

type TestingAWSService struct {
	Log []AWSActionLog
}

// NewStaticAWSService uses the static credential provider (pass in access key ID and secret key)
func NewStaticAWSService(id string, secret string) AWSService {
	s := session.New(&aws.Config{Credentials: credentials.NewStaticCredentials(id, secret, ""), Region: &awsRegion})

	return &RealAWSService{
		elbc: elb.New(s),
		r53c: route53.New(s),
		ec2:  ec2.New(s),
	}
}

// NewAWSService uses the default Environment credential store
func NewAWSService() AWSService {
	s := session.New(&aws.Config{Region: &awsRegion})

	return &RealAWSService{
		elbc: elb.New(s),
		r53c: route53.New(s),
		ec2:  ec2.New(s),
	}
}

// Stupid AWS SDK...
func stringSlicetoStringPointerSlice(s []string) []*string {
	o := []*string{}
	for _, str := range s {
		nstr := str
		o = append(o, &nstr)
	}
	return o
}

func stringPointerSlicetoStringSlice(s []*string) []string {
	o := []string{}
	for _, ptr := range s {
		o = append(o, drefStringPtr(ptr))
	}
	return o
}

func drefStringPtr(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func drefInt64Ptr(ptr *int64) int64 {
	if ptr == nil {
		return int64(0)
	}
	return *ptr
}

func instanceIDSlice(ids []string) []*elb.Instance {
	instances := []*elb.Instance{}
	for _, id := range ids {
		curid := id
		instances = append(instances, &elb.Instance{
			InstanceId: &curid,
		})
	}
	return instances
}
