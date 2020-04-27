package resource

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/jinzhu/copier"
	"math/rand"
	"reflect"
	"sort"
	"time"
)

const (
	generatedClusterNameSuffixLength = 8
	generatedClusterNamePrefix       = "EKS-"
)

var loggingTypes = []string{"api", "audit", "authenticator", "controllerManager", "scheduler"}

func describeClusterToModel(cluster eks.Cluster, model *Model) {
	model.Name = cluster.Name
	model.RoleArn = cluster.RoleArn
	model.Version = cluster.Version
	model.ResourcesVpcConfig = &ResourcesVpcConfig{
		SecurityGroupIds:      aws.StringValueSlice(cluster.ResourcesVpcConfig.SecurityGroupIds),
		SubnetIds:             aws.StringValueSlice(cluster.ResourcesVpcConfig.SubnetIds),
		EndpointPublicAccess:  cluster.ResourcesVpcConfig.EndpointPublicAccess,
		EndpointPrivateAccess: cluster.ResourcesVpcConfig.EndpointPrivateAccess,
		PublicAccessCidrs:     aws.StringValueSlice(cluster.ResourcesVpcConfig.PublicAccessCidrs),
	}
	for _, l := range cluster.Logging.ClusterLogging {
		if *l.Enabled {
			model.EnabledClusterLoggingTypes = aws.StringValueSlice(l.Types)
		}
	}
	if cluster.EncryptionConfig != nil {
		var encryptionConfigs []EncryptionConfig
		for _, e := range cluster.EncryptionConfig {
			encryptionConfigs = append(encryptionConfigs, EncryptionConfig{
				Resources: aws.StringValueSlice(e.Resources),
				Provider: &Provider{
					KeyArn: e.Provider.KeyArn,
				},
			})
		}
		model.EncryptionConfig = encryptionConfigs
	}
	model.Arn = cluster.Arn
	model.CertificateAuthorityData = cluster.CertificateAuthority.Data
	model.ClusterSecurityGroupId = cluster.ResourcesVpcConfig.ClusterSecurityGroupId
	model.Endpoint = cluster.Endpoint
}

func makeCreateClusterInput(model *Model) *eks.CreateClusterInput {
	input := &eks.CreateClusterInput{
		Name: model.Name,
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SubnetIds:             aws.StringSlice(model.ResourcesVpcConfig.SubnetIds),
			EndpointPublicAccess:  model.ResourcesVpcConfig.EndpointPublicAccess,
			EndpointPrivateAccess: model.ResourcesVpcConfig.EndpointPrivateAccess,
		},
		Logging:          createLogging(model),
		RoleArn:          model.RoleArn,
		Version:          model.Version,
		EncryptionConfig: createEncryptionConfig(model),
	}
	if model.ResourcesVpcConfig.SecurityGroupIds != nil {
		input.ResourcesVpcConfig.SecurityGroupIds = aws.StringSlice(model.ResourcesVpcConfig.SecurityGroupIds)
	}
	return input
}

func createEncryptionConfig(model *Model) []*eks.EncryptionConfig {
	var configs []*eks.EncryptionConfig
	for _, c := range model.EncryptionConfig {
		configs = append(configs, &eks.EncryptionConfig{
			Provider:  &eks.Provider{KeyArn: c.Provider.KeyArn},
			Resources: aws.StringSlice(c.Resources),
		})
	}
	return configs
}

func updateVpcConfig(svc eksiface.EKSAPI, model *Model) error {
	input := &eks.UpdateClusterConfigInput{
		Name: model.Name,
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			EndpointPublicAccess:  model.ResourcesVpcConfig.EndpointPublicAccess,
			EndpointPrivateAccess: model.ResourcesVpcConfig.EndpointPrivateAccess,
		},
	}
	if model.ResourcesVpcConfig.PublicAccessCidrs != nil {
		input.ResourcesVpcConfig.PublicAccessCidrs = aws.StringSlice(model.ResourcesVpcConfig.PublicAccessCidrs)
	}
	_, err := svc.UpdateClusterConfig(input)
	return err
}

func createLogging(model *Model) *eks.Logging {
	var logSetups []*eks.LogSetup
	if model.EnabledClusterLoggingTypes != nil {
		if len(model.EnabledClusterLoggingTypes) > 0 {
			logSetups = append(logSetups, &eks.LogSetup{
				Enabled: aws.Bool(true),
				Types:   aws.StringSlice(model.EnabledClusterLoggingTypes),
			})
		}
	}
	disabledTypes := getDisabledLoggingTypes(model.EnabledClusterLoggingTypes)
	if disabledTypes != nil {
		if len(disabledTypes) != 0 {
			logSetups = append(logSetups, &eks.LogSetup{
				Enabled: aws.Bool(false),
				Types:   aws.StringSlice(disabledTypes),
			})
		}
	}
	return &eks.Logging{ClusterLogging: logSetups}
}

func updateLoggingConfig(svc eksiface.EKSAPI, model *Model) error {
	input := &eks.UpdateClusterConfigInput{Name: model.Name, Logging: createLogging(model)}
	_, err := svc.UpdateClusterConfig(input)
	return err
}

func getDisabledLoggingTypes(enabled []string) (disabled []string) {
	if enabled == nil {
		return loggingTypes
	}
	if len(enabled) == 0 {
		return loggingTypes
	}
	for _, d := range loggingTypes {
		isDisabled := true
		for _, e := range enabled {
			if e == d {
				isDisabled = false
				break
			}
		}
		if isDisabled {
			disabled = append(disabled, d)
		}
	}
	return disabled
}

func updateVersionConfig(svc eksiface.EKSAPI, model *Model) error {
	input := &eks.UpdateClusterVersionInput{
		Name:    model.Name,
		Version: model.Version,
	}
	// TODO: version rollback is problematic, see https://github.com/aws/containers-roadmap/issues/497
	_, err := svc.UpdateClusterVersion(input)
	return err
}

func versionChanged(current Model, desired Model) bool {
	if desired.Version == nil {
		return false
	}
	return *current.Version != *desired.Version
}

func vpcChanged(current Model, desired Model) bool {
	desiredVpc := &ResourcesVpcConfig{}
	err := copier.Copy(desiredVpc, desired.ResourcesVpcConfig)
	if err != nil {
		panic(err)
	}
	currentVpc := current.ResourcesVpcConfig
	if desiredVpc.PublicAccessCidrs == nil {
		desiredVpc.PublicAccessCidrs = []string{"0.0.0.0/0"}
	}
	if desiredVpc.EndpointPrivateAccess == nil {
		desiredVpc.EndpointPrivateAccess = aws.Bool(false)
	}
	if desiredVpc.EndpointPublicAccess == nil {
		desiredVpc.EndpointPublicAccess = aws.Bool(true)
	}
	if (!slicesEqual(currentVpc.PublicAccessCidrs, desiredVpc.PublicAccessCidrs)) ||
		(*currentVpc.EndpointPrivateAccess != *desiredVpc.EndpointPrivateAccess) ||
		(*currentVpc.EndpointPublicAccess != *desiredVpc.EndpointPublicAccess) {
		return true
	}
	return false
}

func loggingChanged(current Model, desired Model) bool {
	if !slicesEqual(desired.EnabledClusterLoggingTypes, current.EnabledClusterLoggingTypes) {
		return true
	}
	return false
}

func slicesEqual(s1 []string, s2 []string) bool {
	// evaluates that slices both contain the same strings regardless of order
	sort.StringSlice.Sort(s1)
	sort.StringSlice.Sort(s2)
	return reflect.DeepEqual(s1, s2)
}

func generateClusterName() *string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, generatedClusterNameSuffixLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	generated := generatedClusterNamePrefix + string(b)
	return &generated
}

func matchesAwsErrorCode(err error, code string) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == code {
			return true
		}
	}
	return false
}

func isPrivate(model *Model) bool {
	if !*model.ResourcesVpcConfig.EndpointPublicAccess {
		return true
	}
	cidrs := model.ResourcesVpcConfig.PublicAccessCidrs
	if cidrs == nil {
		return false
	}
	for _, cidr := range cidrs {
		if cidr == "0.0.0.0/0" {
			return false
		}
	}
	return true
}
