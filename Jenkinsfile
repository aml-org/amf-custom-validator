#!groovy
pipeline {
    agent none
    stages {
        stage('Test validator (Go)') {
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target ci-go'
                    registryCredentialsId 'dockerhub-pro-credentials'
                }
            }
            steps {
                echo "Success" // Tests are actually run when building the agent in the Dockerfile
            }
        }
        stage('Test generated WASM (JS)') {
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target ci-js'
                    registryCredentialsId 'dockerhub-pro-credentials'
                }
            }
            steps {
                echo "Success" // Tests are actually run when building the agent in the Dockerfile
            }
        }
        stage('Publish snapshot') {
            when {
                anyOf {
                    branch 'develop'
                    branch 'continuous-delivery'
                }
            }
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target publish-snapshot --build-arg BUILD_NUMBER=${BUILD_NUMBER}'
                    registryCredentialsId 'dockerhub-pro-credentials'
                }
            }
            steps {
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'github-salt', passwordVariable: 'GITHUB_PASS', usernameVariable: 'GITHUB_USER']]) {
                          sh '''#!/bin/bash
                                URL="https://${GITHUB_USER}:${GITHUB_PASS}@github.com/aml-org/amf-custom-validator"
                                cd ./wrappers/js
                                npm-snapshot $BUILD_NUMBER
                                VERSION=$(node -pe "require('./package.json').version")
                                npm publish
                                git tag $VERSION
                                git push $URL $VERSION
                          '''
                }
            }
        }
    }
}
