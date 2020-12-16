# AWSQS::EKS::Cluster

A resource that creates Amazon Elastic Kubernetes Service (Amazon EKS) clusters.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "AWSQS::EKS::Cluster",
    "Properties" : {
        "<a href="#name" title="Name">Name</a>" : <i>String</i>,
        "<a href="#rolearn" title="RoleArn">RoleArn</a>" : <i>String</i>,
        "<a href="#version" title="Version">Version</a>" : <i>String</i>,
        "<a href="#resourcesvpcconfig" title="ResourcesVpcConfig">ResourcesVpcConfig</a>" : <i><a href="resourcesvpcconfig.md">ResourcesVpcConfig</a></i>,
        "<a href="#enabledclusterloggingtypes" title="EnabledClusterLoggingTypes">EnabledClusterLoggingTypes</a>" : <i>[ String, ... ]</i>,
        "<a href="#encryptionconfig" title="EncryptionConfig">EncryptionConfig</a>" : <i>[ <a href="encryptionconfig.md">EncryptionConfig</a>, ... ]</i>,
        "<a href="#kubernetesapiaccess" title="KubernetesApiAccess">KubernetesApiAccess</a>" : <i><a href="kubernetesapiaccess.md">KubernetesApiAccess</a></i>,
    }
}
</pre>

### YAML

<pre>
Type: AWSQS::EKS::Cluster
Properties:
    <a href="#name" title="Name">Name</a>: <i>String</i>
    <a href="#rolearn" title="RoleArn">RoleArn</a>: <i>String</i>
    <a href="#version" title="Version">Version</a>: <i>String</i>
    <a href="#resourcesvpcconfig" title="ResourcesVpcConfig">ResourcesVpcConfig</a>: <i><a href="resourcesvpcconfig.md">ResourcesVpcConfig</a></i>
    <a href="#enabledclusterloggingtypes" title="EnabledClusterLoggingTypes">EnabledClusterLoggingTypes</a>: <i>
      - String</i>
    <a href="#encryptionconfig" title="EncryptionConfig">EncryptionConfig</a>: <i>
      - <a href="encryptionconfig.md">EncryptionConfig</a></i>
    <a href="#kubernetesapiaccess" title="KubernetesApiAccess">KubernetesApiAccess</a>: <i><a href="kubernetesapiaccess.md">KubernetesApiAccess</a></i>
</pre>

## Properties

#### Name

A unique name for your cluster.

_Required_: No

_Type_: String

_Update requires_: [Replacement](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-replacement)

#### RoleArn

Amazon Resource Name (ARN) of the AWS Identity and Access Management (IAM) role. This provides permissions for Amazon EKS to call other AWS APIs.

_Required_: Yes

_Type_: String

_Update requires_: [Replacement](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-replacement)

#### Version

Desired Kubernetes version for your cluster. If you don't specify this value, the cluster uses the latest version from Amazon EKS.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ResourcesVpcConfig

An object that represents the virtual private cloud (VPC) configuration to use for an Amazon EKS cluster.



_Required_: Yes

_Type_: <a href="resourcesvpcconfig.md">ResourcesVpcConfig</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### EnabledClusterLoggingTypes

Enables exporting of logs from the Kubernetes control plane to Amazon CloudWatch Logs. By default, logs from the cluster control plane are not exported to CloudWatch Logs. The valid log types are api, audit, authenticator, controllerManager, and scheduler.

_Required_: No

_Type_: List of String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### EncryptionConfig

Encryption configuration for the cluster.

_Required_: No

_Type_: List of <a href="encryptionconfig.md">EncryptionConfig</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### KubernetesApiAccess

_Required_: No

_Type_: <a href="kubernetesapiaccess.md">KubernetesApiAccess</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Ref

When you pass the resource's logical ID to the intrinsic `Ref` function, `Ref` returns the name.

### Fn::GetAtt

`Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and samples of return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### Arn

ARN of the cluster (e.g., `arn:aws:eks:us-west-2:666666666666:cluster/prod`).

#### Endpoint

Endpoint for your Kubernetes API server (e.g., https://5E1D0CEXAMPLEA591B746AFC5AB30262.yl4.us-west-2.eks.amazonaws.com).

#### ClusterSecurityGroupId

Security group that was created by Amazon EKS for your cluster. Managed-node groups use this security group for control-plane-to-data-plane communications.

#### CertificateAuthorityData

Certificate-authority data for your cluster.

#### EncryptionConfigKeyArn

ARN or alias of the customer master key (CMK).

#### OIDCIssuerURL

Issuer URL for the OpenID Connect identity provider.

