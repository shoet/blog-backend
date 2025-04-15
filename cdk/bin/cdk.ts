#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { BlogBackendAppStack } from "../lib/blog-backend-app-stack";

declare module "aws-cdk-lib" {
  interface StackProps {
    stage: string;
    commitHash?: string;
  }
}

const app = new cdk.App();

const stage = app.node.tryGetContext("STAGE");
if (!stage) {
  throw new Error("STAGE is required");
}
if (!["dev", "prod"].includes(stage)) {
  throw new Error("STAGE must be either dev or prod");
}

const commitHash = app.node.tryGetContext("COMMIT_HASH");

console.log(`Deploying to stage: ${stage}`);

new BlogBackendAppStack(app, `BlogBackendAppStack-${stage}`, {
  stage: stage,
  commitHash: commitHash,
});
