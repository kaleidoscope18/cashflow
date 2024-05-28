package main

import (
	"github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateImage(ctx *pulumi.Context, repositoryUrl pulumi.StringInput) error {
	_, err := ecr.NewImage(ctx, "image", &ecr.ImageArgs{
		RepositoryUrl: repositoryUrl,
		Dockerfile:    pulumi.String("../Dockerfile"),
		Platform:      pulumi.String("linux/arm64"),
	})
	if err != nil {
		return err
	}

	return nil
}
