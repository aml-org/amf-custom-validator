#!groovy
pipeline {
    agent none
    environment {
        DOCKER_BUILDKIT='1' // optimizes target builds for multistage dockerfile
    }
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
        stage ('Validate reports and profiles') {
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target ci-java'
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
        stage('Coverage') {
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target coverage'
                    registryCredentialsId 'dockerhub-pro-credentials'
                }
            }
            steps {
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-official', passwordVariable: 'SONAR_SERVER_TOKEN', usernameVariable: 'SONAR_SERVER_URL']]) {
                    sh '''  #!/bin/bash
                            cp /usr/src/coverage.out ./
                            echo sonar.host.url=${SONAR_SERVER_URL} >> sonar-project.properties
                            echo sonar.login=${SONAR_SERVER_TOKEN} >> sonar-project.properties
                            sonar-scanner
                    '''
                }

            }
        }
        stage('Nexus IQ') {
            when {
                anyOf {
                    branch 'develop'
                    branch 'master'
                }
            }
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target nexus-scan'
                    registryUrl 'https://artifacts.msap.io/'
                    registryCredentialsId 'harbor-docker-registry'
                }
            }
            steps {
                withCredentials([
                    [$class: 'UsernamePasswordMultiBinding', credentialsId: 'nexus-iq', passwordVariable: 'NEXUS_PASS', usernameVariable: 'NEXUS_USER'],
                ]) {
                script {
                    def args = [
                        '--authentication "$NEXUS_USER:$NEXUS_PASS"',
                        "--server-url https://nexusiq.build.msap.io",
                        "--application-id amf-custom-validator",
                        "--fail-on-policy-warnings" // might have to remove this if all policies are taken into account
                     // "--stage $stage",
                    ]
                    def result = sh(returnStatus: true, script: "java -jar /bin/nexusiq-cli ${args.join(' ')} /go.list")
                    if (result != 0) {
                        unstable "Failed Nexus IQ execution"
                    }
                }
                }
            }
        }
        stage('Publish snapshot') {
            when {
                anyOf {
                    branch 'develop'
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
                NPM = credentials('aml-org-bot-npm')
                GITHUB = credentials('github-salt')
            }
            steps {
                sh '''  #!/bin/bash
                        cd /src
                        npm-cli-login -u $NPM_USR -p $NPM_PSW -e als-amf-team@mulesoft.com
                        cd ./wrappers/js
                        npm-snapshot $BUILD_NUMBER
                        VERSION=$(node -pe "require('./package.json').version")
                        npm publish --access public
                        cd ../js-web
                        npm-snapshot $BUILD_NUMBER
                        npm publish --access public
                        URL="https://${GITHUB_USR}:${GITHUB_PSW}@github.com/aml-org/amf-custom-validator"
                        git tag v$VERSION
                        git push $URL v$VERSION
                '''
            }
        }
    }
}

