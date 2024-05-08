package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awselasticache"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awswafv2"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoAwsStackProps struct {
	awscdk.StackProps
}

func NewGoAwsStack(scope constructs.Construct, id string, props *GoAwsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	//create db table here
	userTable := awsdynamodb.NewTable(stack, jsii.String("userTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("userTable"),
	})

	userTable.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("EmailIndex"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("email"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	cosmetologistUserTable := awsdynamodb.NewTable(stack, jsii.String("cosmetologistUserTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("cosmetologistUserTable"),
	})

	cosmetologistUserTable.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("EmailIndex"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("email"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	adminUserTable := awsdynamodb.NewTable(stack, jsii.String("adminUserTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("adminUserTable"),
	})

	adminUserTable.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("EmailIndex"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("email"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	productsTable := awsdynamodb.NewTable(stack, jsii.String("productsTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("productsTable"),
	})

	resetTokensTable := awsdynamodb.NewTable(stack, jsii.String("resetTokensTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("email"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("resetTokensTable"),
	})
	resetTokensTable.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String("TokenIndex"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("token"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	// awselasticache.NewCfnServerlessCache(stack, jsii.String("myElasticache"), &awselasticache.CfnServerlessCacheProps{
	// 	Engine: jsii.String("redis"),
	// 	ServerlessCacheName: jsii.String("myServerlessCache"),
	// })

	myFunction := awslambda.NewFunction(stack, jsii.String("myLambdaFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code: awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
	})

	userTable.GrantReadWriteData(myFunction)
	cosmetologistUserTable.GrantReadWriteData(myFunction)
	productsTable.GrantReadWriteData(myFunction)
	resetTokensTable.GrantReadWriteData(myFunction)

	adminUserTable.GrantReadData(myFunction)

	api := awsapigateway.NewRestApi(stack, jsii.String("myRESTApi"), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
		CloudWatchRole: jsii.Bool(true),
	})

	integration := awsapigateway.NewLambdaIntegration(myFunction, nil)

	// webAcl := awswafv2.NewCfnWebACL(stack, jsii.String("MyWebACL"), &awswafv2.CfnWebACLProps{
  //   Scope: jsii.String("REGIONAL"),
  //   DefaultAction: &awswafv2.CfnWebACL_DefaultActionProperty{
  //       Allow: &awswafv2.CfnWebACL_AllowActionProperty{},
  //   },
  //   VisibilityConfig: &awswafv2.CfnWebACL_VisibilityConfigProperty{
  //       SampledRequestsEnabled: jsii.Bool(true),
  //       CloudWatchMetricsEnabled: jsii.Bool(true),
  //       MetricName: jsii.String("webACL"),
  //   },
  //   Rules: &[]*awswafv2.CfnWebACL_RuleProperty{
  //       {
  //           Name: jsii.String("RateLimitRule"),
  //           Priority: jsii.Number(1),
  //           Action: &awswafv2.CfnWebACL_RuleActionProperty{
  //               Block: &awswafv2.CfnWebACL_BlockActionProperty{}, 
  //           },
  //           Statement: &awswafv2.CfnWebACL_StatementProperty{
  //               RateBasedStatement: &awswafv2.CfnWebACL_RateBasedStatementProperty{
  //                   Limit: jsii.Number(100), 
  //                   AggregateKeyType: jsii.String("IP"), 
	// 									EvaluationWindowSec: jsii.Number(600),
  //               },
  //           },
  //           VisibilityConfig: &awswafv2.CfnWebACL_VisibilityConfigProperty{
  //               SampledRequestsEnabled: jsii.Bool(true),
  //               CloudWatchMetricsEnabled: jsii.Bool(true),
  //               MetricName: jsii.String("RateLimitRule"),
  //           },
  //       },
  //   },
	// })

	// partition := "aws" 
	// region := *stack.Region()
	// restApiId := *api.RestApiId()
	// stageName := "prod" 
	
	// restApiArn := fmt.Sprintf("arn:%s:apigateway:%s::/restapis/%s/stages/%s", partition, region, restApiId, stageName)

	
	// awswafv2.NewCfnWebACLAssociation(stack, jsii.String("WebAclApiGatewayAssociation"), &awswafv2.CfnWebACLAssociationProps{
	// 		ResourceArn: jsii.String(restApiArn), 
	// 		WebAclArn:   webAcl.AttrArn(),
	// })

	// Initialize the root of the API
	apiRoot := api.Root()

	resetPassResource := apiRoot.AddResource(jsii.String("reset-password"), nil)

	resetPassCreateResource := resetPassResource.AddResource(jsii.String("create"), nil)
	resetPassCreateResource.AddMethod(jsii.String("POST"), integration, nil)

	resetPassValidateResource := resetPassResource.AddResource(jsii.String("validate"), nil)
	resetPassValidateResource.AddMethod(jsii.String("POST"), integration, nil)

	registerResource := apiRoot.AddResource(jsii.String("register"), nil)
	registerResource.AddMethod(jsii.String("POST"), integration, nil)

	adminResource := apiRoot.AddResource(jsii.String("admin"), nil)

	adminClientsResource := adminResource.AddResource(jsii.String("clients"), nil)
	adminClientsResource.AddMethod(jsii.String("GET"), integration, nil)

	adminCosmetologistResource := adminResource.AddResource(jsii.String("cosmetologist"), nil)

	adminCosmetologistsResource := adminResource.AddResource(jsii.String("cosmetologists"), nil)
	adminCosmetologistsResource.AddMethod(jsii.String("GET"), integration, nil)

	adminProductsResource := adminResource.AddResource(jsii.String("products"), nil)
	
	adminProductsCreateResource := adminProductsResource.AddResource(jsii.String("create"), nil)
	adminProductsCreateResource.AddMethod(jsii.String("POST"), integration, nil)

	registerCosmetologistResource := adminCosmetologistResource.AddResource(jsii.String("register"), nil)
	registerCosmetologistResource.AddMethod(jsii.String("POST"), integration, nil)

	assignCosmetologistResource := adminCosmetologistResource.AddResource(jsii.String("assign"), nil)
	assignCosmetologistResource.AddMethod(jsii.String("POST"), integration, nil)

	loginResource := apiRoot.AddResource(jsii.String("login"), nil)
	loginResource.AddMethod(jsii.String("POST"), integration, nil)

	cosmetologistResource := apiRoot.AddResource(jsii.String("cosmetologist"), nil)

	loginCosmetologistResource := cosmetologistResource.AddResource(jsii.String("login"), nil)
	loginCosmetologistResource.AddMethod(jsii.String("POST"), integration, nil)

	loginAdminResource := adminResource.AddResource(jsii.String("login"), nil)
	loginAdminResource.AddMethod(jsii.String("POST"), integration, nil)

	protectedUserResource := apiRoot.AddResource(jsii.String("protected-user"), nil)
	protectedUserResource.AddMethod(jsii.String("GET"), integration, nil)

	protectedCosmetologistResource := apiRoot.AddResource(jsii.String("protected-cosmetologist"), nil)
	protectedCosmetologistResource.AddMethod(jsii.String("GET"), integration, nil)

	// Define 'protected-admin' under 'admin'
	protectedAdminResource := apiRoot.AddResource(jsii.String("protected-admin"), nil)
	protectedAdminResource.AddMethod(jsii.String("GET"), integration, nil)
	//If we want to get id from request
	//registerResource := api.Root().AddResource(jsii.String("register/{id}"), nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoAwsStack(app, "GoAwsStack", &GoAwsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
