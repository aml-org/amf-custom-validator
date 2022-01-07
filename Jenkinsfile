#!groovy
pipeline {
    agent none
    environment {
        DOCKER_BUILDKIT='1' // optimizes target builds for multistage dockerfile
    }
    stages {
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
        stage('Test generated WASM (JS Browser)') {
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target ci-browser'
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
                    branch 'apimf-3610'
                }
            }
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target publish-snapshot'
                    registryCredentialsId 'dockerhub-pro-credentials'
                }
            }
            environment {
                NPM_TOKEN = credentials('aml-org-bot-npm-token')
                GITHUB = credentials('github-salt')
            }
            steps {
                sh '''  #!/bin/bash
                        cd /src

                        # Login
                        echo //registry.npmjs.org/:_authToken=${env.NPM_TOKEN} >> ~/.npmrc
                        echo @aml-org:registry=https://registry.npmjs.org/ >> ~/.npmrc
                        npm config set registry https://registry.npmjs.org/
                        npm whoami

                        # Publish
                        cd ./wrappers/js
                        npm-snapshot $BUILD_NUMBER
                        VERSION=$(node -pe "require('./package.json').version")
                        npm publish --access public
                        npm dist-tag add @aml-org/amf-custom-validator@${VERSION} snapshot

                        cd ../js-web
                        npm-snapshot $BUILD_NUMBER
                        npm publish --verbose --access public
                        npm dist-tag add @aml-org/amf-custom-validator-web@${VERSION} snapshot
                        URL="https://${GITHUB_USR}:${GITHUB_PSW}@github.com/aml-org/amf-custom-validator"

                        git tag v$VERSION
                        git push $URL v$VERSION
                '''
            }
        }
    }
}

