package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Containerize(ctx *pulumi.Context) error {
	repo, err := ecr.NewRepository(ctx, "cashflowrepository", &ecr.RepositoryArgs{})
	if err != nil {
		return err
	}
	ctx.Export("ecr repository", repo.RepositoryUrl)

	ecr.NewLifecyclePolicy(ctx, "cashflowrepository-lifecycle-policy", &ecr.LifecyclePolicyArgs{
		Policy: pulumi.String(`
			{
				"rules": [
					{
						"rulePriority": 1,
						"description": "Expire untagged older than 1 days",
						"selection": {
							"tagStatus": "untagged",
							"countType": "sinceImagePushed",
							"countUnit": "days",
							"countNumber": 1
						},
						"action": {
							"type": "expire"
						}
					},
					{
						"rulePriority": 2,
						"description": "Expire images older than 14 days",
						"selection": {
							"tagStatus": "any",
							"countType": "sinceImagePushed",
							"countUnit": "days",
							"countNumber": 14
						},
						"action": {
							"type": "expire"
						}
					}
				]
			}
		`),
		Repository: repo.Name,
	})

	CreateImage(ctx, repo.RepositoryUrl)

	return nil
}
