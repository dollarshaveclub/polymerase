package awsservice

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/elb"
)

type ELBListener struct {
	InstancePort         int64
	LoadBalancerPort     int64
	LoadBalancerProtocol string
	InstanceProtocol     string
	CertificateID        string
}

type LoadBalancerDefinition struct {
	Listeners      []ELBListener
	Name           string
	SecurityGroups []string
	Scheme         string
	Subnets        []string
}

type LBHealthCheck struct {
	Target             string
	Interval           int64
	Timeout            int64
	HealthyThreshold   int64
	UnhealthyThreshold int64
}

type LBInstanceHealth struct {
	ID          string
	Description string
	ReasonCode  string
	State       string
}

type LBInstanceHealthInfo struct {
	LBName    string
	Instances []LBInstanceHealth
}

type LoadBalancerInfo struct {
	Name              string
	Scheme            string
	SecurityGroups    []string
	Subnets           []string
	VPCID             string
	AvailabilityZones []string
	DNSName           string
	Instances         []string
}

func (aws *RealAWSService) CreateLoadBalancer(lbd *LoadBalancerDefinition) (string, error) {
	listeners := []*elb.Listener{}
	for _, l := range lbd.Listeners {
		ip := l.InstancePort // allocate new objects so pointers in struct are unique
		lbp := l.LoadBalancerPort
		pr := l.LoadBalancerProtocol
		ipr := l.InstanceProtocol
		cid := l.CertificateID
		listeners = append(listeners, &elb.Listener{
			InstancePort:     &ip,
			LoadBalancerPort: &lbp,
			Protocol:         &pr,
			InstanceProtocol: &ipr,
			SSLCertificateId: &cid,
		})
	}
	o, err := aws.elbc.CreateLoadBalancer(&elb.CreateLoadBalancerInput{
		Listeners:        listeners,
		LoadBalancerName: &lbd.Name,
		SecurityGroups:   stringSlicetoStringPointerSlice(lbd.SecurityGroups),
		Subnets:          stringSlicetoStringPointerSlice(lbd.Subnets),
	})
	if err != nil {
		return "", err
	}
	return *o.DNSName, nil
}

func (aws *RealAWSService) GetLoadBalancerInfo(n string) (*LoadBalancerInfo, error) {
	dlbi := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: stringSlicetoStringPointerSlice([]string{n}),
	}
	result := &LoadBalancerInfo{}
	res, err := aws.elbc.DescribeLoadBalancers(dlbi)
	if err != nil {
		return result, err
	}
	result.AvailabilityZones = stringPointerSlicetoStringSlice(res.LoadBalancerDescriptions[0].AvailabilityZones)
	result.SecurityGroups = stringPointerSlicetoStringSlice(res.LoadBalancerDescriptions[0].SecurityGroups)
	result.Subnets = stringPointerSlicetoStringSlice(res.LoadBalancerDescriptions[0].Subnets)
	result.DNSName = drefStringPtr(res.LoadBalancerDescriptions[0].DNSName)
	result.Name = drefStringPtr(res.LoadBalancerDescriptions[0].LoadBalancerName)
	result.Scheme = drefStringPtr(res.LoadBalancerDescriptions[0].Scheme)
	result.VPCID = drefStringPtr(res.LoadBalancerDescriptions[0].VPCId)
	il := []string{}
	for _, inst := range res.LoadBalancerDescriptions[0].Instances {
		il = append(il, drefStringPtr(inst.InstanceId))
	}
	result.Instances = il
	return result, nil
}

func (aws *RealAWSService) GetInstanceHealth(n string) (*LBInstanceHealthInfo, error) {
	result := &LBInstanceHealthInfo{
		LBName: n,
	}
	dih := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: &n,
	}
	r, err := aws.elbc.DescribeInstanceHealth(dih)
	if err != nil {
		return result, err
	}
	instances := []LBInstanceHealth{}
	for _, is := range r.InstanceStates {
		inst := LBInstanceHealth{
			ID:          drefStringPtr(is.InstanceId),
			Description: drefStringPtr(is.Description),
			ReasonCode:  drefStringPtr(is.ReasonCode),
			State:       drefStringPtr(is.State),
		}
		instances = append(instances, inst)
	}
	result.Instances = instances
	return result, nil
}

func (aws *RealAWSService) SetHealthCheck(n string, hc *LBHealthCheck) error {
	chk := &elb.ConfigureHealthCheckInput{
		LoadBalancerName: &n,
		HealthCheck: &elb.HealthCheck{
			HealthyThreshold:   &hc.HealthyThreshold,
			Interval:           &hc.Interval,
			Target:             &hc.Target,
			Timeout:            &hc.Timeout,
			UnhealthyThreshold: &hc.UnhealthyThreshold,
		},
	}
	_, err := aws.elbc.ConfigureHealthCheck(chk)
	return err
}

func (aws *RealAWSService) DeleteLoadBalancer(n string) error {
	_, err := aws.elbc.DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{
		LoadBalancerName: &n,
	})
	return err
}

func (aws *RealAWSService) RegisterInstances(n string, ids []string) error {
	_, err := aws.elbc.RegisterInstancesWithLoadBalancer(&elb.RegisterInstancesWithLoadBalancerInput{
		Instances:        instanceIDSlice(ids),
		LoadBalancerName: &n,
	})
	return err
}

func (aws *RealAWSService) DeregisterInstances(n string, ids []string) error {
	_, err := aws.elbc.DeregisterInstancesFromLoadBalancer(&elb.DeregisterInstancesFromLoadBalancerInput{
		Instances:        instanceIDSlice(ids),
		LoadBalancerName: &n,
	})
	return err
}

// Testing mocks

func (aws *TestingAWSService) CreateLoadBalancer(lbd *LoadBalancerDefinition) (string, error) {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "CreateLoadBalancer",
		NotableParams: map[string]string{
			"name":            lbd.Name,
			"security_groups": fmt.Sprintf("%v", lbd.SecurityGroups),
			"scheme":          lbd.Scheme,
			"subnets":         fmt.Sprintf("%v", lbd.Subnets),
			"listeners":       fmt.Sprintf("%v", lbd.Listeners),
		},
	})
	return "", nil
}

func (aws *TestingAWSService) DeleteLoadBalancer(n string) error {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "DeleteLoadBalancer",
		NotableParams: map[string]string{
			"name": n,
		},
	})
	return nil
}

func (aws *TestingAWSService) RegisterInstances(n string, ids []string) error {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "RegisterInstances",
		NotableParams: map[string]string{
			"ids": fmt.Sprintf("%v", ids),
		},
	})
	return nil
}

func (aws *TestingAWSService) DeregisterInstances(n string, ids []string) error {
	aws.Log = append(aws.Log, AWSActionLog{
		Action: "DeregisterInstances",
		NotableParams: map[string]string{
			"ids": fmt.Sprintf("%v", ids),
		},
	})
	return nil
}
