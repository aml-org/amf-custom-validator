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
        stage('Coverage') {
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target coverage'
                    registryCredentialsId 'dockerhub-pro-credentials'
                }
            }
            steps {
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sf-sonarqube-official', passwordVariable: 'SONAR_SERVER_TOKEN', usernameVariable: 'SONAR_SERVER_URL']]) {
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
        stage('Publish artifacts') {
            when {
                anyOf {
                    branch 'develop'
                    branch 'master'
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
                sh '''#!/bin/bash
                      cd /src
                      if [[ ${BRANCH_NAME} = "master" || ${BRANCH_NAME} =~ "release/*" ]]; then
                          IS_SNAPSHOT=false
                      else
                          IS_SNAPSHOT=true
                      fi

                      # Add safe directory exception
                      git config --global --add safe.directory /src

                      # Login
                      echo //registry.npmjs.org/:_authToken=${NPM_TOKEN} >> ~/.npmrc
                      echo @aml-org:registry=https://registry.npmjs.org/ >> ~/.npmrc
                      npm config set registry https://registry.npmjs.org/
                      npm whoami

                      # Publish
                      cd ./wrappers/js
                      if [ "$IS_SNAPSHOT" = true ]; then
                          npm-snapshot $BUILD_NUMBER
                      fi
                      VERSION=$(node -pe "require('./package.json').version")
                      npm publish --access public

                      cd ../js-web
                      if [ "$IS_SNAPSHOT" = true ]; then
                          npm-snapshot $BUILD_NUMBER
                      fi
                      npm publish --access public

                      if [ "$IS_SNAPSHOT" = true ]; then
                          npm dist-tag add @aml-org/amf-custom-validator-web@${VERSION} snapshot
                          npm dist-tag add @aml-org/amf-custom-validator@${VERSION} snapshot
                      else
                          npm dist-tag add @aml-org/amf-custom-validator-web@${VERSION} release
                          npm dist-tag add @aml-org/amf-custom-validator@${VERSION} release
                      fi

                      git tag v$VERSION
                      URL="https://${GITHUB_USR}:${GITHUB_PSW}@github.com/aml-org/amf-custom-validator"
                      git push $URL v$VERSION
                '''
            }
        }
        stage('Publish binaries') {
            when {
                branch 'master'
            }
            agent {
                dockerfile {
                    filename 'Dockerfile'
                    additionalBuildArgs  '--target ci-go'
                    registryCredentialsId 'dockerhub-pro-credentials'
                }
            }
            steps {
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'nexus', usernameVariable: 'NEXUS_USR', passwordVariable: 'NEXUS_PSW']]) {
                    script {
                        def platforms = [ "windows/amd64", "linux/amd64", "darwin/amd64", "darwin/arm64" ]
                        def packageVersion = sh (
                            script: 'grep version ./wrappers/js/package.json | sed \'s/.*"version": "\\(.*\\)".*/\\1/\'',
                            returnStdout: true
                        ).trim()
                        sh (script: 'mkdir -p /tmp/.cache')
                        platforms.each {item ->
                            def nexusRepository = "nexus.build.msap.io/nexus"
                            def repositoryName = "artifacts"
                            def groupId = "aml-org.amf-custom-validator"
                            def packageName = "acv"

                            def i = item.indexOf("/")
                            def goOs = item.substring(0, i)
                            def goArch = item.substring(i + 1, item.length())
                            def artifactId = "acv" + '-' + goOs + '-' + goArch
                            def fileName = packageName + '-' + goOs + '-' + goArch
                            sh (script: "env GOOS=${goOs} GOARCH=${goArch} GOCACHE=/tmp/.cache go build -o ${fileName} ./cmd/main.go", returnStdout: true)

                            println 'nexusRepository=' + nexusRepository
                            println 'repositoryName=' + repositoryName
                            println 'groupId=' + groupId
                            println 'packageVersion=' + packageVersion
                            println 'artifactId=' + artifactId
                            println 'fileName=' + fileName
                            nexusArtifactUploader(
                                nexusUrl: "${nexusRepository}",
                                protocol: 'https',
                                credentialsId: 'nexus',
                                nexusVersion: 'nexus2',
                                repository: "${repositoryName}",
                                groupId: "${groupId}",
                                version: "${packageVersion}",
                                artifacts: [[
                                    artifactId: "${artifactId}",
                                    file: "${fileName}",
                                    classifier: ''
                                ]]
                            )
                        }
                    }
                }
            }
        }
    }
}

