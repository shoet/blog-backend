#!/bin/bash

AWS_ACOUNT_ID=`aws sts get-caller-identity --query 'Account' --output text`
ECR="${AWS_ACOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com"

aws ecr get-login-password | \
  docker login --username AWS --password-stdin ${ECR}

docker tag blog-backend:latest ${ECR}/blog-backend:latest

docker push ${ECR}/blog-backend:latest

